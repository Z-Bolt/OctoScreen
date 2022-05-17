package ui

import (
	"fmt"
	// "os"
	// "strconv"
	// "strings"
	// "sync"
	// "time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type connectionPanel struct {
	CommonPanel

	IsCheckingConnection	bool
	backgroundTask			*utils.BackgroundTask

	// First row
	Logo					*gtk.Image

	// Second row
	Label					*gtk.Label

	// Third row
	ActionBar				*gtk.Box
	RetryButton				*gtk.Button
}

var connectionPanelInstance *connectionPanel

func GetConnectionPanelInstance(
	ui				*UI,
) *connectionPanel {
	if connectionPanelInstance == nil {
		connectionPanelInstance = &connectionPanel {
			CommonPanel: CreateCommonPanel("ConnectionPanel", ui),
			IsCheckingConnection: true,
		}
		connectionPanelInstance.initialize()
		connectionPanelInstance.createBackgroundTask()
	}

	return connectionPanelInstance
}

func (this *connectionPanel) initialize() {
	logger.TraceEnter("ConnectionPanel.initialize()")

	_, windowHeight := this.UI.window.GetSize()
	unscaledLogo := utils.MustImageFromFile("logos/octoscreen-logo.svg")
	pixbuf := unscaledLogo.GetPixbuf()
	width := pixbuf.GetWidth()
	height := pixbuf.GetHeight()

	originalLogoWidth := 154.75
	originalLogoHeight := 103.75
	displayHeight := windowHeight / 2.0

	scaleFactor := float64(displayHeight) / originalLogoHeight
	displayWidth := int(originalLogoWidth * scaleFactor)
	displayHeight = int(originalLogoHeight * scaleFactor)

	this.Logo = utils.MustImageFromFileWithSize("logos/octoscreen-logo.svg", displayWidth, displayHeight)

	pixbuf.ScaleSimple(
		this.UI.scaleFactor * width,
		this.UI.scaleFactor * height,
		gdk.INTERP_NEAREST,
	)

	this.Label = utils.MustLabel("Welcome to OctoScreen")
	this.Label.SetHExpand(true)
	this.Label.SetLineWrap(false)
	this.Label.SetMaxWidthChars(60)

	main := utils.MustBox(gtk.ORIENTATION_VERTICAL, 15)
	main.SetHExpand(true)
	main.SetHAlign(gtk.ALIGN_CENTER)
	main.SetVExpand(true)
	main.SetVAlign(gtk.ALIGN_CENTER)

	main.Add(this.Logo)
	main.Add(this.Label)

	this.createActionBar()

	box := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(main)
	box.Add(this.ActionBar)
	this.Grid().Add(box)

	logger.TraceLeave("ConnectionPanel.initialize()")
}

func (this *connectionPanel) createActionBar() {
	logger.TraceEnter("ConnectionPanel.createActionBar()")

	this.ActionBar = utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	this.ActionBar.SetHAlign(gtk.ALIGN_END)

	this.RetryButton = utils.MustButtonImageStyle("Retry", "refresh.svg", "color-none", this.initializeConnectionState)
	this.RetryButton.SetProperty("width-request", this.Scaled(100))

	this.ActionBar.Add(this.RetryButton)

	this.displayButtons(false)

	logger.TraceLeave("ConnectionPanel.createActionBar()")
}

func (this *connectionPanel) displayButtons(display bool) {
	retryButtonStyleContext, _ := this.RetryButton.GetStyleContext()
	if display {
		retryButtonStyleContext.RemoveClass("hidden")
		this.RetryButton.SetSensitive(true)
	} else {
		retryButtonStyleContext.AddClass("hidden")
		this.RetryButton.SetSensitive(false)
	}
}

func (this *connectionPanel) createBackgroundTask() {
	logger.TraceEnter("ConnectionPanel.createBackgroundTask()")

	this.initializeConnectionState()

	// Default timeout of 5 seconds.
	duration := utils.GetExperimentalFrequency(5, "EXPERIMENTAL_CONNECTION_PANEL_UPDATE_FREQUENCY")
	this.backgroundTask = utils.CreateBackgroundTask(duration, this.update)
	this.backgroundTask.Start()

	logger.TraceLeave("ConnectionPanel.createBackgroundTask()")
}

func (this *connectionPanel) update() {
	logger.TraceEnter("ConnectionPanel.update()")

	connectionManager := utils.GetConnectionManagerInstance(this.UI.Client)
	// connectionManager.UpdateStatus()

	msg := ""
	if connectionManager.IsConnectedToOctoPrint != true {
		if connectionManager.ConnectAttempts >= utils.MAX_CONNECTION_ATTEMPTS {
			msg = fmt.Sprintf("Unable to connect to OctoPrint")
			this.displayButtons(true)
		} else if connectionManager.ConnectAttempts == 0 {
			msg = fmt.Sprintf("Attempting to connect to OctoPrint")
		} else {
			msg = fmt.Sprintf("Attempting to connect to OctoPrint...%d", connectionManager.ConnectAttempts + 1)
		}
	} else if connectionManager.IsConnectedToPrinter != true {
		if connectionManager.ConnectAttempts >= utils.MAX_CONNECTION_ATTEMPTS {
			msg = fmt.Sprintf("Unable to connect to the printer")
			this.displayButtons(true)
		} else if connectionManager.ConnectAttempts == 0 {
			msg = fmt.Sprintf("Attempting to connect to the printer")
		} else {
			msg = fmt.Sprintf("Attempting to connect to the printer...%d", connectionManager.ConnectAttempts + 1)
		}
	}
	
	if msg != "" {
		this.Label.SetText(msg)
		connectionManager.UpdateStatus()
	} else {
		currentPanel := this.UI.PanelHistory.Peek().(interfaces.IPanel)
		if currentPanel.Name() == "ConnectionPanel" {
			this.UI.Update()
			this.UI.GoToPanel(GetIdleStatusPanelInstance(this.UI))
		}
	}

	logger.TraceLeave("ConnectionPanel.update()")
}

func (this *connectionPanel) initializeConnectionState() {
	logger.TraceEnter("ConnectionPanel.initializeConnectionState()")

	this.displayButtons(false)

	this.Label.SetText("Attempting to connect to OctoPrint")
	connectionManager := utils.GetConnectionManagerInstance(this.UI.Client)
	connectionManager.ReInitializeConnectionState()

	logger.TraceLeave("ConnectionPanel.initializeConnectionState()")
}

func (this *connectionPanel) showSystem() {
	logger.TraceEnter("ConnectionPanel.showSystem()")

	this.UI.GoToPanel(GetSystemPanelInstance(this.UI))

	logger.TraceLeave("ConnectionPanel.showSystem()")
}
