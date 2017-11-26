package gui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/OctoPrint-TFT/octoprint"
)

var ImagesFolder string

type GUI struct {
	Current *gtk.Grid
	Printer *octoprint.Printer

	*gtk.Grid
}

func New(endpoint, key string) *GUI {
	gui := &GUI{
		Grid:    MustGrid(),
		Printer: octoprint.NewPrinter(endpoint, key),
	}

	gui.initialize()
	return gui
}

func (g *GUI) initialize() {
	g.ShowMenu()
}

func (g *GUI) ShowMenu() {
	g.Add(NewMenu(g).Grid)
}

func (g *GUI) Add(grid *gtk.Grid) {
	if g.Current != nil {
		g.Current.Destroy()
	}

	g.Current = grid
	g.Attach(g.Current, 1, 0, 1, 1)
	g.ShowAll()
}
