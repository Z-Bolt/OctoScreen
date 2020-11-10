package ui

import (
	"encoding/json"

	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

func getPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
	menuItem		octoprint.MenuItem,
) interfaces.IPanel {
	switch menuItem.Panel {
		// The standard "top four" panels that are in the idleStatus panel
		case "home":
			return HomePanel(ui, parentPanel)

		case "menu":
			fallthrough
		case "custom_items":
			return CustomItemsPanel(ui, parentPanel, menuItem.Items)

		case "filament_multitool":
			fallthrough
		case "filament":
			return FilamentPanel(ui, parentPanel)

		case "configuration":
			return ConfigurationPanel(ui, parentPanel)



		// obsolete and removed, since everything in ExtruderPanel
		// is already in FilamentPanel
		// case "extruder-multitool":
		// 	fallthrough
		// case "extruder":
		// 	return ExtruderPanel(ui, parentPanel)

		case "files":
			return FilesPanel(ui, parentPanel)

		case "temperature":
			return TemperaturePanel(ui, parentPanel)

		case "control":
			return ControlPanel(ui, parentPanel)

		case "network":
			return NetworkPanel(ui, parentPanel)

		case "move":
			return MovePanel(ui, parentPanel)

		case "tool-changer":
			return ToolChangerPanel(ui, parentPanel)

		case "system":
			return SystemPanel(ui, parentPanel)

		case "fan":
			return FanPanel(ui, parentPanel)

		case "bed-level":
			return BedLevelPanel(ui, parentPanel)

		case "nozzle-calibration":
			fallthrough
		case "z-offset-calibration":
			return ZOffsetCalibrationPanel(ui, parentPanel)

		case "print-menu":
			return PrintMenuPanel(ui, parentPanel)

		default:
			return nil
	}
}

func getDefaultMenuItems(client *octoprint.Client) []octoprint.MenuItem {
	defaultMenuItemsForSingleToolhead := `[
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
					"name": "Filament",
					"icon": "filament-spool",
					"panel": "filament"
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

	defaultMenuItemsForMultipleToolheads := `[
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


	// TODO: removed commented code
	// standardLog.Printf("Logger.Level: %q", utils.Logger.GetLogLevel())


	var menuItems []octoprint.MenuItem
	var err error

	toolheadCount := utils.GetToolheadCount(client)
	if toolheadCount > 1 {
		err = json.Unmarshal([]byte(defaultMenuItemsForMultipleToolheads), &menuItems)
	} else {
		err = json.Unmarshal([]byte(defaultMenuItemsForSingleToolhead), &menuItems)
	}

	if err != nil {
		utils.LogError("menu.getDefaultMenuItems()", "json.Unmarshal()", err)
	}
	// utils.LogError("menu.getDefaultMenuItems()", "now leaving", err)


	return menuItems
}
