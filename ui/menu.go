package ui

import (
	"encoding/json"

	"github.com/mcuadros/go-octoprint"
)

func getPanel(ui *UI, parent Panel, item octoprint.MenuItem) Panel {
	switch item.Panel {
		// The standard "top four" panels that are in the idleStatus panel
		case "home":
			return HomePanel(ui, parent)

		case "menu":
			fallthrough
		case "custom_items":
			return CustomItemsPanel(ui, parent, item.Items)

		case "filament_multitool":
			fallthrough
		case "filament":
			return FilamentPanel(ui, parent)

		case "configuration":
			return ConfigurationPanel(ui, parent)



		case "extruder-multitool":
			fallthrough
		case "extruder":
			return ExtruderPanel(ui, parent)

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
			fallthrough
		case "z-offset-calibration":
			return ZOffsetCalibrationPanel(ui, parent)

		case "print-menu":
			return PrintMenuPanel(ui, parent)

		default:
			return nil
	}
}

func getDefaultMenu() []octoprint.MenuItem {

	/*
		// do we need a home_multtool panel?
		// do we need a move_multtool panel?
		// do we need a temperature_multtool panel?


		"name": "Extruder",
					"icon": "extruder-multi",
					"panel": "extruder_multitool"

					"icon": "extruder",
					"panel": "extruder"

		"name": "Filament",
			"icon": "filament-spool",
			"panel": "filament_multitool"

			"icon": "filament-spool",
			"panel": "filament"



		{
			"name": "Configuration",
			"icon": "control",
			"panel": "custom_items",
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


	*/



	defaultMenu := `[
		{
			"name": "Home",
			"icon": "home",
			"panel": "home"
		},
		{
			"name": "Actions",
			"icon": "actions",
			"panel": "custom_items",
			"items": [
				{
					"name": "Move",
					"icon": "move",
					"panel": "move"
				},
				{
					"name": "Extruder",
					"icon": "extruder",
					"panel": "extruder"
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
					"icon": "tool-changer",
					"panel": "tool-changer"
				}



				,{
					"name": "PrintMenu",
					"icon": "tool-changer",
					"panel": "print-menu"
				}
			]
		},
		{
			"name": "Filament",
			"icon": "filament-spool",
			"panel": "filament"
		},
		{
			"name": "Configuration",
			"icon": "control",
			"panel": "configuration"
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

	json.Unmarshal([]byte(defaultMenu), &items)

	return items
}
