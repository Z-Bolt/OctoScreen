package ui

import (
	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/utils"
)


var configurationPanelInstance *configurationPanel

type configurationPanel struct {
	CommonPanel
}

func ConfigurationPanel(
	ui				*UI,
) *configurationPanel {
	if configurationPanelInstance == nil {
		instance := &configurationPanel {
			CommonPanel: NewCommonPanel("ConfigurationPanel", ui),
		}
		instance.initialize()
		configurationPanelInstance = instance
	}

	return configurationPanelInstance
}

func (this *configurationPanel) initialize() {
	defer this.Initialize()

	bedlevelButton := utils.MustButtonImageStyle(
		"Bed Level",
		"bed-level.svg",
		"color1",
		this.showBedLevelPanel,
	)
	this.Grid().Attach(bedlevelButton, 0, 0, 1, 1)

	/*
	TODO: The ZOffsetCalibrationPanel and the buttons/functions within it
	are just too buggy.  Commenting out for now and will look into it later.
	zOffsetCalibrationButton := utils.MustButtonImageStyle(
		"Z-Offset Calibration",
		"z-offset-increase.svg",
		"color2",
		this.showZOffsetCalibrationPanel,
	)
	this.Grid().Attach(zOffsetCalibrationButton, 1, 0, 1, 1)
	*/

	networkButton := utils.MustButtonImageStyle(
		"Network",
		"network.svg",
		"color3",
		this.showNetworkPanel,
	)
	this.Grid().Attach(networkButton, 2, 0, 1, 1)

	systemButton := utils.MustButtonImageStyle(
		"System",
		"info.svg",
		"color4",
		this.showSystemPanel,
	)
	this.Grid().Attach(systemButton, 3, 0, 1, 1)
}

func (this *configurationPanel) showBedLevelPanel() {
	this.UI.GoToPanel(BedLevelPanel(this.UI))
}

func (this *configurationPanel) showZOffsetCalibrationPanel() {
	this.UI.GoToPanel(ZOffsetCalibrationPanel(this.UI))
}

func (this *configurationPanel) showNetworkPanel() {
	this.UI.GoToPanel(NetworkPanel(this.UI))
}

func (this *configurationPanel) showSystemPanel() {
	this.UI.GoToPanel(SystemPanel(this.UI))
}
