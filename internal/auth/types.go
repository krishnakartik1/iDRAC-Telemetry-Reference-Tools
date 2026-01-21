// Licensed to You under the Apache License, Version 2.0.

package auth

import (
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/disc"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

// Service type constants (from disc package)
const (
	UNKNOWN = disc.UNKNOWN
	MSM     = disc.MSM
	EC      = disc.EC
	IDRAC   = disc.IDRAC
	IRC     = disc.IRC
	NVLINK  = disc.NVLINK
)

// Envelope is an alias for wire.Envelope
type Envelope = wire.Envelope

// Service represents a service with authentication information
type Service struct {
	ServiceType int               `json:"serviceType"`
	Ip          string            `json:"ip"`
	AuthType    int               `json:"authType"`
	Auth        map[string]string `json:"auth"`
	State       string            `json:"state"`
}

// ServiceItem represents an item associated with a service
type ServiceItem struct {
	Service
	ServiceIP          string `json:"serviceIp"`
	Systemtypedesc     string `json:"systemtypedesc"`
	IsForcefulshutdown bool   `json:"isforcefulshutdown"`
	Timeout            int    `json:"timeout"`
}

// SplunkConfig holds Splunk HEC configuration
type SplunkConfig struct {
	Url   string `json:"url,omitempty"`
	Key   string `json:"key,omitempty"`
	Index string `json:"index,omitempty"`
}

// SystemType represents a system type identifier
type SystemType string

// ValveState represents the state of valves for an IP
type ValveState struct {
	Ip      string `json:"ip"`
	VState1 string `json:"state1"`
	VState2 string `json:"state2"`
}
