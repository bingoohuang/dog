package main

import (
	"embed"
	"fmt"
	"github.com/bingoohuang/dog"
	"github.com/bingoohuang/gg/pkg/ctl"
	"github.com/bingoohuang/gg/pkg/flagparse"
	"github.com/bingoohuang/golog"
	"runtime"
	"time"
)

func (Config) VersionInfo() string { return "dog v1.0.0 2021-07-14 16:50:32" }

func (c Config) Usage() string {
	return fmt.Sprintf(`Usage of dog:
  -filter value 命令包含，以!开头为不包含，可以多个值
  -kill string 发送信号，eg INT TERM KILL QUIT USR1 USR2 (default "INT")
  -max-mem value 允许最大内存 (默认 0B，不检查内存)
  -max-pcpu int 允许内存最大百分比, eg 1-%d (默认 %d), 0 不检查 CPU
  -max-pmem int 允许CPU最大百分比, eg 1-100 (默认 50)
  -pid int 指定pid
  -ppid int 指定ppid
  -self 是否监控自身
  -span duration 检查时间间隔 (默认 10s)
  -topn int 只取前N个检查
  -v Print version info and exit\n`, runtime.NumCPU()*100, runtime.NumCPU()*50)
}

type Config struct {
	Config string `flag:"c" usage:"yaml config filepath"`
	Init   bool   `usage:"init example dog.yml/ctl and then exit"`

	Topn    int
	Pid     int
	Ppid    int
	Self    bool
	Kill    string        `val:"INT"`
	Span    time.Duration `val:"10s"`
	MaxMem  uint64        `size:"true"`
	MaxPmem int           `val:"50"`
	MaxPcpu int
	Filter  []string

	Version bool `flag:"v" usage:"Print version info and exit"`
}

func (c *Config) PostProcess() {
	if c.MaxPcpu == 0 {
		c.MaxPcpu = runtime.NumCPU() * 50
	}
}

//go:embed initassets
var initAssets embed.FS

func main() {
	c := &Config{}
	flagparse.Parse(c, flagparse.AutoLoadYaml("c", "dog.yml"))
	ctl.Config{Initing: c.Init, InitFiles: initAssets}.ProcessInit()
	golog.SetupLogrus()

	watchConfig := dog.WatchConfig{
		Topn:       c.Topn,
		Pid:        c.Pid,
		Ppid:       c.Ppid,
		Self:       c.Self,
		BitSignals: c.Kill,
		Interval:   c.Span,
		MaxMem:     c.MaxMem,
		MaxPmem:    float32(c.MaxPmem),
		MaxPcpu:    float32(c.MaxPcpu),
		CmdFilter:  c.Filter,
	}

	d := dog.NewDog(dog.WithConfig(watchConfig))
	d.StartWatch()
}
