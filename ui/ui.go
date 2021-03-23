package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/golang-collections/collections/stack"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var Version string
var Build   string

type UI struct {
	sync.Mutex

	PanelHistory				*stack.Stack
	Client						*octoprintApis.Client
	ConnectionState				dataModels.ConnectionState
	Settings					*dataModels.OctoScreenSettingsResponse

	UIState						string

	OctoPrintPluginIsAvailable	bool

	NotificationsBox			*uiWidgets.NotificationsBox

	splashPanel					*SplashPanel
	backgroundTask				*utils.BackgroundTask
	grid						*gtk.Grid
	window						*gtk.Window
	time						time.Time

	width						int
	height						int
	scaleFactor					int
	connectionAttempts			int
}

func New(endpoint, key string, width, height int) *UI {
	logger.TraceEnter("ui.New()")

	if width == 0 || height == 0 {
		width = utils.WindowWidth
		height = utils.WindowHeight
	}

	instance := &UI {
		PanelHistory:				stack.New(),
		Client:						octoprintApis.NewClient(endpoint, key),
		NotificationsBox:			uiWidgets.NewNotificationsBox(),
		OctoPrintPluginIsAvailable:	false,
		Settings:					nil,

		UIState:					"__uninitialized__",

		window:						utils.MustWindow(gtk.WINDOW_TOPLEVEL),
		time:						time.Now(),

		width:						width,
		height:						height,
	}

	instance.window.Connect("configure-event", func(win *gtk.Window) {
		allocatedWidth:= win.GetAllocatedWidth()
		allocatedHeight:= win.GetAllocatedHeight()
		sizeWidth, sizeHeight := win.GetSize()

		if (allocatedWidth > width || allocatedHeight > height) ||
			(sizeWidth > width || sizeHeight > height) {
			logger.Errorf("Window resize went past max size.  allocatedWidth:%d allocatedHeight:%d sizeWidth:%d sizeHeight:%d",
				allocatedWidth,
				allocatedHeight,
				sizeWidth,
				sizeHeight)
			logger.Errorf("Window resize went past max size.  Target width and height: %dx%d",
				width,
				height)
		}
	})

	switch {
		case width > 480:
			instance.scaleFactor = 2

		case width > 1000:
			instance.scaleFactor = 3

		default:
			instance.scaleFactor = 1
	}

	instance.splashPanel = NewSplashPanel(instance)
	instance.backgroundTask = utils.CreateBackgroundTask(time.Second * 10, instance.update)
	instance.initialize()

	logger.TraceLeave("ui.New()")
	return instance
}

func (this *UI) initialize() {
	logger.TraceEnter("ui.initialize()")

	defer this.window.ShowAll()
	this.loadStyle()

	this.window.SetTitle(utils.WindowName)
	this.window.SetDefaultSize(this.width, this.height)
	this.window.SetResizable(false)

	this.window.Connect("show", this.backgroundTask.Start)
	this.window.Connect("destroy", func() {
		logger.Debug("window destroy callback was called, now executing MainQuit()")
		gtk.MainQuit()
	})

	overlay := utils.MustOverlay()
	this.window.Add(overlay)

	this.grid = utils.MustGrid()
	overlay.Add(this.grid)

	this.sdNotify(daemon.SdNotifyReady)

	logger.TraceLeave("ui.initialize()")
}

func (this *UI) loadStyle() {
	logger.TraceEnter("ui.loadStyle()")

	cssProvider := utils.MustCSSProviderFromFile(utils.CSSFilename)

	screenDefault, err := gdk.ScreenGetDefault()
	if err != nil {
		logger.LogError("ui.loadStyle()", "ScreenGetDefault()", err)
		logger.TraceLeave("ui.loadStyle()")
		return
	}

	gtk.AddProviderForScreen(screenDefault, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)

	logger.TraceLeave("ui.loadStyle()")
}

var errMercyPeriod = time.Second * 10

func (this *UI) verifyConnection() {
	logger.TraceEnter("ui.verifyConnection()")

	this.sdNotify(daemon.SdNotifyWatchdog)

	newUIState := "<<uninitialized-state>>"
	splashMessage := "<<uninitialized-message>>"

	connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(this.Client)
	if err == nil {
		logger.Debug("ui.verifyConnection() - ConnectionRequest.Do() succeeded")
		jsonResponse, err := utils.StructToJson(connectionResponse)
		if err != nil {
			logger.Debug("ui.verifyConnection() - utils.StructToJson() failed")
		} else {
			logger.Debugf("ui.verifyConnection() - connectionResponse is: %s", jsonResponse)
		}

		this.ConnectionState = connectionResponse.Current.State
		newUIState, splashMessage = this.getUiStateAndMessageFromConnectionResponse(connectionResponse, newUIState, splashMessage)
	} else {
		logger.LogError("ui.verifyConnection()", "Broke into the else condition because Do(ConnectionRequest) returned an error", err)
		newUIState, splashMessage = this.getUiStateAndMessageFromError(err, newUIState, splashMessage)
	}

	this.splashPanel.Label.SetText(splashMessage)

	defer func() {
		this.setUiState(newUIState, splashMessage)
	}()

	logger.TraceLeave("ui.verifyConnection()")
}


func (this *UI) getUiStateAndMessageFromConnectionResponse(
	connectionResponse *dataModels.ConnectionResponse,
	newUIState string,
	splashMessage string,
) (string, string) {
	logger.TraceEnter("ui.getUiStateAndMessageFromConnectionResponse()")

	strCurrentState := string(connectionResponse.Current.State)
	logger.Debugf("ui.getUiStateAndMessageFromConnectionResponse() - strCurrentState is %s", strCurrentState)

	switch {
		case connectionResponse.Current.State.IsOperational():
			logger.Debug("ui.getUiStateAndMessageFromConnectionResponse() - new state is idle")
			newUIState = "idle"
			splashMessage = "Initializing..."

		case connectionResponse.Current.State.IsPrinting():
			logger.Debug("ui.getUiStateAndMessageFromConnectionResponse() - new state is printing")
			newUIState = "printing"
			splashMessage = "Printing..."

		case connectionResponse.Current.State.IsError():
			logger.Debug("ui.getUiStateAndMessageFromConnectionResponse() - the state has an error")
			fallthrough
		case connectionResponse.Current.State.IsOffline():
			logger.Debug("ui.getUiStateAndMessageFromConnectionResponse() - the state is now offline and displaying the splash panel")
			newUIState = "splash"
			logger.Info("ui.getUiStateAndMessageFromConnectionResponse() - new UI state is 'splash' and is about to call ConnectRequest.Do()")
			if err := (&octoprintApis.ConnectRequest{}).Do(this.Client); err != nil {
				logger.LogError("ui.getUiStateAndMessageFromConnectionResponse()", "s.Current.State is IsOffline, and (ConnectRequest)Do(UI.Client)", err)
				splashMessage = "Loading..."
			} else {
				splashMessage = "Printer is offline, now trying to connect..."
			}

		case connectionResponse.Current.State.IsConnecting():
			logger.Debug("ui.getUiStateAndMessageFromConnectionResponse() - new state is splash (from IsConnecting)")
			newUIState = "splash"
			splashMessage = strCurrentState

		default:
			logger.Debug("ui.getUiStateAndMessageFromConnectionResponse() - the default case was hit")
			switch strCurrentState {
				case "Cancelling":
					newUIState = "idle"

				default:
					logger.Errorf("ui.getUiStateAndMessageFromConnectionResponse() - unknown CurrentState: %q", strCurrentState)
			}
	}

	logger.TraceLeave("ui.getUiStateAndMessageFromConnectionResponse()")
	return newUIState, splashMessage
}


func (this *UI) getUiStateAndMessageFromError(
	err error,
	newUIState string,
	splashMessage string,
) (string, string) {
	logger.TraceEnter("ui.getUiStateAndMessageFromError()")

	logger.Info("ui.getUiStateAndMessageFromError() - now setting newUIState to 'splash'")
	newUIState = "splash"

	if time.Since(this.time) > errMercyPeriod {
		errMessage := this.errToUser(err)

		logger.Info("ui.getUiStateAndMessageFromError() - printer is offline")
		logger.Infof("ui.getUiStateAndMessageFromError() - errMessage is: %q", errMessage)

		if strings.Contains(strings.ToLower(errMessage), "deadline exceeded") {
			splashMessage = "Printer is offline (deadline exceeded), retrying to connect..."
		} else if strings.Contains(strings.ToLower(errMessage), "connection reset by peer") {
			splashMessage = "Printer is offline (peer connection reset), retrying to connect..."
		} else if strings.Contains(strings.ToLower(errMessage), "unexpected status code: 403") {
			splashMessage = "Printer is offline (403), retrying to connect..."
		} else {
			splashMessage = errMessage
		}
	} else {
		splashMessage = "Printer is offline! (retrying to connect...)"
	}

	logger.TraceLeave("ui.getUiStateAndMessageFromError()")
	return newUIState, splashMessage
}


func (this *UI) setUiState(
	newUiState string,
	splashMessage string,
) {
	logger.TraceEnter("ui.setUiState()")

	this.splashPanel.Label.SetText(splashMessage)

	if newUiState == this.UIState {
		logger.Infof("ui.setUiState() - newUiState and ui.UIState are the same (%q)", this.UIState)
		logger.TraceLeave("ui.setUiState()")
		return
	}

	logger.Info("ui.setUiState() - newUiState does not equal ui.UIState")
	logger.Infof("ui.setUiState() - ui.UIState is: %q", this.UIState)
	logger.Infof("ui.setUiState() - newUiState is: %q", newUiState)
	this.UIState = newUiState

	switch newUiState {
		case "idle":
			logger.Info("ui.setUiState() - printer is ready")
			this.GoToPanel(IdleStatusPanel(this))

		case "printing":
			logger.Info("ui.setUiState() - printing a job")
			this.GoToPanel(PrintStatusPanel(this))

		case "splash":
			this.GoToPanel(this.splashPanel)

		default:
			logger.Errorf("ERROR: ui.setUiState() - unknown newUiState case: %q", newUiState)
	}

	logger.TraceLeave("ui.setUiState()")
}


func (this *UI) checkNotification() {
	logger.TraceEnter("ui.checkNotification()")

	if !this.OctoPrintPluginIsAvailable {
		logger.Info("ui.checkNotification() - OctoPrintPluginIsAvailable is false, so not calling GetNotification")
		logger.TraceLeave("ui.checkNotification()")
		return
	}

	notificationResponse, err := (&octoprintApis.NotificationRequest{}).Do(this.Client, this.UIState)
	if err != nil {
		logger.LogError("ui.checkNotification()", "Do(GetNotificationRequest)", err)
		logger.TraceLeave("ui.checkNotification()")
		return
	}

	if notificationResponse != nil && notificationResponse.Message != "" {
		utils.InfoMessageDialogBox(this.window, notificationResponse.Message)
	}

	logger.TraceLeave("ui.checkNotification()")
}

func (this *UI) loadSettings() {
	logger.TraceEnter("ui.loadSettings()")

	settingsResponse, err := (&octoprintApis.OctoScreenSettingsRequest{}).Do(this.Client, this.UIState)
	if err != nil {
		text := err.Error()
		if strings.Contains(strings.ToLower(text), "unexpected status code: 404") {
			// The call to GetSettings is also used to determine whether or not the
			// OctoScreen plug-in is available.  If calling GetSettings returns
			// a 404, the plug-in isn't available.
			logger.Info("The OctoScreen plug-in is not available")
		} else {
			// If we get back any other kind of error, something bad happened, so log an error.
			logger.LogError("ui.loadSettings()", "Do(GetSettingsRequest)", err)
		}

		this.OctoPrintPluginIsAvailable = false

		logger.TraceLeave("ui.loadSettings()")
		return
	} else {
		logger.Info("The call to GetSettings succeeded and the OctoPrint plug-in is available")
		this.OctoPrintPluginIsAvailable = true
	}

	if !this.validateMenuItems(settingsResponse.MenuStructure, "", true) {
		settingsResponse.MenuStructure = nil
	}

	this.Settings = settingsResponse

	logger.TraceLeave("ui.loadSettings()")
}

func (this *UI) validateMenuItems(menuItems []dataModels.MenuItem, name string, isRoot bool) bool {
	logger.TraceEnter("ui.validateMenuItems()")

	if menuItems == nil {
		logger.TraceLeave("ui.validateMenuItems()")
		return true
	}

	maxCount := 11
	if isRoot {
		maxCount = 4
	}

	menuItemsLength := len(menuItems)
	if menuItemsLength > maxCount {
		message := ""
		description := ""
		if isRoot {
			message = fmt.Sprintf("Error!  The custom menu structure can only have %d items\n    at the root level (the idle panel).", maxCount)
			description = fmt.Sprintf("\n    When the MenuStructure was parsed, %d items were found.", menuItemsLength)
		} else {
			message = fmt.Sprintf("Error!  A panel can only have a maximum of %d items.", maxCount)
			description = fmt.Sprintf("\n    When the MenuStructure for '%s' was parsed,\n    %d items were found.", name, menuItemsLength)
		}

		fatalErrorWindow := CreateFatalErrorWindow(
			message,
			description,
		)
		fatalErrorWindow.ShowAll()

		logger.TraceLeave("ui.validateMenuItems()")
		return false
	}

	for i := 0; i < len(menuItems); i++ {
		menuItem := menuItems[i]
		if menuItem.Panel == "menu" {
			if !this.validateMenuItems(menuItem.Items, menuItem.Name, false) {
				logger.TraceLeave("ui.validateMenuItems()")
				return false
			}
		}
	}

	logger.TraceLeave("ui.validateMenuItems()")
	return true
}

func (this *UI) update() {
	logger.TraceEnter("ui.update()")

	if this.connectionAttempts > 8 {
		logger.Info("ui.update() - this.connectionAttempts > 8")
		this.splashPanel.putOnHold()
		this.sdNotify(daemon.SdNotifyWatchdog)

		logger.TraceLeave("ui.update()")
		return
	}

	logger.Infof("ui.update() - this.UIState is: %q", this.UIState)

	if this.UIState == "splash" {
		this.connectionAttempts++
	} else {
		this.connectionAttempts = 0
	}

	if this.Settings == nil {
		this.loadSettings()

		if this.Settings == nil {
			this.Settings = &dataModels.OctoScreenSettingsResponse {
				FilamentInLength: 100,
				FilamentOutLength: 100,
				ToolChanger: false,
				XAxisInverted: false,
				YAxisInverted: false,
				ZAxisInverted: false,
				MenuStructure: nil,
			}
		}
	}

	if this.OctoPrintPluginIsAvailable {
		this.checkNotification()
	}

	this.verifyConnection()

	logger.TraceLeave("ui.update()")
}

func (this *UI) sdNotify(state string) {
	logger.TraceEnter("ui.sdNotify()")

	_, err := daemon.SdNotify(false, state)
	if err != nil {
		logger.Errorf("ui.sdNotify()", "SdNotify()", err)
	}

	logger.TraceLeave("ui.sdNotify()")
}

func (this *UI) GoToPanel(panel interfaces.IPanel) {
	logger.TraceEnter("ui.GoToPanel()")

	this.SetUiToPanel(panel)
	this.PanelHistory.Push(panel)

	logger.TraceLeave("ui.GoToPanel()")
}

func (this *UI) GoToPreviousPanel() {
	logger.TraceEnter("ui.GoToPreviousPanel()")

	stackLength := this.PanelHistory.Len()
	if stackLength < 2 {
		logger.Error("ui.GoToPreviousPanel() - stack does not contain current panel and parent panel")
		logger.TraceLeave("ui.GoToPreviousPanel()")
		return
	}

	if stackLength < 1 {
		logger.Error("ui.GoToPreviousPanel() - GoToPreviousPanel() was called but the stack is empty")
		logger.TraceLeave("ui.GoToPreviousPanel()")
		return
	}

	currentPanel := this.PanelHistory.Pop().(interfaces.IPanel)
	this.RemovePanelFromUi(currentPanel)

	parentPanel := this.PanelHistory.Peek().(interfaces.IPanel)
	this.SetUiToPanel(parentPanel)

	logger.TraceLeave("ui.GoToPreviousPanel()")
}

func (this *UI) SetUiToPanel(panel interfaces.IPanel) {
	logger.TraceEnter("ui.SetUiToPanel()")

	stackLength := this.PanelHistory.Len()
	if stackLength > 0 {
		currentPanel := this.PanelHistory.Peek().(interfaces.IPanel)
		this.RemovePanelFromUi(currentPanel)
	}

	panel.PreShow()
	panel.Show()
	this.grid.Attach(panel.Grid(), 0, 0, 1, 1)
	this.grid.ShowAll()

	logger.TraceLeave("ui.SetUiToPanel()")
}

func (this *UI) RemovePanelFromUi(panel interfaces.IPanel) {
	logger.TraceEnter("ui.RemovePanelFromUi()")

	defer panel.Hide()
	this.grid.Remove(panel.Grid())

	logger.TraceLeave("ui.RemovePanelFromUi()")
}

func (this *UI) errToUser(err error) string {
	logger.TraceEnter("ui.errToUser()")

	text := strings.ToLower(err.Error())
	if strings.Contains(text, "connection refused") {
		logger.TraceLeave("ui.errToUser() - connection refused")
		return "Unable to connect to OctoPrint, check if it running."
	} else if strings.Contains(text, "request canceled") {
		logger.TraceLeave("ui.errToUser() - request canceled")
		return "Loading..."
	} else if strings.Contains(text, "connection broken") {
		logger.TraceLeave("ui.errToUser() - connection broken")
		return "Loading..."
	}

	msg := fmt.Sprintf("ui.errToUser() - unexpected error: %s", text)
	logger.TraceLeave(msg)
	return fmt.Sprintf("Unexpected Error: %s", text)
}
