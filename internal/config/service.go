// Licensed to You under the Apache License, Version 2.0.

package config

import (
	"fmt"
	"log"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/service"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

type ConfigService struct {
	*service.BaseService
	Entries map[string]*ConfigEntry
}

func NewConfigService(bus messagebus.Messagebus, commandQueue string, entries map[string]*ConfigEntry) *ConfigService {
	baseService := service.NewBaseService(bus, commandQueue)
	return &ConfigService{
		BaseService: baseService,
		Entries:     entries,
	}
}

func (d *ConfigService) Run() {
	commands := make(chan Envelope, 10)
	go d.ReceiveCommand(commands)

	for {
		env := <-commands
		switch env.Type {
		default:
			log.Print("Received unknown config command: ", env.Type)
		case GETPROPS:
			d.GetProperties(env)
		case GET:
			d.Get(env)
		case SET:
			d.Set(env)
		case RESET:
			d.Reset(env)
		}
	}
}

func (d *ConfigService) GetProperties(env Envelope) {
	keys := make([]string, 0, len(d.Entries))
	for k := range d.Entries {
		keys = append(keys, k)
	}
	if env.ReplyTo != "" {
		err := d.Reply(env, keys)
		if err != nil {
			log.Printf("Failed to send properties reply: %v", err)
		}
	}
}

func (d *ConfigService) Get(env Envelope) {
	var getReq struct {
		Property string `json:"property"`
	}
	if err := wire.DecodePayload(env.Payload, &getReq); err != nil {
		if env.ReplyTo != "" {
			d.ReplyError(env, fmt.Errorf("Invalid get request payload"))
		}
		return
	}

	entry, ok := d.Entries[getReq.Property]
	if !ok {
		if env.ReplyTo != "" {
			d.ReplyError(env, fmt.Errorf("Could not find property named %s", getReq.Property))
		}
		return
	}

	value, err := entry.Get(getReq.Property)
	if err != nil {
		if env.ReplyTo != "" {
			d.ReplyError(env, err)
		}
		return
	}

	response := struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}{Property: getReq.Property, Value: value}

	if env.ReplyTo != "" {
		err := d.Reply(env, response)
		if err != nil {
			log.Printf("Error sending reply: %v", err)
		}
	}
}

func (d *ConfigService) Set(env Envelope) {
	var setReq struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}
	if err := wire.DecodePayload(env.Payload, &setReq); err != nil {
		if env.ReplyTo != "" {
			d.ReplyError(env, fmt.Errorf("Invalid set request payload"))
		}
		return
	}

	entry, ok := d.Entries[setReq.Property]
	if !ok {
		if env.ReplyTo != "" {
			d.ReplyError(env, fmt.Errorf("Could not find property named %s", setReq.Property))
		}
		return
	}

	err := entry.Set(setReq.Property, setReq.Value)
	if err != nil {
		if env.ReplyTo != "" {
			d.ReplyError(env, err)
		}
		return
	}

	response := struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}{Property: setReq.Property, Value: setReq.Value}

	if env.ReplyTo != "" {
		err := d.Reply(env, response)
		if err != nil {
			log.Printf("Error sending reply: %v", err)
		}
	}
}

func (d *ConfigService) Reset(env Envelope) {
	var resetReq struct {
		Property string `json:"property"`
	}
	if err := wire.DecodePayload(env.Payload, &resetReq); err != nil {
		if env.ReplyTo != "" {
			d.BaseService.ReplyError(env, fmt.Errorf("Invalid reset request payload"))
		}
		return
	}

	entry, ok := d.Entries[resetReq.Property]
	if !ok {
		if env.ReplyTo != "" {
			d.BaseService.ReplyError(env, fmt.Errorf("Could not find property named %s", resetReq.Property))
		}
		return
	}

	err := entry.Set(resetReq.Property, entry.Default)
	if err != nil {
		if env.ReplyTo != "" {
			d.BaseService.ReplyError(env, err)
		}
		return
	}

	response := struct {
		Property string      `json:"property"`
		Value    interface{} `json:"value"`
	}{Property: resetReq.Property, Value: entry.Default}

	if env.ReplyTo != "" {
		err := d.BaseService.Reply(env, response)
		if err != nil {
			log.Printf("Error sending reply: %v", err)
		}
	}
}

// Reply sends a reply to a request
func (d *ConfigService) Reply(req Envelope, payload any) error {
	if req.ReplyTo == "" {
		return nil
	}
	reply, err := wire.ReplyOK(req, payload)
	if err != nil {
		return err
	}
	return d.SendEnvelope(req.ReplyTo, reply)
}

// ReplyError sends an error reply to a request
func (d *ConfigService) ReplyError(req Envelope, replyErr error) error {
	if req.ReplyTo == "" {
		return nil
	}
	reply := wire.ReplyErr(req, replyErr)
	return d.SendEnvelope(req.ReplyTo, reply)
}
