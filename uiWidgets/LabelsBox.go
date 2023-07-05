package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)

/*
type LabelsBox struct {
	*gtk.Box
}

func CreateLabelsBox(
	nameLabel		*gtk.Label,
	infoLabel		*gtk.Label,
) *LabelsBox {
	base := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	instance := &LabelsBox {
		Box:		base,
	}

	if nameLabel != nil {
		instance.Add(nameLabel)
	}

	if infoLabel != nil {
		instance.Add(infoLabel)
	}

	instance.SetVExpand(false)
	instance.SetVAlign(gtk.ALIGN_CENTER)
	instance.SetHAlign(gtk.ALIGN_START)
	labelsBoxStyleContext, _ := instance.GetStyleContext()
	labelsBoxStyleContext.AddClass("labels-box")

	return instance
}
*/

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
