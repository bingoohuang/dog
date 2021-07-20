package dog

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/bingoohuang/gg/pkg/man"
	"github.com/bingoohuang/gg/pkg/ss"
	"github.com/bingoohuang/gg/pkg/timex"
	"github.com/gobars/cmd"
	"github.com/shirou/gopsutil/mem"
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

	Interval      time.Duration
	MaxMem        uint64  // 看住最大内存使用
	MaxPmem       float32 // 看住最大内存占用比例
	MaxPcpu       float32 // 看住最大CPU占用比例
	CmdFilter     []string
	MinFreeMemory uint64
	Whites        []string
	LogItems      []string
	RateConfig    *RateConfig
	limiter       *Limiter
	Jitter        time.Duration
	MaxTime       time.Duration
	MaxTimeEnv    string
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
	BiteForMaxTime                // 超过最大运行时长咬人
	BiteForTopMem                 // 驱逐内存占用第一进程
)

func (b BiteFor) String() string {
	switch b {
	case BiteForMaxMem:
		return "内存超了"
	case BiteForMaxPmem:
		return "内存占比超了"
	case BiteForMaxPcpu:
		return "CPU占比超了"
	case BiteForMaxTime:
		return "运行时长超了"
	case BiteForTopMem:
		return "内存占用第一"
	}

	return "啥都没超"
}

var pid = os.Getpid()

func (d *Dog) watch() {
	c := d.Config

	if c.MinFreeMemory > 0 {
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("get VirtualMemory error: %v", err)
		} else if vmStat.Free < c.MinFreeMemory {
			d.biteTopMem(vmStat)
		}
	}

	items, err := PsAuxTop(c.Topn, 0, PasAuxShell)
	if err != nil {
		log.Printf("ps aux error: %v", err)
		return
	}

	for _, v := range items {
		if d.Filter(v) {
			continue
		}
		if c.Pid > 0 && v.Pid != c.Pid || c.Ppid > 0 && v.Ppid != c.Ppid || c.Self && c.Pid != pid {
			continue
		}
		if !c.Self && v.Pid == pid { // 不看自己，跳过自己
			continue
		}

		biteFor := BiteForNone
		switch {
		case c.MaxMem > 0 && v.Rss > c.MaxMem:
			biteFor = BiteForMaxMem
		case c.MaxPmem > 0 && v.Pmem > c.MaxPmem:
			biteFor = BiteForMaxPmem
		case c.MaxPcpu > 0 && v.Pcpu > c.MaxPcpu:
			biteFor = BiteForMaxPcpu
		case c.MaxTime > 0 && exceedMaxTime(v, `yyyy-MM-dd HH:mm:ss`, c.MaxTime, c.MaxTimeEnv):
			biteFor = BiteForMaxTime
		}
		if biteFor != BiteForNone {
			d.bite(biteFor, v)
		}
	}
}

func exceedMaxTime(item PsAuxItem, layout string, maxTime time.Duration, env string) bool {
	t, err := time.ParseInLocation(timex.ConvertFormat(layout), item.Start, time.Local)
	if err != nil {
		log.Printf("E! failed to parse start time: %v", err)
		return false
	}

	if time.Since(t) < maxTime {
		return false
	}

	return env == "" || strings.Contains(itemEnv(item), env)
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

const TopMemFakePid = -100

func (d *Dog) biteTopMem(vm *mem.VirtualMemoryStat) {
	c := d.Config
	if c.limiter != nil && c.limiter.Allow(TopMemFakePid) {
		log.Printf("Dog barking for low free memory: %s/%s < config min: %s",
			man.Bytes(vm.Free), man.Bytes(vm.Total), man.Bytes(c.MinFreeMemory))
		return
	}

	log.Printf("Dog biting for low free memory: %s/%s < config min: %s",
		man.Bytes(vm.Free), man.Bytes(vm.Total), man.Bytes(c.MinFreeMemory))

	items, err := PsAuxTop(10, 0, PasMemAuxShell)
	if err != nil {
		log.Printf("ps aux error: %v", err)
		return
	}

	for _, v := range items {
		if !d.Whites(v) {
			d.bite(BiteForTopMem, v)
			return
		}
	}

	log.Printf("Dog no biting found for low free memory: %s/%s < config min: %s",
		man.Bytes(vm.Free), man.Bytes(vm.Total), man.Bytes(c.MinFreeMemory))
}

func (d *Dog) bite(biteFor BiteFor, v PsAuxItem) {
	c := d.Config
	if c.limiter != nil && c.limiter.Allow(v.Pid) {
		log.Printf("Dog barking for %s, config:%s, item %+v", biteFor, c.RateConfig, v)
		return
	}

	log.Printf("Dog biting for %s, item %+v", biteFor, v)
	for _, l := range c.LogItems {
		if f, ok := logItemsRegister[l]; ok {
			if m := f(v); m != "" {
				log.Printf("LogItem: %s, Value: %s", l, m)
			}
		}
	}

	for _, s := range c.KillSignals {
		if g, ok := signalMap[s]; ok {
			if err := syscall.Kill(v.Pid, g); err != nil {
				log.Printf("E! Kill %s to %d, err: %v", g, v.Pid, err)
			} else {
				log.Printf("Kill %s to %d succeeded", g, v.Pid)
			}
		}
	}
}

func (d *Dog) Whites(item PsAuxItem) bool {
	for _, cf := range d.Config.Whites {
		if ss.ContainsFold(item.Command, cf) {
			return true // 过滤
		}
	}

	return false
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
		c.limiter = c.RateConfig.NewLimiter()
	}

	return c
}

func itemCwd(item PsAuxItem) (l string) {
	script := fmt.Sprintf(`lsof -p %d | grep cwd | awk '{print $9}'`, item.Pid)
	cmd.BashLiner(script, func(line string) bool { l = line; return false })
	return
}

func itemEnv(item PsAuxItem) (l string) {
	script := fmt.Sprintf(`ps e -ww -p %d | tail -1`, item.Pid)
	cmd.BashLiner(script, func(line string) bool { l = line; return false })
	return
}

var logItemsRegister = map[string]func(PsAuxItem) string{
	"CWD": itemCwd,
	"ENV": itemEnv,
}

type RateConfig struct {
	Times    int
	Duration time.Duration
}

func (r RateConfig) String() string { return fmt.Sprintf("%d/%s", r.Times, r.Duration) }

func (r RateConfig) NewLimiter() *Limiter {
	return &Limiter{RateConfig: r, Series: map[int]*[]time.Time{}}
}

type Limiter struct {
	RateConfig
	Series map[int]*[]time.Time
}

func (l *Limiter) Allow(key int) bool {
	now := time.Now()

	for k, v := range l.Series {
		for i, t := range *v {
			if now.Sub(t) < l.Duration {
				n := copy(*v, (*v)[i:])
				*v = (*v)[:n]
				break
			}
		}
		if len(*v) == 0 {
			delete(l.Series, k)
		}
	}

	if ts, ok := l.Series[key]; ok {
		*ts = append(*ts, now)
		if len(*ts) >= l.Times {
			delete(l.Series, key)
			return false
		}
	} else {
		ts := []time.Time{now}
		l.Series[key] = &ts
	}

	return true
}

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
	times := ss.ParseInt(timesPart)
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
