package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)

type ActionFooter struct {
	*gtk.Box
	refreshButton				*gtk.Button
	backButton					*gtk.Button
}

func CreateActionFooter(
	buttonWidth					int,
	buttonHeight				int,
	refreshClicked				func(),
	backClicked					func(),
) *ActionFooter {
	base := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)

	instance := &ActionFooter {
		Box:				base,
	}

	instance.SetHAlign(gtk.ALIGN_END)
	instance.SetHExpand(true)
	instance.SetMarginTop(5)
	instance.SetMarginBottom(5)
	instance.SetMarginEnd(5)

	instance.refreshButton = instance.createRefreshButton(buttonWidth, buttonHeight, refreshClicked)
	instance.Add(instance.refreshButton)

	instance.backButton = instance.createBackButton(buttonWidth, buttonHeight, backClicked)
	instance.Add(instance.backButton)

	return instance
}

func (this *ActionFooter) createRefreshButton(buttonWidth int, buttonHeight int, clicked func()) *gtk.Button {
	image := utils.MustImageFromFileWithSize("refresh.svg", buttonWidth, buttonHeight)
	return utils.MustButton(image, clicked)
}

func (this *ActionFooter) createBackButton(buttonWidth int, buttonHeight int, clicked func()) *gtk.Button {
	image := utils.MustImageFromFileWithSize("back.svg", buttonWidth, buttonHeight)
	return utils.MustButton(image, clicked)
}
