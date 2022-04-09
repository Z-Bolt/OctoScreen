package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type bedLevelPanel struct {
	CommonPanel

	innerGrid			*gtk.Grid
	points				map[string][]float64
	homed				bool
}

var bedLevelPanelInstance *bedLevelPanel

func GetBedLevelPanelInstance(
	ui				*UI,
) *bedLevelPanel {
	if bedLevelPanelInstance == nil {
		instance := &bedLevelPanel {
			CommonPanel: CreateCommonPanel("BedLevelPanel", ui),
		}
		instance.initialize()
		bedLevelPanelInstance = instance
		bedLevelPanelInstance.homed = false
	}

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
	logger.TraceEnter("BedLevelPanel.defineLevelingPoints()")

	connectionResponse, err := (&octoprintApis.ConnectionRequest{}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("BedLevelPanel.defineLevelingPoints()", "version.Get()", err)
		logger.TraceLeave("BedLevelPanel.defineLevelingPoints()")
		return
	}

	logger.Info(connectionResponse.Current.PrinterProfile)

	printerProfile, err := (&octoprintApis.PrinterProfilesRequest{Id: connectionResponse.Current.PrinterProfile}).Do(this.UI.Client)
	if err != nil {
		logger.LogError("BedLevelPanel.defineLevelingPoints()", "Do(PrinterProfilesRequest)", err)
		logger.TraceLeave("BedLevelPanel.defineLevelingPoints()")
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

	logger.TraceLeave("BedLevelPanel.defineLevelingPoints()")
}

func (this *bedLevelPanel) createLevelButton(placement string) *gtk.Button {
	imageFileName := fmt.Sprintf("bed-level-parts/bed-corner-%s.svg", placement)
	noLabel := ""
	button := utils.MustButtonImage(noLabel, imageFileName, func() {
		this.goHomeIfRequired()

		cmd := &octoprintApis.CommandRequest{}
		cmd.Commands = []string {
			"G0 Z10 F2000",
			fmt.Sprintf("G0 X%f Y%f F10000", this.points[placement][0], this.points[placement][1]),
			"G0 Z0 F400",
		}

		if err := cmd.Do(this.UI.Client); err != nil {
			logger.LogError("BedLevelPanel.createLevelButton()", "Do(CommandRequest)", err)
			return
		}
	})

	return button
}

func (this *bedLevelPanel) goHomeIfRequired() {
	if this.homed {
		return
	}

	cmd := &octoprintApis.CommandRequest{}
	cmd.Commands = []string{
		"G28",
	}

	if err := cmd.Do(this.UI.Client); err != nil {
		logger.LogError("BedLevelPanel.goHomeIfRequire()", "Do(CommandRequest)", err)
		return
	}

	this.homed = true
}

func (this *bedLevelPanel) createAutoLevelButton(gcode string) *gtk.Button {
	button := utils.MustButtonImage("Auto Level", "bed-level.svg", func() {
		cmd := &octoprintApis.CommandRequest{}
		cmd.Commands = []string{
			gcode,
		}

		if err := cmd.Do(this.UI.Client); err != nil {
			logger.LogError("BedLevelPanel.createAutoLevelButton()", "Do(CommandRequest)", err)
			return
		}
	})

	return button
}
