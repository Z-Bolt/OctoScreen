package interfaces

import (
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
)

type ITemperatureDataDisplay interface {
	UpdateTemperatureData(temperatureData map[string]octoprintApis.TemperatureData)
}
