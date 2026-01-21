// Licensed to You under the Apache License, Version 2.0.

package databus

import (
	"log"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/service"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

type DataBusService struct {
	*service.BaseService
	Recievers []string
}

// NewDataBusService creates a new DataBusService using BaseService
func NewDataBusService(bus messagebus.Messagebus) *DataBusService {
	baseService := service.NewBaseService(bus, CommandQueue)
	return &DataBusService{
		BaseService: baseService,
		Recievers:   make([]string, 0),
	}
}

func (d *DataBusService) Broadcast(msgType string, payload any) {
	env, err := wire.NewRequest(msgType, payload, "")
	if err != nil {
		log.Printf("Failed to create broadcast envelope: %v", err)
		return
	}
	data, err := wire.EncodeEnvelope(env)
	if err != nil {
		log.Printf("Failed to encode broadcast envelope: %v", err)
		return
	}
	for _, queue := range d.Recievers {
		err := d.Bus.SendMessage(data, queue)
		if err != nil {
			log.Printf("Failed to send broadcast to %s: %v", queue, err)
		}
	}
}

func (d *DataBusService) SendGroup(group DataGroup) {
	d.Broadcast(SUBSCRIBE, group)
}

func (d *DataBusService) SendGroupToQueue(group DataGroup, replyTo string, correlationID string) error {
	req := Envelope{
		ReplyTo:       replyTo,
		CorrelationID: correlationID,
	}
	return d.Reply(req, group)
}

func (d *DataBusService) SendProducersToQueue(producers []*DataProducer, replyTo string, correlationID string) error {
	req := Envelope{
		ReplyTo:       replyTo,
		CorrelationID: correlationID,
	}
	return d.Reply(req, producers)
}

func (d *DataBusService) ReceiveCommand(commands chan<- Envelope) error {
	allCommands := make(chan Envelope, 10)
	go d.ListenToQueue(CommandQueue, allCommands)

	for env := range allCommands {
		if env.Type == SUBSCRIBE {
			var subscribeReq struct {
				ReceiveQueue string `json:"ReceiveQueue"`
			}
			if err := wire.DecodePayload(env.Payload, &subscribeReq); err == nil {
				found := false
				for _, rec := range d.Recievers {
					if rec == subscribeReq.ReceiveQueue {
						found = true
						break
					}
				}
				if !found {
					d.Recievers = append(d.Recievers, subscribeReq.ReceiveQueue)
				}
			}
		} else {
			commands <- env
		}
	}
	return nil
}

// Reply sends a reply to a request (overrides BaseService.Reply for databus-specific handling)
func (d *DataBusService) Reply(req Envelope, payload any) error {
	if req.ReplyTo == "" {
		return nil
	}
	reply, err := wire.ReplyOK(req, payload)
	if err != nil {
		return err
	}
	return d.SendEnvelope(req.ReplyTo, reply)
}
