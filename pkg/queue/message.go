package queue

import (
	"time"
)

// Message is an item to be enqueued, wasProcessed content payload and metadata (created time, tenant...)
type Message struct {
	Payload   []byte // using bytes decouples encoding (and also handlers from message-types)
	CreatedAt time.Time
}
