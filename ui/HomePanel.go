package ui

import (
	// "github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

var homePanelInstance *homePanel

type homePanel struct {
	CommonPanel
}

func HomePanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *homePanel {
	if homePanelInstance == nil {
		instance := &homePanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		instance.initialize()
		homePanelInstance = instance
	}

	return homePanelInstance
}

func (this *homePanel) initialize() {
	defer this.Initialize()

	homeXButton := uiWidgets.CreateHomeButton(this.UI.Printer, "Home X", "home-x.svg", octoprint.XAxis)
	this.Grid().Attach(homeXButton, 2, 1, 1, 1)

	homeYButton := uiWidgets.CreateHomeButton(this.UI.Printer, "Home Y", "home-y.svg", octoprint.YAxis)
	this.Grid().Attach(homeYButton, 1, 0, 1, 1)

	homeZButton := uiWidgets.CreateHomeButton(this.UI.Printer, "Home Z", "home-z.svg", octoprint.ZAxis)
	this.Grid().Attach(homeZButton, 1, 1, 1, 1)

	homeAllButton := uiWidgets.CreateHomeAllButton(this.UI.Printer)
	this.Grid().Attach(homeAllButton, 2, 0, 1, 1)
}
