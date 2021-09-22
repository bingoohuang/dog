package dog

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/bingoohuang/gg/pkg/man"
	"github.com/bingoohuang/gg/pkg/ss"
	"github.com/gobars/cmd"
)

// CpuPercent 得到当前系统的CPU百分比使用率，值在0-100之间
func CpuPercent() (v float32) {
	c := cpuUsageCmd()
	_, _ = cmd.BashLiner(c, func(line string) bool {
		v = ss.ParseFloat32(line)
		return false
	})

	return
}

// PsAuxItem 是一条PS命令输出的各个列信息.
type PsAuxItem struct {
	User    string
	Pid     int
	Ppid    int
	Pcpu    float32
	Pmem    float32
	Vsz     uint64
	Rss     uint64
	Tty     string
	Stat    string
	Start   string
	Time    string
	Command string
}

// PsAuxRawItem 是一条PS命令输出的各个列原始信息.
type PsAuxRawItem struct {
	User    string
	Pid     string
	Ppid    string
	Pcpu    string
	Pmem    string
	Vsz     string
	Rss     string
	Tty     string
	Stat    string
	Start   string
	Time    string
	Command string
}

func (p PsAuxItem) String() string {
	return fmt.Sprintf("User: %s Pid: %d Ppid: %d %%cpu: %s %%mem: %s VSZ: %s, RSS: %s Tty: %s Stat: %s Start: %s Time: %s Command: %s",
		p.User, p.Pid, p.Ppid,
		strconv.FormatFloat(float64(p.Pcpu), 'f', -1, 32),
		strconv.FormatFloat(float64(p.Pmem), 'f', -1, 32),
		man.Bytes(p.Vsz), man.Bytes(p.Rss), p.Tty, p.Stat, p.Start, p.Time, p.Command)
}

// PsAuxTop ...
func PsAuxTop(n, printN int, psFn func(topN int, heading bool) string) ([]PsAuxRawItem, error) {
	i := 0
	auxItems := make([]PsAuxRawItem, 0)
	_, status := cmd.BashLiner(psFn(n, false), func(line string) bool {
		item := splitLine(line)
		if i < printN {
			log.Printf("%s", item)
		}
		i++

		auxItems = append(auxItems, item)
		return true
	})

	return auxItems, status.Error
}

func (p PsAuxRawItem) ToPsAuxItem() *PsAuxItem {
	return &PsAuxItem{
		User:    p.User,
		Pid:     ss.ParseInt(p.Pid),
		Ppid:    ss.ParseInt(p.Ppid),
		Pcpu:    ss.ParseFloat32(p.Pcpu),
		Pmem:    ss.ParseFloat32(p.Pmem),
		Vsz:     ss.ParseUint64(p.Vsz) * 1024,
		Rss:     ss.ParseUint64(p.Rss) * 1024,
		Tty:     p.Tty,
		Stat:    p.Stat,
		Start:   p.Start,
		Time:    p.Time,
		Command: p.Command,
	}
}

var Pid = os.Getpid()

// ExecPsAuxSelf execute ps command with specific self pid.
func ExecPsAuxSelf() (item *PsAuxRawItem, err error) {
	return ExecPsAuxPid(Pid)
}

var re = regexp.MustCompile(`\s+`)

// ExecPsAuxPid execute ps command with specific pid.
func ExecPsAuxPid(pid int) (item *PsAuxRawItem, err error) {
	shell := PsAuxPid(0, pid, false)
	_, status := cmd.BashLiner(shell, func(line string) bool {
		a := splitLine(line)
		item = &a
		return true
	})

	return item, status.Error
}

// PsAuxPid ...
func PsAuxPid(topN, pid int, heading bool) string {
	return prefix + ss.If(heading, "", noheading) + pidPostfix + fmt.Sprintf("%d", pid) + psAuxTopOpt(topN) + fixedLtime
}

// PasAuxShell ...
func PasAuxShell(topN int, heading bool) string {
	return prefix + ss.If(heading, "", noheading) + psAuxTopOpt(topN) + fixedLtime
}

// PasCpuAuxShell ...
func PasCpuAuxShell(topN int, heading bool) string {
	return cpuPrefix + ss.If(heading, "", noheading) + psAuxCpuTopOpt(topN) + fixedLtime
}

// PasMemAuxShell ...
func PasMemAuxShell(topN int, heading bool) string {
	return memPrefix + ss.If(heading, "", noheading) + psAuxMemTopOpt(topN) + fixedLtime
}

func splitLine(line string) PsAuxRawItem {
	f := re.Split(line, 13)
	return PsAuxRawItem{
		User:    f[2],
		Pid:     f[3],
		Ppid:    f[4],
		Pcpu:    f[5],
		Pmem:    f[6],
		Vsz:     f[7],
		Rss:     f[8],
		Tty:     f[9],
		Stat:    f[10],
		Start:   f[0] + ` ` + f[1],
		Time:    f[11],
		Command: f[12],
	}
}

/*
ps是linux系统的进程管理工具，般来说，ps aux命令执行结果的几个列的信息分别是：

USER 进程所属用户
PID 进程ID
%CPU 进程占用CPU百分比
%MEM 进程占用内存百分比
VSZ 虚拟内存占用大小
RSS 实际内存占用大小
TTY 终端类型
STAT 进程状态
START 进程启动时刻
TIME 进程运行时长
COMMAND 启动进程的命令

https://superuser.com/a/117921

USER = user owning the process
PID = process ID of the process
%CPU = It is the CPU time used divided by the time the process has been running.
%MEM = ratio of the process’s resident set size to the physical memory on the machine
VSZ = virtual memory usage of entire process (in KiB)
RSS = resident set size, the non-swapped physical memory that a task has used (in KiB)
TTY = controlling tty (terminal)
STAT = multi-character process state
START = starting time or date of the process
TIME = cumulative CPU time
COMMAND = command with all its arguments

https://alvinalexander.com/linux/unix-linux-process-memory-sort-ps-command-cpu/
vsz        VSZ      virtual memory size of the process in KiB (1024-byte units). Device mappings are currently excluded; this is subject to change.
rss        RSS      resident set size, the non-swapped physical memory that a task has used (in kiloBytes).
*/
