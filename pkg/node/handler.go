package node

import (
	"context"
	"io"
	"net"

	"github.com/varas/distributed-queue-essay/pkg/cluster"
	"github.com/varas/distributed-queue-essay/pkg/errhandler"
	"github.com/varas/distributed-queue-essay/pkg/queue"
)

type connHandler struct {
	errHandle  errhandler.ErrHandler
	conns      <-chan net.Conn
	cmdHandler *queue.CmdHandler
	publisher  cluster.Publisher
}

func newConnHandler(
	errHandle errhandler.ErrHandler,
	conns <-chan net.Conn,
	cmdHandler *queue.CmdHandler,
	clusterPublisher cluster.Publisher,
) *connHandler {
	return &connHandler{
		errHandle:  errHandle,
		conns:      conns,
		cmdHandler: cmdHandler,
		publisher:  clusterPublisher,
	}
}

func (r *connHandler) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case c, open := <-r.conns:
			if !open {
				return
			}
			r.handle(ctx, c)
		}
	}
}

// context unhandled here to avoid data loss, as client has no guarantees of sent data is processed on service stop
func (r *connHandler) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	commandSerializer := queue.NewConnSerializer(conn)

	for {
		cmd, err := commandSerializer.ReadCommand()
		if err == io.EOF {
			return
		}
		if err != nil {
			r.errHandle(err)
			continue
		}

		result := r.cmdHandler.Handle(ctx, cmd)

		if result.Err != nil {
			r.errHandle(r.publisher.Publish(eventFromResult(result)))
		}

		err = commandSerializer.WriteResponse(result)
		if err != nil {
			// TODO retry/rollback policy
			r.errHandle(err)
		}
	}
}

// TODO
func eventFromResult(result queue.CommandResult) cluster.QueueEvent {
	return cluster.QueueEvent{}
}
