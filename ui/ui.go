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
	"github.com/sirupsen/logrus"
)

var (
	StylePath    string
	WindowName   = "OctoScreen"
	WindowHeight = 480
	WindowWidth  = 800
)

const (
	ImageFolder = "images"
	CSSFilename = "style.css"
)

type UI struct {
	Current  Panel
	Printer  *octoprint.Client
	State    octoprint.ConnectionState
	Settings *octoprint.GetSettingsResponse
	UIState  string

	OctoPrintPlugin bool

	Notifications *Notifications

	s *SplashPanel
	b *BackgroundTask
	g *gtk.Grid
	w *gtk.Window
	t time.Time

	width, height      int
	scaleFactor        int
	connectionAttempts int

	sync.Mutex
}

func New(endpoint, key string, width, height int) *UI {
	if width == 0 || height == 0 {
		width = WindowWidth
		height = WindowHeight
	}

	ui := &UI{
		Printer:         octoprint.NewClient(endpoint, key),
		Notifications:   NewNotifications(),
		OctoPrintPlugin: true,
		Settings:        nil,

		w: MustWindow(gtk.WINDOW_TOPLEVEL),
		t: time.Now(),

		width:  width,
		height: height,
	}

	switch {
	case width > 480:
		ui.scaleFactor = 2
	case width > 1000:
		ui.scaleFactor = 3
	default:
		ui.scaleFactor = 1
	}

	ui.s = NewSplashPanel(ui)
	ui.b = NewBackgroundTask(time.Second*2, ui.update)
	ui.initialize()
	return ui
}

func (ui *UI) initialize() {
	defer ui.w.ShowAll()
	ui.loadStyle()

	ui.w.SetTitle(WindowName)
	ui.w.SetDefaultSize(ui.width, ui.height)
	ui.w.SetResizable(false)

	ui.w.Connect("show", ui.b.Start)
	ui.w.Connect("destroy", func() {
		gtk.MainQuit()
	})

	o := MustOverlay()
	ui.w.Add(o)

	ui.g = MustGrid()
	o.Add(ui.g)

	ui.sdNotify("READY=1")
}

func (ui *UI) loadStyle() {
	p := MustCSSProviderFromFile(CSSFilename)

	s, err := gdk.ScreenGetDefault()
	if err != nil {
		logrus.Errorf("Error getting GDK screen: %s", err)
		return
	}

	gtk.AddProviderForScreen(s, p, gtk.STYLE_PROVIDER_PRIORITY_USER)
}

var errMercyPeriod = time.Second * 10

func (ui *UI) verifyConnection() {

	ui.sdNotify("WATCHDOG=1")

	newUiState := "splash"
	splashMessage := "Initializing..."

	s, err := (&octoprint.ConnectionRequest{}).Do(ui.Printer)
	if err == nil {
		ui.State = s.Current.State
		switch {
		case s.Current.State.IsOperational():
			newUiState = "idle"
		case s.Current.State.IsPrinting():
			newUiState = "printing"
		case s.Current.State.IsError():
			fallthrough //Sigh
		case s.Current.State.IsOffline():
			newUiState = "splash"
			splashMessage = "Loading..."
		case s.Current.State.IsConnecting():
			splashMessage = string(s.Current.State)
		}
	} else {
		if time.Since(ui.t) > errMercyPeriod {
			splashMessage = ui.errToUser(err)
		}

		newUiState = "splash"
		Logger.Debugf("Unexpected error: %s", err)
	}

	defer func() { ui.UIState = newUiState }()

	ui.s.Label.SetText(splashMessage)

	if newUiState == ui.UIState {
		return
	}

	switch newUiState {
	case "idle":
		Logger.Info("Printer is ready")
		ui.Add(IdleStatusPanel(ui))
	case "printing":
		Logger.Info("Printing a job")
		ui.Add(PrintStatusPanel(ui))
	case "splash":
		ui.Add(ui.s)
	}
}

func (m *UI) checkNotification() {
	n, err := (&octoprint.GetNotificationRequest{}).Do(m.Printer)
	if err != nil {
		text := err.Error()
		if strings.Contains(text, "unexpected status code: 404") {
			m.OctoPrintPlugin = false
		}
		return
	}

	if n.Message != "" {
		MessageDialog(m.w, n.Message)
	}
}

func (m *UI) loadSettings() {
	n, err := (&octoprint.GetSettingsRequest{}).Do(m.Printer)
	if err != nil {
		Logger.Error(err)
		return
	}
	m.Settings = n
}

func (m *UI) update() {
	if m.connectionAttempts > 8 {
		m.sdNotify("WATCHDOG=1")
		m.s.putOnHold()
		return
	} else if m.UIState == "splash" {
		m.connectionAttempts += 1
	} else {
		m.connectionAttempts = 0
	}

	if m.OctoPrintPlugin {
		m.checkNotification()
		m.loadSettings()
	}

	m.verifyConnection()
}

func (ui *UI) sdNotify(m string) {
	_, err := daemon.SdNotify(false, m)

	if err != nil {
		logrus.Errorf("Error sending notification: %s", err)
		return
	}
}

func (ui *UI) Add(p Panel) {
	if ui.Current != nil {
		ui.Remove(ui.Current)
	}

	ui.Current = p
	ui.Current.Show()
	ui.g.Attach(ui.Current.Grid(), 1, 0, 1, 1)
	ui.g.ShowAll()
}

func (ui *UI) Remove(p Panel) {
	defer p.Hide()
	ui.g.Remove(p.Grid())
}

func (ui *UI) GoHistory() {
	ui.Add(ui.Current.Parent())
}

func (ui *UI) errToUser(err error) string {
	text := err.Error()
	if strings.Contains(text, "connection refused") {
		return "Unable to connect to OctoPrint, check if it running."
	} else if strings.Contains(text, "request canceled") {
		return "Loading..."
	} else if strings.Contains(text, "connection broken") {
		return "Loading..."
	}

	return fmt.Sprintf("Unexpected error: %s", err)
}
