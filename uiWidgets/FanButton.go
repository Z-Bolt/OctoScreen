package uiWidgets

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type FanButton struct {
	*gtk.Button

	client				*octoprintApis.Client
	amount				int
}

func CreateFanButton(
	client				*octoprintApis.Client,
	amount				int,
) *FanButton {
	var (
		label string
		image string
	)

	if amount == 0 {
		label = "Fan Off"
		image = "fan-off.svg"
	} else {
		label = fmt.Sprintf("%d %%", amount)
		image = "fan.svg"
	}

	base := utils.MustButtonImageStyle(label, image, "", nil)
	instance := &FanButton {
		Button:				base,
		client:				client,
		amount:				amount,
	}
	instance.Button.Connect("clicked", instance.handleClicked)

	return instance
}

func (this *FanButton) handleClicked() {
	cmd := &octoprintApis.CommandRequest{}
	cmd.Commands = []string{
		fmt.Sprintf("M106 S%d", (255 * this.amount / 100)),
	}

	err := cmd.Do(this.client)
	if err != nil {
		logger.LogError("FanButton.handleClicked()", "Do(CommandRequest)", err)
	}
}
