package tsender

import (
	"sync"
	"time"
)

// chat keeps a previous key(based on timestamp) and latency
type chat struct {
	// chat latency depends on whether it is user or group
	latency int64
	sync.Mutex
	// keeps previous key
	prev int64
}

// key returns next valid key for chat, the key is based on a timestamp (the number of nanoseconds)
// it is calculated according to latency and the previous key
// the difference between calling `key()` function twice must be more or equal previous key + latency
// e.g key() - key() >= latency
func (c *chat) key() int64 {
	c.Lock()
	defer c.Unlock()

	// generates new possible key value
	var key = time.Now().UnixNano()

	// generated key become actual if the previous key is absent
	if c.prev == 0 {
		c.prev = key
		// returns actual key
		return c.prev
	}

	// calculates minimum next possible key value
	c.prev = c.prev + c.latency

	// generated key become actual if generated key greater than the minimum possible key
	if key > c.prev {
		c.prev = key
	}

	// returns actual key
	return c.prev
}

// newChat returns latency based chat structure
func newChat(l int64) *chat {
	return &chat{latency: l}
}
