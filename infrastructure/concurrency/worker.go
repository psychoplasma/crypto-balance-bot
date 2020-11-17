package concurrency

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Job represents a handle to a function to be executed in parallel with some others in a Worker instance
type Job struct {
	id    string
	done  chan struct{}
	s     chan os.Signal // Channel to capture system interrupt signal
	Error error
}

func generateJid() (string, error) {
	// Return 12 random bytes as 24 character hex
	b := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// NewJob creates a new instance of Job
func NewJob() (*Job, error) {
	jid, err := generateJid()
	if err != nil {
		return nil, fmt.Errorf("cannot generate job id")
	}

	j := &Job{
		id:   jid,
		done: make(chan struct{}),
		s:    make(chan os.Signal),
	}

	signal.Notify(j.s, syscall.SIGINT, syscall.SIGTERM)

	return j, nil
}

// Run runs the given function
func (j *Job) Run(f func()) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); !ok {
				j.Error = fmt.Errorf("%s", r)
			} else {
				j.Error = err
			}
		}

		j.done <- struct{}{}
	}()

	f()
}

// Wait waits for the corresponding job to be completed and returns any error if there is
func (j *Job) Wait(finalizer func()) error {
	defer finalizer()
	// defer close(j.done)
	// defer close(j.s)

	select {
	case <-j.done:
		return j.Error
	case <-j.s:
		j.Error = errors.New("interrupt signal has been received")
		return j.Error
	}
}

// Worker represents a manager to execute jobs in parallel
type Worker struct {
	m              *sync.Mutex
	jobs           map[string]*Job
	maxParallelism int
	waitDelay      time.Duration
	exitTimeout    time.Duration
	isStopped      bool
}

// NewWorker creates a new instance of Worker
func NewWorker(maxParallelism int, exitTimeout time.Duration) *Worker {
	w := &Worker{
		m:              &sync.Mutex{},
		jobs:           make(map[string]*Job),
		maxParallelism: maxParallelism,
		isStopped:      false,
		waitDelay:      time.Millisecond * 500,
		exitTimeout:    exitTimeout,
	}

	return w
}

func (w *Worker) add() (*Job, error) {
	if w.isStopped {
		return nil, fmt.Errorf("cannot add a new job. worker is exiting")
	}

	w.m.Lock()
	defer w.m.Unlock()
	if len(w.jobs) >= w.maxParallelism {
		return nil, fmt.Errorf("parallel execution limit(%d) has already been reached", w.maxParallelism)
	}

	j, err := NewJob()
	if err != nil {
		return nil, err
	}

	w.jobs[j.id] = j

	return j, nil
}

func (w *Worker) remove(j *Job) {
	w.m.Lock()
	delete(w.jobs, j.id)
	w.m.Unlock()
}

func (w *Worker) await(timeout time.Duration) error {
	elapsedTime := time.Millisecond * 0
	for len(w.jobs) > 0 {
		if timeout > 0 && elapsedTime >= timeout {
			return fmt.Errorf("exit timeout exceeded")
		}
		time.Sleep(w.waitDelay)
		elapsedTime += w.waitDelay
	}

	return nil
}

// Run runs a job in parallel with other jobs in this worker instance
func (w *Worker) Run(f func()) (*Job, error) {
	j, err := w.add()
	if err != nil {
		return nil, err
	}

	go j.Run(f)
	go j.Wait(func() { w.remove(j) })

	return j, nil
}

// WaitAll waits for all jobs to finish
func (w *Worker) WaitAll() {
	w.await(0)
}

// Stop stops the worker instance and wait until the in-progress jobs to finish or exits with an error upon exit timeout
func (w *Worker) Stop() error {
	w.isStopped = true
	return w.await(w.exitTimeout)
}

// JobsInProgress returns the number of jobs in-progress
func (w *Worker) JobsInProgress() int {
	return len(w.jobs)
}

// IsJobAlive returns  whether the given job is alive(runnig) or not
func (w *Worker) IsJobAlive(j *Job) bool {
	_, ok := w.jobs[j.id]
	return ok
}
