// Licensed to You under the Apache License, Version 2.0.

package config

import (
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

// Envelope is an alias for wire.Envelope
type Envelope = wire.Envelope

type SetFunc func(name string, value interface{}) error
type GetFunc func(name string) (interface{}, error)

type ConfigEntry struct {
	Set     SetFunc
	Get     GetFunc
	Default interface{}
}

type Command struct {
	Command       string      `json:"command"`
	ResponseQueue string      `json:"ReceiveQueue"`
	Property      string      `json:"property,omitempty"`
	Value         interface{} `json:"value,omitempty"`
}

type Response struct {
	Command  string      `json:"command"`
	Property string      `json:"property,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	Error    error       `json:"error,omitempty"`
}
