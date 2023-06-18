package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/utils"
)

type ListBoxRow struct {
	*gtk.ListBoxRow
	rowContentsBox		*gtk.Box
}

func CreateListBoxRow(
	index				int,
	padding				int,
) *ListBoxRow {
	base := createListBoxRow(index)

	rowContentsBox := createRowContentsBox(padding)
	base.Add(rowContentsBox)

	instance := &ListBoxRow {
		ListBoxRow:				base,
		rowContentsBox:			rowContentsBox,
	}

	return instance
}

// A gtk.ListBoxRow can only have one child.
// A uiWidgets.ListBoxRow will often have many children in it,
// so a contents box is added to the gtk.ListBoxRow, and when a
// child widget is added to uiWidgets.ListBoxRow, it will actually
// be added to rowContentsBox.
//
// 		Object hierarchy:
//			gtk.ListBoxRow (base)
//				gtk.Box (rowContentsBox)
													// TODO: remove this func (this *ListBoxRow) AddContent(widget gtk.IWidget) {
func (this *ListBoxRow) Add(widget gtk.IWidget) {
	this.rowContentsBox.Add(widget)
}



// Internal functions
func createListBoxRow(
	index			int,
) *gtk.ListBoxRow {
	listBoxRow, _ := gtk.ListBoxRowNew()
	listBoxRowStyleContext, _ := listBoxRow.GetStyleContext()
	listBoxRowStyleContext.AddClass("list-box-row")
	if index % 2 != 0 {
		listBoxRowStyleContext.AddClass("list-item-nth-child-even-background-color")
	} else {
		listBoxRowStyleContext.AddClass("list-item-nth-child-odd-background-color")
	}

	return listBoxRow
}

func createRowContentsBox(padding int) *gtk.Box {
	rowContentsBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 0)
	rowContentsBox.SetMarginTop(padding)
	rowContentsBox.SetMarginBottom(padding)
	rowContentsBox.SetMarginStart(padding)
	rowContentsBox.SetMarginEnd(padding)
	rowContentsBox.SetHAlign(gtk.ALIGN_START)
	// base.SetHExpand(true) // TODO: should this be set?

	return rowContentsBox
}
