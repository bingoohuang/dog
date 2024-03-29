package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/bingoohuang/gg/pkg/v"

	"github.com/bingoohuang/dog"
	"github.com/bingoohuang/gg/pkg/flagparse"
	"github.com/bingoohuang/gg/pkg/ss"
	"github.com/bingoohuang/golog"
)

func (*Config) VersionInfo() string { return v.Version() }

func (c *Config) Usage() string {
	return fmt.Sprintf(`Usage of %s:
  -filter value 命令包含，以!开头为不包含，可以多个值
  -cond string 发送条件，默认触发1次就发信号，eg.3/30s，在30s内发生3次，则触发 
  -kill string 发送信号，多个逗号分隔，eg. INT,TERM,KILL,QUIT,USR1,USR2 (默认 INT)
  -log  string 记录日志信息，多个逗号分隔，eg. ENV,CWD
  -max-time value 允许最大启动时长 (默认 0，不检查启动时长)
  -max-time-env value 允许最大启动时长包含的环境变量
  -max-mem value 允许最大内存 (默认 0B，不检查内存)
  -max-pcpu int 允许CPU最大百分比, eg. 1-%d (默认 %d), 0 不查 CPU
  -max-pmem int 允许内存最大百分比, eg. 1-100 (默认 50)
  -min-available-memory 允许最小总可用内存 (默认 0B，不检查此项)
  -max-host-cpu 允许最大机器CPU百分比（0-100） (默认 0，不检查此项)
  -whites value 总最小内存/最大机器CPU百分比触发时，驱逐进程命令行包含白名单，可以多个值
  -pid int 指定pid
  -ppid int 指定ppid
  -self 是否监控自身
  -span duration 检查时间间隔 (默认 10s)
  -jitter duration 最大抖动 (默认 1s)
  -topn int 只取前N个检查
  -v Print version info and exit`,
		os.Args[0], runtime.NumCPU()*100, runtime.NumCPU()*50)
}

type Config struct {
	Config string `flag:"c" usage:"yaml config filepath"`
	Init   bool   `usage:"init example dogwatch.yml/ctl and then exit"`

	Topn       int
	Pid        int
	Ppid       int
	Self       bool
	Kill       string `val:"INT"`
	Log        string
	Cond       string
	Span       time.Duration `val:"10s"`
	Jitter     time.Duration `val:"1s"`
	MaxTime    time.Duration
	MaxTimeEnv string
	MaxMem     uint64 `size:"true" yaml:",label=size"`
	MaxPmem    int    `val:"50"`
	MaxPcpu    int
	Filter     []string

	MinAvailableMemory uint64 `size:"true" yaml:",label=size"` // 最小整个机器可用内存阈值
	MaxHostCpu         int
	Whites             []string // 驱逐白名单

	Version    bool `flag:"v" usage:"Print version info and exit"`
	rateConfig *dog.RateConfig
}

func (c *Config) PostProcess() {
	if c.MaxPcpu == 0 {
		c.MaxPcpu = runtime.NumCPU() * 50
	}

	var err error
	if c.rateConfig, err = dog.ParseRateConfig(c.Cond); err != nil {
		log.Fatalf("ParseRateConfig error: %v", err)
	}
}

//go:embed initassets
var initAssets embed.FS

func main() {
	c := &Config{}
	flagparse.Parse(c,
		flagparse.AutoLoadYaml("c", "dogwatch.yml"),
		flagparse.ProcessInit(&initAssets))
	golog.Setup()
	splitter := func(s string) []string {
		return ss.Split(s, ss.WithCase(ss.CaseUpper), ss.WithIgnoreEmpty(true), ss.WithTrimSpace(true))
	}
	watchConfig := &dog.WatchConfig{
		Topn:               c.Topn,
		Pid:                c.Pid,
		Ppid:               c.Ppid,
		Self:               c.Self,
		KillSignals:        splitter(c.Kill),
		LogItems:           splitter(c.Log),
		Interval:           c.Span,
		Jitter:             c.Jitter,
		MaxTime:            c.MaxTime,
		MaxTimeEnv:         c.MaxTimeEnv,
		MaxMem:             c.MaxMem,
		MaxPmem:            float32(c.MaxPmem),
		MaxPcpu:            float32(c.MaxPcpu),
		CmdFilter:          c.Filter,
		MinAvailableMemory: c.MinAvailableMemory,
		MaxHostPcpu:        float32(c.MaxHostCpu),
		Whites:             c.Whites,
		RateConfig:         c.rateConfig,
	}

	d := dog.NewDog(dog.WithConfig(watchConfig))
	d.StartWatch()
}
