package dog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPsauxSelf(t *testing.T) {
	p := PsauxSelf()
	assert.Equal(t, os.Getpid(), p.Pid)
}
