package ui

import (
	"github.com/mcuadros/go-octoprint"
)

type customItemsPanel struct {
	CommonPanel
	items []octoprint.MenuItem
}

func CustomItemsPanel(ui *UI, parent Panel, items []octoprint.MenuItem) Panel {
	m := &customItemsPanel{
		CommonPanel: NewCommonPanel(ui, parent),
		items:       items,
	}

	m.initialize()
	return m
}

func (m *customItemsPanel) initialize() {
	defer m.Initialize()
	m.arrangeMenuItems(m.g, m.items, 4)
}
