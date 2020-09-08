package ui

import (
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
)

type customItemsPanel struct {
	CommonPanel
	items			[]octoprint.MenuItem
}

func CustomItemsPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
	items			[]octoprint.MenuItem,
) *customItemsPanel {
	instance := &customItemsPanel {
		CommonPanel: NewCommonPanel(ui, parentPanel),
		items:       items,
	}
	instance.initialize()

	return instance
}

func (this *customItemsPanel) initialize() {
	defer this.Initialize()
	this.arrangeMenuItems(this.grid, this.items, 4)
}
