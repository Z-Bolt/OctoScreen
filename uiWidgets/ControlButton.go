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
	base := utils.MustButtonImage(utils.StrEllipsisLen(controlDefinition.Name, 16), iconName + ".svg", nil)
	instance := &ControlButton {
		Button:				base,
		client:				client,
		parentWindow:		parentWindow,
		controlDefinition:	controlDefinition,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		logger.LogError("PANIC!!! - CreateControlButton()", "instance.Button.Connect()", err)
		panic(err)
	}

	return instance
}

func (this *ControlButton) handleClicked() {
	if len(this.controlDefinition.Confirm) != 0 {
		utils.MustConfirmDialogBox(
			this.parentWindow,
			this.controlDefinition.Confirm,
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
