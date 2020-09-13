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
	this.AddButton(homeXButton)

	homeYButton := uiWidgets.CreateHomeButton(this.UI.Printer, "Home Y", "home-y.svg", octoprint.YAxis)
	this.AddButton(homeYButton)

	homeZButton := uiWidgets.CreateHomeButton(this.UI.Printer, "Home Z", "home-z.svg", octoprint.ZAxis)
	this.AddButton(homeZButton)

	homeAllButton := uiWidgets.CreateHomeAllButton(this.UI.Printer)
	this.AddButton(homeAllButton)
}
