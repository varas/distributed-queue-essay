package queue

import (
	"net"
)

// CommandSerializer reads message commands to enqueue or dequeue and provides execution confirmation (Ack/Nack)
type CommandSerializer interface {
	ReadCommand() (Command, error)
	WriteResponse(CommandResult) error
}

// NewConnSerializer creates a Message reader from a connection
// Illustrative, not implemented
func NewConnSerializer(conn net.Conn) CommandSerializer {
	return &noopSerializer{}
}

type noopSerializer struct{}

func (r *noopSerializer) ReadCommand() (cmd Command, err error)   { return }
func (r *noopSerializer) WriteResponse(CommandResult) (err error) { return }
