package ui

import (
	"encoding/json"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	// "github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

func getPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
	menuItem		dataModels.MenuItem,
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


		case "filament_multitool":
			fallthrough
		case "extrude_multitool":
			fallthrough
		case "extruder":
			logger.Warnf("WARNING! the '%s' panel has been deprecated.  Please use the 'filament' panel instead.", menuItem.Panel)
			logger.Warnf("Support for the %s panel remains in this release, but will be removed in a future.", menuItem.Panel)
			logger.Warn("Please update the custom menu structure in your OctoScreen settings in OctoPrint.")
			return FilamentPanel(ui, parentPanel)

		case "toolchanger":
			logger.Warn("WARNING! the 'toolchanger' panel has been renamed to 'tool-changer'.  Please use the 'tool-changer' panel instead.")
			logger.Warnf("Support for the %s panel remains in this release, but will be removed in a future.", menuItem.Panel)
			logger.Warn("Please update the custom menu structure in your OctoScreen settings in OctoPrint.")
			return ToolChangerPanel(ui, parentPanel)

		case "nozzle-calibration":
			logger.Warn("WARNING! the 'nozzle-calibration' panel has been deprecated.  Please use the 'z-offset-calibration' panel instead.")
			logger.Warn("Support for the nozzle-calibration panel remains in this release, but will be removed in a future.")
			logger.Warn("Please update the custom menu structure in your OctoScreen settings in OctoPrint.")
			return ZOffsetCalibrationPanel(ui, parentPanel)

		default:
			logLevel := logger.LogLevel()
			if logLevel == "debug" {
				logger.Fatalf("menu.getPanel() - unknown menuItem.Panel: %q", menuItem.Panel)
			}

			return nil
	}
}

func getDefaultMenuItems(client *octoprintApis.Client) []dataModels.MenuItem {
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


	var menuItems []dataModels.MenuItem
	var err error

	hotendCount := utils.GetHotendCount(client)
	if hotendCount > 1 {
		err = json.Unmarshal([]byte(defaultMenuItemsForMultipleToolheads), &menuItems)
	} else {
		err = json.Unmarshal([]byte(defaultMenuItemsForSingleToolhead), &menuItems)
	}

	if err != nil {
		logger.LogError("menu.getDefaultMenuItems()", "json.Unmarshal()", err)
	}

	return menuItems
}
