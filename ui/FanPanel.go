package ui

import (
	// "fmt"

	// "github.com/gotk3/gotk3/gtk"
	// "github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

var fanPanelInstance *fanPanel

type fanPanel struct {
	CommonPanel
}

func FanPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *fanPanel {
	if fanPanelInstance == nil {
		instance := &fanPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		instance.initialize()
		fanPanelInstance = instance
	}

	return fanPanelInstance
}

func (this *fanPanel) initialize() {
	defer this.Initialize()

	// First row
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Printer, 25),  0, 0, 1, 1)
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Printer, 50),  1, 0, 1, 1)
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Printer, 75),  2, 0, 1, 1)
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Printer, 100), 3, 0, 1, 1)

	// Second row
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Printer, 0),   0, 1, 1, 1)
}
