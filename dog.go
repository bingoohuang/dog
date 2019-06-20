package dog

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Dog 表示 看门狗
type Dog struct {
	MaxMemKib     uint32 // 看住最大内存使用
	MaxMemPercent uint32 // 看住最大内存占用比例(1-99)
	MaxCpuPercent uint32 // 看住最大CPU占用比例(1-99)
	BiteLive      bool   // 是否咬了不死

	TimerDuration time.Duration // 狗巡视周期

	paused bool
	free   bool
	cmd    chan CmdType

	biteListeners []BiteListener
	pid           int
}

// CmdType 命令类型
type CmdType int

const (
	CmdNoop   CmdType = iota
	CmdCaging  // 收狗进狗笼
)

// BiteListener 咬人监听器
type BiteListener interface {
	Biting(barkType BiteFor, threshold, real uint32)
}

// ListenBiting 监听狗咬事件
func (d *Dog) ListenBiting(l BiteListener) {
	d.biteListeners = append(d.biteListeners, l)
}

// SetBite4Dead 设置是否直接咬死
func (d *Dog) SetBite4Dead(bite4Dead bool) {
	d.BiteLive = bite4Dead
}

// CageDog 收狗进狗笼
func (d *Dog) CageDog() {
	d.cmd <- CmdCaging
}

// FreeDog 开始放狗看门
func (d *Dog) FreeDog() {
	d.FreeDog4Pid(os.Getpid())
}

// FreeDog4Pid 开始放狗看门
func (d *Dog) FreeDog4Pid(pid int) {
	if d.free {
		return
	}

	d.pid = pid
	d.free = true
	d.paused = false
	d.cmd = make(chan CmdType)
	go d.watching()
}

// PauseWatching 暂停看门
func (d *Dog) PauseWatching() {
	d.paused = true
}

// ResumeWatching 继续看门
func (d *Dog) ResumeWatching() {
	d.paused = false
}

func (d *Dog) watching() {
	if d.TimerDuration == 0 {
		d.TimerDuration = 60 * time.Second
	}
	timer := time.NewTimer(d.TimerDuration)
	defer timer.Stop()

	for {
		select {
		case cmd := <-d.cmd:
			switch cmd {
			case CmdCaging:
				d.free = false
				return
			case CmdNoop:
				// noop!
			}
		case <-timer.C:
			timer.Reset(d.TimerDuration)
			if !d.paused {
				d.watch()
			}
		}
	}
}

// BiteFor 咬人原因
type BiteFor int

const (
	BiteForMaxMem        BiteFor = iota + 1 // 超过最大内存咬人
	BiteForMaxMemPercent                    // 超过最大内存占比咬人
	BiteForMaxCpuPercent                    // 超过最大CPU占比咬人
)

func (d *Dog) watch() {
	s := Psaux(d.pid)
	logrus.Debugf("dog is watching %+v", s)

	if d.MaxMemKib > 0 && s.RssKib > d.MaxMemKib {
		d.bite(BiteForMaxMem, d.MaxMemKib, s.RssKib)
	}

	if d.MaxCpuPercent > 0 && s.Pcpu > d.MaxCpuPercent {
		d.bite(BiteForMaxCpuPercent, d.MaxCpuPercent, s.Pcpu)
	}

	if d.MaxMemPercent > 0 && s.Pmem > d.MaxMemPercent {
		d.bite(BiteForMaxMemPercent, d.MaxMemPercent, s.Pmem)
	}
}

func (d *Dog) bite(biteFor BiteFor, threshold, real uint32) {
	logrus.Warnf("Dog biting for %v, threshold %v, real %v", biteFor, threshold, real)

	for _, l := range d.biteListeners {
		l.Biting(biteFor, threshold, real)
	}

	if !d.BiteLive {
		logrus.Panicf("Dog biting for %v, threshold %v, real %v", biteFor, threshold, real)
	}
}
