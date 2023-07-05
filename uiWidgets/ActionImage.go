package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)


func CreateActionImage(
	imageFileName		string,
	buttonWidth			int,
	buttonHeight		int,
	colorClass			string,
) *gtk.Image {
	image := utils.MustImageFromFileWithSize(
		imageFileName,
		buttonWidth,
		buttonHeight,
	)

	imageStyleContext, _ := image.GetStyleContext()
	imageStyleContext.AddClass(colorClass)

	return image
}
