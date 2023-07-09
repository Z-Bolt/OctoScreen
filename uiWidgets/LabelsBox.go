package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)

func CreateLabelsBox(
	nameLabel		*gtk.Label,
	infoLabel		*gtk.Label,
) *gtk.Box {
	box := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)

	if nameLabel != nil {
		box.Add(nameLabel)
	}

	if infoLabel != nil {
		box.Add(infoLabel)
	}

	box.SetVExpand(false)
	box.SetVAlign(gtk.ALIGN_CENTER)
	box.SetHAlign(gtk.ALIGN_START)
	labelsBoxStyleContext, _ := box.GetStyleContext()
	labelsBoxStyleContext.AddClass("labels-box")

	return box
}
