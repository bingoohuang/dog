package dog

import (
	"strconv"

	"github.com/bingoohuang/gg/pkg/ss"
)

func cpuUsageCmd() string { return `top -l 1 -n 0 -F|egrep -o '\S+ idle'|awk '{print 100-$1}'` }

func psAuxMemTopOpt(n int) string { return ss.If(n > 0, `|head -n `+strconv.Itoa(n), ``) }
func psAuxCpuTopOpt(n int) string { return ss.If(n > 0, `|head -n `+strconv.Itoa(n), ``) }
func psAuxTopOpt(n int) string    { return ss.If(n > 0, `|head -n `+strconv.Itoa(n), ``) }

const (
	memPrefix  = `(export TZ=UTC0 LC_ALL=C; ps xo lstart,user,pid,ppid,pcpu,pmem,vsz,rss,tt,stat,time,args -m`
	cpuPrefix  = `(export TZ=UTC0 LC_ALL=C; ps xo lstart,user,pid,ppid,pcpu,pmem,vsz,rss,tt,stat,time,args -r`
	prefix     = `(export TZ=UTC0 LC_ALL=C; ps xo lstart,user,pid,ppid,pcpu,pmem,vsz,rss,tt,stat,time,args`
	noheading  = `|tail -n +2`
	pidPostfix = ` -p `
)

// nolint
const fixedLtime = `|awk '{c="date -jf \"%a %b %e %T %Y\" \""$1 FS $2 FS $3 FS $4 FS $5"\" +\047%Y-%m-%d %H:%M:%S\047"; c|getline d; close(c); $1=$2=$3=$4=$5=""; printf "%s\n",d$0 }' )`

// nolint
/*
https://unix.stackexchange.com/questions/115736/unable-to-use-date-command-to-change-specific-date-format-in-bash-shell-on-os-x
➜  sysinfo git:(master) ✗ (export TZ=UTC0 LC_ALL=C; ps axo lstart,user,pid,ppid,pcpu,pmem,vsz,rss,tt,stat,time,args |  head -n 3 )
STARTED                      USER               PID  PPID  %cpu %MEM      VSZ    RSS   TT  STAT      TIME ARGS
Mon Jul 22 01:13:22 2019     root                 1     0   0.0  0.1  4417876  13896   ??  Ss     8:44.17 /sbin/launchd
Mon Jul 22 01:13:30 2019     root                40     1   0.0  0.0  4395956   1448   ??  Ss     0:05.95 /usr/sbin/syslogd
➜  sysinfo git:(master) ✗ (export TZ=UTC0 LC_ALL=C;date -jf "%a %b %e %T %Y" "Mon Jul 22 01:13:22 2019" +"%Y-%m-%d %H:%M:%S")
2019-07-22 01:13:22
➜  sysinfo git:(master) ✗ (export TZ=UTC0 LC_ALL=C; ps axo lstart,user,pid,ppid,pcpu,pmem,vsz,rss,tt,stat,time,args | tail -n +2|head -n 1|awk '{c="date -jf \"%a %b %e %T %Y\" \""$1 FS $2 FS $3 FS $4 FS $5"\" +\047%Y-%m-%d %H:%M:%S\047"; c|getline d; close(c); $1=$2=$3=$4=$5=""; printf "%s\n",d$0 }' )
2019-07-22 01:13:22     root 1 0 0.1 0.1 4418400 13912 ?? Ss 8:56.17 /sbin/launchd
*/
