package broadcaster

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	b := NewBroadcaster()

	done := make(chan bool, 2)
	b.Go(func() {
		log.Println("function 1 finished")
		done <- true
	})
	b.Go(func() {
		log.Println("function 2 finished")
		done <- true
	})

	if len(done) != 0 {
		t.Fatalf("Failed to wait for broadcast")
	}

	b.Broadcast()
	<-done
	<-done
}

func TestLateBroadcast(t *testing.T) {
	b := NewBroadcaster()

	done := make(chan bool, 2)
	b.Go(func() {
		log.Println("function 1 finished")
		done <- true
	})
	b.Go(func() {
		log.Println("function 2 finished")
		done <- true
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		b.Broadcast()
	}()

	<-done
	<-done
}

func TestBroadcastWaitDone(t *testing.T) {
	b := NewBroadcaster()

	var counter atomic.Uint32

	b.Go(func() {
		log.Println("function 1 finished")
		counter.Add(1)
	})
	b.Go(func() {
		log.Println("function 2 finished")
		counter.Add(1)
	})

	b.Broadcast()
	b.Wait()

	if counter.Load() != 2 {
		t.Fatalf("Failed to wait for all functions to finish")
	}
}
