package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
)

var ImagesFolder string

type UI struct {
	Current Panel
	Printer *octoprint.Client

	*gtk.Grid
}

func New(endpoint, key string) *UI {
	ui := &UI{
		Grid:    MustGrid(),
		Printer: octoprint.NewClient(endpoint, key),
	}

	ui.initialize()
	return ui
}

func (ui *UI) initialize() {
	ui.ShowDefaultPanel()
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
