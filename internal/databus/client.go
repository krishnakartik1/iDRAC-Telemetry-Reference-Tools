// Licensed to You under the Apache License, Version 2.0.

package databus

import (
	"context"
	"log"
	"time"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/service"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

type DataBusClient struct {
	*service.BaseClient
}

// NewDataBusClient creates a new DataBusClient using BaseClient
func NewDataBusClient(bus messagebus.Messagebus) *DataBusClient {
	baseClient := service.NewBaseClient(bus, CommandQueue, ReplyQueuePrefix, "databus", 30*time.Second)
	return &DataBusClient{
		BaseClient: baseClient,
	}
}

// NewDataBusClientWithBus creates a new DataBusClient with the provided message bus
func NewDataBusClientWithBus(bus messagebus.Messagebus) *DataBusClient {
	return NewDataBusClient(bus)
}

func (d *DataBusClient) Get(queue string) error {
	payload := struct {
		ReceiveQueue string `json:"ReceiveQueue"`
	}{ReceiveQueue: queue}
	return d.BaseClient.Send(GET, payload)
}

func (d *DataBusClient) Subscribe(queue string) error {
	payload := struct {
		ReceiveQueue string `json:"ReceiveQueue"`
	}{ReceiveQueue: queue}
	return d.Send(SUBSCRIBE, payload)
}

func (d *DataBusClient) ReadOneMessage(queue string) string {
	messages := make(chan string)
	sub, err := d.Bus.ReceiveMessage(messages, queue)
	if err != nil {
		log.Println("Error receiving message: ", err)
		return ""
	}
	message := <-messages
	sub.Close()
	return message
}

func (d *DataBusClient) GetResponse(queue string) (Envelope, error) {
	message := d.ReadOneMessage(queue)
	env, err := wire.DecodeEnvelope([]byte(message))
	if err != nil {
		log.Printf("Error reading response queue: %v", err)
		return Envelope{}, err
	}
	return env, nil
}

func (d *DataBusClient) GetProducers(queue string) ([]DataProducer, error) {
	payload := struct {
		ReceiveQueue string `json:"ReceiveQueue"`
	}{ReceiveQueue: queue}

	var producers []DataProducer
	err := d.Call(GETPRODUCERS, payload, &producers)
	if err != nil {
		return nil, err
	}
	return producers, nil
}

func (d *DataBusClient) DeleteProducer(queue string, serviceIP string) error {
	payload := struct {
		ReceiveQueue string `json:"ReceiveQueue"`
		ServiceIP    string `json:"serviceIP"`
	}{ReceiveQueue: queue, ServiceIP: serviceIP}
	return d.Send(DELETEPRODUCER, payload)
}

func (d *DataBusClient) GetGroup(ctx context.Context, groups chan<- *DataGroup, queue string) {
	subscribeEnvelopes := make(chan service.Envelope, 10)
	getEnvelopes := make(chan service.Envelope, 10)

	d.ListenToQueueFiltered(ctx, queue, SUBSCRIBE, subscribeEnvelopes)
	d.ListenToQueueFiltered(ctx, queue, GET, getEnvelopes)

	go func() {
		for env := range subscribeEnvelopes {
			group := new(DataGroup)
			if err := wire.DecodePayload(env.Payload, group); err == nil {
				groups <- group
			} else {
				log.Printf("Error decoding DataGroup payload from SUBSCRIBE: %v", err)
			}
		}
	}()

	for env := range getEnvelopes {
		group := new(DataGroup)
		if err := wire.DecodePayload(env.Payload, group); err == nil {
			groups <- group
		} else {
			log.Printf("Error decoding DataGroup payload from GET: %v", err)
		}
	}
}
