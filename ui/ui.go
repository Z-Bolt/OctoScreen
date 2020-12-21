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
	"github.com/mcuadros/go-octoprint"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type UI struct {
	sync.Mutex

	PanelHistory				*stack.Stack
	Client						*octoprint.Client
	ConnectionState				octoprint.ConnectionState
	Settings					*octoprint.GetSettingsResponse

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
	utils.Logger.Debug("entering ui.New()")

	if width == 0 || height == 0 {
		width = utils.WindowWidth
		height = utils.WindowHeight
	}

	instance := &UI{
		PanelHistory:				stack.New(),
		Client:						octoprint.NewClient(endpoint, key),
		NotificationsBox:			uiWidgets.NewNotificationsBox(),
		OctoPrintPluginIsAvailable:	true,
		Settings:					nil,

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
			utils.Logger.Errorf("Widow resize went past max size.  allocatedWidth:%d allocatedHeight:%d sizeWidth:%d sizeHeight:%d",
				allocatedWidth,
				allocatedHeight,
				sizeWidth,
				sizeHeight)
			utils.Logger.Errorf("Widow resize went past max size.  Target width and height: %dx%d",
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

	utils.Logger.Debug("leaving ui.New()")
	return instance
}

func (this *UI) initialize() {
	utils.Logger.Debug("entering ui.initialize()")

	defer this.window.ShowAll()
	this.loadStyle()

	this.window.SetTitle(utils.WindowName)
	this.window.SetDefaultSize(this.width, this.height)
	this.window.SetResizable(false)

	this.window.Connect("show", this.backgroundTask.Start)
	this.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	overlay := utils.MustOverlay()
	this.window.Add(overlay)

	this.grid = utils.MustGrid()
	overlay.Add(this.grid)

	this.sdNotify("READY=1")

	utils.Logger.Debug("leaving ui.initialize()")
}

func (this *UI) loadStyle() {
	utils.Logger.Debug("entering ui.loadStyle()")

	cssProvider := utils.MustCSSProviderFromFile(utils.CSSFilename)

	screenDefault, err := gdk.ScreenGetDefault()
	if err != nil {
		utils.LogError("ui.loadStyle()", "ScreenGetDefault()", err)

		utils.Logger.Debug("leaving ui.loadStyle()")
		return
	}

	gtk.AddProviderForScreen(screenDefault, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)

	utils.Logger.Debug("leaving ui.loadStyle()")
}

var errMercyPeriod = time.Second * 10

func (this *UI) verifyConnection() {
	utils.Logger.Debug("entering ui.verifyConnection()")

	this.sdNotify("WATCHDOG=1")

	newUIState := "<<uninitialized-state>>"
	splashMessage := "<<uninitialized-message>>"

	connectionResponse, err := (&octoprint.ConnectionRequest{}).Do(this.Client)
	if err == nil {
		this.ConnectionState = connectionResponse.Current.State
		strCurrentState := string(connectionResponse.Current.State)

		switch {
			case connectionResponse.Current.State.IsOperational():
				newUIState = "idle"
				splashMessage = "Initializing..."

			case connectionResponse.Current.State.IsPrinting():
				newUIState = "printing"
				splashMessage = "Printing..."

			case connectionResponse.Current.State.IsError():
				fallthrough
			case connectionResponse.Current.State.IsOffline():
				newUIState = "splash"
				if err := (&octoprint.ConnectRequest{}).Do(this.Client); err != nil {
					utils.LogError("ui.verifyConnection()", "s.Current.State is IsOffline, and (ConnectRequest)Do(UI.Client)", err)
					splashMessage = "Loading..."
				} else {
					// Use 'Offline' here and 'offline' later.  Having different variations may help in
					// troubleshooting any issues around this state.
					splashMessage = "Printer is Offline."
				}

			case connectionResponse.Current.State.IsConnecting():
				newUIState = "splash"
				splashMessage = strCurrentState

			default:
				switch strCurrentState {
					case "Cancelling":
						newUIState = "idle"

					default:
						logLevel := utils.LowerCaseLogLevel()
						if logLevel == "debug" {
							utils.Logger.Fatalf("menu.getPanel() - unknown CurrentState: %q", strCurrentState)
						}
				}
		}
	} else {
		utils.LogError("ui.verifyConnection()", "Broke into the else condition because Do(ConnectionRequest)", err)
		utils.Logger.Info("ui.verifyConnection() - now setting newUIState to 'splash'")
		newUIState = "splash"

		if time.Since(this.time) > errMercyPeriod {
			errMessage := this.errToUser(err)

			utils.Logger.Info("ui.verifyConnection() - printer is offline")
			utils.Logger.Infof("ui.verifyConnection() - errMessage is: %q", errMessage)

			if strings.Contains(strings.ToLower(errMessage), "deadline exceeded") {
				// Use 'offline' here, but no ending period.
				splashMessage = "Printer is offline"
			} else {
				splashMessage = errMessage
			}
		} else {
			// Use 'offline.' here and 'offline' above.  Having different variations may help in
			// troubleshooting any issues around this state.
			splashMessage = "Printer is offline."
		}
	}

	defer func() {
		this.UIState = newUIState
	}()

	this.splashPanel.Label.SetText(splashMessage)

	if newUIState == this.UIState {
		utils.Logger.Infof("ui.verifyConnection() - newUIState equals ui.UIState and is: %q", this.UIState)
		utils.Logger.Debug("leaving ui.verifyConnection()")
		return
	}

	utils.Logger.Info("ui.verifyConnection() - newUIState does not equals ui.UIState")
	utils.Logger.Infof("ui.verifyConnection() - ui.UIState is: %q", this.UIState)
	utils.Logger.Infof("ui.verifyConnection() - newUIState is: %q", newUIState)

	switch newUIState {
		case "idle":
			utils.Logger.Info("ui.verifyConnection() - printer is ready")
			this.GoToPanel(IdleStatusPanel(this))

		case "printing":
			utils.Logger.Info("ui.verifyConnection() - printing a job")
			this.GoToPanel(PrintStatusPanel(this))

		case "splash":
			this.GoToPanel(this.splashPanel)

		default:
			logLevel := utils.LowerCaseLogLevel()
			if logLevel == "debug" {
				utils.Logger.Fatalf("ui.verifyConnection() - unknown switch of newUIState: %q", newUIState)
			}
	}

	utils.Logger.Debug("leaving ui.verifyConnection()")
}

func (this *UI) checkNotification() {
	utils.Logger.Debug("entering ui.checkNotification()")

	notificationRespone, err := (&octoprint.GetNotificationRequest{}).Do(this.Client)
	if err != nil {
		text := err.Error()
		if strings.Contains(strings.ToLower(text), "unexpected status code: 404") {
			this.OctoPrintPluginIsAvailable = false
		}

		utils.LogError("ui.checkNotification()", "Do(GetNotificationRequest)", err)

		utils.Logger.Debug("leaving ui.checkNotification()")
		return
	}

	if notificationRespone.Message != "" {
		utils.InfoMessageDialogBox(this.window, notificationRespone.Message)
	}

	utils.Logger.Debug("leaving ui.checkNotification()")
}

func (this *UI) loadSettings() {
	utils.Logger.Debug("entering ui.loadSettings()")

	settingsResponse, err := (&octoprint.GetSettingsRequest{}).Do(this.Client)
	if err != nil {
		utils.LogError("ui.loadSettings()", "Do(GetSettingsRequest)", err)
		utils.Logger.Error("leaving ui.loadSettings() - Do(GetSettingsRequest) returned an error")
		return
	}

	if !this.validateMenuItems(settingsResponse.MenuStructure, "", true) {
		settingsResponse.MenuStructure = nil
	}

	this.Settings = settingsResponse

	utils.Logger.Debug("leaving ui.loadSettings()")
}

func (this *UI) validateMenuItems(menuItems []octoprint.MenuItem, name string, isRoot bool) bool {
	if menuItems == nil {
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

		return false
	}

	for i := 0; i < len(menuItems); i++ {
		menuItem := menuItems[i]
		if menuItem.Panel == "menu" {
			if !this.validateMenuItems(menuItem.Items, menuItem.Name, false) {
				return false
			}
		}
	}

	return true
}

func (this *UI) update() {
	utils.Logger.Debug("entering ui.update()")

	if this.connectionAttempts > 8 {
		this.splashPanel.putOnHold()

		utils.Logger.Debug("leaving ui.update() - connectionAttempts > 8")
		return
	}

	utils.Logger.Infoln("ui.update() - thus.UIState is: ", this.UIState)

	if this.UIState == "splash" {
		this.connectionAttempts++
	} else {
		this.connectionAttempts = 0
	}

	if this.OctoPrintPluginIsAvailable {
		this.checkNotification()
		this.loadSettings()
	}

	this.verifyConnection()

	utils.Logger.Debug("leaving ui.update()")
}

func (this *UI) sdNotify(state string) {
	utils.Logger.Debug("entering ui.sdNotify()")

	_, err := daemon.SdNotify(false, state)
	if err != nil {
		utils.Logger.Errorf("ui.sdNotify()", "SdNotify()", err)
		utils.Logger.Debug("leaving ui.sdNotify()")
		return
	}

	utils.Logger.Debug("leaving ui.sdNotify()")
}

func (this *UI) GoToPanel(panel interfaces.IPanel) {
	utils.Logger.Debug("entering ui.GoToPanel()")

	this.SetUiToPanel(panel)
	this.PanelHistory.Push(panel)

	utils.Logger.Debug("leaving ui.GoToPanel()")
}

func (this *UI) GoToPreviousPanel() {
	utils.Logger.Debug("entering ui.GoToPreviousPanel()")

	stackLength := this.PanelHistory.Len()
	if stackLength < 2 {
		utils.Logger.Error("ui.GoToPreviousPanel() - stack does not contain current panel and parent panel")

		utils.Logger.Debug("leaving ui.GoToPreviousPanel()")
		return
	}

	if stackLength < 1 {
		utils.Logger.Error("ui.GoToPreviousPanel() - GoToPreviousPanel() was called but the stack is empty")

		utils.Logger.Debug("leaving ui.GoToPreviousPanel()")
		return
	}

	currentPanel := this.PanelHistory.Pop().(interfaces.IPanel)
	this.RemovePanelFromUi(currentPanel)

	parentPanel := this.PanelHistory.Peek().(interfaces.IPanel)
	this.SetUiToPanel(parentPanel)

	utils.Logger.Debug("leaving ui.GoToPreviousPanel()")
}

func (this *UI) SetUiToPanel(panel interfaces.IPanel) {
	utils.Logger.Debug("entering ui.SetUiToPanel()")

	stackLength := this.PanelHistory.Len()
	if stackLength > 0 {
		currentPanel := this.PanelHistory.Peek().(interfaces.IPanel)
		this.RemovePanelFromUi(currentPanel)
	}

	panel.PreShow()
	panel.Show()
	this.grid.Attach(panel.Grid(), 0, 0, 1, 1)
	this.grid.ShowAll()

	utils.Logger.Debug("leaving ui.SetUiToPanel()")
}

func (this *UI) RemovePanelFromUi(panel interfaces.IPanel) {
	utils.Logger.Debug("entering ui.RemovePanelFromUi()")

	defer panel.Hide()
	this.grid.Remove(panel.Grid())

	utils.Logger.Debug("leaving ui.RemovePanelFromUi()")
}

func (this *UI) errToUser(err error) string {
	utils.Logger.Debug("entering ui.errToUser()")

	text := strings.ToLower(err.Error())
	if strings.Contains(text, "connection refused") {
		utils.Logger.Debug("leaving ui.errToUser() - connection refused")
		return "Unable to connect to OctoPrint, check if it running."
	} else if strings.Contains(text, "request canceled") {
		utils.Logger.Debug("leaving ui.errToUser() - request canceled")
		return "Loading..."
	} else if strings.Contains(text, "connection broken") {
		utils.Logger.Debug("leaving ui.errToUser() - connection broken")
		return "Loading..."
	}

	utils.Logger.Debugf("leaving ui.errToUser() - unexpected error: %q", text)
	return fmt.Sprintf("Unexpected error: %s", err)
}
