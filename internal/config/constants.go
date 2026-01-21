// Licensed to You under the Apache License, Version 2.0.

package config

// Config message types
const (
	GETPROPS = "getprops"
	GET      = "get"
	SET      = "set"
	RESET    = "reset"
)

// Config queue constants
const (
	DefaultCommandQueue = "/queue/config.command.v1"
	ReplyQueuePrefix    = "/queue/config.reply."
)
