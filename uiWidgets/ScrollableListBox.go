package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)

type ScrollableListBox struct {
	*gtk.ScrolledWindow
	ListBox					*gtk.Box
}

func CreateScrollableListBox() *ScrollableListBox {
	base, _ := gtk.ScrolledWindowNew(nil, nil)

	instance := &ScrollableListBox {
		ScrolledWindow:		base,
	}

	instance.SetProperty("overlay-scrolling", false)
	// ctx1, _ := instance.GetStyleContext()
	// ctx1.AddClass("red-background")

	instance.ListBox = utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	instance.ListBox.SetVExpand(true)
	// ctx2, _ := instance.ListBox.GetStyleContext()
	// ctx2.AddClass("green-background")

	instance.ScrolledWindow.Add(instance.ListBox)

	return instance
}

func (this *ScrollableListBox) Add(widget gtk.IWidget) {
	this.ListBox.Add(widget)
}

func (this *ScrollableListBox) ListBoxContainer() *gtk.Container {
	return &this.ListBox.Container
}

func (this *ScrollableListBox) ShowAll() {
	// Might need to also call this.ScrolledWindow.ShowAll().
	// Add it in later if needed.
	this.ListBox.ShowAll()
}
