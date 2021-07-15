package main

import (
	"fmt"
	"github.com/bingoohuang/dog"
	flag "github.com/bingoohuang/gg/pkg/fla9"
	"github.com/bingoohuang/gg/pkg/man"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// https://github.com/vikyd/go-cpu-load
func main() {
	var cores int
	var percentage int
	var duration time.Duration
	var memory string

	numCPU := runtime.NumCPU()
	flag.IntVar(&cores, "c", numCPU, "")
	flag.IntVar(&percentage, "p", 100, "")
	flag.StringVar(&memory, "m", "", "")
	flag.DurationVar(&duration, "d", 0, "")
	flag.Usage = func() {
		fmt.Printf(`Usage of busy:
  -c int 使用核数，默认 %d
  -d duration 跑多久，默认一直跑
  -m string 总内存,增量, eg. 1) 10M 直接达到10M 2) 10M,1K/10s 总用量10M,每10秒增加1K
  -p int 每核CPU百分比 (默认 100)
`, runtime.NumCPU())
	}
	flag.Parse()

	if cores < 1 || cores > numCPU {
		log.Fatalf("cores %d is not between 1 - %d", cores, numCPU)
	}
	if percentage > 100 {
		log.Fatalf("percentage %d is invalid, should be between 0 and 100", percentage)
	}

	log.Printf("busy starting, pid %d", os.Getpid())

	sig := make(chan os.Signal, 1)
	// syscall.SIGINT: ctl + c, syscall.SIGTERM: kill pid
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		s := <-sig
		log.Printf("received signal %s, exiting", s)
		os.Exit(-1)
	}()

	if memory != "" {
		go controlMem(memory)
	}

	if percentage > 0 {
		log.Printf(" run %d%% of %d/%d CPU cores %s.", percentage, cores, numCPU, printDuration(duration))
		RunCPULoad(cores, duration, percentage)
	} else {
		select {}
	}
}

func controlMem(memory string) {
	total := memory
	incr := ""
	if p := strings.Index(memory, ","); p > 0 {
		total = memory[:p]
		incr = memory[p+1:]
	}

	var per time.Duration
	if incr != "" {
		if p := strings.Index(incr, "/"); p > 0 {
			per, _ = time.ParseDuration(incr[p+1:])
			incr = incr[:p]
		}
	}

	totalMem, _ := man.ParseBytes(total)
	incrMem, _ := man.ParseBytes(incr)
	pid := os.Getpid()

	for {
		item := findItem(pid)
		if item.Rss >= totalMem {
			mem = nil
			runtime.GC()
		}

		if per == 0 {
			mem = make([]byte, totalMem-item.Rss)
			runtime.GC()
			time.Sleep(1 * time.Second)
		} else {
			mem = append(mem, make([]byte, incrMem)...)
			time.Sleep(per)
		}
	}
}

var mem []byte

func findItem(pid int) *dog.PsAuxItem {
	items, _ := dog.PsAuxTop(0, 0)
	for _, item := range items {
		if item.Pid == pid {
			return &item
		}
	}

	return nil
}

func printDuration(d time.Duration) string {
	if d == 0 {
		return "forever"
	}

	return "for " + d.String()
}

// RunCPULoad run CPU load in specify cores count and percentage
func RunCPULoad(coresCount int, timeSeconds time.Duration, percentage int) {
	runtime.GOMAXPROCS(coresCount)

	// second     ,s  * 1
	// millisecond,ms * 1000
	// microsecond,μs * 1000 * 1000
	// nanosecond ,ns * 1000 * 1000 * 1000

	// every loop : run + sleep = 1 unit

	// 1 unit = 100 ms may be the best
	const unitHundredOfMs = 1000
	runMicrosecond := unitHundredOfMs * percentage
	sleepMicrosecond := unitHundredOfMs*100 - runMicrosecond
	for i := 0; i < coresCount; i++ {
		go func() {
			runtime.LockOSThread()
			// endless loop
			for {
				begin := time.Now()
				for {
					// run 100%
					if time.Now().Sub(begin) > time.Duration(runMicrosecond)*time.Microsecond {
						break
					}
				}
				// sleep
				time.Sleep(time.Duration(sleepMicrosecond) * time.Microsecond)
			}
		}()
	}
	// how long
	if timeSeconds > 0 {
		time.Sleep(timeSeconds)
	} else {
		select {}
	}
}

func xxx() {
	duration := flag.Duration("d", 1*time.Minute, "duration")
	flag.Parse()

	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	quit := make(chan bool)

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-quit:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(*duration)
	for i := 0; i < n; i++ {
		quit <- true
	}
}
