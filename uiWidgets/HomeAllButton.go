package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type HomeAllButton struct {
	*gtk.Button

	client				*octoprint.Client
}

func CreateHomeAllButton(
	client				*octoprint.Client,
) *HomeAllButton {
	base := utils.MustButtonImageStyle("Home All", "home.svg", "", nil)

	instance := &HomeAllButton {
		Button:				base,
		client:				client,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		panic(err)
	}

	return instance
}

func (this *HomeAllButton) handleClicked() {
	utils.Logger.Infof("Homing the print head")

	// Version A:
	axes := []octoprint.Axis {
		octoprint.XAxis,
		octoprint.YAxis,
		octoprint.ZAxis,
	}
	cmd := &octoprint.PrintHeadHomeRequest{Axes: axes}
	err := cmd.Do(this.client);
	if err != nil {
		utils.LogError("HomeAllButton.handleClicked()", "Do(PrintHeadHomeRequest)", err)
	}


	/*
	// If there are issues with version A, there's also version B:
	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		"G28 Z",
		"G28 X",
		"G28 Y",
	}

	if err := cmd.Do(m.UI.Client); err != nil {
		utils.LogError("HomeAllButton.handleClicked()", "Do(CommandRequest)", err)
	}
	*/
}
