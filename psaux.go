package dog

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

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
	return Psaux(os.Getpid())
}

func Psaux(pid int) PsauxOut {
	ps := fmt.Sprintf(`ps aux |awk '$2 == %d { print $0 }'`, pid)

	output, err := ExecShell(ps)
	if err != nil {
		logrus.Errorf("exec %s error %v", ps, err)
		return PsauxOut{}
	}

	n := strings.Fields(output)
	return PsauxOut{
		Pid:    ParseInt(n[1]),
		Pcpu:   uint32(ParseFloat32(n[2])),
		Pmem:   uint32(ParseFloat32(n[3])),
		VszKib: ParseUint32(n[4]),
		RssKib: ParseUint32(n[5]),
		Line:   output,
	}
}

func ParseUint32(s string) uint32 {
	return uint32(ParseInt(s))
}

func ParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func ParseFloat32(s string) float32 {
	v, _ := strconv.ParseFloat(s, 32)
	return float32(v)
}

func ExecShell(s string) (string, error) {
	cmd := exec.Command("bash", "-c", s)
	output, err := cmd.CombinedOutput()

	return string(output), err
}
