package tsender

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := newQueue()
	require.NotNil(t, q)

	now := time.Now().UnixNano()

	// create two messages with the same time
	msg1 := newMessage(now, "msg1")
	msg2 := newMessage(now, "msg2")
	// create one more message with different time
	msg3 := newMessage(now+1, "msg3")

	// should be the last message in Queue
	q.push(msg3)
	// messages msg1 and msg2 should be in Queue consistently
	// should be the first message in Queue
	q.push(msg1)
	// should be the second message in Queue
	q.push(msg2)

	// gets the first message from Queue
	qMsg, ok := q.shift()
	require.True(t, ok)
	require.Equal(t, msg1, qMsg)

	// gets the second message from Queue
	qMsg, ok = q.shift()
	require.True(t, ok)
	require.Equal(t, msg2, qMsg)

	// gets the last message from Queue
	qMsg, ok = q.shift()
	require.True(t, ok)
	require.Equal(t, msg3, qMsg)

	// Queue is empty
	qMsg, ok = q.shift()
	require.False(t, ok)
	require.Nil(t, qMsg)
}
