package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Dequeue(t *testing.T) {
	q := NewQueue()

	msg1 := Message{CreatedAt: time.Now()}
	msg2 := Message{CreatedAt: time.Now()}

	q.Queue(msg1)
	q.Queue(msg2)

	dequeued, err := q.Dequeue()

	assert.NoError(t, err)
	assert.Equal(t, msg1, dequeued)

	dequeued, err = q.Dequeue()

	assert.NoError(t, err)
	assert.Equal(t, msg2, dequeued)
}

func TestQueue_DequeueOnEmpty(t *testing.T) {
	q := NewQueue()

	_, err := q.Dequeue()

	assert.Error(t, err)
	assert.Equal(t, ErrEmptyQueue, err)
}
