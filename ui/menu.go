package ui

import (
	"encoding/json"

	"github.com/mcuadros/go-octoprint"
)

func getPanel(ui *UI, parent Panel, item octoprint.MenuItem) Panel {
	switch item.Panel {
	case "menu":
		return MenuPanel(ui, parent, item.Items)
	case "home":
		return HomePanel(ui, parent)
	case "filament":
		return FilamentPanel(ui, parent)
	case "filament_multitool":
		return FilamentMultitoolPanel(ui, parent)
	case "extrude":
		return ExtrudePanel(ui, parent)
	case "extrude_multitool":
		return ExtrudeMultitoolPanel(ui, parent)
	case "files":
		return FilesPanel(ui, parent)
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

func getDefaultMenu() []octoprint.MenuItem {
	default_menu := `[
		{
			"name": "Home",
			"icon": "home",
			"panel": "home"
		},
		{
			"name": "Actions",
			"icon": "actions",
			"panel": "menu",
			"items": [
				{
					"name": "Move",
					"icon": "move",
					"panel": "move"
				},
				{
					"name": "Extrude",
					"icon": "filament",
					"panel": "extrude_multitool"
				},
				{
					"name": "Fan",
					"icon": "fan",
					"panel": "fan"
				},
				{
					"name": "Temperature",
					"icon": "heat-up",
					"panel": "temperature"
				},
				{
					"name": "Control",
					"icon": "control",
					"panel": "control"
				},
				{
					"name": "ToolChanger",
					"icon": "toolchanger",
					"panel": "toolchanger"
				}
			]
		},
		{
			"name": "Filament",
			"icon": "filament",
			"panel": "filament_multitool"
		},
		{
			"name": "Configuration",
			"icon": "control",
			"panel": "menu",
			"items": [
				{
					"name": "Bed Level",
					"icon": "bed-level",
					"panel": "bed-level"
				},
				{
					"name": "ZOffsets",
					"icon": "z-offset-increase",
					"panel": "nozzle-calibration"
				},
				{
					"name": "Network",
					"icon": "network",
					"panel": "network"
				},
				{
					"name": "System",
					"icon": "info",
					"panel": "system"
				}
			]
		}
	]`

	// filePath := filepath.Join(os.Getenv("OCTOSCREEN_STYLE_PATH"), "default_menu.json")
	// // filePath := "/etc/octoscreen/config/default_menu.json"
	// jsonFile, err := os.Open(filePath)

	// if err != nil {
	// 	Logger.Info(err)
	// }

	// defer jsonFile.Close()

	// byteValue, err := ioutil.ReadAll(jsonFile)
	// if err != nil {
	// 	Logger.Info("Error in default_menu.json")
	// 	Logger.Info(err)
	// 	return items
	// }

	var items []octoprint.MenuItem

	json.Unmarshal([]byte(default_menu), &items)

	return items
}
