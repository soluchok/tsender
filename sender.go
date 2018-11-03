package tsender

import (
	"sync"
	"time"
)

const (
	duration        = time.Second / 30
	latencyOrdinary = int64(time.Second)
	latencyGroup    = int64(time.Minute / 20)
)

// Provider ensure implementation for sending message
type Provider interface {
	Send(msg interface{})
}

type Sender struct {
	done       chan struct{}
	provider   Provider
	chats      sync.Map
	distribute chan *message
	wg         sync.WaitGroup
	queue      *queue
}

func (s *Sender) getChat(ID int64) *chat {
	latency := latencyOrdinary
	if ID < 0 {
		latency = latencyGroup
	}

	raw, _ := s.chats.LoadOrStore(ID, newChat(latency))
	return raw.(*chat)
}

func (s *Sender) worker() {
	defer s.wg.Done()

	for msg := range s.distribute {
		s.provider.Send(msg.payload)
	}
}

// Run workers to distribute messages
func (s *Sender) Run(workers int) {
	if workers <= 0 {
		workers = 1
	}

	s.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go s.worker()
	}

	for {
		select {
		case <-s.done:
			close(s.distribute)
			s.wg.Wait()
			return
		default:
			if msg, ok := s.queue.shift(); ok {
				s.distribute <- msg
			}
			time.Sleep(duration)
		}
	}
}

// Stop distributing messages at all
func (s *Sender) Stop() {
	close(s.done)
}

// Send a message to a particular chat or group thread-safe
func (s *Sender) Send(chatID int64, msg interface{}) {
	s.queue.push(newMessage(s.getChat(chatID).key(), msg))
}

// NewSender returns *Sender using time.Second / 30 strategy
// to distribute messages also avoid sending
// more than one message per second to a particular chat
// and sending more than 20 messages per minute to the same group.
func NewSender(p Provider) *Sender {
	return &Sender{
		provider:   p,
		queue:      newQueue(),
		distribute: make(chan *message),
		done:       make(chan struct{}),
	}
}
