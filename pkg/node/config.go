package node

// Default config values
const (
	DefaultHandlersAmount = 4
)

type config struct {
	port           int
	handlersAmount int // I assume nodes would do some kind of message processing, so several handlers leverage this
}

func newConfig(port int) *config {
	return &config{
		port:           port,
		handlersAmount: DefaultHandlersAmount,
	}
}
