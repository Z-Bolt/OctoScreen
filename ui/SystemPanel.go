package ui

import (
	// "time"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type systemPanel struct {
	CommonPanel

	// First row
	octoPrintInfoBox			*uiWidgets.OctoPrintInfoBox
	octoScreenInfoBox			*uiWidgets.OctoScreenInfoBox
	octoScreenPluginInfoBox		*uiWidgets.OctoScreenPluginInfoBox

	// Second row
	systemInformationInfoBox	*uiWidgets.SystemInformationInfoBox

	// Third row
	shutdownSystemButton		*uiWidgets.SystemCommandButton
	rebootSystemButton			*uiWidgets.SystemCommandButton
	restartOctoPrintButton		*uiWidgets.SystemCommandButton
}

var systemPanelInstance *systemPanel = nil

func GetSystemPanelInstance(
	ui				*UI,
) *systemPanel {
	if systemPanelInstance == nil {
		instance := &systemPanel {
			CommonPanel: CreateCommonPanel("SystemPanel", ui),
		}
		instance.initialize()
		instance.preShowCallback = instance.refreshSystemInformationInfoBox
		systemPanelInstance = instance
	}

	return systemPanelInstance
}

func (this *systemPanel) initialize() {
	logger.TraceEnter("SystemPanel.initialize()")

	defer this.Initialize()

	// First row
	logoWidth := this.Scaled(52)
	this.octoPrintInfoBox = uiWidgets.CreateOctoPrintInfoBox(this.UI.Client, logoWidth)
	this.Grid().Attach(this.octoPrintInfoBox,        0, 0, 1, 1)

	this.octoScreenInfoBox = uiWidgets.CreateOctoScreenInfoBox(this.UI.Client, utils.OctoScreenVersion)
	this.Grid().Attach(this.octoScreenInfoBox,       1, 0, 2, 1)

	this.octoScreenPluginInfoBox = uiWidgets.CreateOctoScreenPluginInfoBox(this.UI.Client, this.UI.UIState, this.UI.OctoPrintPluginIsAvailable)
	this.Grid().Attach(this.octoScreenPluginInfoBox, 3, 0, 1, 1)


	// Second row
	this.systemInformationInfoBox = uiWidgets.CreateSystemInformationInfoBox(this.UI.window, this.UI.scaleFactor)
	this.Grid().Attach(this.systemInformationInfoBox, 0, 1, 4, 1)


	// Third row
	this.shutdownSystemButton = uiWidgets.CreateSystemCommandButton(
		this.UI.Client,
		this.UI.window,
		"Shutdown System",
		"shutdown",
		"color-warning-sign-yellow",
	)
	this.Grid().Attach(this.shutdownSystemButton,    0, 2, 1, 1)

	this.rebootSystemButton = uiWidgets.CreateSystemCommandButton(
		this.UI.Client,
		this.UI.window,
		"Reboot System",
		"reboot",
		"color-warning-sign-yellow",
	)
	this.Grid().Attach(this.rebootSystemButton,      1, 2, 1, 1)

	this.restartOctoPrintButton = uiWidgets.CreateSystemCommandButton(
		this.UI.Client,
		this.UI.window,
		"Restart OctoPrint",
		"restart",
		"color-warning-sign-yellow",
	)
	this.Grid().Attach(this.restartOctoPrintButton,  2, 2, 1, 1)

	logger.TraceLeave("SystemPanel.initialize()")
}


func (this *systemPanel) refreshSystemInformationInfoBox() {
	this.systemInformationInfoBox.Refresh()
}
