// TODO: This file should now be obsolete, and should be deleted


package utils

import (
	"time"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// TODO: convert this into a singleton

var UpdateTemperaturesBackgroundTask *BackgroundTask = nil
var temperatureDataDisplays []interfaces.ITemperatureDataDisplay
var registeredClient *octoprintApis.Client = nil

func CreateUpdateTemperaturesBackgroundTask(
	temperatureDataDisplay interfaces.ITemperatureDataDisplay,
	client *octoprintApis.Client,
) {
	if UpdateTemperaturesBackgroundTask != nil {
		logger.Error("UpdateTemperaturesBackgroundTask.CreateUpdateTemperaturesBackgroundTask() - updateTemperaturesBackgroundTask has already been set")
		return
	}

	UpdateTemperaturesBackgroundTask = CreateBackgroundTask(time.Second * 3, updateTemperaturesCallback)
	RegisterTemperatureStatusBox(temperatureDataDisplay, client)
	UpdateTemperaturesBackgroundTask.Start()
}

func RegisterTemperatureStatusBox(temperatureDataDisplay interfaces.ITemperatureDataDisplay, client *octoprintApis.Client) {
	temperatureDataDisplays = append(temperatureDataDisplays, temperatureDataDisplay)
	registeredClient = client
}

func updateTemperaturesCallback() {
	logger.TraceEnter("UpdateTemperaturesBackgroundTask.updateTemperaturesCallback()")

	// TODO: add guard if printer isn't connected
	// can't do right now due to circular dependency:
	//		TemperatureStatusBox creates the background task...
	//		background task needs UI.UIState or UIState.connectionAttempts
	//		UI has panel
	//		panel has TemperatureStatusBox

	currentTemperatureData, err := GetCurrentTemperatureData(registeredClient)
	if err != nil {
		logger.LogError("UpdateTemperaturesBackgroundTask.updateTemperaturesCallback()", "GetCurrentTemperatureData(client)", err)
	} else {
		for _, temperatureDataDisplay := range temperatureDataDisplays {
			temperatureDataDisplay.UpdateTemperatureData(currentTemperatureData)
		}
	}

	logger.TraceLeave("UpdateTemperaturesBackgroundTask.updateTemperaturesCallback()")
}
