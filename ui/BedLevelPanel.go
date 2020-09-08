package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var bedLevelPanelInstance *bedLevelPanel

type bedLevelPanel struct {
	CommonPanel

	innerGrid			*gtk.Grid
	points				map[string][]float64
	homed				bool
}

func BedLevelPanel(
	ui				*UI,
	parentPanel		interfaces.IPanel,
) *bedLevelPanel {
	if bedLevelPanelInstance == nil {
		instance := &bedLevelPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		instance.initialize()
		bedLevelPanelInstance = instance
	}

	bedLevelPanelInstance.homed = false
	return bedLevelPanelInstance
}

func (this *bedLevelPanel) initialize() {
	defer this.Initialize()

	this.defineLevelingPoints()

	this.innerGrid = utils.MustGrid()
	this.innerGrid.SetRowHomogeneous(true)
	this.innerGrid.SetColumnHomogeneous(true)
	this.Grid().Attach(this.innerGrid, 1, 0, 2, 3)


	this.addBedLevelCornerButton(true, true)
	this.addBedLevelCornerButton(false, true)
	this.addBedLevelCornerButton(true, false)
	this.addBedLevelCornerButton(false, false)

	if this.UI.Settings != nil && this.UI.Settings.GCodes.AutoBedLevel != "" {
		autoLevelButton := this.createAutoLevelButton(this.UI.Settings.GCodes.AutoBedLevel)
		this.Grid().Attach(autoLevelButton, 3, 0, 1, 1)
	}
}

func (this *bedLevelPanel) addBedLevelCornerButton(isLeft, isTop bool) {
	x := 0
	y := 1
	placement := "b-"
	if isTop {
		placement = "t-"
		y = 0
	}

	if isLeft {
		placement += "l"
	} else {
		placement += "r"
		x = 1
	}

	button := this.createLevelButton(placement)
	if isLeft {
		button.SetHAlign(gtk.ALIGN_END)
	} else {
		button.SetHAlign(gtk.ALIGN_START)
	}

	if isTop {
		button.SetVAlign(gtk.ALIGN_END)
		button.SetImagePosition(gtk.POS_BOTTOM)
	} else {
		button.SetVAlign(gtk.ALIGN_START)
		button.SetImagePosition(gtk.POS_TOP)
	}

	styleContext, _ := button.GetStyleContext()
	styleContext.AddClass("no-margin")
	styleContext.AddClass("no-border")
	styleContext.AddClass("no-padding")

	if isLeft {
		styleContext.AddClass("padding-left-20")
	} else {
		styleContext.AddClass("padding-right-20")
	}

	this.innerGrid.Attach(button, x, y, 1, 1)
}

func (this *bedLevelPanel) defineLevelingPoints() {
	connectionRequest, err := (&octoprint.ConnectionRequest{}).Do(this.UI.Printer)
	if err != nil {
		utils.LogError("BedLevelPanel.defineLevelingPoints()", "Do(ConnectionRequest)", err)
		return
	}

	utils.Logger.Info(connectionRequest.Current.PrinterProfile)

	printerProfile, err := (&octoprint.PrinterProfilesRequest{Id: connectionRequest.Current.PrinterProfile}).Do(this.UI.Printer)
	if err != nil {
		utils.LogError("BedLevelPanel.defineLevelingPoints()", "Do(PrinterProfilesRequest)", err)
		return
	}

	xMax := printerProfile.Volume.Width
	yMax := printerProfile.Volume.Depth
	xOffset := xMax * 0.1
	yOffset := yMax * 0.1

	this.points = map[string][]float64 {
		"t-l": {xOffset, yMax - yOffset},
		"t-r": {xMax - xOffset, yMax - yOffset},
		"b-l": {xOffset, yOffset},
		"b-r": {xMax - xOffset, yOffset},
	}
}

func (this *bedLevelPanel) createLevelButton(placement string) *gtk.Button {
	imageFileName := fmt.Sprintf("bed-level-parts/bed-level-%s.svg", placement)
	noLabel := ""
	button := utils.MustButtonImage(noLabel, imageFileName, func() {
		this.goHomeIfRequired()

		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string {
			"G0 Z10 F2000",
			fmt.Sprintf("G0 X%f Y%f F10000", this.points[placement][0], this.points[placement][1]),
			"G0 Z0 F400",
		}

		if err := cmd.Do(this.UI.Printer); err != nil {
			utils.LogError("BedLevelPanel.createLevelButton()", "Do(CommandRequest)", err)
			return
		}
	})

	return button
}

func (this *bedLevelPanel) goHomeIfRequired() {
	if this.homed {
		return
	}

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{
		"G28",
	}

	if err := cmd.Do(this.UI.Printer); err != nil {
		utils.LogError("BedLevelPanel.goHomeIfRequire()", "Do(CommandRequest)", err)
		return
	}

	this.homed = true
}

func (this *bedLevelPanel) createAutoLevelButton(gcode string) *gtk.Button {
	button := utils.MustButtonImage("Auto Level", "bed-level.svg", func() {
		cmd := &octoprint.CommandRequest{}
		cmd.Commands = []string{
			gcode,
		}

		if err := cmd.Do(this.UI.Printer); err != nil {
			utils.LogError("BedLevelPanel.createAutoLevelButton()", "Do(CommandRequest)", err)
			return
		}
	})

	return button
}
