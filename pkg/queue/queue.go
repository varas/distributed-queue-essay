package queue

import (
	"sync"

	"github.com/pkg/errors"
)

// Errors
var (
	ErrEmptyQueue     = errors.New("empty queue")
	ErrUnknownCommand = errors.New("unknown command")
)

// Queue in-memory store for messages
// All its operations should be either rollbackable, either 2 phase commited (2PC).
// With 2PC we can avoid to commit the operation if cluster didn't acknowledge it.
// I didn't make that due to time constraints.
// A timeout should be added on all public methods
type Queue struct {
	messages []Message // could be a channel without mutex if we wouldn't want to replicate state across nodes
	mutex    sync.Mutex
}

// NewQueue creates a new node internal in-memory queue
func NewQueue() *Queue {
	return &Queue{}
}

// Queue adds a message to the queue
// Returns if action was done to avoid broadcast an already taken action
func (r *Queue) Queue(msg Message) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.messages = append(r.messages, msg)

	return
}

// Dequeue extract a message from the queue
func (r *Queue) Dequeue() (Message, error) {
	if len(r.messages) == 0 {
		return Message{}, ErrEmptyQueue
	}

	msg := r.messages[0]
	r.messages = r.messages[1:]

	return msg, nil
}

// ReplicateDequeue extract a given message from the queue
func (r *Queue) ReplicateDequeue(cmd ReplicateDequeueRequest) error {
	// TODO

	return nil
}

// ReplicateEnqueue replicates a enqueue command
func (r *Queue) ReplicateEnqueue(cmd ReplicateEnqueueRequest) error {
	// TODO

	return nil
}
