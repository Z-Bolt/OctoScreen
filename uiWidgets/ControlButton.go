package uiWidgets

import (
	// "fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type ControlButton struct {
	*gtk.Button

	client				*octoprintApis.Client
	parentWindow		*gtk.Window
	controlDefinition	*dataModels.ControlDefinition
}

func CreateControlButton(
	client				*octoprintApis.Client,
	parentWindow		*gtk.Window,
	controlDefinition	*dataModels.ControlDefinition,
	iconName			string,
) *ControlButton {
	style := ""
	if controlRequiresConfirmation(controlDefinition) {
		style = "color-warning-sign-yellow"
	}
	base := utils.MustButtonImageStyle(utils.StrEllipsisLen(controlDefinition.Name, 16), iconName + ".svg", style, nil)

	instance := &ControlButton {
		Button:				base,
		client:				client,
		parentWindow:		parentWindow,
		controlDefinition:	controlDefinition,
	}
	instance.Button.Connect("clicked", instance.handleClicked)

	return instance
}

func controlRequiresConfirmation(controlDefinition *dataModels.ControlDefinition) bool {
	return controlDefinition != nil && len(controlDefinition.Confirm) > 0
}

func (this *ControlButton) handleClicked() {
	if controlRequiresConfirmation(this.controlDefinition) {
		utils.MustConfirmDialogBox(
			this.parentWindow,
			this.controlDefinition.Confirm + "\n\nAre you sure you want to proceed?",
			this.sendCommand,
		)()
	} else {
		this.sendCommand()
	}
}

func (this *ControlButton) sendCommand() {
	logger.Infof("ControlButton.sendCommand(), now sending command %q", this.controlDefinition.Name)

	commandRequest := &octoprintApis.CommandRequest{
		Commands: this.controlDefinition.Commands,
	}

	if len(this.controlDefinition.Command) != 0 {
		commandRequest.Commands = []string{this.controlDefinition.Command}
	}

	err := commandRequest.Do(this.client)
	if err != nil {
		logger.LogError("ControlButton.sendCommand()", "Do(CommandRequest)", err)
	}
}
