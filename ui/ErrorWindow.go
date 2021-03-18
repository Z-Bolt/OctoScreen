package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

func CreateErrorWindow(message string, description string) *gtk.Window {
	window, error := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if error != nil {
		logger.LogError("ErrorWindow.CreateErrorWindow()", "WindowNew()", error)
	}

	window.SetTitle("Error")
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	window.SetDefaultSize(800, 480)

	// Create a new label widget to show in the window.
	label, error := gtk.LabelNew("\n    " + message + "\n    " + description)
	if error != nil {
		logger.LogError("ErrorWindow.CreateErrorWindow()", "LabelNew()", error)
	}

	label.SetHAlign(gtk.ALIGN_START)
	label.SetVAlign(gtk.ALIGN_START)
	window.Add(label)

	return window
}
