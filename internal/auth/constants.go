// Licensed to You under the Apache License, Version 2.0.

package auth

// Authentication type constants
const (
	AuthTypeUsernamePassword = 1
	AuthTypeXAuthToken       = 2
	AuthTypeBearerToken      = 3
	ReadTimeout              = 5
)

// Service states
const (
	STARTING     = "Starting"
	RUNNING      = "Running"
	RUNNINGWOTEL = "Running Only Alerts"
	TELNOTFOUND  = "Telemetry Service Not Found"
	CONNFAILED   = "Connection Failed"
	LEAKED       = "Leak Detected"
	MONITORING   = "Monitoring"
)

// ServiceItem states
const (
	SHUTDOWNSENT   = "Shutdown Sent"
	SHUTDOWNFAILED = "Shutdown Failed"
	POWERSTATEON   = "Power Status On"
	POWERSTATEOFF  = "Power Status Off"
)

// Message type constants
const (
	RESEND        = "auth.resend"
	ADDSERVICE    = "auth.add_service"
	DELETESERVICE = "auth.delete_service"
	UPDATESERVICE = "auth.update_service"
	GETSERVICE    = "auth.get_service"
	TERMINATE     = "auth.terminate"
	SPLUNKADDHEC  = "auth.splunk_add_hec"
	GETHECCONFIG  = "auth.get_hec_config"

	GETALLSERVICES    = "auth.get_all_services"
	ADDSERVICEITEM    = "auth.add_service_item"
	DELETESERVICEITEM = "auth.delete_service_item"
	GETSERVICEITEMS   = "auth.get_service_items"
	GETSERVICEITEM    = "auth.get_service_item"
	UPDATESERVICEITEM = "auth.update_service_item"
	UPDATEVALVESTATE  = "auth.update_valve_state"
	GETVALVESTATE     = "auth.get_valve_state"
	ADDSYSTEMTYPE     = "auth.add_system_type"
	GETSYSTEMTYPES    = "auth.get_system_types"

	UPDATESERVICEX    = "auth.update_service_x"
	SERVICEEVENT      = "auth.service_event"
	SPLUNKCONFIGEVENT = "auth.splunk_config_event"
)

// Queue constants
const (
	CommandQueue     = "/queue/auth.command.v1"
	EventQueue       = "/topic/auth.events.v1"
	ReplyQueuePrefix = "/queue/auth.reply."
)
