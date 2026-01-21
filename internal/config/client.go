// Licensed to You under the Apache License, Version 2.0.

package config

import (
	"log"
	"time"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/service"
)

type ConfigClient struct {
	*service.BaseClient
	CommandQueue  string
	ResponseQueue string
}

func NewConfigClient(bus messagebus.Messagebus, commandQueue string, responseQueue string) *ConfigClient {
	baseClient := service.NewBaseClient(bus, commandQueue, "/queue/config.reply.", "config", 30*time.Second)
	return &ConfigClient{
		BaseClient:    baseClient,
		CommandQueue:  commandQueue,
		ResponseQueue: responseQueue,
	}
}

// NewConfigClientWithBus creates a new ConfigClient with the provided message bus
func NewConfigClientWithBus(bus messagebus.Messagebus) *ConfigClient {
	return NewConfigClient(bus, "/queue/config.command.v1", "/queue/config.reply")
}

func (d *ConfigClient) ReadOneMessage(queue string) string {
	messages := make(chan string, 1)
	sub, err := d.Bus.ReceiveMessage(messages, queue)
	if err != nil {
		log.Printf("Error receiving message: %v", err)
		return ""
	}
	message := <-messages
	sub.Close()
	return message
}

func (d *ConfigClient) GetProperties() ([]string, error) {
	var props []string
	err := d.Call(GETPROPS, struct{}{}, &props)
	return props, err
}

func (d *ConfigClient) Get(name string) (interface{}, error) {
	payload := struct {
		Property string `json:"property"`
	}{Property: name}

	var response struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}
	err := d.Call(GET, payload, &response)
	if err != nil {
		return nil, err
	}
	return response.Value, nil
}

func (d *ConfigClient) Set(name string, value interface{}) error {
	payload := struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}{Property: name, Value: value}

	var response struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}
	return d.Call(SET, payload, &response)
}

func (d *ConfigClient) Reset(name string) error {
	payload := struct {
		Property string `json:"property"`
	}{Property: name}

	var response struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}
	return d.Call(RESET, payload, &response)
}

// Backward-compatible wrapper methods for old API signatures
func (d *ConfigClient) GetOld(name string) (*Response, error) {
	value, err := d.Get(name)
	if err != nil {
		return &Response{
			Command:  GET,
			Property: name,
			Error:    err,
		}, err
	}
	return &Response{
		Command:  GET,
		Property: name,
		Value:    value,
	}, nil
}

func (d *ConfigClient) SetOld(name string, value interface{}) (*Response, error) {
	err := d.Set(name, value)
	if err != nil {
		return &Response{
			Command:  SET,
			Property: name,
			Value:    value,
			Error:    err,
		}, err
	}
	return &Response{
		Command:  SET,
		Property: name,
		Value:    value,
	}, nil
}
