package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/mock"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

func receiveMessage(t *testing.T, ch <-chan string) string {
	t.Helper()
	select {
	case msg := <-ch:
		return msg
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for message")
		return ""
	}
}

func TestAddServiceItemSendsEnvelope(t *testing.T) {
	mb := mock.NewMockMessageBus()
	authClient := NewAuthorizationClient(mb)
	commandCh := make(chan string, 1)
	sub, err := mb.ReceiveMessage(commandCh, CommandQueue)
	if err != nil {
		t.Fatalf("ReceiveMessage() error = %v", err)
	}
	defer sub.Close()

	expected := ServiceItem{
		Service:   Service{Ip: "foo"},
		ServiceIP: "sip1",
	}
	if err := authClient.AddServiceItem(expected); err != nil {
		t.Fatalf("AddServiceItem() error = %v", err)
	}

	msg := receiveMessage(t, commandCh)
	env, err := wire.DecodeEnvelope([]byte(msg))
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	if env.Type != ADDSERVICEITEM {
		t.Fatalf("env.Type = %q, want %q", env.Type, ADDSERVICEITEM)
	}
	var got ServiceItem
	if err := wire.DecodePayload(env.Payload, &got); err != nil {
		t.Fatalf("DecodePayload() error = %v", err)
	}
	if got.ServiceIP != expected.ServiceIP || got.Ip != expected.Ip {
		t.Fatalf("payload = %+v, want %+v", got, expected)
	}
}

func TestGetServiceItemsCall(t *testing.T) {
	mb := mock.NewMockMessageBus()
	authClient := NewAuthorizationClient(mb)
	commandCh := make(chan string, 1)
	sub, err := mb.ReceiveMessage(commandCh, CommandQueue)
	if err != nil {
		t.Fatalf("ReceiveMessage() error = %v", err)
	}
	defer sub.Close()

	expectedItems := []ServiceItem{
		{Service: Service{Ip: "foo"}, ServiceIP: "sip1"},
		{Service: Service{Ip: "bar"}, ServiceIP: "sip2"},
	}
	resultCh := make(chan []ServiceItem, 1)
	go func() {
		resultCh <- authClient.GetServiceItems("sip1")
	}()

	msg := receiveMessage(t, commandCh)
	req, err := wire.DecodeEnvelope([]byte(msg))
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	if req.Type != GETSERVICEITEMS {
		t.Fatalf("req.Type = %q, want %q", req.Type, GETSERVICEITEMS)
	}
	var payload ServiceItem
	if err := wire.DecodePayload(req.Payload, &payload); err != nil {
		t.Fatalf("DecodePayload() error = %v", err)
	}
	if payload.ServiceIP != "sip1" {
		t.Fatalf("payload.ServiceIP = %q, want %q", payload.ServiceIP, "sip1")
	}

	if req.ReplyTo == "" {
		t.Fatalf("req.ReplyTo is empty")
	}
	reply, err := wire.ReplyOK(req, expectedItems)
	if err != nil {
		t.Fatalf("ReplyOK() error = %v", err)
	}
	data, err := wire.EncodeEnvelope(reply)
	if err != nil {
		t.Fatalf("EncodeEnvelope() error = %v", err)
	}
	if err := mb.SendMessage(data, req.ReplyTo); err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}

	select {
	case got := <-resultCh:
		if len(got) != len(expectedItems) {
			t.Fatalf("GetServiceItems() len = %d, want %d", len(got), len(expectedItems))
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for GetServiceItems response")
	}
}

func TestCallReturnsErrorStatus(t *testing.T) {
	mb := mock.NewMockMessageBus()
	authClient := NewAuthorizationClient(mb)
	commandCh := make(chan string, 1)
	sub, err := mb.ReceiveMessage(commandCh, CommandQueue)
	if err != nil {
		t.Fatalf("ReceiveMessage() error = %v", err)
	}
	defer sub.Close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- authClient.Call("auth.test", nil, &struct{}{})
	}()

	msg := receiveMessage(t, commandCh)
	req, err := wire.DecodeEnvelope([]byte(msg))
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	reply := wire.ReplyErr(req, errors.New("boom"))
	data, err := wire.EncodeEnvelope(reply)
	if err != nil {
		t.Fatalf("EncodeEnvelope() error = %v", err)
	}
	if err := mb.SendMessage(data, req.ReplyTo); err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}

	select {
	case err := <-errCh:
		if err == nil || err.Error() != "boom" {
			t.Fatalf("Call() error = %v, want boom", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for Call error")
	}
}

func TestGetServiceReceivesEvent(t *testing.T) {
	mb := mock.NewMockMessageBus()
	authClient := NewAuthorizationClient(mb)
	serviceCh := make(chan *Service, 1)
	ctx := context.Background()
	go authClient.GetService(ctx, serviceCh)

	expected := Service{Ip: "10.0.0.1"}
	env, err := wire.NewEnvelope(SERVICEEVENT, expected)
	if err != nil {
		t.Fatalf("NewEnvelope() error = %v", err)
	}
	data, err := wire.EncodeEnvelope(env)
	if err != nil {
		t.Fatalf("EncodeEnvelope() error = %v", err)
	}
	if err := mb.SendMessage(data, EventQueue); err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}

	select {
	case got := <-serviceCh:
		if got == nil || got.Ip != expected.Ip {
			t.Fatalf("GetService() = %+v, want %+v", got, expected)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for service event")
	}
}

func TestAuthorizationServiceSendEvent(t *testing.T) {
	mb := mock.NewMockMessageBus()
	authService := NewAuthorizationService(mb)
	messageCh := make(chan string, 1)
	sub, err := mb.ReceiveMessage(messageCh, EventQueue)
	if err != nil {
		t.Fatalf("ReceiveMessage() error = %v", err)
	}
	defer sub.Close()

	expected := Service{Ip: "10.0.0.2"}
	if err := authService.SendEvent(SERVICEEVENT, expected); err != nil {
		t.Fatalf("SendEvent() error = %v", err)
	}

	msg := receiveMessage(t, messageCh)
	env, err := wire.DecodeEnvelope([]byte(msg))
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	if env.Type != SERVICEEVENT {
		t.Fatalf("env.Type = %q, want %q", env.Type, SERVICEEVENT)
	}
	var got Service
	if err := wire.DecodePayload(env.Payload, &got); err != nil {
		t.Fatalf("DecodePayload() error = %v", err)
	}
	if got.Ip != expected.Ip {
		t.Fatalf("payload.Ip = %q, want %q", got.Ip, expected.Ip)
	}
}

func TestAuthorizationServiceReply(t *testing.T) {
	mb := mock.NewMockMessageBus()
	authService := NewAuthorizationService(mb)
	replyQueue := "/queue/reply.test"
	messageCh := make(chan string, 1)
	sub, err := mb.ReceiveMessage(messageCh, replyQueue)
	if err != nil {
		t.Fatalf("ReceiveMessage() error = %v", err)
	}
	defer sub.Close()

	req := Envelope{Type: GETSERVICE, CorrelationID: "corr-1", ReplyTo: replyQueue}
	expected := Service{Ip: "10.0.0.3"}
	if err := authService.Reply(req, expected); err != nil {
		t.Fatalf("Reply() error = %v", err)
	}

	msg := receiveMessage(t, messageCh)
	env, err := wire.DecodeEnvelope([]byte(msg))
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	if env.Status != wire.StatusOK {
		t.Fatalf("env.Status = %q, want %q", env.Status, wire.StatusOK)
	}
	if env.CorrelationID != req.CorrelationID {
		t.Fatalf("env.CorrelationID = %q, want %q", env.CorrelationID, req.CorrelationID)
	}
	var got Service
	if err := wire.DecodePayload(env.Payload, &got); err != nil {
		t.Fatalf("DecodePayload() error = %v", err)
	}
	if got.Ip != expected.Ip {
		t.Fatalf("payload.Ip = %q, want %q", got.Ip, expected.Ip)
	}
}

func TestAuthorizationServiceReplyError(t *testing.T) {
	mb := mock.NewMockMessageBus()
	authService := NewAuthorizationService(mb)
	replyQueue := "/queue/reply.error"
	messageCh := make(chan string, 1)
	sub, err := mb.ReceiveMessage(messageCh, replyQueue)
	if err != nil {
		t.Fatalf("ReceiveMessage() error = %v", err)
	}
	defer sub.Close()

	req := Envelope{Type: GETSERVICE, CorrelationID: "corr-2", ReplyTo: replyQueue}
	if err := authService.ReplyError(req, errors.New("nope")); err != nil {
		t.Fatalf("ReplyError() error = %v", err)
	}

	msg := receiveMessage(t, messageCh)
	env, err := wire.DecodeEnvelope([]byte(msg))
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	if env.Status != wire.StatusError {
		t.Fatalf("env.Status = %q, want %q", env.Status, wire.StatusError)
	}
	if env.Error != "nope" {
		t.Fatalf("env.Error = %q, want %q", env.Error, "nope")
	}
}
