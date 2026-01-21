// Licensed to You under the Apache License, Version 2.0.

package auth

import (
	"fmt"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/service"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

// AuthorizationService handles authorization-related service operations
type AuthorizationService struct {
	*service.BaseService
}

// NewAuthorizationService creates a new AuthorizationService using BaseService
func NewAuthorizationService(bus messagebus.Messagebus) *AuthorizationService {
	baseService := service.NewBaseService(bus, CommandQueue)
	return &AuthorizationService{
		BaseService: baseService,
	}
}

// SendWithQueue sends a message to a specific queue
func (as *AuthorizationService) SendWithQueue(msgType string, payload any, queue string) error {
	env, err := wire.NewEnvelope(msgType, payload)
	if err != nil {
		return err
	}
	return as.SendEnvelope(queue, env)
}

// SendEvent sends an event to the event queue
func (as *AuthorizationService) SendEvent(msgType string, payload any) error {
	return as.SendWithQueue(msgType, payload, EventQueue)
}

// Reply sends a reply to a request
func (as *AuthorizationService) Reply(req Envelope, payload any) error {
	if req.ReplyTo == "" {
		return fmt.Errorf("reply queue missing")
	}
	env, err := wire.ReplyOK(req, payload)
	if err != nil {
		return err
	}
	return as.SendEnvelope(req.ReplyTo, env)
}

// ReplyError sends an error reply to a request
func (as *AuthorizationService) ReplyError(req Envelope, replyErr error) error {
	if req.ReplyTo == "" {
		return fmt.Errorf("reply queue missing")
	}
	env := wire.ReplyErr(req, replyErr)
	return as.SendEnvelope(req.ReplyTo, env)
}

// SendService sends a service event
func (as *AuthorizationService) SendService(service Service) error {
	return as.SendEvent(SERVICEEVENT, service)
}

// SendServiceWithQ sends a service event to a specific queue
func (as *AuthorizationService) SendServiceWithQ(service Service, queue string) error {
	return as.SendWithQueue(SERVICEEVENT, service, queue)
}

// Sendconfig sends a Splunk configuration event
func (as *AuthorizationService) Sendconfig(config SplunkConfig) error {
	return as.SendEvent(SPLUNKCONFIGEVENT, config)
}
