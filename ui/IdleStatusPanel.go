package ui

import (
	// "encoding/json"
	// "fmt"
	// "os"
	// "strconv"
	// "sync"
	// "time"

	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type idleStatusPanel struct {
	CommonPanel

	tool0Button			*uiWidgets.ToolButton
	tool1Button			*uiWidgets.ToolButton
	tool2Button			*uiWidgets.ToolButton
	tool3Button			*uiWidgets.ToolButton
	bedButton			*uiWidgets.ToolButton

	backgroundTask		*utils.BackgroundTask
}

var idleStatusPanelInstance *idleStatusPanel

func GetIdleStatusPanelInstance(ui *UI) *idleStatusPanel {
	if idleStatusPanelInstance == nil {
		idleStatusPanelInstance = &idleStatusPanel{
			CommonPanel: CreateTopLevelCommonPanel("IdleStatusPanel", ui),
		}
		idleStatusPanelInstance.initialize()
		idleStatusPanelInstance.createBackgroundTask()
	}

	return idleStatusPanelInstance
}

func (this *idleStatusPanel) initialize() {
	logger.TraceEnter("IdleStatusPanel.initialize()")

	defer this.Initialize()

	logger.Info("IdleStatusPanel.initialize() - settings are:")
	if this.UI == nil {
		logger.Error("IdleStatusPanel.initialize() - this.UI is nil")
	} else if this.UI.Settings == nil {
		logger.Error("IdleStatusPanel.initialize() - this.UI.Settings is nil")
	} else {
		logger.Info("struct values:")
		logger.Info(this.UI.Settings)

		jsonStr, err := utils.StructToJson(this.UI.Settings)
		if err == nil {
			logger.Info("JSON:")
			logger.Info(jsonStr)
		}
	}

	var menuItems []dataModels.MenuItem
	if this.UI.Settings == nil || this.UI.Settings.MenuStructure == nil || len(this.UI.Settings.MenuStructure) < 1 {
		logger.Info("IdleStatusPanel.initialize() - Loading default menu")
		this.UI.Settings.MenuStructure = getDefaultMenuItems(this.UI.Client)
	} else {
		logger.Info("IdleStatusPanel.initialize() - Loading octo menu")
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

	logger.TraceLeave("IdleStatusPanel.initialize()")
}

func (this *idleStatusPanel) showFiles() {
	logger.TraceEnter("IdleStatusPanel.showFiles()")

	this.UI.GoToPanel(GetFilesPanelInstance(this.UI))

	logger.TraceLeave("IdleStatusPanel.showFiles()")
}

func (this *idleStatusPanel) showTools() {
	logger.TraceEnter("IdleStatusPanel.showTools()")

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

	logger.TraceLeave("IdleStatusPanel.showTools()")
}

func (this *idleStatusPanel) createBackgroundTask() {
	logger.TraceEnter("IdleStatusPanel.createBackgroundTask()")

	// Default timeout of 1 second.
	duration := utils.GetExperimentalFrequency(1, "EXPERIMENTAL_IDLE_UPDATE_FREQUENCY")
	this.backgroundTask = utils.CreateBackgroundTask(duration, this.update)
	// Update the UI every second, but the data is only updated once every 10 seconds.
	// See OctoPrintResponseManager.update(). 
	this.backgroundTask.Start()

	logger.TraceLeave("IdleStatusPanel.createBackgroundTask()")
}

func (this *idleStatusPanel) update() {
	logger.TraceEnter("IdleStatusPanel.update()")

	this.updateTemperature()

	logger.TraceLeave("IdleStatusPanel.update()")
}

func (this *idleStatusPanel) updateTemperature() {
	logger.TraceEnter("IdleStatusPanel.updateTemperature()")

	octoPrintResponseManager := GetOctoPrintResponseManagerInstance(this.UI)
	if octoPrintResponseManager.IsConnected() != true {
		// If not connected, do nothing and leave.
		logger.TraceLeave("IdleStatusPanel.updateTemperature() (not connected)")
		return
	}

	for tool, currentTemperatureData := range octoPrintResponseManager.FullStateResponse.Temperature.CurrentTemperatureData {
		switch tool {
			case "bed":
				logger.Debug("Updating the UI's bed temp")
				this.bedButton.SetTemperatures(currentTemperatureData)

			case "tool0":
				logger.Debug("Updating the UI's tool0 temp")
				this.tool0Button.SetTemperatures(currentTemperatureData)

			case "tool1":
				logger.Debug("Updating the UI's tool1 temp")
				this.tool1Button.SetTemperatures(currentTemperatureData)

			case "tool2":
				logger.Debug("Updating the UI's tool2 temp")
				this.tool2Button.SetTemperatures(currentTemperatureData)

			case "tool3":
				logger.Debug("Updating the UI's tool3 temp")
				this.tool3Button.SetTemperatures(currentTemperatureData)

			default:
				logger.Errorf("IdleStatusPanel.updateTemperature() - GetOctoPrintResponseManagerInstance() returned an unknown tool: %q", tool)
		}
	}

	logger.TraceLeave("IdleStatusPanel.updateTemperature()")
}
