package ui

import (
	"github.com/mcuadros/go-octoprint"
)

func getPanel(ui *UI, parent Panel, item octoprint.MenuItem) Panel {

	switch item.Panel {
	case "menu":
		return MenuPanel(ui, parent, item.Items)
	case "home":
		return HomePanel(ui, parent)
	case "files":
		return FilesPanel(ui, parent)
	case "filament":
		return FilamentPanel(ui, parent)
	case "temperature":
		return TemperaturePanel(ui, parent)
	case "control":
		return ControlPanel(ui, parent)
	case "network":
		return NetworkPanel(ui, parent)
	case "move":
		return MovePanel(ui, parent)
	case "toolchanger":
		return ToolchangerPanel(ui, parent)
	case "system":
		return SystemPanel(ui, parent)
	case "fan":
		return FanPanel(ui, parent)
	case "extrude":
		return ExtrudePanel(ui, parent)
	case "bed-level":
		return BedLevelPanel(ui, parent)
	case "nozzle-calibration":
		return NozzleCalibrationPanel(ui, parent)
	default:
		return nil
	}
}

type menuPanel struct {
	CommonPanel
	items []octoprint.MenuItem
}

func MenuPanel(ui *UI, parent Panel, items []octoprint.MenuItem) Panel {
	m := &menuPanel{
		CommonPanel: NewCommonPanel(ui, parent),
		items:       items,
	}

	m.panelH = 1 + len(items)/4

	m.initialize()
	return m
}

func (m *menuPanel) initialize() {
	defer m.Initialize()
	m.arrangeMenuItems(m.g, m.items, 4)
}
