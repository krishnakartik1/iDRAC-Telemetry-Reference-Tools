// Licensed to You under the Apache License, Version 2.0.

package databus

import (
	"time"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

// Envelope is an alias for wire.Envelope
type Envelope = wire.Envelope

type DataValue struct {
	ID        string
	Context   string
	Label     string
	Value     string
	System    string
	HostName  string
	Timestamp string
}

type EventValue struct {
	EventType         string
	EventId           string
	EventTimestamp    string
	MemberId          string
	MessageSeverity   string
	Message           string
	MessageId         string
	MessageArgs       []string
	OriginOfCondition string
}

type DataGroup struct {
	HostID    string
	ID        string
	Label     string
	Sequence  string
	System    string
	HostName  string
	Model     string
	SKU       string
	FQDN      string
	FwVer     string
	ImgID     string
	Timestamp string
	Values    []DataValue
	Events    []EventValue
}

type DataProducer struct {
	Hostname  string
	Username  string
	State     string
	LastEvent time.Time
}

type Command struct {
	Command      string `json:"command"`
	ReceiveQueue string `json:"ReceiveQueue"`
	ReportData   string `json:"reportdata,omitempty"`
	ServiceIP    string `json:"serviceIP,omitempty"`
}

type Response struct {
	Command  string      `json:"command"`
	DataType string      `json:"dataType"`
	Data     interface{} `json:"data"`
}
