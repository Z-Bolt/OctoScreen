package gui

import (
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

// MustGrid returns a new gtk.Grid, if error panics.
func MustGrid() *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		panic(err)
	}

	return grid
}

// MustButtonImage returns a new gtk.Button with the given label, image and
// clicked callback. If error panics.
func MustButtonImage(label, img string, clicked func()) *gtk.Button {
	i, err := gtk.ImageNewFromFile(filepath.Join(ImagesFolder, img))
	if err != nil {
		panic(err)
	}

	b, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		panic(err)
	}

	b.SetImage(i)
	b.SetAlwaysShowImage(true)
	b.SetImagePosition(gtk.POS_TOP)
	b.SetVExpand(true)
	b.SetHExpand(true)

	if clicked != nil {
		b.Connect("clicked", clicked)
	}

	return b
}
