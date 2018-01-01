package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/sirupsen/logrus"
)

var (
	StylePath    string
	WindowName   = "OctoPrint-TFT"
	WindowHeight = 320
	WindowWidth  = 480
)

const (
	ImageFolder = "images"
	CSSFilename = "style.css"
)

type UI struct {
	Current       Panel
	Printer       *octoprint.Client
	State         octoprint.ConnectionState
	Notifications *Notifications

	b *BackgroundTask
	g *gtk.Grid
	o *gtk.Overlay
	w *gtk.Window
}

func New(endpoint, key string) *UI {
	ui := &UI{
		Printer:       octoprint.NewClient(endpoint, key),
		Notifications: NewNotifications(),

		w: MustWindow(gtk.WINDOW_TOPLEVEL),
	}

	ui.b = NewBackgroundTask(time.Second*5, ui.verifyConnection)
	ui.initialize()
	return ui
}

func (ui *UI) initialize() {
	defer ui.w.ShowAll()
	ui.loadStyle()

	ui.w.SetTitle(WindowName)
	ui.w.SetDefaultSize(WindowWidth, WindowHeight)
	ui.w.Connect("show", ui.b.Start)
	ui.w.Connect("destroy", func() {
		gtk.MainQuit()
	})

	ui.o = MustOverlay()
	ui.w.Add(ui.o)

	ui.g = MustGrid()
	ui.o.Add(ui.g)
	ui.o.AddOverlay(ui.Notifications)
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

func (ui *UI) verifyConnection() {
	splash := NewSplashPanel(ui)

	s, err := (&octoprint.ConnectionRequest{}).Do(ui.Printer)
	if err != nil {

		splash.Label.SetText(ui.errToUser(err))
		ui.Add(splash)

		// Isn't an error since, error is being displayed already on the panel.
		Logger.Debugf("Unexpected error: %s", err)
		return
	}

	defer func() { ui.State = s.Current.State }()

	switch {
	case s.Current.State.IsOperational():
		if !ui.State.IsOperational() && !ui.State.IsPrinting() {
			Logger.Info("Printer is ready")
			ui.Add(NewDefaultPanel(ui))
		}
		return
	case s.Current.State.IsPrinting():
		if !ui.State.IsPrinting() {
			Logger.Info("Printing a job")
			ui.Add(NewStatusPanel(ui))
		}
		return
	case s.Current.State.IsError():
		fallthrough
	case s.Current.State.IsOffline():
		Logger.Infof("Connection offline, connecting: %s", s.Current.State)
		if err := (&octoprint.ConnectRequest{}).Do(ui.Printer); err != nil {
			splash.Label.SetText(fmt.Sprintf("Error connecting to printer: %s", err))
		}
	case s.Current.State.IsConnecting():
		Logger.Infof("Waiting for connection: %s", s.Current.State)
		splash.Label.SetText(string(s.Current.State))
	}

	ui.Add(splash)
}

func (ui *UI) ShowDefaultPanel() {
	ui.Add(NewDefaultPanel(ui))
}

func (ui *UI) Add(p Panel) {
	if ui.Current != nil {
		ui.Current.Destroy()
	}

	ui.Current = p
	ui.g.Attach(ui.Current.Grid(), 1, 0, 1, 1)
	ui.g.ShowAll()
}

func (ui *UI) errToUser(err error) string {
	text := err.Error()
	if strings.Contains(text, "connection refused") {
		return fmt.Sprintf(
			"Unable to connect to %q (API: %v), maybe OctoPrint not running?",
			ui.Printer.Endpoint, ui.Printer.Endpoint == "",
		)
	}

	return fmt.Sprintf("Unexpected error: %s", err)
}
