package main

import (
	"context"
	// "encoding/csv"
	// "fmt"
	"log"
	// "net/http"
	// "os"
	// "path/filepath"
	"strconv"
	"time"

	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/auth"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/config"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/databus"
	"github.com/dell/iDRAC-Telemetry-Reference-Tools/internal/messagebus/stomp"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var configStrings = map[string]string{
	"mbhost":   "activemq",
	"mbip":     "172.18.0.2",
	"mbport":   "61613",
	"httpport": "8082",
}

type SystemHandler struct {
	AuthClient *auth.AuthorizationClient
	DataBus    *databus.DataBusClient
	ConfigBus  *config.ConfigClient
}

func splunkAddContainer(cli *client.Client, ctx context.Context, SplunkConfig auth.SplunkConfig) error {
	env := []string{
		"MESSAGEBUS_HOST: " + configStrings["mbhost"],
		"MESSAGEBUS_PORT: " + configStrings["mbport"],
		"SPLUNK_HEC_URL: " + SplunkConfig.Url,
		"SPLUNK_HEC_KEY: " + SplunkConfig.Key,
		"SPLUNK_HEC_INDEX: " + SplunkConfig.Index,
	}

	imageName := "idrac-telemetry-reference-tools/splunkpump:latest"

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
	}, nil, nil, nil, "idrac-telemetry-reference-tools-splunkpump")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return err
}

func main() {

	systemHandler := new(SystemHandler)
	systemHandler.AuthClient = new(auth.AuthorizationClient)
	systemHandler.DataBus = new(databus.DataBusClient)
	systemHandler.ConfigBus = new(config.ConfigClient)

	//Initialize messagebus
	for {
		stompPort, _ := strconv.Atoi(configStrings["mbport"])
		mb, err := stomp.NewStompMessageBus(configStrings["mbip"], stompPort)
		if err != nil {
			log.Printf("Could not connect to message bus: %s", err)
			time.Sleep(5 * time.Second)
		} else {
			systemHandler.AuthClient.Bus = mb
			systemHandler.DataBus.Bus = mb
			systemHandler.ConfigBus.Bus = mb
			defer mb.Close()
			break
		}
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	//Setu authorization service
	authorizationService := new(auth.AuthorizationService)

	commands := make(chan *auth.Command)
	go authorizationService.ReceiveCommand(commands) //nolint: errcheck
	for {
		command := <-commands
		log.Printf("Received command in dockermanage: %s", command.Command)
		switch command.Command {
		case auth.SPLUNKADDCONTAINER:
			err := splunkAddContainer(cli, ctx, command.SplunkConfig)
			if err != nil {
				log.Print("Failed to add container: ", err)
				break
			}
		}
	}
}
