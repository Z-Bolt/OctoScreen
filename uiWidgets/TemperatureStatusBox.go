package uiWidgets

import (
	// "time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type TemperatureStatusBox struct {
	*gtk.Box
	interfaces.ITemperatureDataDisplay

	client						*octoprintApis.Client
	labelWithImages				map[string]*utils.LabelWithImage
}

func CreateTemperatureStatusBox(
	client						*octoprintApis.Client,
	includeHotends				bool,
	includeBed					bool,
) *TemperatureStatusBox {
	if !includeHotends && !includeBed {
		logger.Error("TemperatureStatusBox.CreateTemperatureStatusBox() - both includeToolheads and includeBed are false, but at least one needs to be true")
		return nil
	}

	currentTemperatureData, err := utils.GetCurrentTemperatureData(client)
	if err != nil {
		logger.LogError("TemperatureStatusBox.CreateTemperatureStatusBox()", "GetCurrentTemperatureData(client)", err)
		return nil
	}

	base := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)

	instance := &TemperatureStatusBox{
		Box:						base,
		client:						client,
		labelWithImages:			map[string]*utils.LabelWithImage{},
	}

	instance.SetVAlign(gtk.ALIGN_CENTER)
	instance.SetHAlign(gtk.ALIGN_CENTER)

	var bedTemperatureData *dataModels.TemperatureData = nil
	var hotendIndex int = 0
	var hotendCount int = utils.GetHotendCount(client)
	for key, temperatureData := range currentTemperatureData {
		if key == "bed" {
			bedTemperatureData = &temperatureData
		} else {
			hotendIndex++

			if includeHotends {
				if hotendIndex <= hotendCount {
					strImageFileName := utils.GetNozzleFileName(hotendIndex, hotendCount)
					instance.labelWithImages[key] = utils.MustLabelWithImage(strImageFileName, "")
					instance.Add(instance.labelWithImages[key])
				}
			}
		}
	}

	if bedTemperatureData != nil {
		if includeBed {
			instance.labelWithImages["bed"] = utils.MustLabelWithImage("bed.svg", "")
			instance.Add(instance.labelWithImages["bed"])
		}
	}

	return instance
}

// interfaces.ITemperatureDataDisplay
func (this *TemperatureStatusBox) UpdateTemperatureData(currentTemperatureData map[string]dataModels.TemperatureData) {
	for key, temperatureData := range currentTemperatureData {
		if labelWithImage, ok := this.labelWithImages[key]; ok {
			temperatureDataString := utils.GetTemperatureDataString(temperatureData)
			labelWithImage.Label.SetText(temperatureDataString)
			labelWithImage.ShowAll()
		}
	}
}
