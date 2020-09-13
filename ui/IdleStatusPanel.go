package ui

import (
	// "fmt"
	// "sync"
	"time"

	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var idleStatusPanelInstance *idleStatusPanel

type idleStatusPanel struct {
	CommonPanel

	tool0				*uiWidgets.ToolHeatup
	tool1				*uiWidgets.ToolHeatup
	tool2				*uiWidgets.ToolHeatup
	tool3				*uiWidgets.ToolHeatup
	bed					*uiWidgets.ToolHeatup
}

func IdleStatusPanel(ui *UI) *idleStatusPanel {
	if idleStatusPanelInstance == nil {
		instance := &idleStatusPanel{
			CommonPanel: NewTopLevelCommonPanel(ui, nil),
		}
		instance.backgroundTask = utils.CreateBackgroundTask(time.Second * 2, instance.update)
		instance.initialize()

		idleStatusPanelInstance = instance
	}

	return idleStatusPanelInstance
}

func (this *idleStatusPanel) initialize() {
	defer this.Initialize()

	utils.Logger.Info(this.UI.Settings)

	var menuItems []octoprint.MenuItem
	if this.UI.Settings == nil || len(this.UI.Settings.MenuStructure) == 0 {
		utils.Logger.Info("Loading default menu")
		this.UI.Settings.MenuStructure = getDefaultMenuItems(this.UI.Printer)
	} else {
		utils.Logger.Info("Loading octo menu")
	}

	menuItems = this.UI.Settings.MenuStructure

	menuGrid := utils.MustGrid()
	menuGrid.SetRowHomogeneous(true)
	menuGrid.SetColumnHomogeneous(true)
	this.Grid().Attach(menuGrid, 2, 0, 2, 2)
	this.arrangeMenuItems(menuGrid, menuItems, 2)

	printButton := utils.MustButtonImageStyle("Print", "print.svg", "color2", this.showFiles)
	this.Grid().Attach(printButton, 2, 2, 2, 1)

	this.showTools()
}

func (this *idleStatusPanel) showFiles() {
	this.UI.GoToPanel(FilesPanel(this.UI, this))
}

func (this *idleStatusPanel) update() {
	this.updateTemperature()
}

func (this *idleStatusPanel) showTools() {
	toolheadCount := utils.GetToolheadCount(this.UI.Printer)

	if toolheadCount == 1 {
		this.tool0 = uiWidgets.CreteToolHeatupButton(0, this.UI.Printer)
	} else {
		this.tool0 = uiWidgets.CreteToolHeatupButton(1, this.UI.Printer)
	}
	this.tool1 = uiWidgets.CreteToolHeatupButton(2, this.UI.Printer)
	this.tool2 = uiWidgets.CreteToolHeatupButton(3, this.UI.Printer)
	this.tool3 = uiWidgets.CreteToolHeatupButton(4, this.UI.Printer)
	this.bed = uiWidgets.CreteToolHeatupButton( -1, this.UI.Printer)

	switch toolheadCount {
		case 1:
			grid := utils.MustGrid()
			grid.SetRowHomogeneous(true)
			grid.SetColumnHomogeneous(true)
			this.Grid().Attach(grid, 0, 0, 2, 3)
			grid.Attach(this.tool0,  0, 0, 2, 1)
			grid.Attach(this.bed,    0, 1, 2, 1)

		case 2:
			this.Grid().Attach(this.tool0, 0, 0, 2, 1)
			this.Grid().Attach(this.tool1, 0, 1, 2, 1)
			this.Grid().Attach(this.bed,   0, 2, 2, 1)

		case 3:
			this.Grid().Attach(this.tool0, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1, 1, 0, 1, 1)
			this.Grid().Attach(this.tool2, 0, 1, 2, 1)
			this.Grid().Attach(this.bed,   0, 2, 2, 1)

		case 4:
			this.Grid().Attach(this.tool0, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1, 1, 0, 1, 1)
			this.Grid().Attach(this.tool2, 0, 1, 1, 1)
			this.Grid().Attach(this.tool3, 1, 1, 1, 1)
			this.Grid().Attach(this.bed,   0, 2, 2, 1)
	}


	// if toolheadCount == 1 {
	// 	this.tool0 = creteToolHeatupButton(0, this.UI.Printer)
	// } else {
	// 	this.tool0 = creteToolHeatupButton(1, this.UI.Printer)
	// }

	// this.tool1 = creteToolHeatupButton(2, this.UI.Printer)
	// this.tool2 = creteToolHeatupButton(3, this.UI.Printer)
	// this.tool3 = creteToolHeatupButton(4, this.UI.Printer)
	// this.bed   = creteToolHeatupButton(-1, this.UI.Printer)

	// this.Grid().Attach(this.tool0, 0, 0, 1, 1)
	// if toolheadCount >= 2 {
	// 	this.Grid().Attach(this.tool1, 1, 0, 1, 1)
	// }

	// if toolheadCount >= 3 {
	// 	this.Grid().Attach(this.tool2, 0, 1, 1, 1)
	// }

	// if toolheadCount >= 4 {
	// 	this.Grid().Attach(this.tool3, 1, 1, 1, 1)
	// }

	// this.Grid().Attach(this.bed, 0, 2, 1, 1)
}

func (this *idleStatusPanel) updateTemperature() {
	fullStateResponse, err := (&octoprint.FullStateRequest{Exclude: []string{"sd"}}).Do(this.UI.Printer)
	if err != nil {
		utils.LogError("idle_status.updateTemperature()", "Do(StateRequest)", err)
		return
	}

	for tool, currentTemperatureData := range fullStateResponse.Temperature.CurrentTemperatureData {
		switch tool {
			case "bed":
				this.bed.SetTemperatures(currentTemperatureData)

			case "tool0":
				this.tool0.SetTemperatures(currentTemperatureData)

			case "tool1":
				this.tool1.SetTemperatures(currentTemperatureData)

			case "tool2":
				this.tool2.SetTemperatures(currentTemperatureData)

			case "tool3":
				this.tool3.SetTemperatures(currentTemperatureData)
		}
	}
}
