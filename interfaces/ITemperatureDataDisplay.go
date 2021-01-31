package interfaces

import (
	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)

type ITemperatureDataDisplay interface {
	UpdateTemperatureData(temperatureData map[string]dataModels.TemperatureData)
}
