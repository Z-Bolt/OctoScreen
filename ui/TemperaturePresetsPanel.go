package ui

import (
	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	// "github.com/Z-Bolt/OctoScreen/utils"
)


type temperaturePresetsPanel struct {
	CommonPanel

	selectHotendStepButton	*uiWidgets.SelectToolStepButton
}

var temperaturePresetsPanelInstance *temperaturePresetsPanel

func GetTemperaturePresetsPanelInstance(
	ui						*UI,
	selectHotendStepButton	*uiWidgets.SelectToolStepButton,
) *temperaturePresetsPanel {
	if temperaturePresetsPanelInstance == nil {
		instance := &temperaturePresetsPanel {
			CommonPanel:			CreateCommonPanel("temperaturePresetsPanel", ui),
			selectHotendStepButton:	selectHotendStepButton,
		}
		instance.initialize()
		temperaturePresetsPanelInstance = instance
	}

	return temperaturePresetsPanelInstance
}

func (this *temperaturePresetsPanel) initialize() {
	defer this.Initialize()
	this.createAllOffButton()
	this.createTemperaturePresetButtons()
}

func (this *temperaturePresetsPanel) createAllOffButton() {
	allOffButton := uiWidgets.CreateCoolDownButton(this.UI.Client, this.UI.GoToPreviousPanel)
	this.AddButton(allOffButton)
}

func (this *temperaturePresetsPanel) createTemperaturePresetButtons() {
	settings, err := (&octoprintApis.SettingsRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("TemperaturePresetsPanel.getTemperaturePresets()", "Do(SettingsRequest)", err)
		return
	}

	// 12 (max) - Back button - All Off button = 10 available slots to display.
	const maxSlots = 10

	count := 0
	for _, temperaturePreset := range settings.Temperature.TemperaturePresets {
		if count < maxSlots {
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
