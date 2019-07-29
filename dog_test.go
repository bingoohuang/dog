package dog

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
		MaxMemKib: 100,
		BiteLive:  true,
	}
	dog.TimerDuration = 100 * time.Millisecond

	my := &MyListerner{}
	dog.ListenBiting(my)
	dog.FreeDog()

	time.Sleep(300 * time.Millisecond)

	assert.Equal(t, my.barkType, BiteForMaxMem)
}

func Test4Demo(t *testing.T) {
	dog := &Dog{
		MaxMemKib:     100 * 1024, // 最大内存 100M
		MaxCPUPercent: 90,         // 最大CPU占比90%
		MaxMemPercent: 90,         // 最大内存占比90%
		BiteLive:      false,      // 咬了不活，直接死掉
	}

	dog.FreeDog() // 放狗看门（不阻塞)

	dog.CageDog() // 收狗不看门

}
