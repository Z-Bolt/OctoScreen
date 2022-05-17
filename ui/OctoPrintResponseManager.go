package ui

import (
	// "os"
	// "time"
	// "strconv"
	// "sync"

	// "github.com/gotk3/gotk3/glib"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

// const MAX_CONNECTION_ATTEMPTS = 8
// const MAX_CONNECTION_ATTEMPTS = 4

type octoPrintResponseManager struct {
	UI							*UI
	//Client						*octoprintApis.Client
	// // IsRunning					bool
	// ConnectAttempts				int
	// IsConnectedToOctoPrint		bool
	// IsConnectedToPrinter		bool
	backgroundTask				*utils.BackgroundTask
	FullStateResponse			dataModels.FullStateResponse
}

var octoPrintResponseManagerInstance *octoPrintResponseManager

//func GetOctoPrintResponseManagerInstance(UI *UI, client *octoprintApis.Client) (*octoPrintResponseManager) {
func GetOctoPrintResponseManagerInstance(ui *UI) (*octoPrintResponseManager) {
	if octoPrintResponseManagerInstance == nil {
		if ui == nil {
			panic("GetOctoPrintResponseManagerInstance() was called for the first time, but the ui object passed in was nil")
		}

		octoPrintResponseManagerInstance = &octoPrintResponseManager{
			UI: ui,
			// Client: client,
			// // IsRunning: false,
			// ConnectAttempts: 0,
			// IsConnectedToOctoPrint: false,
			// IsConnectedToPrinter: false,
		}

		octoPrintResponseManagerInstance.createBackgroundTask()
	}

	return octoPrintResponseManagerInstance
}

func (this *octoPrintResponseManager) createBackgroundTask() {
	logger.TraceEnter("OctoPrintResponseManager.createBackgroundTask()")

	// Default timeout of 10 seconds.
	duration := utils.GetExperimentalFrequency(10, "EXPERIMENTAL_OCTO_PRINT_RESPONSE_MANGER_UPDATE_FREQUENCY")
	this.backgroundTask = utils.CreateBackgroundTask(duration, this.update)
	this.backgroundTask.Start()

	logger.TraceLeave("OctoPrintResponseManager.createBackgroundTask()")
}

func (this *octoPrintResponseManager) update() {
	logger.TraceEnter("OctoPrintResponseManager.update()")

	connectionManager := utils.GetConnectionManagerInstance(this.UI.Client)
	if connectionManager.IsConnected() != true {
		// If not connected, do nothing and leave.
		logger.TraceLeave("OctoPrintResponseManager.update()")
		return
	}


	// call APIs
	// uiWidgets/TemperatureStatusBox.go
	// utils.CreateUpdateTemperaturesBackgroundTask(instance, client)
	// UpdateTemperaturesBackgroundTask.updateTemperaturesCallback()
	// 	Tools.GetCurrentTemperatureData()
	// 	target: /api/printer?history=false&exclude=sd,state  file=logger.go func=logger.Debugf line=141
	// 	temperatureDataResponse, err := (&octoprintApis.TemperatureDataRequest{}).Do(client)
	// 	JSON response:
	// 	{
	// 		"temperature": {
	// 			"bed": {
	// 				"actual": 22.1,
	// 				"offset": 0,
	// 				"target": 0.0
	// 			},
	// 			"tool0": {
	// 				"actual": 15.8,
	// 				"offset": 0,
	// 				"target": 0.0
	// 			}
	// 		}
	// 	}

	// IdleStatusPanel.updateTemperature()
	// 	target: /api/printer?history=false&limit=0&exclude=sd  file=logger.go func=logger.Debugf line=141
	// 	fullStateResponse, err := (&octoprintApis.FullStateRequest{Exclude: []string{"sd"}}).Do(this.UI.Client)
	// 	JSON response:
	// 	{
	// 		"state": {
	// 			"error": "",
	// 			"flags": {
	// 				"cancelling": false,
	// 				"closedOrError": false,
	// 				"error": false,
	// 				"finishing": false,
	// 				"operational": true,
	// 				"paused": false,
	// 				"pausing": false,
	// 				"printing": false,
	// 				"ready": true,
	// 				"resuming": false,
	// 				"sdReady": true
	// 			},
	// 			"text": "Operational"
	// 		},
	// 		"temperature": {
	// 			"bed": {
	// 				"actual": 22.1,
	// 				"offset": 0,
	// 				"target":0.0
	// 			},
	// 			"tool0": {
	// 				"actual": 15.8,
	// 				"offset": 0,
	// 				"target": 0.0
	// 			}
	// 		}
	// 	}

	// 	fullStateResponse, err := (&octoprintApis.FullStateRequest{}).Do(this.UI.Client)
	// 	JSON response:
	// {
	// 	"sd": {
	// 		"ready": true
	// 	},
	// 	"state": {
	// 		"error": "",
	// 		"flags": {
	// 			"cancelling": false,
	// 			"closedOrError": false,
	// 			"error": false,
	// 			"finishing": false,
	// 			"operational": true,
	// 			"paused": false,
	// 			"pausing": false,
	// 			"printing": false,
	// 			"ready": true,
	// 			"resuming": false,
	// 			"sdReady": true
	// 		},
	// 		"text": "Operational"
	// 	},
	// 	"temperature": {
	// 		"bed": {
	// 			"actual": 22.3,
	// 			"offset": 0,
	// 			"target": 0.0
	// 		},
	// 		"tool0": {
	// 			"actual": 16.2,
	// 			"offset": 0,
	// 			"target": 0.0
	// 		}
	// 	}
	// }


	fullStateResponse, err := (&octoprintApis.FullStateRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("OctoPrintResponseManager.update()", "Do(FullStateRequest)", err)
	
		connectionManager.ReInitializeConnectionState()
		this.UI.GoToPanel(GetConnectionPanelInstance(this.UI))
	
		logger.TraceLeave("OctoPrintResponseManager.update()")
		return
	}

	if fullStateResponse == nil /*|| fullStateResponse.Temperature == nil*/ || fullStateResponse.Temperature.CurrentTemperatureData == nil {
		logger.Error("OctoPrintResponseManager.update() - fullStateResponse.Temperature.CurrentTemperatureData is invalid")
	
		connectionManager.ReInitializeConnectionState()
		this.UI.GoToPanel(GetConnectionPanelInstance(this.UI))
	
		logger.TraceLeave("OctoPrintResponseManager.update()")
		return
	}


	this.FullStateResponse = *fullStateResponse

	// connectionManager.UpdateStatus()


	logger.TraceLeave("OctoPrintResponseManager.update()")
}


func (this *octoPrintResponseManager) IsConnected() bool {
	// TODO: should this be named IsFullyConnected()?
	connectionManager := utils.GetConnectionManagerInstance(this.UI.Client)
	return connectionManager.IsConnected()
}


// func (this *connectionManager) ReInitializeConnectionState() {
// 	// this.IsRunning = true
// 	this.ConnectAttempts = 0
// 	this.IsConnectedToOctoPrint = false
// 	this.IsConnectedToPrinter = false
// }

// func (this *connectionManager) UpdateStatus() {
// 	logger.TraceEnter("ConnectionManager.UpdateStatus()")

// 	if this.IsConnected() != true {
// 		if this.ConnectAttempts > MAX_CONNECTION_ATTEMPTS {
// 			// verify
// 			this.ConnectAttempts++
// 		}

// 		logger.Debug("ConnectionManager.UpdateStatus() - about to call ConnectionRequest.Do()")
// 		t1 := time.Now()
// 		connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(this.Client)
// 		t2 := time.Now()
// 		logger.Debug("ConnectionManager.UpdateStatus() - finished calling ConnectionRequest.Do()")
// 		logger.Debugf("time elapsed: %q", t2.Sub(t1))
		
// 		if err != nil {
// 			logger.LogError("ConnectionManager.UpdateStatus()", "ConnectionRequest.Do()", err)
// 								// newUIState, splashMessage = this.getUiStateAndMessageFromError(err, newUIState, splashMessage)
// 								// logger.Debugf("ConnectionManager.UpdateStatus() - newUIState is now: %s", newUIState)
// 			this.IsConnectedToOctoPrint = false
// 			logger.Debug("ConnectionManager.UpdateStatus() - Connection state: IsConnectedToOctoPrint is now false")
// 			logger.TraceLeave("ConnectionManager.UpdateStatus()")
// 			return
// 		}
		
// 		logger.Debug("ConnectionManager.UpdateStatus() - ConnectionRequest.Do() succeeded")
		
// 		this.IsConnectedToOctoPrint = true
		
// 		jsonResponse, err := StructToJson(connectionResponse)
// 		if err != nil {
// 			logger.LogError("ConnectionManager.UpdateStatus()", "StructToJson()", err)
// 			// If there's an error here, it's with the serialization of the object to JSON.
// 			// This is just for debugging, so don't return if there's an issue, and just
// 			// carry on (and hopefully connectionResponse isn't corrupted)
// 		} else {
// 			logger.Debugf("ConnectionManager.UpdateStatus() - connectionResponse is: %s", jsonResponse)
// 		}
		
// 		/*
// 		Example JSON response:
// 		{
// 			"Current": {
// 				"state": "Operational",
// 				"port": "/dev/ttyACM0",
// 				"baudrate": 115200,
// 				"printerProfile": "_default"
// 			},
// 			"Options": {
// 				"ports": [
// 					"/dev/ttyACM0"
// 				],
// 				"baudrates": [
// 					250000,
// 					230400,
// 					115200,
// 					57600,
// 					38400,
// 					19200,
// 					9600
// 				],
// 				"printerProfiles": [
// 					{
// 						"id": "_default",
// 						"name": "name-of-the-printer"
// 					}
// 				],
// 				"portPreference": "",
// 				"baudratePreference": 0,
// 				"printerProfilePreference": "_default",
// 				"autoconnect": false
// 			}
// 		}
// 		*/
		
// 		printerConnectionState := connectionResponse.Current.State
// 		if printerConnectionState.IsOffline() || printerConnectionState.IsError() {
// 			this.IsConnectedToPrinter = false
// 		} else {
// 			this.IsConnectedToPrinter = true
// 		}
// 	}

// 	logger.TraceLeave("ConnectionManager.UpdateStatus()")
// }

// func (this *connectionManager) IsConnected() bool {
// 	// TODO: should this be named IsFullyConnected?
// 	return this.IsConnectedToOctoPrint == true && this.IsConnectedToPrinter == true;
// }
