package gui

import (
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

var ImagesFolder string

type GUI struct {
	Current *gtk.Grid

	*gtk.Grid
}

func New() *GUI {
	grid, _ := gtk.GridNew()

	gui := &GUI{Grid: grid}
	gui.ShowHomeMenu()
	return gui
}

func (g *GUI) ShowHomeMenu() {
	g.Add(NewHomeMenu(g).Grid)
}

func (g *GUI) Add(grid *gtk.Grid) {
	if g.Current != nil {
		g.Current.Destroy()
	}

	g.Current = grid
	g.Attach(g.Current, 1, 0, 1, 1)
	g.ShowAll()
}

func NewButtonImage(label, img string, clicked func()) gtk.IWidget {
	i, err := gtk.ImageNewFromFile(filepath.Join(ImagesFolder, img))
	if err != nil {
		panic(err)
	}

	b, _ := gtk.ButtonNewWithLabel(label)
	b.SetImage(i)
	b.SetAlwaysShowImage(true)
	b.SetImagePosition(gtk.POS_TOP)
	b.SetVExpand(true)
	b.SetHExpand(true)
	b.Connect("clicked", clicked)

	return b
}
