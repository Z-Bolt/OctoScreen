package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/utils"
)

type ClickableListBoxRow struct {
	gtk.ListBoxRow

	ListItemButton		*gtk.Button
	ContentsBox			*gtk.Box
	RowIndex			int
}

func CreateClickableListBoxRow(
	rowIndex			int,
	padding				int,
	rowClickHandler		func (button *gtk.Button, rowIndex int),
) *ClickableListBoxRow {
	/*
		Object hierarchy:
		(for a static, non-clickable list item)

		ScrollableListBox (ScrolledWindow + ListBox)
			ListBoxRow
				ContentsBox (to layout the objects)


		Object hierarchy:
		(for a clickable list item/button)

		ScrollableListBox (ScrolledWindow + ListBox)
			ClickableListBoxRow
				listItemButton (to handle to click for the entire item amd all of the child controls)
					ContentsBox (to layout the objects)
	*/

	base := createListBoxRow(rowIndex)
	// ctx1, _ := base.GetStyleContext()
	// ctx1.AddClass("blue-background")

	listItemButton := createListItemButton(rowIndex, rowClickHandler)
	base.Add(listItemButton)

	contentsBox := createContentsBox(padding)
	// ctx3, _ := contentsBox.GetStyleContext()
	// ctx3.AddClass("magenta-background")
	listItemButton.Add(contentsBox)


	instance := &ClickableListBoxRow {
		ListBoxRow:				*base,
		ListItemButton:			listItemButton,
		ContentsBox:			contentsBox,
		RowIndex:				rowIndex,
	}

	return instance
}

// A gtk.ListBoxRow can only have one child.
// A uiWidgets.ListBoxRow will often have many children in it,
// so a contents box is added to the gtk.ListBoxRow, and when a
// child widget is added to uiWidgets.ListBoxRow, it will actually
// be added to contentsBox.
//
// 	Object hierarchy:
//		gtk.ListBoxRow (base)
//			gtk.Box (ContentsBox)
func (this *ClickableListBoxRow) Add(widget gtk.IWidget) {
	this.ContentsBox.Add(widget)
}

func createListItemButton(
	rowIndex			int,
	rowClickHandler		func (button *gtk.Button, rowIndex int),
) *gtk.Button{
	listItemButton, _ := gtk.ButtonNew()
	listItemButton.Connect("clicked", rowClickHandler, rowIndex)

	listItemButtonStyleContext, _ := listItemButton.GetStyleContext()
	listItemButtonStyleContext.AddClass("list-item-button")

	// buttonStyleContext.AddClass("cyan-background")

	if rowIndex % 2 == 0 {
		listItemButtonStyleContext.AddClass("list-item-nth-child-odd-background-color")
	} else {
		listItemButtonStyleContext.AddClass("list-item-nth-child-even-background-color")
	}

	return listItemButton
}
