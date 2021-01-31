package utils

import (
	"time"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// TODO: convert this into a singleton

var UpdateTemperaturesBackgroundTask *BackgroundTask = nil
var temperatureDataDisplays []interfaces.ITemperatureDataDisplay
var registeredClient *octoprintApis.Client = nil

func CreateUpdateTemperaturesBackgroundTask(temperatureDataDisplay interfaces.ITemperatureDataDisplay, client *octoprintApis.Client) {
	if UpdateTemperaturesBackgroundTask != nil {
		Logger.Error("UpdateTemperaturesBackgroundTask.CreateUpdateTemperaturesBackgroundTask() - updateTemperaturesBackgroundTask has already been set")
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
	currentTemperatureData, err := GetCurrentTemperatureData(registeredClient)
	if err != nil {
		LogError("UpdateTemperaturesBackgroundTask.updateTemperaturesCallback()", "GetCurrentTemperatureData(client)", err)
		return
	}

	for _, temperatureDataDisplay := range temperatureDataDisplays {
		temperatureDataDisplay.UpdateTemperatureData(currentTemperatureData)
	}
}
