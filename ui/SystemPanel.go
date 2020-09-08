package ui

import (
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
)

var systemPanelInstance *systemPanel = nil

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

func SystemPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *systemPanel {
	if systemPanelInstance == nil {
		instance := &systemPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		instance.initialize()
		systemPanelInstance = instance
	} else {
		systemPanelInstance.parentPanel = parentPanel
	}

	return systemPanelInstance
}

func (this *systemPanel) initialize() {
	defer this.Initialize()

	// First row
	logoWidth := this.Scaled(52)
	this.octoPrintInfoBox = uiWidgets.CreateOctoPrintInfoBox(this.UI.Printer, logoWidth)
	this.Grid().Attach(this.octoPrintInfoBox,        0, 0, 1, 1)

	this.octoScreenInfoBox = uiWidgets.CreateOctoScreenInfoBox(this.UI.Printer, OctoScreenVersion)
	this.Grid().Attach(this.octoScreenInfoBox,       1, 0, 2, 1)

	this.octoScreenPluginInfoBox = uiWidgets.CreateOctoScreenPluginInfoBox(this.UI.Printer, this.UI.OctoPrintPlugin)
	this.Grid().Attach(this.octoScreenPluginInfoBox, 3, 0, 1, 1)


	// Second row
	this.systemInformationInfoBox = uiWidgets.CreateSystemInformationInfoBox(this.UI.Printer)
	this.Grid().Attach(this.systemInformationInfoBox, 0, 1, 4, 1)


	// Third row
	this.shutdownSystemButton = uiWidgets.CreateSystemCommandButton(
		this.UI.Printer,
		this.UI.window,
		"Shutdown System",
		"shutdown",
		"color-warning-sign-yellow",
	)
	this.Grid().Attach(this.shutdownSystemButton,    0, 2, 1, 1)

	this.rebootSystemButton = uiWidgets.CreateSystemCommandButton(
		this.UI.Printer,
		this.UI.window,
		"Reboot System",
		"reboot",
		"color-warning-sign-yellow",
	)
	this.Grid().Attach(this.rebootSystemButton,      1, 2, 1, 1)

	this.restartOctoPrintButton = uiWidgets.CreateSystemCommandButton(
		this.UI.Printer,
		this.UI.window,
		"Restart OctoPrint",
		"restart",
		"color-warning-sign-yellow",
	)
	this.Grid().Attach(this.restartOctoPrintButton,  2, 2, 1, 1)
}
