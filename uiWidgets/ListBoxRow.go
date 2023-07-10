package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)

type ListBoxRow struct {
	gtk.ListBoxRow

	ContentsBox			*gtk.Box
	RowIndex			int
}

func CreateListBoxRow(
	rowIndex			int,
	padding				int,
) *ListBoxRow {
	/*
		Object hierarchy (for a static, non-clickable list item):
		ScrollableListBox (ScrolledWindow + ListBox)
			ListBoxRow
				ContentsBox (to layout the objects)

		Object hierarchy (for a clickable list item/button):
		ScrollableListBox (ScrolledWindow + ListBox)
			ClickableListBoxRow
				listItemButton (to handle to click for the entire item amd all of the child controls)
	*/

	base := createListBoxRow(rowIndex)
	// ctx1, _ := base.GetStyleContext()
	// ctx1.AddClass("blue-background")

	contentsBox := createContentsBox(padding) // jab rename
	// ctx2, _ := contentsBox.GetStyleContext()
	// ctx2.AddClass("cyan-background-JAB-NO-USE")

	base.Add(contentsBox)


	instance := &ListBoxRow {
		ListBoxRow:				*base,
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
func (this *ListBoxRow) Add(widget gtk.IWidget) {
	this.ContentsBox.Add(widget)
}


// Internal functions
func createListBoxRow(
	rowIndex		int,
) *gtk.ListBoxRow {
	listBoxRow, _ := gtk.ListBoxRowNew()
	listBoxRowStyleContext, _ := listBoxRow.GetStyleContext()
	listBoxRowStyleContext.AddClass("list-box-row")
	if rowIndex % 2 != 0 {
		listBoxRowStyleContext.AddClass("list-item-nth-child-even-background-color")
	} else {
		listBoxRowStyleContext.AddClass("list-item-nth-child-odd-background-color")
	}

	return listBoxRow
}

func createContentsBox(padding int) *gtk.Box {
	contentsBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 0)
	contentsBox.SetMarginTop(padding)
	contentsBox.SetMarginBottom(padding)
	contentsBox.SetMarginStart(padding)
	contentsBox.SetMarginEnd(padding)
	contentsBox.SetHAlign(gtk.ALIGN_START)
	// contentsBox.SetHExpand(true) // TODO: should this be set?

	return contentsBox
}
