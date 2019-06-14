package dog

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MyListerner struct {
	barkType BiteFor
}

func (m *MyListerner) Biting(barkType BiteFor, threshold, real uint32) {
	m.barkType = barkType
}

func TestDog_FreeDog4Pid(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	dog := &Dog{
		MaxMem:   1,
		BiteLive: true,
	}
	dog.TimerDuration = 100 * time.Millisecond

	my := &MyListerner{}
	dog.ListenBiting(my)
	dog.FreeDog()

	time.Sleep(300 * time.Millisecond)

	assert.Equal(t, my.barkType, BiteForMaxMem)
}
