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

		case "filament":
			return FilamentPanel(ui, parentPanel)

		case "configuration":
			return ConfigurationPanel(ui, parentPanel)



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

		case "z-offset-calibration":
			return ZOffsetCalibrationPanel(ui, parentPanel)

		case "print-menu":
			return PrintMenuPanel(ui, parentPanel)

		default:
			logLevel := utils.LowerCaseLogLevel()
			if logLevel == "debug" {
				utils.Logger.Fatalf("menu.getPanel() - unknown menuItem.Panel: %q", menuItem.Panel)
			}

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
				},
				{
					"name": "Tool Changer",
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

	return menuItems
}
