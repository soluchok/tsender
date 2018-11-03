package tsender

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockProvider struct {
	send func(msg interface{})
}

func (p *mockProvider) Send(msg interface{}) {
	if p.send != nil {
		p.send(msg)
	}
}

func TestNewSender(t *testing.T) {
	s := NewSender(nil)
	require.NotNil(t, s)
	require.Nil(t, s.provider)
	require.NotNil(t, s.queue)
	require.NotNil(t, s.distribute)
	require.NotNil(t, s.done)
}

func TestSender_Run(t *testing.T) {
	t.Parallel()
	s := NewSender(nil)
	done := make(chan struct{})
	go func() {
		s.Run(0)
		close(done)
	}()
	s.Stop()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}

func TestSender_Send(t *testing.T) {
	t.Parallel()
	execute(t, 30, int64(duration))
}

func TestSender_SendOrdinary(t *testing.T) {
	t.Parallel()
	execute(t, 3, latencyOrdinary)
}

func TestSender_SendGroup(t *testing.T) {
	t.Parallel()
	execute(t, 2, latencyGroup)
}

func execute(t *testing.T, msgCount int, latency int64) {
	var (
		prevTime int64
		count    int
		done     = make(chan struct{})
	)
	mock := &mockProvider{
		send: func(msg interface{}) {
			msgTime := time.Now().UnixNano()
			diff := msgTime - (prevTime + int64(duration))
			approximateError := int64(time.Millisecond / 2)
			minNextTime := prevTime + latency - approximateError
			require.GreaterOrEqual(t, msgTime, minNextTime, diff)
			prevTime = msgTime
			count++
			if count == msgCount {
				close(done)
			}
		},
	}

	s := NewSender(mock)
	go s.Run(1)

	for i := 0; i < msgCount; i++ {
		id := i
		switch latency {
		case latencyGroup:
			id = -1
		case latencyOrdinary:
			id = 1
		}
		s.Send(int64(id), nil)
	}

	select {
	case <-done:
		s.Stop()
	case <-time.After(time.Duration(latency * int64(msgCount))):
		t.Error("timeout")
		return
	}
}
