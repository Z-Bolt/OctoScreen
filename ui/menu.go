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
			return GetHomePanelInstance(ui)

		case "menu":
			fallthrough
		case "custom_items":
			return CreateCustomItemsPanel(ui, menuItem.Items)

		case "filament":
			return GetFilamentPanelInstance(ui)

		case "configuration":
			return GetConfigurationPanelInstance(ui)



		case "files":
			return GetFilesPanelInstance(ui)

		case "temperature":
			return GetTemperaturePanelInstance(ui)

		case "control":
			return GetControlPanelInstance(ui)

		case "network":
			return GetNetworkPanelInstance(ui)

		case "move":
			return GetMovePanelInstance(ui)

		case "tool-changer":
			return GetToolChangerPanelInstance(ui)

		case "system":
			return GetSystemPanelInstance(ui)

		case "fan":
			return GetFanPanelInstance(ui)

		case "bed-level":
			return GetBedLevelPanelInstance(ui)

		case "z-offset-calibration":
			return GetZOffsetCalibrationPanelInstance(ui)

		case "print-menu":
			return GetPrintMenuPanelInstance(ui)


		case "filament_multitool":
			fallthrough
		case "extrude_multitool":
			fallthrough
		case "extruder":
			logger.Warnf("WARNING! the '%s' panel has been deprecated.  Please use the 'filament' panel instead.", menuItem.Panel)
			logger.Warnf("Support for the %s panel remains in this release, but will be removed in a future.", menuItem.Panel)
			logger.Warn("Please update the custom menu structure in your OctoScreen settings in OctoPrint.")
			return GetFilamentPanelInstance(ui)

		case "toolchanger":
			logger.Warn("WARNING! the 'toolchanger' panel has been renamed to 'tool-changer'.  Please use the 'tool-changer' panel instead.")
			logger.Warnf("Support for the %s panel remains in this release, but will be removed in a future.", menuItem.Panel)
			logger.Warn("Please update the custom menu structure in your OctoScreen settings in OctoPrint.")
			return GetToolChangerPanelInstance(ui)

		case "nozzle-calibration":
			logger.Warn("WARNING! the 'nozzle-calibration' panel has been deprecated.  Please use the 'z-offset-calibration' panel instead.")
			logger.Warn("Support for the nozzle-calibration panel remains in this release, but will be removed in a future.")
			logger.Warn("Please update the custom menu structure in your OctoScreen settings in OctoPrint.")
			return GetZOffsetCalibrationPanelInstance(ui)

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
