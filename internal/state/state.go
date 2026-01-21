package state

import (
	"time"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/auth"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/service"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

// State message types
const (
	UPDATESERVICESTATUS     = "updateservicestatus"
	UPDATESERVICEITEMSTATUS = "updateserviceitemstatus"
)

// State queue constants
const (
	CommandQueue     = "/status/command"
	EventQueue       = "/status"
	ReplyQueuePrefix = "/queue/state.reply."
)

type Envelope = wire.Envelope

type Command struct {
	Command string `json:"command"`
}

type StateService struct {
	*service.BaseService
}

// NewStateService creates a new StateService using BaseService
func NewStateService(bus messagebus.Messagebus) *StateService {
	baseService := service.NewBaseService(bus, CommandQueue)
	return &StateService{
		BaseService: baseService,
	}
}

type StateClient struct {
	*service.BaseClient
}

// NewStateClient creates a new StateClient using BaseClient
func NewStateClient(bus messagebus.Messagebus) *StateClient {
	baseClient := service.NewBaseClient(bus, CommandQueue, ReplyQueuePrefix, "state", 30*time.Second)
	return &StateClient{
		BaseClient: baseClient,
	}
}

func (s *StateClient) SendCommandString(command string) error {
	return s.Send(command, struct{}{})
}

func (s *StateClient) UpdateServiceStatus(service auth.Service) error {
	return s.Send(UPDATESERVICESTATUS, service)
}

func (s *StateClient) UpdateServiceItemStatus(si auth.ServiceItem) error {
	return s.Send(UPDATESERVICEITEMSTATUS, si)
}
