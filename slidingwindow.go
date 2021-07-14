package dog

import (
	"sync"
	"time"
)

// Window represents a fixed-window.
type Window interface {
	// Start returns the start boundary.
	Start() time.Time

	// Count returns the accumulated count.
	Count() int64

	// AddCount increments the accumulated count by n.
	AddCount(n int64)

	// Reset sets the state of the window with the given settings.
	Reset(s time.Time, c int64)

	// Sync tries to exchange data between the window and the central
	// datastore at time now, to keep the window's count up-to-date.
	Sync(now time.Time)
}

// StopFunc stops the window's sync behaviour.
type StopFunc func()

// NewWindow creates a new window, and returns a function to stop
// the possible sync behaviour within it.
type NewWindow func() (Window, StopFunc)

type Limiter struct {
	size  time.Duration
	limit int64

	mu sync.Mutex

	curr Window
	prev Window
}

// NewLimiter creates a new limiter, and returns a function to stop
// the possible sync behaviour within the current window.
func NewLimiter(size time.Duration, limit int64, newWindow NewWindow) (*Limiter, StopFunc) {
	currWin, currStop := newWindow()

	// The previous window is static (i.e. no add changes will happen within it),
	// so we always create it as an instance of LocalWindow.
	//
	// In this way, the whole limiter, despite containing two windows, now only
	// consumes at most one goroutine for the possible sync behaviour within
	// the current window.
	prevWin, _ := NewLocalWindow()

	lim := &Limiter{
		size:  size,
		limit: limit,
		curr:  currWin,
		prev:  prevWin,
	}

	return lim, currStop
}

// Size returns the time duration of one window size. Note that the size
// is defined to be read-only, if you need to change the size,
// create a new limiter with a new size instead.
func (l *Limiter) Size() time.Duration {
	return l.size
}

// Limit returns the maximum events permitted to happen during one window size.
func (l *Limiter) Limit() int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit
}

// SetLimit sets a new Limit for the limiter.
func (l *Limiter) SetLimit(newLimit int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.limit = newLimit
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (l *Limiter) Allow() bool { return l.AllowN(time.Now(), 1) }

// AllowN reports whether n events may happen at time now.
func (l *Limiter) AllowN(now time.Time, n int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.advance(now)

	elapsed := now.Sub(l.curr.Start())
	weight := float64(l.size-elapsed) / float64(l.size)
	count := int64(weight*float64(l.prev.Count())) + l.curr.Count()

	// Trigger the possible sync behaviour.
	defer l.curr.Sync(now)

	if count+n > l.limit {
		return false
	}

	l.curr.AddCount(n)
	return true
}

// advance updates the current/previous windows resulting from the passage of time.
func (l *Limiter) advance(now time.Time) {
	// Calculate the start boundary of the expected current-window.
	newCurrStart := now.Truncate(l.size)

	diffSize := newCurrStart.Sub(l.curr.Start()) / l.size
	if diffSize >= 1 {
		// The current-window is at least one-window-size behind the expected one.

		newPrevCount := int64(0)
		if diffSize == 1 {
			// The new previous-window will overlap with the old current-window,
			// so it inherits the count.
			//
			// Note that the count here may be not accurate, since it is only a
			// SNAPSHOT of the current-window's count, which in itself tends to
			// be inaccurate due to the asynchronous nature of the sync behaviour.
			newPrevCount = l.curr.Count()
		}
		l.prev.Reset(newCurrStart.Add(-l.size), newPrevCount)

		// The new current-window always has zero count.
		l.curr.Reset(newCurrStart, 0)
	}
}

// LocalWindow represents a window that ignores sync behavior entirely
// and only stores counters in memory.
type LocalWindow struct {
	// The start boundary (timestamp in nanoseconds) of the window.
	// [start, start + size)
	start int64

	// The total count of events happened in the window.
	count int64
}

func NewLocalWindow() (*LocalWindow, StopFunc) { return &LocalWindow{}, func() {} }

func (w *LocalWindow) Start() time.Time { return time.Unix(0, w.start) }
func (w *LocalWindow) Count() int64     { return w.count }
func (w *LocalWindow) AddCount(n int64) { w.count += n }

func (w *LocalWindow) Reset(s time.Time, c int64) {
	w.start = s.UnixNano()
	w.count = c
}

func (w *LocalWindow) Sync(_ time.Time) {}

type (
	SyncRequest struct {
		Key     string
		Start   int64
		Count   int64
		Changes int64
	}

	SyncResponse struct {
		// Whether the synchronization succeeds.
		OK    bool
		Start int64
		// The changes accumulated by the local limiter.
		Changes int64
		// The total changes accumulated by all the other limiters.
		OtherChanges int64
	}

	MakeFunc   func() SyncRequest
	HandleFunc func(SyncResponse)
)

type Synchronizer interface {
	// Start starts the synchronization goroutine, if any.
	Start()

	// Stop stops the synchronization goroutine, if any, and waits for it to exit.
	Stop()

	// Sync sends a synchronization request.
	Sync(time.Time, MakeFunc, HandleFunc)
}

// SyncWindow represents a window that will sync counter data to the
// central datastore asynchronously.
//
// Note that for the best coordination between the window and the synchronizer,
// the synchronization is not automatic but is driven by the call to Sync.
type SyncWindow struct {
	LocalWindow
	changes int64

	key    string
	syncer Synchronizer
}

// NewSyncWindow creates an instance of SyncWindow with the given synchronizer.
func NewSyncWindow(key string, syncer Synchronizer) (*SyncWindow, StopFunc) {
	w := &SyncWindow{
		key:    key,
		syncer: syncer,
	}

	w.syncer.Start()
	return w, w.syncer.Stop
}

func (w *SyncWindow) AddCount(n int64) {
	w.changes += n
	w.LocalWindow.AddCount(n)
}

func (w *SyncWindow) Reset(s time.Time, c int64) {
	// Clear changes accumulated within the OLD window.
	//
	// Note that for simplicity, we do not sync remaining changes to the
	// central datastore before the reset, thus let the periodic synchronization
	// take full charge of the accuracy of the window's count.
	w.changes = 0

	w.LocalWindow.Reset(s, c)
}

func (w *SyncWindow) makeSyncRequest() SyncRequest {
	return SyncRequest{
		Key:     w.key,
		Start:   w.LocalWindow.start,
		Count:   w.LocalWindow.count,
		Changes: w.changes,
	}
}

func (w *SyncWindow) handleSyncResponse(resp SyncResponse) {
	if resp.OK && resp.Start == w.LocalWindow.start {
		// Update the state of the window, only when it has not been reset
		// during the latest sync.

		// Take the changes accumulated by other limiters into consideration.
		w.LocalWindow.count += resp.OtherChanges

		// Subtract the amount that has been synced from existing changes.
		w.changes -= resp.Changes
	}
}

func (w *SyncWindow) Sync(now time.Time) {
	w.syncer.Sync(now, w.makeSyncRequest, w.handleSyncResponse)
}
