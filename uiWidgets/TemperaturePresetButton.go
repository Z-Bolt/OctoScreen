package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type TemperaturePresetButton struct {
	*gtk.Button

	client						*octoprintApis.Client
	selectHotendStepButton		*SelectToolStepButton
	imageFileName				string
	temperaturePreset			*dataModels.TemperaturePreset
	callback					func()
}

func CreateTemperaturePresetButton(
	client						*octoprintApis.Client,
	selectHotendStepButton		*SelectToolStepButton,
	imageFileName				string,
	temperaturePreset			*dataModels.TemperaturePreset,
	callback					func(),
) *TemperaturePresetButton {
	presetName := utils.StrEllipsisLen(temperaturePreset.Name, 10)
	base := utils.MustButtonImage(presetName, imageFileName, nil)

	instance := &TemperaturePresetButton{
		Button:						base,
		client:						client,
		selectHotendStepButton:		selectHotendStepButton,
		imageFileName:				imageFileName,
		temperaturePreset:			temperaturePreset,
		callback:					callback,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		logger.LogError("PANIC!!! - CreateTemperaturePresetButton()", "instance.Button.Connect()", err)
		panic(err)
	}

	return instance
}

func (this *TemperaturePresetButton) handleClicked() {
	logger.Infof("TemperaturePresetButton.handleClicked() - setting temperature to preset %s.", this.temperaturePreset.Name)
	logger.Infof("TemperaturePresetButton.handleClicked() - setting hotend temperature to %.0f.", this.temperaturePreset.Extruder)
	logger.Infof("TemperaturePresetButton.handleClicked() - setting bed temperature to %.0f.", this.temperaturePreset.Bed)

	currentTool := this.selectHotendStepButton.Value()
	if currentTool == "" {
		logger.Error("TemperaturePresetButton.handleClicked() - currentTool is invalid (blank), defaulting to tool0")
		currentTool = "tool0"
	}

	/*
	CreateTemperaturePresetButton is used by TemperaturePresetsPanel.  Strictly speaking,
	CreateTemperaturePresetButton should only set the temperature of one device at at time,
	but that's a lousy UX.  Imagine being in the TemperaturePanel... with the tool set to
	the hotend, click the More button (and go to the TemperaturePresetsPanel), then
	clicking PLA (and get taken back to the TemperaturePanel), __THEN__ have to click the
	tool button to change to the bed, and then repeat the process over again.

	So, instead, the temperature of both the bed and the selected tool (or tool0 if the bed
	is selected) are set.

	NOTE: This only changes the temperature of the bed and the currently selected hotend
	(which is passed into the TemperaturePresetsPanel, and then passed into
	CreateTemperaturePresetButton).  The code could be changed so it sets the temperature
	of every hotend, but this is problematic if one is using different materials with
	different temperature characteristics.
	*/

	// Set the bed's temp.
	bedTargetRequest := &octoprintApis.BedTargetRequest{Target: this.temperaturePreset.Bed}
	err := bedTargetRequest.Do(this.client)
	if err != nil {
		logger.LogError("TemperaturePresetButton.handleClicked()", "Do(BedTargetRequest)", err)
		return
	}

	// Set the hotend's temp.
	var toolTargetRequest *octoprintApis.ToolTargetRequest
	if currentTool == "bed" {
		// If current tool is set to "bed", use tool0.
		toolTargetRequest = &octoprintApis.ToolTargetRequest {
			Targets: map[string]float64 {
				"tool0": this.temperaturePreset.Extruder,
			},
		}
	} else {
		toolTargetRequest = &octoprintApis.ToolTargetRequest {
			Targets: map[string]float64 {
				currentTool: this.temperaturePreset.Extruder,
			},
		}
	}

	err = toolTargetRequest.Do(this.client)
	if err != nil {
		logger.LogError("TemperaturePresetButton.handleClicked()", "Do(ToolTargetRequest)", err)
	}

	if this.callback != nil {
		this.callback()
	}
}
