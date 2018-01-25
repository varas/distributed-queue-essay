package node

import (
	"context"

	"github.com/pkg/errors"
	"github.com/varas/distributed-queue-essay/pkg/cluster"
	"github.com/varas/distributed-queue-essay/pkg/errhandler"
)

// Node tcp server that store unique numbers writen
// It works as a bg daemon, so its API is based on channels to trigger graceful stop and wait for state completions
type Node struct {
	config           config
	runtime          *runtime
	errHandle        errhandler.ErrHandler
	clusterPublisher cluster.Publisher
	// channel api to wait (on channel close) for status change
	Ready   chan struct{} // enables to wait until ready
	Stop    chan struct{} // enables to gracefully stop the server
	Stopped chan struct{} // enables to wait until stopped
}

// New generates a new queue node
func New(publicPort int, publisher cluster.Publisher) *Node {
	errHandle := errhandler.Logger("[error] ")

	return &Node{
		config:           *newConfig(publicPort),
		runtime:          &runtime{}, // stateless runtime, loses queue items on restart
		errHandle:        errHandle,
		clusterPublisher: publisher,
		Ready:            make(chan struct{}),
	}
}

// Run bootstraps the runtime so resilience could be added via recover, and runs the app
// context cancellation is aimed for fast teardown, for graceful stop use Stop channel instead
func (n *Node) Run(ctx context.Context) {
	n.Stop = make(chan struct{})
	n.Stopped = make(chan struct{})

	err := n.runtime.start(ctx, n.config, n.errHandle, n.clusterPublisher)
	if err != nil {
		n.errHandle(errors.Wrap(err, "error on start"))
		return
	}

	go n.waitForContextTermination(ctx)
	go n.waitForClientStop()

	close(n.Ready)

	<-n.runtime.stopped
	n.stop()
}

func (n *Node) stop() {
	if n.runtime.isUp.IsSet() {
		n.runtime.stop()
		close(n.Stopped)
	}
}

func (n *Node) waitForClientStop() {
	<-n.Stop
	n.stop()
}

func (n *Node) waitForContextTermination(ctx context.Context) {
	<-ctx.Done()
	n.stop()
}
