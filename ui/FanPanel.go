package ui

import (
	// "fmt"

	// "github.com/gotk3/gotk3/gtk"
	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	// "github.com/Z-Bolt/OctoScreen/utils"
)


type fanPanel struct {
	CommonPanel
}

var fanPanelInstance *fanPanel

func GetFanPanelInstance(
	ui				*UI,
) *fanPanel {
	if fanPanelInstance == nil {
		instance := &fanPanel {
			CommonPanel: CreateCommonPanel("FanPanel", ui),
		}
		instance.initialize()
		fanPanelInstance = instance
	}

	return fanPanelInstance
}

func (this *fanPanel) initialize() {
	defer this.Initialize()

	// First row
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Client, 25),  0, 0, 1, 1)
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Client, 50),  1, 0, 1, 1)
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Client, 75),  2, 0, 1, 1)
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Client, 100), 3, 0, 1, 1)

	// Second row
	this.Grid().Attach(uiWidgets.CreateFanButton(this.UI.Client, 0),   0, 1, 1, 1)
}
