package ui

import (
	"strings"

	// "github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
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

func (this *controlPanel) getDefaultControls() []*octoprint.ControlDefinition {
	var controlDefinitions = []*octoprint.ControlDefinition{{
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

func (this *controlPanel) getCustomControls() []*octoprint.ControlDefinition {
	control := []*octoprint.ControlDefinition{}

	utils.Logger.Info("control.getCustomControl() - Retrieving custom controls")
	r, err := (&octoprint.CustomCommandsRequest{}).Do(this.UI.Client)
	if err != nil {
		utils.LogError("control.getCustomControl()", "Do(ControlDefinition)", err)
		return control
	}

	for _, c := range r.Controls {
		for _, cc := range c.Children {
			if cc.Command != "" || cc.Script != "" {
				control = append(control, cc)
			}
		}
	}

	return control
}

func (this *controlPanel) getCommands() []*octoprint.CommandDefinition {
	utils.Logger.Info("Retrieving custom commands")
	r, err := (&octoprint.SystemCommandsRequest{}).Do(this.UI.Client)
	if err != nil {
		utils.LogError("control.getCommands()", "Do(SystemCommandsRequest)", err)
		return nil
	}

	return r.Custom
}
