package dog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type PsauxOut struct {
	Pid  int
	Pcpu uint32
	Pmem uint32
	Vsz  uint32
	Rss  uint32
	Line string
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
		Pid:  ParseInt(n[1]),
		Pcpu: uint32(ParseFloat32(n[2])),
		Pmem: uint32(ParseFloat32(n[3])),
		Vsz:  ParseKib(n[4]),
		Rss:  ParseKib(n[5]),
		Line: output,
	}
}

func ParseKib(s string) uint32 {
	return uint32(ParseFloat32(s) * 1024)
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
