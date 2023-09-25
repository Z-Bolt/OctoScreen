package uiWidgets

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)


func CreateNameLabel(name string) *gtk.Label {
	label := utils.MustLabel(name)
	truncatedName := utils.TruncateString(name, 28)
	markup := fmt.Sprintf("<big>%s</big>", truncatedName)
	label.SetMarkup(markup)
	label.SetHExpand(true)
	label.SetHAlign(gtk.ALIGN_START)

	return label
}
