package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bingoohuang/cmd"
	"github.com/bingoohuang/dog"
	"github.com/bingoohuang/gou/str"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/docker/go-units"
	"github.com/mitchellh/go-homedir"
)

// App 表示应用
type App struct {
	Conf       Conf
	ListenAddr string      // 监听端口，为空表示不启动监听。例如:9900
	Gin        *gin.Engine // Gin引擎

	ProgramLock sync.Mutex
	Programs    map[string]Program

	PoolsLock sync.Mutex
	Pools     map[string]*Pool
}

type StateCode int

const (
	Succ StateCode = iota
	Timeout
	Error
)

type CompleteState struct {
	StateCode StateCode
	Data      string
}

type ProgramIn struct {
	In           string
	CompleteChan chan CompleteState
}

type Pool struct {
	program   Program
	inChan    chan ProgramIn
	duration  time.Duration
	maxMemKib int64
}

func (p *Pool) start() {
	for i := 0; i < p.program.PoolSize; i++ {
		go p.poolLoop()
	}
}

func (p *Pool) poolLoop() {
	var c *cmd.Cmd
	d := 10 * time.Second
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case in := <-p.inChan:
			if c == nil {
				c = p.createCmd()
			}

			c.Stdin <- in.In
			if p.waitResult(c, in) == ResultRestart {
				c = nil
			}
		case <-timer.C:
			timer.Reset(d)
			if c != nil {
				_ = c.Stop()
				c = nil
			}
		}

	}

}

type ResultState int

const (
	ResultOK   ResultState = iota
	ResultWait ResultState = iota
	ResultRestart
)

func (p *Pool) waitResult(c *cmd.Cmd, in ProgramIn) ResultState {
	timer := time.NewTimer(p.duration)
	defer timer.Stop()

	for {
		select {
		case out := <-c.Stdout:
			switch p.processOut(out, c, in) {
			case ResultOK:
				return ResultOK
			case ResultWait:
				continue
			case ResultRestart:
				return ResultRestart
			}
		case <-timer.C:
			p.restart(c, in)
			return ResultRestart
		}
	}
}

func (p *Pool) restart(c *cmd.Cmd, in ProgramIn) {
	logrus.Warnf("PID:%d, %s timeout %s in=%s kill and restart",
		c.Status().PID, p.program.Bash, p.program.Timeout, in.In)
	in.CompleteChan <- CompleteState{
		StateCode: Timeout,
		Data:      "timeout in " + p.duration.String(),
	}
	_ = c.Stop()
}

func (p *Pool) processOut(out string, c *cmd.Cmd, in ProgramIn) ResultState {
	if p.maxMemKib > 0 {
		ps := dog.Psaux(uint32(c.Status().PID))
		if int64(ps.RssKib) > p.maxMemKib {
			logrus.Warnf("PID:%d, %s reached maxMem %s, real %dKIB in=%s kill and restart",
				c.Status().PID, p.program.Bash, p.program.MaxMem, ps.RssKib, in.In)
			_ = c.Stop()
			in.CompleteChan <- CompleteState{StateCode: Error,
				Data: fmt.Sprintf("reached max memory, %s  > %s ",
					units.HumanSize(float64(ps.RssKib*1024)), p.program.MaxMem)}
			return ResultRestart
		}
	}

	ok := p.program.ExpectPrefix == "" || strings.HasPrefix(out, p.program.ExpectPrefix)
	if ok {
		in.CompleteChan <- CompleteState{StateCode: Succ, Data: out}

		logrus.Infof("PID:%d, %s processed in=%s, out=%s",
			c.Status().PID, p.program.Bash, in.In, out)
		return ResultOK
	}

	return ResultWait
}

func (p *Pool) createCmd() *cmd.Cmd {
	var err error
	if p.program.Timeout != "" {
		p.duration, err = time.ParseDuration(p.program.Timeout)
		if err != nil {
			logrus.Warnf("bad format for timeout %s, error %v", p.program.Timeout, err)
		}
	}
	if p.duration == 0 {
		p.duration = 10 * time.Second
	}

	if p.program.MaxMem != "" {
		maxMem, err := units.FromHumanSize(p.program.MaxMem)
		if err != nil {
			logrus.Warnf("bad format for maxMem %s, error %v", p.program.MaxMem, err)
		}

		p.maxMemKib = maxMem / 1024
	}

	cmdparts := strings.Fields(p.program.Bash)
	c := cmd.NewCmd(cmdparts...)
	c.Options(cmd.Stdin(), cmd.Streaming(), cmd.Buffered(false))
	c.Start()
	return c
}

// CreateAgApp 创建AgApp应用。
func CreateAgApp() *App {
	app := &App{
		ListenAddr: viper.GetString("addr"),
	}

	gin.SetMode(gin.ReleaseMode)
	app.Gin = gin.Default()

	return app
}

// GoStart 异步启动应用
func (a *App) GoStart() {
	conf, _ := homedir.Expand(viper.GetString("conf"))
	a.Conf = MustLoadConf(conf)

	a.Programs = make(map[string]Program)
	for key, p := range a.Conf.Programs {
		a.Programs[key] = p
	}

	a.Pools = make(map[string]*Pool)

	go a.setupRoutes()
}

func (a *App) setupRoutes() {
	r := a.Gin

	r.GET("/reg/:key", a.Register)
	r.GET("/run/:key", a.Exec)

	logrus.Infof("start to run at address %s", a.ListenAddr)
	if err := r.Run(a.ListenAddr); err != nil {
		logrus.Warnf("fail to start at %s, error %v", a.ListenAddr, err)
	}
}

func (a *App) Register(c *gin.Context) {
	a.ProgramLock.Lock()
	defer a.ProgramLock.Unlock()

	key := c.Param("key")
	a.Programs[key] = Program{
		Bash:     c.Query("bash"),
		MaxMem:   c.Query("maxMem"),
		Timeout:  c.Query("timeout"),
		PoolSize: str.ParseInt(c.Query("poolSize")),
	}
}

func (a *App) Exec(c *gin.Context) {
	key := c.Param("key")
	pg, ok := a.tryProgram(key, c)
	if !ok {
		return
	}

	in := c.Query("in")
	if pg.PoolSize == 0 {
		a.noPoolExec(pg, c, in)
		return
	}

	result := <-a.tryPool(key, pg, in)
	switch result.StateCode {
	case Succ:
		c.String(200, result.Data)
	default:
		c.String(500, result.Data)
	}
}

func (a *App) noPoolExec(pg Program, c *gin.Context, in string) {
	cmdparts := strings.Fields(pg.Bash)
	p := cmd.NewCmd(cmdparts...)
	p.Options(cmd.Stdin())
	chanStatuses := p.Start()

	p.Stdin <- in
	time.Sleep(100 * time.Millisecond)
	_ = p.Stop()
	status := <-chanStatuses
	errInfo := ""
	for _, stderr := range status.Stdout {
		if errInfo != "" {
			errInfo += "\n"
		}
		errInfo += stderr
	}
	if status.Error != nil {
		if errInfo != "" {
			errInfo += "\n"
		}
		errInfo += status.Error.Error()
	}
	for _, stderr := range status.Stderr {
		if errInfo != "" {
			errInfo += "\n"
		}
		errInfo += stderr
	}
	c.String(200, errInfo)
}

func (a *App) tryProgram(key string, c *gin.Context) (Program, bool) {
	a.ProgramLock.Lock()
	defer a.ProgramLock.Unlock()
	pg, ok := a.Programs[key]
	if !ok {
		c.Status(http.StatusNotFound)
		return Program{}, false
	}

	return pg, true
}

func (a *App) tryPool(key string, program Program, in string) chan CompleteState {
	a.PoolsLock.Lock()
	defer a.PoolsLock.Unlock()

	pool, ok := a.Pools[key]
	if !ok {
		pool = &Pool{
			program: program,
			inChan:  make(chan ProgramIn, program.PoolSize),
		}
		go pool.start()
		a.Pools[key] = pool
	}

	completeChan := make(chan CompleteState, 2)
	pool.inChan <- ProgramIn{
		In:           in,
		CompleteChan: completeChan,
	}

	return completeChan
}
