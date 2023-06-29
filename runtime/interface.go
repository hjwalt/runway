package runtime

type Runtime interface {
	Start() error
	Stop()
	SetController(Controller)
}

type Controller interface {
	Started()    // denotes that a runtime element has started
	Stopped()    // denotes that a runtime element has stopped
	Error(error) // denotes that a runtime element has errored
	Wait()
}
