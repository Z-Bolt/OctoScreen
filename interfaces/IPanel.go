package interfaces

import (
	"github.com/gotk3/gotk3/gtk"
)


type IPanel interface {
	// Initialize()
	// AddButton()
	Name() string
	Grid() *gtk.Grid
	PreShow()
	Show()
	Hide()
	// maybe add PostShow(), PreHide(), and PostHide() ?
	// Scaled(s int) int
	// arrangeMenuItems(
	// 	grid			*gtk.Grid,
	// 	items			[]dataModels.MenuItem,
	// 	cols			int,
	// )
	// command(gcode string) error
}
