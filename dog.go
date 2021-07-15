package dog

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/bingoohuang/gg/pkg/ss"
	"github.com/gobars/cmd"
	"log"
	"math/big"
	"os"
	"strings"
	"syscall"
	"time"
)

// Dog 表示 看门狗
type Dog struct {
	Config *WatchConfig
	stop   chan interface{}
}

func WithConfig(v WatchConfig) WatchOption   { return func(c *WatchConfig) { *c = v } }
func WithTopn(v int) WatchOption             { return func(c *WatchConfig) { c.Topn = v } }
func WithPid(v int) WatchOption              { return func(c *WatchConfig) { c.Pid = v } }
func WithPpid(v int) WatchOption             { return func(c *WatchConfig) { c.Ppid = v } }
func WithWatchSelf(v bool) WatchOption       { return func(c *WatchConfig) { c.Self = v } }
func WithKillSignals(v []string) WatchOption { return func(c *WatchConfig) { c.KillSignals = v } }
func WithMaxMem(v uint64) WatchOption        { return func(c *WatchConfig) { c.MaxMem = v } }
func WithMaxPmem(v float32) WatchOption      { return func(c *WatchConfig) { c.MaxPmem = v } }
func WithMaxPcpu(v float32) WatchOption      { return func(c *WatchConfig) { c.MaxPcpu = v } }
func WithCmdFilter(v ...string) WatchOption  { return func(c *WatchConfig) { c.CmdFilter = v } }

func NewDog(options ...WatchOption) *Dog {
	c := createWatchConfig(options)
	log.Printf("dog with config: %+v created", c)
	return &Dog{Config: c}
}

// BiteListener 咬人监听器
type BiteListener interface {
	Biting(barkType BiteFor, threshold, real uint32)
}

type WatchConfig struct {
	Topn        int
	Pid         int
	Ppid        int
	Self        bool
	KillSignals []string

	Interval   time.Duration
	MaxMem     uint64  // 看住最大内存使用
	MaxPmem    float32 // 看住最大内存占用比例
	MaxPcpu    float32 // 看住最大CPU占用比例
	CmdFilter  []string
	LogItems   []string
	RateConfig *RateConfig
	limiter    *Limiter
	Jitter     time.Duration
}

type WatchOption func(*WatchConfig)

func (d *Dog) Stop() {
	d.stop <- struct{}{}
}

// StartWatch 开始放狗看门.
func (d *Dog) StartWatch() {
	ticker := time.NewTicker(d.Config.Interval)
	defer ticker.Stop()

	d.watch()

	for {
		RandomSleep(d.Config.Jitter)

		select {
		case <-d.stop:
			return
		case <-ticker.C:
			d.watch()
		}
	}
}

// BiteFor 咬人原因
type BiteFor int

const (
	BiteForNone    BiteFor = iota // 不咬
	BiteForMaxMem                 // 超过最大内存咬人
	BiteForMaxPmem                // 超过最大内存占比咬人
	BiteForMaxPcpu                // 超过最大CPU占比咬人
)

func (d *Dog) watch() {
	c := d.Config
	items, err := PsAuxTop(c.Topn, 0)
	if err != nil {
		log.Printf("ps aux error: %v", err)
		return
	}

	pid := os.Getpid()

	for _, item := range items {
		if d.Filter(item) {
			continue
		}
		if c.Pid > 0 && item.Pid != c.Pid || c.Ppid > 0 && item.Ppid != c.Ppid || c.Self && c.Pid != pid {
			continue
		}
		if !c.Self && c.Pid == pid { // 不看自己，跳过自己
			continue
		}

		biteFor := BiteForNone
		switch {
		case c.MaxMem > 0 && item.Rss > c.MaxMem:
			biteFor = BiteForMaxMem
		case c.MaxPmem > 0 && item.Pmem > c.MaxPmem:
			biteFor = BiteForMaxPmem
		case c.MaxPcpu > 0 && item.Pcpu > c.MaxPcpu:
			biteFor = BiteForMaxPcpu
		}
		if biteFor != BiteForNone {
			d.bite(biteFor, item)
		}
	}
}

// Ctrl+C - SIGINT
// Ctrl+\ - SIGQUIT
// Ctrl+Z - SIGTSTP
var signalMap = map[string]syscall.Signal{
	"INT":  syscall.SIGINT,
	"TERM": syscall.SIGTERM,
	"QUIT": syscall.SIGQUIT,
	"KILL": syscall.SIGKILL,
	"USR1": syscall.SIGUSR1,
	"USR2": syscall.SIGUSR2,
}

func (d *Dog) bite(biteFor BiteFor, item PsAuxItem) {
	c := d.Config
	if c.limiter != nil && c.limiter.Allow() {
		log.Printf("Dog barking for %v, config:%s, item %+v", biteFor, c.RateConfig, item)
		return
	}

	log.Printf("Dog biting for %v, item %+v", biteFor, item)
	for _, v := range c.LogItems {
		if f, ok := logItemsRegister[v]; ok {
			if m := f(item); m != "" {
				log.Printf("LogItem: %s, Value: %s", v, m)
			}
		}
	}

	for _, s := range c.KillSignals {
		if v, ok := signalMap[s]; ok {
			if err := syscall.Kill(item.Pid, v); err != nil {
				log.Printf("E! Kill %s to %d, err: %v", v, item.Pid, err)
			} else {
				log.Printf("Kill %s to %d succeeded", v, item.Pid)
			}
		}
	}
}

func (d *Dog) Filter(item PsAuxItem) bool {
	for _, cf := range d.Config.CmdFilter {
		if strings.HasPrefix(cf, "!") {
			if ss.ContainsFold(item.Command, cf[1:]) {
				return true // 配置不能包含，但是包含，过滤掉
			}
		} else {
			if !ss.ContainsFold(item.Command, cf) {
				return true // 配置包含，但是不包含，过滤掉
			}
		}
	}

	return false
}

func createWatchConfig(options []WatchOption) *WatchConfig {
	c := &WatchConfig{}
	for _, option := range options {
		option(c)
	}
	if c.Interval == 0 {
		c.Interval = 10 * time.Second
	}
	if len(c.KillSignals) == 0 {
		c.KillSignals = []string{"INT"}
	}

	if c.RateConfig != nil {
		c.limiter, _ = NewLimiter(c.RateConfig.Duration, c.RateConfig.Times, func() (Window, StopFunc) {
			// NewLocalWindow returns an empty stop function, so it's
			// unnecessary to call it later.
			return NewLocalWindow()
		})
	}

	return c
}

var logItemsRegister = map[string]func(PsAuxItem) string{
	"CWD": func(item PsAuxItem) (l string) {
		script := fmt.Sprintf(`lsof -p %d | grep cwd | awk '{print $9}'`, item.Pid)
		cmd.BashLiner(script, func(line string) bool { l = line; return false })
		return
	},
	"ENV": func(item PsAuxItem) (l string) {
		script := fmt.Sprintf(`ps e -ww -p %d | tail -1`, item.Pid)
		cmd.BashLiner(script, func(line string) bool { l = line; return false })
		return
	},
}

type RateConfig struct {
	Times    int64
	Duration time.Duration
}

func (r RateConfig) String() string { return fmt.Sprintf("%d/%s", r.Times, r.Duration) }

var ErrBadRateConfig = errors.New("bad format for rate config, eg 10/30s")

func ParseRateConfig(expr string) (*RateConfig, error) {
	if expr == "" {
		return nil, nil
	}

	pos := strings.Index(expr, "/")
	if pos < 0 {
		return nil, ErrBadRateConfig
	}

	timesPart := expr[:pos]
	times := ss.ParseInt64(timesPart)
	if times <= 0 {
		return nil, ErrBadRateConfig
	}

	durationPart := expr[pos+1:]
	duration, err := time.ParseDuration(durationPart)
	if err != nil {
		return nil, ErrBadRateConfig
	}

	return &RateConfig{Times: times, Duration: duration}, nil
}

// RandomSleep will sleep for a random amount of time up to max.
// If the shutdown channel is closed, it will return before it has finished
// sleeping.
func RandomSleep(max time.Duration) {
	if max == 0 {
		return
	}

	var sleepns int64
	maxSleep := big.NewInt(max.Nanoseconds())
	if j, err := rand.Int(rand.Reader, maxSleep); err == nil {
		sleepns = j.Int64()
	}

	t := time.NewTimer(time.Nanosecond * time.Duration(sleepns))
	<-t.C
}
