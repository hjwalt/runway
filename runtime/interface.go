package runtime

type Runtime interface {
	Start() error
	Stop()
}

type Controller interface {
	Started() // denotes that a runtime element has started
	Stopped() // denotes that a runtime element has stopped
	Wait()
}
