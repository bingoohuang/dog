package dog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bingoohuang/gou/str"
	"github.com/gobars/cmd"
	"github.com/sirupsen/logrus"
)

type PsauxOut struct {
	Pid    int
	Pcpu   uint32
	Pmem   uint32
	VszKib uint32
	RssKib uint32
	Line   string
}

func PsauxSelf() PsauxOut {
	return Psaux(uint32(os.Getpid()))
}

func Psaux(pid uint32) PsauxOut {
	ps := fmt.Sprintf(`ps aux |awk '$2 == %d { print $0 }'`, pid)

	_, stat := cmd.Bash(ps, cmd.Timeout(1*time.Second))
	if stat.Error != nil {
		logrus.Errorf("exec %s error %v", ps, stat)
		return PsauxOut{}
	}

	output := stat.Stdout[0]
	n := strings.Fields(output)
	return PsauxOut{
		Pid:    str.ParseInt(n[1]),
		Pcpu:   str.ParseUint32(n[2]),
		Pmem:   str.ParseUint32(n[3]),
		VszKib: str.ParseUint32(n[4]),
		RssKib: str.ParseUint32(n[5]),
		Line:   output,
	}
}
