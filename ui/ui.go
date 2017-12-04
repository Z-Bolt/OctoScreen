package ui

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var ImagesFolder string

type UI struct {
	Current Panel
	Printer *octoprint.Client
	State   octoprint.ConnectionState

	b *BackgroundTask
	*gtk.Grid
}

func New(endpoint, key string) *UI {
	ui := &UI{
		Grid:    MustGrid(),
		Printer: octoprint.NewClient(endpoint, key),
	}

	ui.b = NewBackgroundTask(time.Second*5, ui.verifyConnection)
	ui.Connect("show", ui.b.Start)
	return ui
}

func (ui *UI) verifyConnection() {
	splash := NewSplashPanel(ui)

	s, err := (&octoprint.ConnectionRequest{}).Do(ui.Printer)
	if err != nil {
		splash.Label.SetText(fmt.Sprintf("Unexpected error: %s", err))
		return
	}

	defer func() { ui.State = s.Current.State }()

	switch {
	case s.Current.State.IsOperational():
		Logger.Debug("Printer is ready")
		if !ui.State.IsOperational() && !ui.State.IsPrinting() {
			ui.Add(NewDefaultPanel(ui))
		}
		return
	case s.Current.State.IsPrinting():
		Logger.Warning("Printing a job")
		if !ui.State.IsPrinting() {
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
	ui.Attach(ui.Current.Grid(), 1, 0, 1, 1)
	ui.ShowAll()
}
