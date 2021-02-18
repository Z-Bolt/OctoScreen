package ui

import (
	// "encoding/json"
	// "fmt"
	// "sync"
	"time"

	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


var idleStatusPanelInstance *idleStatusPanel

type idleStatusPanel struct {
	CommonPanel

	tool0Button			*uiWidgets.ToolButton
	tool1Button			*uiWidgets.ToolButton
	tool2Button			*uiWidgets.ToolButton
	tool3Button			*uiWidgets.ToolButton
	bedButton			*uiWidgets.ToolButton
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
	utils.Logger.Debug("entering IdleStatusPanel.initialize()")

	defer this.Initialize()

	utils.Logger.Info("IdleStatusPanel.initialize() - settings are:")
	if this.UI == nil {
		utils.Logger.Error("IdleStatusPanel.initialize() - this.UI is nil")
	} else if this.UI.Settings == nil {
		utils.Logger.Error("IdleStatusPanel.initialize() - this.UI.Settings is nil")
	} else {
		utils.Logger.Info("struct values:")
		utils.Logger.Info(this.UI.Settings)

		jsonStr, err := utils.StructToJson(this.UI.Settings)
		if err == nil {
			utils.Logger.Info("JSON:")
			utils.Logger.Info(jsonStr)
		}
	}
	utils.Logger.Info("")
	utils.Logger.Info("")

	var menuItems []dataModels.MenuItem
	if this.UI.Settings == nil || this.UI.Settings.MenuStructure == nil || len(this.UI.Settings.MenuStructure) < 1 {
		utils.Logger.Info("IdleStatusPanel.initialize() - Loading default menu")
		this.UI.Settings.MenuStructure = getDefaultMenuItems(this.UI.Client)
	} else {
		utils.Logger.Info("IdleStatusPanel.initialize() - Loading octo menu")
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

	utils.Logger.Debug("leaving IdleStatusPanel.initialize()")
}

func (this *idleStatusPanel) showFiles() {
	utils.Logger.Debug("entering IdleStatusPanel.showFiles()")

	this.UI.GoToPanel(FilesPanel(this.UI, this))

	utils.Logger.Debug("leaving IdleStatusPanel.showFiles()")
}

func (this *idleStatusPanel) update() {
	utils.Logger.Debug("entering IdleStatusPanel.update()")

	this.updateTemperature()

	utils.Logger.Debug("leaving IdleStatusPanel.update()")
}

func (this *idleStatusPanel) showTools() {
	utils.Logger.Debug("entering IdleStatusPanel.showTools()")

	// Note: The creation and initialization of the tool buttons in IdleStatusPanel and
	// PrintStatusPanel look similar, but there are subtle differences between the two
	// and they can't be reused.
	hotendCount := utils.GetHotendCount(this.UI.Client)
	if hotendCount == 1 {
		this.tool0Button = uiWidgets.CreateToolButton(0, this.UI.Client)
	} else {
		this.tool0Button = uiWidgets.CreateToolButton(1, this.UI.Client)
	}
	this.tool1Button = uiWidgets.CreateToolButton( 2, this.UI.Client)
	this.tool2Button = uiWidgets.CreateToolButton( 3, this.UI.Client)
	this.tool3Button = uiWidgets.CreateToolButton( 4, this.UI.Client)
	this.bedButton   = uiWidgets.CreateToolButton(-1, this.UI.Client)

	switch hotendCount {
		case 1:
			toolGrid := utils.MustGrid()
			toolGrid.SetRowHomogeneous(true)
			toolGrid.SetColumnHomogeneous(true)
			this.Grid().Attach(toolGrid, 0, 0, 2, 3)
			toolGrid.Attach(this.tool0Button, 0, 0, 2, 1)
			toolGrid.Attach(this.bedButton,   0, 1, 2, 1)

		case 2:
			this.Grid().Attach(this.tool0Button, 0, 0, 2, 1)
			this.Grid().Attach(this.tool1Button, 0, 1, 2, 1)
			this.Grid().Attach(this.bedButton,   0, 2, 2, 1)

		case 3:
			this.Grid().Attach(this.tool0Button, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1Button, 1, 0, 1, 1)
			this.Grid().Attach(this.tool2Button, 0, 1, 2, 1)
			this.Grid().Attach(this.bedButton,   0, 2, 2, 1)

		case 4:
			this.Grid().Attach(this.tool0Button, 0, 0, 1, 1)
			this.Grid().Attach(this.tool1Button, 1, 0, 1, 1)
			this.Grid().Attach(this.tool2Button, 0, 1, 1, 1)
			this.Grid().Attach(this.tool3Button, 1, 1, 1, 1)
			this.Grid().Attach(this.bedButton,   0, 2, 2, 1)
	}

	utils.Logger.Debug("leaving IdleStatusPanel.showTools()")
}

func (this *idleStatusPanel) updateTemperature() {
	utils.Logger.Debug("entering IdleStatusPanel.updateTemperature()")

	fullStateResponse, err := (&octoprintApis.FullStateRequest{Exclude: []string{"sd"}}).Do(this.UI.Client)
	if err != nil {
		utils.LogError("IdleStatusPanel.updateTemperature()", "Do(StateRequest)", err)

		utils.Logger.Debug("leaving IdleStatusPanel.updateTemperature()")
		return
	}

	for tool, currentTemperatureData := range fullStateResponse.Temperature.CurrentTemperatureData {
		switch tool {
			case "bed":
				this.bedButton.SetTemperatures(currentTemperatureData)

			case "tool0":
				this.tool0Button.SetTemperatures(currentTemperatureData)

			case "tool1":
				this.tool1Button.SetTemperatures(currentTemperatureData)

			case "tool2":
				this.tool2Button.SetTemperatures(currentTemperatureData)

			case "tool3":
				this.tool3Button.SetTemperatures(currentTemperatureData)
		}
	}

	utils.Logger.Debug("leaving IdleStatusPanel.updateTemperature()")
}
