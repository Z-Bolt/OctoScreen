package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type HomeButton struct {
	*gtk.Button

	client				*octoprint.Client
	axes				[]octoprint.Axis
}

func CreateHomeButton(
	client				*octoprint.Client,
	label				string,
	image				string,
	axes				...octoprint.Axis,
) *HomeButton {
	base := utils.MustButtonImageStyle(label, image, "", nil)

	instance := &HomeButton {
		Button:				base,
		client:				client,
		axes:				axes,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		panic(err)
	}

	return instance
}

func (this *HomeButton) handleClicked() {
	cmd := &octoprint.PrintHeadHomeRequest{Axes: this.axes}
	utils.Logger.Infof("Homing the print head in %s axes", this.axes)
	err := cmd.Do(this.client);
	if err != nil {
		utils.LogError("HomeButton.handleClicked()", "Do(PrintHeadHomeRequest)", err)
	}
}
