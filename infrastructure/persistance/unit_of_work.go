package persistance

import "sync"

// UnitOfWork implements a set of functionalities for Unit-of-Work design pattern
type UnitOfWork struct {
	m sync.Mutex
}

func (w UnitOfWork) Register() {
	w.m.Lock()
}

func (w UnitOfWork) Finalize() error {
	defer w.m.Unlock()

	return nil
}
