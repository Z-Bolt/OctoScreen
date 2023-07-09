package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/utils"
)

func  CreateListItemButton(
	index				int,
) *gtk.Button {
	listItemButton, _ := gtk.ButtonNew()
	listItemButtonStyleContext, _ := listItemButton.GetStyleContext()
	listItemButtonStyleContext.AddClass("list-item-button")
	if index % 2 != 0 {
		listItemButtonStyleContext.AddClass("list-item-nth-child-even")
	}

	return listItemButton
}
