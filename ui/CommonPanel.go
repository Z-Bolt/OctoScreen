package ui

import (
	"fmt"
	"math"
	// "strings"
	// "sync"
	// "time"

	// "github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type CommonPanel struct {
	UI					*UI
	grid				*gtk.Grid
	preShowCallback		func()
	backgroundTask		*utils.BackgroundTask
	// parentPanel			interfaces.IPanel
	panelWidth			int
	panelHeight			int
	includeBackButton	bool
	backButton			*gtk.Button
	buttons				[]gtk.IWidget
}

func NewCommonPanel(ui *UI, parentPanel interfaces.IPanel) CommonPanel {
	return newPanel(ui, parentPanel, true)
}

func NewTopLevelCommonPanel(ui *UI, parentPanel interfaces.IPanel) CommonPanel {
	return newPanel(ui, parentPanel, false)
}

func newPanel(ui *UI, parentPanel interfaces.IPanel, includeBackButton bool) CommonPanel {
	grid := utils.MustGrid()
	grid.SetRowHomogeneous(true)
	grid.SetColumnHomogeneous(true)

	return CommonPanel {
		UI:					ui,
		grid:				grid,
		// parentPanel:		parentPanel,
		panelWidth:			4,
		panelHeight:		3,
		includeBackButton:	includeBackButton,
	}
}

func (this *CommonPanel) Initialize() {
	last := this.panelWidth * this.panelHeight
	if last < len(this.buttons) {
		cols := math.Ceil(float64(len(this.buttons)) / float64(this.panelWidth))
		last = int(cols) * this.panelWidth
	}

	for i := len(this.buttons) + 1; i < last; i++ {
		box := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 0)
		this.AddButton(box)
	}

	this.backButton = utils.MustButtonImageStyle("Back", "back.svg", "color-none", this.UI.GoToPreviousPanel)
	if this.includeBackButton {
		this.AddButton(this.backButton)
	}
}

func (this *CommonPanel) AddButton(button gtk.IWidget) {
	x := len(this.buttons) % this.panelWidth
	y := len(this.buttons) / this.panelWidth
	this.grid.Attach(button, x, y, 1, 1)
	this.buttons = append(this.buttons, button)
}

func (this *CommonPanel) PreShow() {
	if this.preShowCallback != nil {
		this.preShowCallback()
	}
}

func (this *CommonPanel) Show() {
	if this.backgroundTask != nil {
		this.backgroundTask.Start()
	}
}

func (this *CommonPanel) Hide() {
	if this.backgroundTask != nil {
		this.backgroundTask.Close()
	}
}

func (this *CommonPanel) Grid() *gtk.Grid {
	return this.grid
}

func (this *CommonPanel) Scaled(s int) int {
	return s * this.UI.scaleFactor
}

func (this *CommonPanel) arrangeMenuItems(
	grid			*gtk.Grid,
	items			[]dataModels.MenuItem,
	cols			int,
) {
	for i, item := range items {
		panel := getPanel(this.UI, this, item)
		if panel != nil {
			color := fmt.Sprintf("color%d", (i % 4) + 1)
			icon := fmt.Sprintf("%s.svg", item.Icon)
			button := utils.MustButtonImageStyle(item.Name, icon, color, func() {
				this.UI.GoToPanel(panel)
			})
			grid.Attach(button, (i % cols), i / cols, 1, 1)
		}
	}
}

func (this *CommonPanel) command(gcode string) error {
	cmd := &octoprintApis.CommandRequest{}
	cmd.Commands = []string{gcode}
	return cmd.Do(this.UI.Client)
}
