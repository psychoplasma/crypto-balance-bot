package concurrency_test

import (
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
)

func TestJob_RunAndWait(t *testing.T) {
	err_msg := "test_panic"
	j1, err := concurrency.NewJob()
	if err != nil {
		t.Fatal(err)
	}

	// failing job
	go j1.Run(func() {
		time.Sleep(time.Millisecond * 100)
		panic(err_msg)
	})

	if err := j1.Wait(func() {}); err == nil || err.Error() != err_msg {
		t.Fatalf("\ngot:%s\nwant:%s", err, err_msg)
	}

	j2, err := concurrency.NewJob()
	if err != nil {
		t.Fatal(err)
	}

	// succeeding job
	go j2.Run(func() {
		time.Sleep(time.Millisecond * 100)
	})

	if err := j2.Wait(func() {}); err != nil {
		t.Fatal(err)
	}
}

func TestWorker_Run(t *testing.T) {
	numberOfJobs := 2
	err_msg := "test_panic"

	w := concurrency.NewWorker(
		numberOfJobs,
		time.Second*1,
	)

	// succeeding job
	j1, err := w.Run(func() {
		time.Sleep(time.Millisecond * 100)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !w.IsJobAlive(j1) {
		t.Fatalf("job is supposed to be alive")
	}

	// failing job
	j2, err := w.Run(func() {
		time.Sleep(time.Millisecond * 200)
		panic(err_msg)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !w.IsJobAlive(j2) {
		t.Fatalf("job is supposed to be alive")
	}

	if w.JobsInProgress() != numberOfJobs {
		t.Fatalf("Jobs in progress\ngot:%d\nwant:%d", w.JobsInProgress(), numberOfJobs)
	}

	w.WaitAll()
	if j2.Error == nil || j2.Error.Error() != err_msg {
		t.Fatalf("\ngot:%s\nwant:%s", j2.Error, err_msg)
	}

	if w.IsJobAlive(j1) || w.IsJobAlive(j2) {
		t.Fatalf("all jobs are supposed to be done")
	}
}

func TestWorker_Run_WithMoreThanParalellismLimit(t *testing.T) {
	parallelismLimit := 2
	expected_err := fmt.Sprintf("parallel execution limit(%d) has already been reached", parallelismLimit)

	w := concurrency.NewWorker(
		parallelismLimit,
		time.Second*1,
	)

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 100)
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 200)
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 300)
	}); err == nil || err.Error() != expected_err {
		t.Fatalf("\ngot: %s\nwant: %s", err, expected_err)
	}
}

func TestWorker_Run_WhenStopped(t *testing.T) {
	expected_err := "cannot add a new job. worker is exiting"

	w := concurrency.NewWorker(
		2,
		time.Second*1,
	)

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 100)
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Stop(); err != nil {
		t.Fatal(err)
	}

	// Try to run a job after stopped
	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 100)
	}); err == nil || err.Error() != expected_err {
		t.Fatalf("\ngot: %s\nwant: %s", err, expected_err)
	}
}

func TestWorker_Stop(t *testing.T) {
	w := concurrency.NewWorker(
		2,
		time.Second*1,
	)

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 200)
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 400)
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Stop(); err != nil {
		t.Fatal(err)
	}

	if w.JobsInProgress() > 0 {
		t.FailNow()
		t.Fatalf("Jobs in progress\ngot:%d\nwant:%d", w.JobsInProgress(), 0)
	}
}

func TestWorker_Stop_WithTimeout(t *testing.T) {
	expected_err := "exit timeout exceeded"
	w := concurrency.NewWorker(
		2,
		time.Millisecond*500,
	)

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 200)
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := w.Run(func() {
		time.Sleep(time.Millisecond * 600)
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Stop(); err == nil || err.Error() != expected_err {
		t.Fatalf("\ngot: %s\nwant: %s", err, expected_err)
	}
}

func TestWorker_Run_WhenInterrupted(t *testing.T) {
	err_msg := "interrupt signal has been received"
	w := concurrency.NewWorker(
		2,
		time.Second*1,
	)

	j1, err1 := w.Run(func() {
		time.Sleep(time.Millisecond * 500)
	})
	if err1 != nil {
		t.Fatal(err1)
	}

	j2, err2 := w.Run(func() {
		time.Sleep(time.Millisecond * 700)
	})
	if err2 != nil {
		t.Fatal(err2)
	}

	time.Sleep(time.Millisecond * 200)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	w.WaitAll()

	if j1.Error == nil || j1.Error.Error() != err_msg {
		t.Fatalf("\ngot: %s\nwant: %s", j1.Error, err_msg)
	}

	if j2.Error == nil || j2.Error.Error() != err_msg {
		t.Fatalf("\ngot: %s\nwant: %s", j2.Error, err_msg)
	}
}
