package queue

import (
	"context"

	"sync"
)

// CmdHandler idempotent queue-operation request handler
type CmdHandler struct {
	queue          *Queue
	processedCmds  map[CmdID]struct{} // enables idempotency (we should use a hash list for storage reduction, but CmdID kept for readability)
	processedMutex sync.Mutex
}

// NewCmdHandler initializes a command handler
func NewCmdHandler(
	queue *Queue,
) *CmdHandler {
	return &CmdHandler{
		queue:         queue,
		processedCmds: make(map[CmdID]struct{}),
	}
}

// Handle handles a queue operation request
// TODO handle context and consistency
func (h *CmdHandler) Handle(ctx context.Context, cmd Command) (result CommandResult) {
	result.ID = cmd.ID

	if h.wasProcessed(cmd) {
		return
	}

	switch c := cmd.Request.(type) {

	case EnqueueRequest:
		h.queue.Queue(c.Msg)

	case DequeueRequest:
		result.DequeuedMsg, result.Err = h.queue.Dequeue()

	case ReplicateDequeueRequest:
		result.Err = h.queue.ReplicateDequeue(c)

	case ReplicateEnqueueRequest:
		result.Err = h.queue.ReplicateEnqueue(c)

	default:
		result.Err = ErrUnknownCommand
	}

	return result
}

func (h *CmdHandler) wasProcessed(cmd Command) bool {
	h.processedMutex.Lock()
	defer h.processedMutex.Unlock()

	if _, exists := h.processedCmds[cmd.ID]; exists {
		return true
	}

	h.processedCmds[cmd.ID] = struct{}{}

	return false
}
