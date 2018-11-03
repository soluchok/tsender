package tsender

// keeps a key (unix time, the number of nanoseconds)
// and payload which should be delivered
type message struct {
	key     int64
	payload interface{}
}

// newMessage returns a new message
func newMessage(k int64, data interface{}) *message {
	return &message{
		key:     k,
		payload: data,
	}
}
