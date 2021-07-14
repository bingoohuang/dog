package dog

import (
	"github.com/bingoohuang/gg/pkg/ss"
	"log"
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

func WithConfig(v WatchConfig) WatchOption  { return func(c *WatchConfig) { *c = v } }
func WithTopn(v int) WatchOption            { return func(c *WatchConfig) { c.Topn = v } }
func WithPid(v int) WatchOption             { return func(c *WatchConfig) { c.Pid = v } }
func WithPpid(v int) WatchOption            { return func(c *WatchConfig) { c.Ppid = v } }
func WithWatchSelf(v bool) WatchOption      { return func(c *WatchConfig) { c.Self = v } }
func WithBitSignals(v string) WatchOption   { return func(c *WatchConfig) { c.BitSignals = v } }
func WithMaxMem(v uint64) WatchOption       { return func(c *WatchConfig) { c.MaxMem = v } }
func WithMaxPmem(v float32) WatchOption     { return func(c *WatchConfig) { c.MaxPmem = v } }
func WithMaxPcpu(v float32) WatchOption     { return func(c *WatchConfig) { c.MaxPcpu = v } }
func WithCmdFilter(v ...string) WatchOption { return func(c *WatchConfig) { c.CmdFilter = v } }

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
	Topn       int
	Pid        int
	Ppid       int
	Self       bool
	BitSignals string

	Interval  time.Duration
	MaxMem    uint64  // 看住最大内存使用
	MaxPmem   float32 // 看住最大内存占用比例
	MaxPcpu   float32 // 看住最大CPU占用比例
	CmdFilter []string
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
	log.Printf("Dog biting for %v, item %+v", biteFor, item)

	for k, v := range signalMap {
		if ss.ContainsFold(d.Config.BitSignals, k) {
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
	c := &WatchConfig{Interval: 10 * time.Second, BitSignals: "INT"}
	for _, option := range options {
		option(c)
	}
	return c
}
