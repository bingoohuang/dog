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

const ver = "v1.0.1 2021-07-22 09:32:19"

// https://github.com/vikyd/go-cpu-load
func main() {
	var cores int
	var p int
	var lockOsThread bool
	var duration time.Duration
	var memory string
	version := false

	numCPU := runtime.NumCPU()
	flag.IntVar(&cores, "c", numCPU, "")
	flag.IntVar(&p, "p", 100, "")
	flag.BoolVar(&lockOsThread, "l", false, "")

	flag.BoolVar(&version, "v", false, "")
	flag.StringVar(&memory, "m", "", "")
	flag.DurationVar(&duration, "d", 0, "")
	flag.Usage = func() {
		fmt.Printf(`Usage of busy (`+ver+`):
  -c int      使用核数，默认 %d
  -p int      每核 CPU 百分比 (默认 100), 0 时不开启 CPU 耗用
  -l          是否在 CPU 耗用时锁定 OS 线程
  -m string   总内存耗用，默认不开启, eg. 1) 10M 直达10M 2) 10M,1K/10s 总10M,每10秒加1K
  -d duration 跑多久，默认一直跑
  -v          看下版本号
`, runtime.NumCPU())
	}
	flag.Parse()

	if version {
		fmt.Printf(ver)
		os.Exit(0)
	}

	if cores < 1 || cores > numCPU {
		log.Fatalf("cores %d is not between 1 - %d", cores, numCPU)
	}
	if p > 100 {
		log.Fatalf("percentage %d is invalid, should be between 0 and 100", p)
	}

	log.Printf("busy starting, pid %d", os.Getpid())

	setupSignals()

	if memory != "" {
		go controlMem(memory)
	}

	if p > 0 {
		log.Printf("run %d%% of %d/%d CPU cores %s.", p, cores, numCPU, printDuration(duration))
		go RunCPULoad(cores, p, lockOsThread)
	}

	if memory == "" && p == 0 {
		log.Printf("not busy for memory or cpu, please adjust arguments")
		return
	}

	// how long
	if duration > 0 {
		time.Sleep(duration)
		log.Printf("already runned %s exiting", duration)
	} else {
		select {}
	}
}

func setupSignals() {
	sig := make(chan os.Signal, 1)
	// syscall.SIGINT: ctl + c, syscall.SIGTERM: kill pid
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			s := <-sig
			log.Printf("received signal %s", s)
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Printf("exiting")
				os.Exit(0)
			}
		}
	}()
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
	items, _ := dog.PsAuxTop(0, 0, func(topN int, heading bool) string {
		return dog.PasAuxPid(topN, pid, heading)
	})
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
func RunCPULoad(coresCount, percentage int, lockOsThread bool) {
	runtime.GOMAXPROCS(coresCount)

	// second     ,s  * 1
	// millisecond,ms * 1000
	// microsecond,μs * 1000 * 1000
	// nanosecond ,ns * 1000 * 1000 * 1000

	// every loop : run + sleep = 1 unit

	// 1 unit = 100 ms may be the best
	const unitHundredOfMs = 1000
	runMs := unitHundredOfMs * percentage
	sleepMs := unitHundredOfMs*100 - runMs
	runDuration := time.Duration(runMs) * time.Microsecond
	sleepDuration := time.Duration(sleepMs) * time.Microsecond
	for i := 0; i < coresCount; i++ {
		go func() {
			if lockOsThread {
				// https://github.com/golang/go/wiki/LockOSThread
				// Some libraries—especially graphical frameworks and libraries like Cocoa, OpenGL, and libSDL—use thread-local state and can require functions to be called only
				// from a specific OS thread, typically the 'main' thread. Go provides the runtime. LockOSThread function for this, but it's notoriously difficult to use correctly.
				// https://stackoverflow.com/a/25362395
				// With the Go threading model, calls to C code, assembler code, or blocking system calls occur in the same thread as the calling Go code, which is managed by the Go runtime scheduler.
				// The os.LockOSThread() mechanism is mostly useful when Go has to interface with some foreign library (a C library for instance). It guarantees that several successive calls to this library will be done in the same thread.
				// This is interesting in several situations:
				//   1. a number of graphic libraries (OS X Cocoa, OpenGL, SDL, ...) require all the calls to be done on a specific thread (or the main thread in some cases).
				//   2. some foreign libraries are based on thread local storage (TLS) facilities. They store some context in a data structure attached to the thread. Or some functions of the API provide results whose memory lifecycle is attached to the thread. This concept is used in both Windows and Unix-like systems. A typical example is the errno global variable commonly used in C libraries to store error codes. On systems supporting multi-threading, errno is generally defined as a thread-local variable.
				//   3. more generally, some foreign libraries may use a thread identifier to index/manage internal resources.
				//   4. doing any sort of linux namespace switch (e.g. unsharing a network or process namespace) is also bound to a thread, so if you don't lock the OS thread before you might get part of your code randomly scheduled into a different network/process namespace.
				runtime.LockOSThread()
				// runtime.UnlockOSThread()
			}
			for { // endless loop
				begin := time.Now()
				for { // run 100%
					if time.Now().Sub(begin) > runDuration {
						break
					}
				}
				time.Sleep(sleepDuration)
			}
		}()
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
