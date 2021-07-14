package dog

import (
	"fmt"
	"time"
)

func ExampleLocalWindow() {
	lim, _ := NewLimiter(time.Second, 10, func() (Window, StopFunc) {
		// NewLocalWindow returns an empty stop function, so it's
		// unnecessary to call it later.
		return NewLocalWindow()
	})

	ok := lim.Allow()
	fmt.Printf("ok: %v\n", ok)

	// Output:
	// ok: true
}
