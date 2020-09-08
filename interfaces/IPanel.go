package interfaces

import (
	"github.com/gotk3/gotk3/gtk"
)

type IPanel interface {
	Grid() *gtk.Grid
	Show()
	Hide()
	ParentPanel() IPanel
}
