package dog

import (
	"fmt"
	"github.com/bingoohuang/gg/pkg/man"
	"github.com/bingoohuang/gou/str"
	"log"
	"regexp"
	"strconv"

	"github.com/bingoohuang/gg/pkg/ss"
	"github.com/gobars/cmd"
)

// PsAuxItem ...
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

func (p PsAuxItem) String() string {
	return fmt.Sprintf("User: %s Pid: %d Ppid: %d %%cpu: %s %%mem: %s VSZ: %s, RSS: %s Tty: %s Stat: %s Start: %s Time: %s Command: %s",
		p.User, p.Pid, p.Ppid,
		strconv.FormatFloat(float64(p.Pcpu), 'f', -1, 32),
		strconv.FormatFloat(float64(p.Pmem), 'f', -1, 32),
		man.Bytes(p.Vsz), man.Bytes(p.Rss), p.Tty, p.Stat, p.Start, p.Time, p.Command)
}

// PsAuxTop ...
func PsAuxTop(n, printN int, psFn func(topN int, heading bool) string) ([]PsAuxItem, error) {
	auxItems := make([]PsAuxItem, 0)
	re := regexp.MustCompile(`\s+`)
	i := 0
	_, status := cmd.BashLiner(psFn(n, false), func(line string) bool {
		f := re.Split(line, 13)
		item := PsAuxItem{
			User:    f[2],
			Pid:     ss.ParseInt(f[3]),
			Ppid:    ss.ParseInt(f[4]),
			Pcpu:    ss.ParseFloat32(f[5]),
			Pmem:    ss.ParseFloat32(f[6]),
			Vsz:     ss.ParseUint64(f[7]) * 1024,
			Rss:     ss.ParseUint64(f[8]) * 1024,
			Tty:     f[9],
			Stat:    f[10],
			Start:   f[0] + ` ` + f[1],
			Time:    f[11],
			Command: f[12],
		}

		if i < printN {
			log.Printf("%s", item)
		}
		i++

		auxItems = append(auxItems, item)
		return true
	})

	return auxItems, status.Error
}

// PasAuxPid ...
func PasAuxPid(topN, pid int, heading bool) string {
	return prefix + str.If(heading, "", noheading) + pidPostfix + fmt.Sprintf("%d", pid) + psAuxTopOpt(topN) + fixedLtime
}

const pidPostfix = ` -p `

// PasAuxShell ...
func PasAuxShell(topN int, heading bool) string {
	return prefix + str.If(heading, "", noheading) + psAuxTopOpt(topN) + fixedLtime
}

// PasMemAuxShell ...
func PasMemAuxShell(topN int, heading bool) string {
	return memPrefix + str.If(heading, "", noheading) + psAuxMemTopOpt(topN) + fixedLtime
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
