package interfaces

import (
	"github.com/mcuadros/go-octoprint"
)

type ITemperatureDataDisplay interface {
	UpdateTemperatureData(temperatureData map[string]octoprint.TemperatureData)
}
