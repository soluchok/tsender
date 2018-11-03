package tsender

import (
	"sort"
	"sync"
	"time"
)

// queue keeps messages which should be delivered to the user or group
// messages in queue are already sorted by key
type queue struct {
	mu   sync.Mutex
	data []*message
}

// newQueue returns new queue
func newQueue() *queue {
	return &queue{}
}

// index finds the index for the adding message
func (q *queue) index(x int64) int {
	return sort.Search(len(q.data), func(i int) bool { return q.data[i].key > x })
}

// checks whether a message is ready for delivery
func (q *queue) isAvailable() bool {
	return len(q.data) > 0 && q.data[0].key <= time.Now().UnixNano()
}

// push adds a message to the queue (thread-safe)
// the message must be added in proper order (by message key)
func (q *queue) push(msg *message) {
	q.mu.Lock()
	// finds index according to key (timestamp)
	i := q.index(msg.key)
	// appends message to queue
	q.data = append(q.data, msg)
	// shifts rest of messages after index by one
	copy(q.data[i+1:], q.data[i:])
	// sets a message to the specific position
	q.data[i] = msg
	q.mu.Unlock()
}

// shift returns available message from the queue and found flag (thread-safe)
// available message is a message which can be delivered to the user or group
func (q *queue) shift() (*message, bool) {
	var msg *message
	q.mu.Lock()
	ok := q.isAvailable()
	if ok {
		msg, q.data = q.data[0], q.data[1:]
	}
	q.mu.Unlock()
	return msg, ok
}
