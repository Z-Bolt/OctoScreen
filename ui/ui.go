package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"

	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type UI struct {
	sync.Mutex

	CurrentPanel			interfaces.IPanel
	Printer					*octoprint.Client
	State					octoprint.ConnectionState
	Settings				*octoprint.GetSettingsResponse
	UIState					string

	OctoPrintPlugin			bool

	NotificationsBox		*uiWidgets.NotificationsBox

	splashPanel				*SplashPanel
	backgroundTask			*utils.BackgroundTask
	grid					*gtk.Grid
	window					*gtk.Window
	time					time.Time

	width					int
	height					int
	scaleFactor				int
	connectionAttempts		int
}

func New(endpoint, key string, width, height int) *UI {
	utils.Logger.Info("entering ui.New()")

	if width == 0 || height == 0 {
		width = utils.WindowWidth
		height = utils.WindowHeight
	}

	ui := &UI{
		Printer:				octoprint.NewClient(endpoint, key),
		NotificationsBox:		uiWidgets.NewNotificationsBox(),
		OctoPrintPlugin:		true,
		Settings:				nil,

		window:					utils.MustWindow(gtk.WINDOW_TOPLEVEL),
		time:					time.Now(),

		width:					width,
		height:					height,
	}

	ui.window.Connect("configure-event", func(win *gtk.Window) {
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
			ui.scaleFactor = 2

		case width > 1000:
			ui.scaleFactor = 3

		default:
			ui.scaleFactor = 1
	}

	ui.splashPanel = NewSplashPanel(ui)
	ui.backgroundTask = utils.CreateBackgroundTask(time.Second * 2, ui.update)
	ui.initialize()

	utils.Logger.Info("leaving ui.New()")
	return ui
}

func (ui *UI) initialize() {
	utils.Logger.Info("entering ui.initialize()")

	defer ui.window.ShowAll()
	ui.loadStyle()

	ui.window.SetTitle(utils.WindowName)
	ui.window.SetDefaultSize(ui.width, ui.height)
	ui.window.SetResizable(false)

	ui.window.Connect("show", ui.backgroundTask.Start)
	ui.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	o := utils.MustOverlay()
	ui.window.Add(o)

	ui.grid = utils.MustGrid()
	o.Add(ui.grid)

	ui.sdNotify("READY=1")

	utils.Logger.Info("leaving ui.initialize()")
}

func (ui *UI) loadStyle() {
	utils.Logger.Info("entering ui.loadStyle()")

	p := utils.MustCSSProviderFromFile(utils.CSSFilename)

	s, err := gdk.ScreenGetDefault()
	if err != nil {
		utils.LogError("ui.loadStyle()", "ScreenGetDefault()", err)
		utils.Logger.Error("leaving ui.loadStyle()")
		return
	}

	gtk.AddProviderForScreen(s, p, gtk.STYLE_PROVIDER_PRIORITY_USER)

	utils.Logger.Info("leaving ui.loadStyle()")
}

var errMercyPeriod = time.Second * 10

func (ui *UI) verifyConnection() {
	utils.Logger.Info("entering ui.verifyConnection()")

	ui.sdNotify("WATCHDOG=1")

	newUIState := "<<uninitialized-state>>"
	splashMessage := "<<uninitialized-message>>"

	s, err := (&octoprint.ConnectionRequest{}).Do(ui.Printer)
	if err == nil {
		ui.State = s.Current.State
		strCurrentState := string(s.Current.State)

		switch {
			case s.Current.State.IsOperational():
				newUIState = "idle"
				splashMessage = "Initializing..."

			case s.Current.State.IsPrinting():
				newUIState = "printing"
				splashMessage = "Printing..."

			case s.Current.State.IsError():
				fallthrough
			case s.Current.State.IsOffline():
				newUIState = "splash"
				if err := (&octoprint.ConnectRequest{}).Do(ui.Printer); err != nil {
					utils.LogError("ui.verifyConnection()", "s.Current.State is IsOffline, and (ConnectRequest)Do(ui.Printer)", err)
					splashMessage = "Loading..."
				} else {
					// Use 'Offline' here and 'offline' later.  Having different variations may help in
					// troubleshooting any issues around this state.
					splashMessage = "Printer is Offline."
				}

			case s.Current.State.IsConnecting():
				newUIState = "splash"
				splashMessage = strCurrentState

			default:
				utils.Logger.Fatalf("ui.verifyConnection() - unknown switch of s.Current.State: %q", strCurrentState)
		}
	} else {
		utils.LogError("ui.verifyConnection()", "Broke into the else condition because Do(ConnectionRequest)", err)
		utils.Logger.Info("ui.verifyConnection() - now setting newUIState to 'splash'")
		newUIState = "splash"

		if time.Since(ui.time) > errMercyPeriod {
			errMessage := ui.errToUser(err)

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

	defer func() { ui.UIState = newUIState }()

	ui.splashPanel.Label.SetText(splashMessage)

	if newUIState == ui.UIState {
		utils.Logger.Infof("ui.verifyConnection() - newUIState equals ui.UIState and is: %q", ui.UIState)
		utils.Logger.Info("leaving ui.verifyConnection()")
		return
	}

	utils.Logger.Info("ui.verifyConnection() - newUIState does not equals ui.UIState")
	utils.Logger.Infof("ui.verifyConnection() - ui.UIState is: %q", ui.UIState)
	utils.Logger.Infof("ui.verifyConnection() - newUIState is: %q", newUIState)

	switch newUIState {
		case "idle":
			utils.Logger.Info("ui.verifyConnection() - printer is ready")
			ui.Add(IdleStatusPanel(ui))

		case "printing":
			utils.Logger.Info("ui.verifyConnection() - printing a job")
			ui.Add(PrintStatusPanel(ui))

		case "splash":
			ui.Add(ui.splashPanel)

		default:
			utils.Logger.Fatalf("ui.verifyConnection() - unknown switch of newUIState: %q", newUIState)
	}

	utils.Logger.Info("leaving ui.verifyConnection()")
}

func (m *UI) checkNotification() {
	utils.Logger.Info("entering ui.checkNotification()")

	n, err := (&octoprint.GetNotificationRequest{}).Do(m.Printer)
	if err != nil {
		text := err.Error()
		if strings.Contains(strings.ToLower(text), "unexpected status code: 404") {
			m.OctoPrintPlugin = false
		}

		utils.LogError("ui.checkNotification()", "Do(GetNotificationRequest)", err)
		utils.Logger.Error("leaving ui.checkNotification()")
		return
	}

	if n.Message != "" {
		utils.InfoMessageDialogBox(m.window, n.Message)
	}

	utils.Logger.Info("leaving ui.checkNotification()")
}

func (m *UI) loadSettings() {
	utils.Logger.Info("entering ui.loadSettings()")

	n, err := (&octoprint.GetSettingsRequest{}).Do(m.Printer)
	if err != nil {
		utils.LogError("ui.loadSettings()", "Do(GetSettingsRequest)", err)
		utils.Logger.Error("leaving ui.loadSettings() - Do(GetSettingsRequest) returned an error")
		return
	}

	m.Settings = n

	utils.Logger.Info("leaving ui.loadSettings()")
}

func (m *UI) update() {
	utils.Logger.Info("entering ui.update()")

	if m.connectionAttempts > 8 {
		m.splashPanel.putOnHold()
		utils.Logger.Info("leaving ui.update() - connectionAttempts > 8")
		return
	}

	utils.Logger.Infoln("ui.update() - m.UIState is: ", m.UIState)

	if m.UIState == "splash" {
		m.connectionAttempts++
	} else {
		m.connectionAttempts = 0
	}

	if m.OctoPrintPlugin {
		m.checkNotification()
		m.loadSettings()
	}

	m.verifyConnection()

	utils.Logger.Info("leaving ui.update()")
}

func (ui *UI) sdNotify(m string) {
	utils.Logger.Info("entering ui.sdNotify()")

	_, err := daemon.SdNotify(false, m)
	if err != nil {
		utils.Logger.Errorf("ui.sdNotify()", "SdNotify()", err)
		utils.Logger.Error("leaving ui.sdNotify()")
		return
	}

	utils.Logger.Info("leaving ui.sdNotify()")
}

func (ui *UI) Add(panel interfaces.IPanel) {
	utils.Logger.Info("entering ui.Add()")

	if ui.CurrentPanel != nil {
		ui.Remove(ui.CurrentPanel)
	}

	ui.CurrentPanel = panel
	ui.CurrentPanel.Show()
	ui.grid.Attach(ui.CurrentPanel.Grid(), 0, 0, 1, 1)
	ui.grid.ShowAll()

	utils.Logger.Info("leaving ui.Add()")
}

func (ui *UI) Remove(panel interfaces.IPanel) {
	utils.Logger.Info("entering ui.Remove()")

	defer panel.Hide()
	ui.grid.Remove(panel.Grid())

	utils.Logger.Info("leaving ui.Remove()")
}

func (ui *UI) GoHistory() {
	utils.Logger.Info("entering ui.GoHistory()")

	ui.Add(ui.CurrentPanel.ParentPanel())

	utils.Logger.Info("entering ui.GoHistory()")
}

func (ui *UI) errToUser(err error) string {
	utils.Logger.Info("entering ui.errToUser()")

	text := strings.ToLower(err.Error())
	if strings.Contains(text, "connection refused") {
		utils.Logger.Error("leaving ui.errToUser() - connection refused")
		return "Unable to connect to OctoPrint, check if it running."
	} else if strings.Contains(text, "request canceled") {
		utils.Logger.Error("leaving ui.errToUser() - request canceled")
		return "Loading..."
	} else if strings.Contains(text, "connection broken") {
		utils.Logger.Error("leaving ui.errToUser() - connection broken")
		return "Loading..."
	}

	utils.Logger.Errorf("leaving ui.errToUser() - unexpected error: %q", text)
	return fmt.Sprintf("Unexpected error: %s", err)
}
