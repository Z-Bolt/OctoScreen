package uiWidgets

import (
	// "fmt"

	"github.com/gotk3/gotk3/gtk"
)


func CreatePrintImage(
	buttonWidth			int,
	buttonHeight		int,
) *gtk.Image {
	return CreateActionImage("print.svg", buttonWidth, buttonHeight, "color-warning-sign-yellow")
}
