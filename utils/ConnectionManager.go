package utils

import (
	"time"
	// "sync"

	// "github.com/gotk3/gotk3/glib"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/logger"
)

const MAX_CONNECTION_ATTEMPTS = 10

type connectionManager struct {
	Client						*octoprintApis.Client
	ConnectAttempts				int
	IsConnectedToOctoPrint		bool
	IsConnectedToPrinter		bool
}

var connectionManagerInstance *connectionManager

func GetConnectionManagerInstance(client *octoprintApis.Client) (*connectionManager) {
	if connectionManagerInstance == nil {
		if client == nil {
			panic("GetConnectionManagerInstance() was called for the first time, but the client passed was nil")
		}

		connectionManagerInstance = &connectionManager{
			Client: client,
			ConnectAttempts: 0,
			IsConnectedToOctoPrint: false,
			IsConnectedToPrinter: false,
		}
	}

	return connectionManagerInstance
}

func (this *connectionManager) ReInitializeConnectionState() {
	this.ConnectAttempts = 0
	this.IsConnectedToOctoPrint = false
	this.IsConnectedToPrinter = false
}

func (this *connectionManager) UpdateStatus() {
	logger.TraceEnter("ConnectionManager.UpdateStatus()")

	logger.Infof("ConnectAttempts: %d", this.ConnectAttempts)

	// If OctoScreen is connected to OctoPrint,
	// and OctoPrint is connected to the printer,
	// don't bother checking again.
	if this.IsConnected() == true {
		logger.TraceLeave("ConnectionManager.UpdateStatus()")
		return
	}

	// Continue on if OctoScreen isn't connected...

	// If the maximum number of attempts have already been made, don't bother trying agin.
	if this.ConnectAttempts >= MAX_CONNECTION_ATTEMPTS {
		logger.TraceLeave("ConnectionManager.UpdateStatus()")
		return
	}

	this.ConnectAttempts++

	logger.Debug("ConnectionManager.UpdateStatus() - about to call ConnectionRequest.Do()")
	t1 := time.Now()
	connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(this.Client)
	t2 := time.Now()
	logger.Debug("ConnectionManager.UpdateStatus() - finished calling ConnectionRequest.Do()")
	logger.Debugf("time elapsed: %q", t2.Sub(t1))
	
	if err != nil {
		logger.LogError("ConnectionManager.UpdateStatus()", "ConnectionRequest.Do()", err)
							// newUIState, splashMessage = this.getUiStateAndMessageFromError(err, newUIState, splashMessage)
							// logger.Debugf("ConnectionManager.UpdateStatus() - newUIState is now: %s", newUIState)
		this.IsConnectedToOctoPrint = false
		logger.Debug("ConnectionManager.UpdateStatus() - Connection state: IsConnectedToOctoPrint is now false")
		logger.TraceLeave("ConnectionManager.UpdateStatus()")
		return
	}
	
	logger.Debug("ConnectionManager.UpdateStatus() - ConnectionRequest.Do() succeeded")
	
	this.IsConnectedToOctoPrint = true
	
	jsonResponse, err := StructToJson(connectionResponse)
	if err != nil {
		logger.LogError("ConnectionManager.UpdateStatus()", "StructToJson()", err)
		// If there's an error here, it's with the serialization of the object to JSON.
		// This is just for debugging, so don't return if there's an issue, and just
		// carry on (and hopefully connectionResponse isn't corrupted)
	} else {
		logger.Debugf("ConnectionManager.UpdateStatus() - connectionResponse is: %s", jsonResponse)
	}
	
	/*
	Example JSON response:
	{
		"Current": {
			"state": "Operational",
			"port": "/dev/ttyACM0",
			"baudrate": 115200,
			"printerProfile": "_default"
		},
		"Options": {
			"ports": [
				"/dev/ttyACM0"
			],
			"baudrates": [
				250000,
				230400,
				115200,
				57600,
				38400,
				19200,
				9600
			],
			"printerProfiles": [
				{
					"id": "_default",
					"name": "name-of-the-printer"
				}
			],
			"portPreference": "",
			"baudratePreference": 0,
			"printerProfilePreference": "_default",
			"autoconnect": false
		}
	}
	*/
	
	printerConnectionState := connectionResponse.Current.State
	if printerConnectionState.IsOffline() || printerConnectionState.IsError() {
		this.IsConnectedToPrinter = false
	} else {
		this.IsConnectedToPrinter = true
	}

	logger.TraceLeave("ConnectionManager.UpdateStatus()")
}

func (this *connectionManager) IsConnected() bool {
	// TODO: should this be named IsFullyConnected()?
	return this.IsConnectedToOctoPrint == true && this.IsConnectedToPrinter == true;
}
