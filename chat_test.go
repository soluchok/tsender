package tsender

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChat(t *testing.T) {
	// no previous key (latencyOrdinary)
	c := newChat(latencyOrdinary)
	first := c.key()
	second := c.key()
	require.Equal(t, latencyOrdinary, second-first)

	// with previous key (latencyOrdinary)
	c.prev = time.Now().UnixNano() - latencyOrdinary
	first = c.key()
	second = c.key()
	require.Equal(t, latencyOrdinary, second-first)

	// no previous key (latencyGroup)
	c = newChat(latencyGroup)
	first = c.key()
	second = c.key()
	require.Equal(t, latencyGroup, second-first)

	// with previous key (latencyGroup)
	c.prev = time.Now().UnixNano() - latencyGroup
	first = c.key()
	second = c.key()
	require.Equal(t, latencyGroup, second-first)
}
