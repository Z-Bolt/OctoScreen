package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)


func CreateVerticalLayoutBox() *gtk.Box {
	verticalLayoutBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	verticalLayoutBox.SetMarginTop(0)
	verticalLayoutBox.SetMarginBottom(0)
	verticalLayoutBox.SetMarginStart(0)
	verticalLayoutBox.SetMarginEnd(0)

	return verticalLayoutBox
}
