package uiWidgets

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type SystemCommandButton struct {
	*gtk.Button
}

func CreateSystemCommandButton(
	client				*octoprintApis.Client,
	parentWindow		*gtk.Window,
	name				string,
	action				string,
	style				string,
) *SystemCommandButton {
	systemCommandsResponse, err := (&octoprintApis.SystemCommandsRequest{}).Do(client)
	if err != nil {
		logger.LogError("PANIC!!! - CreateSystemCommandButton()", "SystemCommandsRequest.Do()", err)
		panic(err)
	}

	var cmd *dataModels.CommandDefinition
	var cb func()

	for _, commandDefinition := range systemCommandsResponse.Core {
		if commandDefinition.Action == action {
			cmd = commandDefinition
		}
	}

	if cmd != nil {
		do := func() {
			systemExecuteCommandRequest := &octoprintApis.SystemExecuteCommandRequest{
				Source: dataModels.Core,
				Action: cmd.Action,
			}

			if err := systemExecuteCommandRequest.Do(client); err != nil {
				logger.LogError("system.createCommandButton()", "Do(SystemExecuteCommandRequest)", err)
				return
			}
		}

		confirmationMessage := ""
		if len(cmd.Confirm) != 0 {
			confirmationMessage = fmt.Sprintf("%s\n\nDo you wish to proceed?", cmd.Confirm)
		} else if len(name) != 0 {
			confirmationMessage = fmt.Sprintf("Do you wish to %s?", name)
		} else {
			confirmationMessage = "Do you wish to proceed?"
		}

		cb = utils.MustConfirmDialogBox(parentWindow, confirmationMessage, do)
	}

	base := utils.MustButtonImageStyle(name, action + ".svg", style, cb)
	ctx, _ := base.GetStyleContext()
	ctx.AddClass("font-size-19")

	instance := &SystemCommandButton {
		Button:				base,
	}

	if cmd == nil {
		instance.SetSensitive(false)
	}

	return instance
}
