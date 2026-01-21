package mock

import (
	"context"
	"sync"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
)

type MockSubscription struct {
	ctxCancel context.CancelFunc
}

func (ms MockSubscription) Close() error {
	// No-op for mock
	ms.ctxCancel()
	return nil
}

// MockMessageBus is a mock that implements the Messagebus interface.
//
// Messages are tracked per queue, and ReceiveMessage registers subscribers
// that receive new messages for the queue.
type MockMessageBus struct {
	Messages    map[string][]string
	subscribers map[string][]chan<- string
	mu          sync.Mutex
	sendError   error
}

// NewMockMessageBus creates a new MockMessageBus instance
func NewMockMessageBus() *MockMessageBus {
	return &MockMessageBus{
		Messages:    make(map[string][]string),
		subscribers: make(map[string][]chan<- string),
	}
}

// SetSendError allows setting an error to be returned by SendMessage for testing
func (mb *MockMessageBus) SetSendError(err error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.sendError = err
}

// GetMessages returns all messages sent to a specific queue
func (mb *MockMessageBus) GetMessages(queue string) []string {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	return append([]string{}, mb.Messages[queue]...)
}

// ClearMessages clears all messages for a specific queue
func (mb *MockMessageBus) ClearMessages(queue string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.Messages[queue] = nil
}

func (mb *MockMessageBus) SendMessage(message []byte, queue string) error {
	mb.mu.Lock()
	if mb.sendError != nil {
		err := mb.sendError
		mb.mu.Unlock()
		return err
	}

	if mb.Messages == nil {
		mb.Messages = make(map[string][]string)
	}
	mb.Messages[queue] = append(mb.Messages[queue], string(message))
	subs := append([]chan<- string(nil), mb.subscribers[queue]...)
	mb.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- string(message):
		default:
			// Channel full, skip
		}
	}
	return nil
}

func (mb *MockMessageBus) ReceiveMessage(message chan<- string, queue string) (messagebus.Subscription, error) {
	mb.mu.Lock()
	if mb.subscribers == nil {
		mb.subscribers = make(map[string][]chan<- string)
	}
	if mb.Messages == nil {
		mb.Messages = make(map[string][]string)
	}
	mb.subscribers[queue] = append(mb.subscribers[queue], message)
	queued := append([]string(nil), mb.Messages[queue]...)
	mb.Messages[queue] = nil
	mb.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for _, msg := range queued {
			select {
			case <-ctx.Done():
				return
			case message <- msg:
			}
		}
		<-ctx.Done()
	}()
	msg := &MockSubscription{
		ctxCancel: cancel,
	}
	return msg, nil
}

func (mb *MockMessageBus) Close() error {
	// No-op for mock
	return nil
}

func (mb *MockMessageBus) SendMessageWithHeaders(message []byte, queue string, headers map[string]string) error {
	// For mock, just delegate to SendMessage
	return mb.SendMessage(message, queue)
}
