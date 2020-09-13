package ui

import (
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var temperaturePresetsPanelInstance *temperaturePresetsPanel

type temperaturePresetsPanel struct {
	CommonPanel
	selectToolStepButton	*uiWidgets.SelectToolStepButton
}

func TemperaturePresetsPanel(
	ui						*UI,
	parentPanel				interfaces.IPanel,
	selectToolStepButton	*uiWidgets.SelectToolStepButton,
) *temperaturePresetsPanel {
	if temperaturePresetsPanelInstance == nil {
		instance := &temperaturePresetsPanel {
			CommonPanel:			NewCommonPanel(ui, parentPanel),
			selectToolStepButton:	selectToolStepButton,
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
	settings, err := (&octoprint.SettingsRequest{}).Do(this.UI.Printer)
	if err != nil {
		utils.LogError("TemperaturePresetsPanel.getTemperaturePresets()", "Do(SettingsRequest)", err)
		return
	}

	count := 0
	for _, temperaturePreset := range settings.Temperature.TemperaturePresets {
		if count < 10 {
			temperaturePresetButton := uiWidgets.CreateTemperaturePresetButton(
				this.UI.Printer,
				this.selectToolStepButton,
				"heat-up.svg",
				temperaturePreset,
				this.UI.GoToPreviousPanel,
			)
			this.AddButton(temperaturePresetButton)
			count++
		}
	}

	coolDownTemperaturePreset := octoprint.TemperaturePreset{
		Name:		"Cool Down",
		Bed:		0.0,
		Extruder:	0.0,
	}
	coolDownButton := uiWidgets.CreateTemperaturePresetButton(
		this.UI.Printer,
		this.selectToolStepButton,
		"cool-down.svg",
		&coolDownTemperaturePreset,
		this.UI.GoToPreviousPanel,
	)
	this.AddButton(coolDownButton)
}
