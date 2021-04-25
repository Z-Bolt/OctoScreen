package interfaces

import (
	"github.com/gotk3/gotk3/gtk"
)

type IPanel interface {
	Name() string
	Grid() *gtk.Grid
	PreShow()
	Show()
	Hide()
	// maybe add PostShow(), PreHide(), and PostHide()
}
