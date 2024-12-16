// Package broadcaster provides some concurrency
// tools to block and unblock execution across
// many goroutines based on broadcasting a signal.
package broadcaster

import (
	"sync"
)

type Broadcaster struct {
	cond     *sync.Cond
	signaled bool
	wg       sync.WaitGroup
}

// NewBroadcaster creates a new broadcaster.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		cond:     sync.NewCond(&sync.Mutex{}),
		signaled: false,
	}
}

// Go waits until something is broadcasted,
// and runs the given function in a new
// goroutine.
func (b *Broadcaster) Go(fn func()) {
	b.wg.Add(1)
	go func() {
		b.cond.L.Lock()

		for !b.signaled {
			b.cond.Wait()
		}
		fn()

		b.cond.L.Unlock()
		b.wg.Done()
	}()
}

// Broadcast broadcasts a signal to all
// waiting function and unblocks them.
func (b *Broadcaster) Broadcast() {
	b.cond.L.Lock()
	b.signaled = true
	b.cond.L.Unlock()

	b.cond.Broadcast()
}

// Wait waits for all functions to finish.
func (b *Broadcaster) Wait() {
	b.wg.Wait()
}
