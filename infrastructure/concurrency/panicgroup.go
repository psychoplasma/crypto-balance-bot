package concurrency

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// PanicGroup runs functions in go-routines and returns upon
// the first error panicking from either of go-routines or
// system interrupt signal or completion of all the go-routines
type PanicGroup struct {
	i       int // Number of go-routines running
	m       sync.Mutex
	errOnce sync.Once
	err     chan error     // Channel to capture the first panicking error
	s       chan os.Signal // Channel to capture system interrupt signal
	done    chan struct{}  // Channel to signal that all go-routines are done
}

// NewPanicGroup creates a new instance of PanicGroup
func NewPanicGroup() *PanicGroup {
	pg := &PanicGroup{
		i:    0,
		done: make(chan struct{}),
		m:    sync.Mutex{},
		s:    make(chan os.Signal),
		err:  make(chan error),
	}

	signal.Notify(pg.s, syscall.SIGINT, syscall.SIGTERM)
	return pg
}

// Go runs the given function in a go-routine
// The first panicking error will be caught and
// the panic group instance will be halted with the error
func (pg *PanicGroup) Go(f func()) {
	pg.add(1)
	go pg.catch(f)
}

// Wait waits until either an panicking error is caugth or
// interrupt signal is received or all the go-routines created by this group return
func (pg *PanicGroup) Wait() error {
	// If the Wait function is called without run any function trough Go()
	// then return immediately without waiting
	if pg.i <= 0 {
		return nil
	}

	for {
		select {
		case <-pg.s:
			return errors.New("interrupt signal has been received")
		case err := <-pg.err:
			return err
		case <-pg.done:
			return nil
		}
	}
	// We're not closing the channels here because not all the go-routines may finish yet
	// and if any remaining go-routine try to write any thing to a closed channel will
	// program to panic. Rather we leave it to the caller to clear the open channels
	// according to its logic.
}

// Clear closes any remainging open channels
func (pg *PanicGroup) Clear() {
	close(pg.err)
	close(pg.s)
	close(pg.done)
}

func (pg *PanicGroup) add(i int) {
	pg.m.Lock()
	pg.i += i
	pg.m.Unlock()
}

func (pg *PanicGroup) remove() {
	pg.m.Lock()
	pg.i--
	if pg.i < 1 {
		pg.done <- struct{}{}
	}
	pg.m.Unlock()
}

func (pg *PanicGroup) catch(f func()) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				// Capture the first panicking error
				pg.errOnce.Do(func() {
					pg.err <- err
				})
			} else {
				pg.err <- fmt.Errorf("Cannot recover panicking error. %+v", r)
			}
		} else {
			pg.remove()
		}
	}()

	f()
}
