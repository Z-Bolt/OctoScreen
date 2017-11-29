package ui

import (
	"fmt"
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

// MustBox returns a new gtk.Box, with the given configuration, if err panics.
func MustBox(o gtk.Orientation, spacing int) *gtk.Box {
	box, err := gtk.BoxNew(o, spacing)
	if err != nil {
		panic(err)
	}

	return box
}

// MustLabel returns a new gtk.Label, if err panics.
func MustLabel(label string, args ...interface{}) *gtk.Label {
	l, err := gtk.LabelNew(fmt.Sprintf(label, args...))
	if err != nil {
		panic(err)
	}

	//l.SetVExpand(true)
	return l
}

// MustButtonImage returns a new gtk.Button with the given label, image and
// clicked callback. If error panics.
func MustButtonImage(label, img string, clicked func()) *gtk.Button {
	b, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		panic(err)
	}

	b.SetImage(MustImageFromFile(img))
	b.SetAlwaysShowImage(true)
	b.SetImagePosition(gtk.POS_TOP)
	b.SetVExpand(true)
	b.SetHExpand(true)

	if clicked != nil {
		b.Connect("clicked", clicked)
	}

	return b
}

func MustImageFromFile(img string) *gtk.Image {
	i, err := gtk.ImageNewFromFile(filepath.Join(ImagesFolder, img))
	if err != nil {
		panic(err)
	}

	return i
}
