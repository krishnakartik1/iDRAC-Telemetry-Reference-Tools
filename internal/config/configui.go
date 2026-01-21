package config

import (
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/auth"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/databus"
)

type SystemHandler struct {
	AuthClient *auth.AuthorizationClient
	DataBus    *databus.DataBusClient
	ConfigBus  *ConfigClient
}

func NewSystemHandler(authClient *auth.AuthorizationClient, dataBus *databus.DataBusClient, configBus *ConfigClient) *SystemHandler {
	return &SystemHandler{
		AuthClient: authClient,
		DataBus:    dataBus,
		ConfigBus:  configBus,
	}
}
