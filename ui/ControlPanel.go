package ui

import (
	"strings"

	// "github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

var controlPanelInstance *controlPanel

type controlPanel struct {
	CommonPanel
}

func ControlPanel(
	ui 				*UI,
	parentPanel		interfaces.IPanel,
) *controlPanel {
	if controlPanelInstance == nil {
		instance := &controlPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		instance.initialize()
		controlPanelInstance = instance
	}

	return controlPanelInstance
}

func (this *controlPanel) initialize() {
	defer this.Initialize()

	for _, controlDefinition := range this.getDefaultControls() {
		icon := strings.ToLower(strings.Replace(controlDefinition.Name, " ", "-", -1))
		button := uiWidgets.CreateControlButton(this.UI.Client, this.UI.window, controlDefinition, icon)
		this.AddButton(button)
	}

	for _, controlDefinition := range this.getCustomControls() {
		button := uiWidgets.CreateControlButton(this.UI.Client, this.UI.window, controlDefinition, "custom-script")
		this.AddButton(button)
	}

	for _, commandDefinition := range this.getCommands() {
		button := uiWidgets.CreateCommandButton(this.UI.Client, this.UI.window, commandDefinition, "custom-script")
		this.AddButton(button)
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
			if childControl.Command != "" || childControl.Script != "" {
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
