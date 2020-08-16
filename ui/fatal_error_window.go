package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

func CreateFatalErrorWindow(message string, description string) *gtk.Window {
	window, error := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if error != nil {
		Logger.Fatalln("Unable to create window: ", error)
	}

	window.SetTitle("Fatal Error")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	window.SetDefaultSize(800, 480)

	// Create a new label widget to show in the window.
	label, error := gtk.LabelNew("\n    " + message + "\n    " + description)
	if error != nil {
		Logger.Fatalln("Unable to create label: ", error)
	}

	label.SetHAlign(gtk.ALIGN_START)
	label.SetVAlign(gtk.ALIGN_START)
	window.Add(label)

	return window
}
