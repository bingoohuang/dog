package dog

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPsauxSelf(t *testing.T) {
	p := PsauxSelf()
	assert.Equal(t, os.Getpid(), p.Pid)
}
