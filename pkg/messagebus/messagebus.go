// Licensed to You under the Apache License, Version 2.0.

package messagebus

type Subscription interface {
	Close() error
}

type Messagebus interface {
	SendMessage(message []byte, queue string) error
	SendMessageWithHeaders(message []byte, queue string, headers map[string]string) error
	ReceiveMessage(message chan<- string, queue string) (Subscription, error)
	Close() error
}

// MessagebusFactory is a function type for creating Messagebus instances
type MessagebusFactory func(host string, port int) (Messagebus, error)

// DefaultFactory holds the default messagebus factory (set by stomp package init or manually)
var DefaultFactory MessagebusFactory
