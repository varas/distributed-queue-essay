package node

import (
	"context"
	"net"

	"sync"

	"github.com/pkg/errors"
	"github.com/tevino/abool"
	"github.com/varas/distributed-queue-essay/pkg/cluster"
	"github.com/varas/distributed-queue-essay/pkg/errhandler"
	"github.com/varas/distributed-queue-essay/pkg/queue"
)

// runtime node runtime context
type runtime struct {
	isUp           *abool.AtomicBool
	stopped        chan struct{}
	errHandle      errhandler.ErrHandler
	cancelListener context.CancelFunc
	cancelHandlers context.CancelFunc
	wgHandlers     sync.WaitGroup
}

// passing config on start enables hot config-reloading
func (r *runtime) start(ctx context.Context, c config, errHandle errhandler.ErrHandler, publisher cluster.Publisher) (err error) {
	r.stopped = make(chan struct{})
	r.errHandle = errHandle

	conns := make(chan net.Conn)
	listener, err := NewListener(c.port, conns)
	if err != nil {
		return errors.Wrap(err, "cannot create connection listener")
	}

	// stop runtime in order
	var ctxListener, ctxHandlers context.Context
	ctxListener, r.cancelListener = context.WithCancel(ctx)
	ctxHandlers, r.cancelHandlers = context.WithCancel(ctx)

	// stop bg jobs: listener
	go func() {
		listerErr := listener.Listen(ctxListener)
		// avoid connection closed as error on teardown
		if r.isUp.IsSet() {
			r.errHandle(listerErr)
		}
	}()

	cmdHandler := queue.NewCmdHandler(queue.NewQueue())
	connHandler := newConnHandler(errHandle, conns, cmdHandler, publisher)

	r.wgHandlers = sync.WaitGroup{}
	r.wgHandlers.Add(c.handlersAmount)
	for w := c.handlersAmount; w > 0; w-- {
		go func() {
			connHandler.run(ctxHandlers)
			r.wgHandlers.Done()
		}()
	}

	go r.waitForContextTermination(ctx)

	r.isUp = abool.NewBool(true)

	return nil
}

func (r *runtime) stop() {
	wasStopped := r.isUp.SetToIf(true, false)
	if !wasStopped {
		return
	}

	r.cancelListener()

	// stop conn handlers
	r.cancelHandlers()
	r.wgHandlers.Wait()

	close(r.stopped)
}

func (r *runtime) waitForContextTermination(ctx context.Context) {
	<-ctx.Done()
	r.stop()
}
