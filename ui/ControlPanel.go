package ui

import (
	"strings"

	// "github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	// "github.com/Z-Bolt/OctoScreen/utils"
)


type controlPanel struct {
	CommonPanel
}

var controlPanelInstance *controlPanel

func GetControlPanelInstance(
	ui 				*UI,
) *controlPanel {
	if controlPanelInstance == nil {
		instance := &controlPanel {
			CommonPanel: CreateCommonPanel("ControlPanel", ui),
		}
		instance.initialize()
		controlPanelInstance = instance
	}

	return controlPanelInstance
}

func (this *controlPanel) initialize() {
	defer this.Initialize()

	defaultControls := this.getDefaultControls()
	for _, controlDefinition := range defaultControls {
		icon := strings.ToLower(strings.Replace(controlDefinition.Name, " ", "-", -1))
		button := uiWidgets.CreateControlButton(this.UI.Client, this.UI.window, controlDefinition, icon)
		this.AddButton(button)
	}


	// 12 (max) - Back button = 11 available slots to display.
	const maxSlots = 11
	currentButtonCount := len(defaultControls)

	for _, controlDefinition := range this.getCustomControls() {
		if currentButtonCount < maxSlots {
			button := uiWidgets.CreateControlButton(this.UI.Client, this.UI.window, controlDefinition, "custom-script")
			this.AddButton(button)
			currentButtonCount++
		}
	}

	for _, commandDefinition := range this.getCommands() {
		if currentButtonCount < maxSlots {
			button := uiWidgets.CreateCommandButton(this.UI.Client, this.UI.window, commandDefinition, "custom-script")
			this.AddButton(button)
			currentButtonCount++
		}
	}
}

func (this *controlPanel) getDefaultControls() []*dataModels.ControlDefinition {
	var controlDefinitions = []*dataModels.ControlDefinition{{
		Name:    "Motor Off",
		Command: "M18",			// Disable all stepper motors immediately
	}, {
		Name:    "Motor On",
		Command: "M17",			// Enable all stepper motors
	}, {
		Name:    "Fan Off",
		Command: "M106 S0",		// Sets the fan speed to off
	}, {
		Name:    "Fan On",
		Command: "M106",		// Sets the fan speed to full
	}}

	return controlDefinitions
}

func (this *controlPanel) getCustomControls() []*dataModels.ControlDefinition {
	controlDefinitions := []*dataModels.ControlDefinition{}

	logger.Info("control.getCustomControl() - Retrieving custom controls")
	response, err := (&octoprintApis.CustomCommandsRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("control.getCustomControl()", "Do(ControlDefinition)", err)
		return controlDefinitions
	}

	for _, control := range response.Controls {
		for _, childControl := range control.Children {
			if childControl.Command != "" || childControl.Script != "" || childControl.Commands != nil {
				controlDefinitions = append(controlDefinitions, childControl)
			}
		}
	}

	return controlDefinitions
}

func (this *controlPanel) getCommands() []*dataModels.CommandDefinition {
	logger.Info("Retrieving custom commands")
	response, err := (&octoprintApis.SystemCommandsRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("control.getCommands()", "Do(SystemCommandsRequest)", err)
		return nil
	}

	return response.Custom
}
