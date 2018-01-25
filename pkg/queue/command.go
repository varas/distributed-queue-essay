package queue

// CmdID is a unique identifier for a command (may be a UUID)
// Enables idempotency (duplicated commands can be discarded) and tracing
type CmdID string

// Command command to enqueue a message or to dequeue a message
type Command struct {
	ID      CmdID
	Request interface{}
}

// CommandResult command result
type CommandResult struct {
	ID          CmdID
	Err         error
	DequeuedMsg Message
}

// EnqueueRequest enqueue request
type EnqueueRequest struct {
	Msg Message
}

// DequeueRequest dequeue request
type DequeueRequest struct {
}

// ReplicateCommand command to replicate a queue operation processed out of the node
type ReplicateCommand struct {
	ID         CmdID
	Replicated Command
	Request    interface{}
}

// ReplicateEnqueueRequest replicate enqueue request
type ReplicateEnqueueRequest struct {
	Order int // to place message in expected order
}

// ReplicateDequeueRequest replicate dequeue request
type ReplicateDequeueRequest struct {
}
