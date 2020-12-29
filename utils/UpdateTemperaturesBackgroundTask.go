package utils

import (
	"time"

	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
)

// TODO: convert this into a singleton class

var UpdateTemperaturesBackgroundTask *BackgroundTask = nil
var temperatureDataDisplays []interfaces.ITemperatureDataDisplay
var registeredClient *octoprint.Client = nil

func CreateUpdateTemperaturesBackgroundTask(temperatureDataDisplay interfaces.ITemperatureDataDisplay, client *octoprint.Client) {
	if UpdateTemperaturesBackgroundTask != nil {
		Logger.Error("UpdateTemperaturesBackgroundTask.CreateUpdateTemperaturesBackgroundTask() - updateTemperaturesBackgroundTask has already been set")
		return
	}

	UpdateTemperaturesBackgroundTask = CreateBackgroundTask(time.Second * 3, updateTemperaturesCallback)
	RegisterTemperatureStatusBox(temperatureDataDisplay, client)
	UpdateTemperaturesBackgroundTask.Start()
}

func RegisterTemperatureStatusBox(temperatureDataDisplay interfaces.ITemperatureDataDisplay, client *octoprint.Client) {
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
