package ui

import (
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

var temperaturePresetsPanelInstance *temperaturePresetsPanel

type temperaturePresetsPanel struct {
	CommonPanel

	selectHotendStepButton	*uiWidgets.SelectToolStepButton

}

func TemperaturePresetsPanel(
	ui						*UI,
	parentPanel				interfaces.IPanel,
	selectHotendStepButton	*uiWidgets.SelectToolStepButton,
) *temperaturePresetsPanel {
	if temperaturePresetsPanelInstance == nil {
		instance := &temperaturePresetsPanel {
			CommonPanel:			NewCommonPanel(ui, parentPanel),
			selectHotendStepButton:	selectHotendStepButton,
		}
		instance.initialize()
		temperaturePresetsPanelInstance = instance
	}

	return temperaturePresetsPanelInstance
}

func (this *temperaturePresetsPanel) initialize() {
	defer this.Initialize()
	this.createTemperaturePresetButtons()
}

func (this *temperaturePresetsPanel) createTemperaturePresetButtons() {
	settings, err := (&octoprintApis.SettingsRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("TemperaturePresetsPanel.getTemperaturePresets()", "Do(SettingsRequest)", err)
		return
	}

	this.createAllOffButton()

	count := 0
	for _, temperaturePreset := range settings.Temperature.TemperaturePresets {
		if count < 10 {
			temperaturePresetButton := uiWidgets.CreateTemperaturePresetButton(
				this.UI.Client,
				this.selectHotendStepButton,
				"heat-up.svg",
				temperaturePreset,
				this.UI.GoToPreviousPanel,
			)
			this.AddButton(temperaturePresetButton)
			count++
		}
	}

}

func (this *temperaturePresetsPanel) createAllOffButton() {
	allOffButton := uiWidgets.CreateCoolDownButton(this.UI.Client, this.UI.GoToPreviousPanel)
	this.AddButton(allOffButton)
}
