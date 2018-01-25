package queue

import (
	"context"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestCmdHandler_Handle(t *testing.T) {
	q := NewQueue()

	handler := NewCmdHandler(q)

	msg := Message{
		CreatedAt: time.Now(),
	}

	result := handler.Handle(context.Background(), Command{
		ID:      CmdID("some-id"),
		Request: EnqueueRequest{Msg: msg},
	})

	assert.NoError(t, result.Err)
	assert.Equal(t, CmdID("some-id"), result.ID)

	result = handler.Handle(context.Background(), Command{
		ID:      CmdID("other-id"),
		Request: DequeueRequest{},
	})

	assert.NoError(t, result.Err)
	assert.Equal(t, CmdID("other-id"), result.ID)
	assert.Equal(t, msg, result.DequeuedMsg)
}
