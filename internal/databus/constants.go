// Licensed to You under the Apache License, Version 2.0.

package databus

// DataBus message types
const (
	GET            = "get"
	SUBSCRIBE      = "subscribe"
	GETPRODUCERS   = "getproducers"
	DELETEPRODUCER = "deleteproducers"
	TERMINATE      = "terminate"
)

// DataBus queue constants
const (
	CommandQueue     = "/databus"
	ReplyQueuePrefix = "/queue/databus.reply."
)
