// Licensed to You under the Apache License, Version 2.0.

package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/messagebus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/service"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/pkg/wire"
)

// AuthClientInterface defines the minimal interface for authorization clients in iDRAC
type AuthClientInterface interface {
	AddService(service Service) error
	DeleteService(service Service) error
	GetService(ctx context.Context, services chan<- *Service)
	ResendAll()
	SplunkAddHEC(SplunkHttp SplunkConfig) error
	UpdateServiceState(state string, sip string) error
}

// AuthorizationClient handles authorization-related client operations
type AuthorizationClient struct {
	*service.BaseClient
}

// NewAuthorizationClient creates a new AuthorizationClient using BaseClient
func NewAuthorizationClient(bus messagebus.Messagebus) *AuthorizationClient {
	baseClient := service.NewBaseClient(bus, CommandQueue, ReplyQueuePrefix, "auth", ReadTimeout*time.Second)
	return &AuthorizationClient{
		BaseClient: baseClient,
	}
}

// GetHECConfig requests HEC configuration
func (ac *AuthorizationClient) GetHECConfig() {
	if err := ac.Send(GETHECCONFIG, nil); err != nil {
		log.Printf("Failed to send GetHECConfig command: %v", err)
	}
}

// ResendAll requests all services to be resent
func (ac *AuthorizationClient) ResendAll() {
	if err := ac.Send(RESEND, nil); err != nil {
		log.Printf("Failed to send resend command: %v", err)
	}
}

// SplunkAddHEC adds a Splunk HEC configuration
func (ac *AuthorizationClient) SplunkAddHEC(SplunkHttp SplunkConfig) error {
	return ac.Send(SPLUNKADDHEC, SplunkHttp)
}

// AddService adds a new service
func (ac *AuthorizationClient) AddService(service Service) error {
	return ac.Send(ADDSERVICE, service)
}

// DeleteService deletes a service
func (ac *AuthorizationClient) DeleteService(service Service) error {
	return ac.Send(DELETESERVICE, service)
}

// GetService listens for service events
func (ac *AuthorizationClient) GetService(ctx context.Context, services chan<- *Service) {
	envelopes := make(chan service.Envelope, 10)

	// Use the standardized filtered listening method from BaseClient with context
	ac.ListenToQueueFiltered(ctx, EventQueue, SERVICEEVENT, envelopes)

	for env := range envelopes {
		svc := new(Service)
		if err := wire.DecodePayload(env.Payload, svc); err != nil {
			log.Print("Error decoding service payload: ", err)
			continue
		}
		services <- svc
	}
}

// UpdateService updates a service
func (ac *AuthorizationClient) UpdateService(s Service) error {
	return ac.Send(UPDATESERVICE, s)
}

// UpdateServiceState updates a service's state
func (ac *AuthorizationClient) UpdateServiceState(state string, sip string) error {
	switch state {
	case CONNFAILED, STARTING, RUNNING, TELNOTFOUND, RUNNINGWOTEL, LEAKED, MONITORING:
		ac.UpdateService(
			Service{
				Ip:    sip,
				State: state,
			},
		)
	default:
		return fmt.Errorf("invalid state %s", state)
	}
	return nil
}

// UpdateValveState updates valve state for an IP
func (ac *AuthorizationClient) UpdateValveState(ip string, state1 string, state2 string) error {
	return ac.Send(UPDATEVALVESTATE, ValveState{
		Ip:      ip,
		VState1: state1,
		VState2: state2,
	})
}

// GetAllServices retrieves all configured services using request/reply.
func (ac *AuthorizationClient) GetAllServices() []Service {
	services := []Service{}
	err := ac.Call(GETALLSERVICES, nil, &services)
	if err != nil {
		log.Print("Error getting all services: ", err)
		return []Service{}
	}
	return services
}

// GetAllSystemTypes retrieves all configured system types using request/reply.
func (ac *AuthorizationClient) GetAllSystemTypes() []SystemType {
	systemtype := []SystemType{}
	err := ac.Call(GETSYSTEMTYPES, nil, &systemtype)
	if err != nil {
		log.Print("Error getting all services: ", err)
		return []SystemType{}
	}
	return systemtype
}

// GetValveStatus retrieves the current valve status using request/reply.
func (ac *AuthorizationClient) GetValveStatus() []ValveState {
	valvestatus := []ValveState{}
	err := ac.Call(GETVALVESTATE, nil, &valvestatus)
	if err != nil {
		log.Print("Error getting valve status: ", err)
		return []ValveState{}
	}
	log.Print("valvestatus: ", valvestatus)
	return valvestatus
}

// GetServiceWithIP retrieves one service by IP using request/reply.
func (ac *AuthorizationClient) GetServiceWithIP(ip string) Service {
	svc := Service{}
	err := ac.Call(GETSERVICE, Service{Ip: ip}, &svc)
	if err != nil {
		log.Print("Error getting service with ip: ", ip, " err: ", err)
		return Service{}
	}
	return svc
}

// AddServiceItem adds a service item
func (ac *AuthorizationClient) AddServiceItem(si ServiceItem) error {
	return ac.Send(ADDSERVICEITEM, si)
}

// DeleteServiceItem deletes a service item
func (ac *AuthorizationClient) DeleteServiceItem(si ServiceItem) error {
	return ac.Send(DELETESERVICEITEM, si)
}

// UpdateServiceX updates a service (extended)
func (ac *AuthorizationClient) UpdateServiceX(s Service) error {
	return ac.Send(UPDATESERVICEX, s)
}

// UpdateServiceXState updates a service's state (extended)
func (ac *AuthorizationClient) UpdateServiceXState(state string, ip string) error {
	return ac.UpdateServiceX(
		Service{
			Ip:    ip,
			State: state,
		},
	)
}

// UpdateServiceItem updates a service item
func (ac *AuthorizationClient) UpdateServiceItem(si ServiceItem) error {
	return ac.Send(UPDATESERVICEITEM, si)
}

// UpdateServiceItemState updates a service item's state
func (ac *AuthorizationClient) UpdateServiceItemState(state string, siip string) error {
	switch state {
	case SHUTDOWNSENT, SHUTDOWNFAILED, CONNFAILED, RUNNING, POWERSTATEON, POWERSTATEOFF:
		ac.UpdateServiceItem(
			ServiceItem{
				Service: Service{
					Ip:    siip,
					State: state,
				},
			},
		)
	default:
		return fmt.Errorf("invalid state %s", state)
	}
	return nil
}

// GetServiceItems retrieves the associated systems for a service IP using request/reply.
func (ac *AuthorizationClient) GetServiceItems(sip string) []ServiceItem {
	serviceItems := []ServiceItem{}
	err := ac.Call(GETSERVICEITEMS, ServiceItem{ServiceIP: sip}, &serviceItems)
	if err != nil {
		log.Print("Error reading service items: ", err)
		return nil
	}
	fmt.Println("Get associated systems", serviceItems)
	return serviceItems
}

// GetServiceItemWithIP retrieves exactly one service item by IP using request/reply.
func (ac *AuthorizationClient) GetServiceItemWithIP(siip string) ServiceItem {
	serviceItems := []ServiceItem{}
	err := ac.Call(GETSERVICEITEM, ServiceItem{Service: Service{Ip: siip}}, &serviceItems)
	if err != nil {
		log.Print("Error reading service items: ", err)
		return ServiceItem{}
	}
	if len(serviceItems) != 1 {
		log.Print("Error reading service items: Found multiple items for ip ", siip)
		return ServiceItem{}
	}
	return serviceItems[0]
}

// AddSystemType adds a system type
func (ac *AuthorizationClient) AddSystemType(sysType string) error {
	return ac.Send(ADDSYSTEMTYPE, SystemType(sysType))
}
