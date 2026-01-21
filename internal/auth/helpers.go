// Licensed to You under the Apache License, Version 2.0.

package auth

import (
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
)

// NewEnvelopeChan creates a new channel for receiving envelopes
func NewEnvelopeChan() chan Envelope {
	return make(chan Envelope, 10)
}

// NewAuthServiceWithBus is an alias for NewAuthorizationService for compatibility
func NewAuthServiceWithBus(bus messagebus.Messagebus) *AuthorizationService {
	return NewAuthorizationService(bus)
}

// NewAuthClientWithBus is an alias for NewAuthorizationClient for compatibility
func NewAuthClientWithBus(bus messagebus.Messagebus) *AuthorizationClient {
	return NewAuthorizationClient(bus)
}
