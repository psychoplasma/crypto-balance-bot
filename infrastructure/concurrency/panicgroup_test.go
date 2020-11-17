package concurrency_test

import (
	"errors"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
)

func TestGo_WithoutError(t *testing.T) {
	pg := concurrency.NewPanicGroup()

	pg.Go(func() {
		time.Sleep(time.Millisecond * 100)
	})

	pg.Go(func() {
		time.Sleep(time.Millisecond * 300)
	})

	if err := pg.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestGo_CatchingFirstPanickingError(t *testing.T) {
	errMsg := "test_panicking_error"
	pg := concurrency.NewPanicGroup()

	pg.Go(func() {
		time.Sleep(time.Millisecond * 100)
	})

	pg.Go(func() {
		time.Sleep(time.Millisecond * 300)
	})

	pg.Go(func() {
		time.Sleep(time.Millisecond * 150)
		panic(errors.New(errMsg))
	})

	if err := pg.Wait(); err != nil {
		if err.Error() != errMsg {
			t.Fatalf("got: %s\nwant: %s", err.Error(), errMsg)
		}
	}
}

func TestGo_CatchingFirstPanic_WithoutErrorStruct(t *testing.T) {
	errMsg := "test_panicking_error"
	expectedMsg := fmt.Sprintf("Cannot recover panicking error. %+v", errMsg)
	pg := concurrency.NewPanicGroup()

	pg.Go(func() {
		time.Sleep(time.Millisecond * 100)
	})

	pg.Go(func() {
		// We are paniking with a string on purpose
		// to see if it will catch even without error struct
		// which imitates that recover() cannot asserts to error struct
		panic(errMsg)
	})

	if err := pg.Wait(); err != nil {
		if err.Error() != expectedMsg {
			t.Fatalf("got: %s\nwant: %s", err.Error(), expectedMsg)
		}
	}
}

func TestGo_CathingInterruptSignal(t *testing.T) {
	errMsg := "interrupt signal has been received"
	pg := concurrency.NewPanicGroup()

	pg.Go(func() {
		time.Sleep(time.Millisecond * 2000)
	})

	time.Sleep(time.Millisecond * 500)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	if err := pg.Wait(); err != nil {
		if err.Error() != errMsg {
			t.Fatalf("got: %s\nwant: %s", err.Error(), errMsg)
		}
	}
}
