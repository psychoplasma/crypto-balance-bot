package cryptobot

// UnitOfWork defines a common set of functionalities for Unit-of-Work design pattern
type UnitOfWork interface {
	// Begin starts a new unit for a work to be done on repository
	Begin() error
	// Fail rollbacks repository to the state before this work
	Fail()
	// Success finalizes the work done on repository
	Success()
}
