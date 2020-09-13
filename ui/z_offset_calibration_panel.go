package ui

import (
	"fmt"
	// "math"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/interfaces"
	"github.com/Z-Bolt/OctoScreen/uiWidgets"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var zOffsetCalibrationPanelInstance *zOffsetCalibrationPanel

type pointCoordinates struct {
	x float64
	y float64
	z float64
}

type zOffsetCalibrationPanel struct {
	CommonPanel
	zCalibrationMode				bool
	activeTool						int
	cPoint							pointCoordinates
	zOffset							float64

	// First row
	selectToolStepButton			*uiWidgets.SelectToolStepButton
	decreaseZOffsetButton			*uiWidgets.IncreaseZOffsetButton
	increaseZOffsetButton			*uiWidgets.IncreaseZOffsetButton

	// Second row
	zOffsetLabel					*gtk.Label

	// Third row
	manualZCalibrationStepButton	*uiWidgets.ManualZCalibrationStepButton
}

func ZOffsetCalibrationPanel(
	ui 					*UI,
	parentPanel 		interfaces.IPanel,
) *zOffsetCalibrationPanel {
	if zOffsetCalibrationPanelInstance == nil {
		instane := &zOffsetCalibrationPanel {
			CommonPanel: NewCommonPanel(ui, parentPanel),
		}
		instane.cPoint = pointCoordinates {
			x: 20,
			y: 20,
			z: 0,
		}
		instane.initialize()

		zOffsetCalibrationPanelInstance = instane
	}

	return zOffsetCalibrationPanelInstance
}

func (this *zOffsetCalibrationPanel) initialize() {
	defer this.Initialize()

	// First row
	this.CreateSelectToolStepButton()
	this.CreateDecreaseZOffsetButton()
	this.CreateIncreaseZOffsetButton()


	// Second row
	this.zOffsetLabel = this.CreateZOffsetLabel()
	this.Grid().Attach(this.zOffsetLabel, 1, 1, 2, 1)


	// Third row
	// Start Manual
	// Z Calibration
	this.CreateManualZCalibrationStepButton()

	// Auto Z Calibration
	this.Grid().Attach(this.CreateAutoZCalibrationButton(), 1, 2, 2, 1)
}




// First row
func (this *zOffsetCalibrationPanel) CreateSelectToolStepButton() {
	this.selectToolStepButton = uiWidgets.CreateSelectToolStepButton(this.UI.Printer, false)
	_, err := this.selectToolStepButton.Connect("clicked", this.selectToolStepButtonHandleClick)
	if err != nil {
		panic(err)
	}

	toolheadCount := utils.GetToolheadCount(this.UI.Printer)
	if toolheadCount > 1 {
		// Only display the select tool button if there are multiple toolheads.
		this.Grid().Attach(this.selectToolStepButton, 0, 0, 1, 1)
	}
}

func (this *zOffsetCalibrationPanel) selectToolStepButtonHandleClick() {
	toolheadIndex := this.selectToolStepButton.Index()
	utils.Logger.Infof("Changing tool to tool%d", toolheadIndex)

	gcode := fmt.Sprintf("T%d", toolheadIndex)

	if this.zCalibrationMode {
		this.activeTool = toolheadIndex
		this.command(fmt.Sprintf("G0 Z%f", 5.0))
		this.command(gcode)
		time.Sleep(time.Second * 1)
		this.command(fmt.Sprintf("G0 X%f Y%f F10000", this.cPoint.x, this.cPoint.y))

		cmd := &octoprint.GetZOffsetRequest{Tool: this.activeTool}
		response, err := cmd.Do(this.UI.Printer)
		if err != nil {
			utils.LogError("z_offset_calibration.setToolheadButtonClickHandler()", "Do(GetZOffsetRequest)", err)
			return
		}

		this.updateZOffset(response.Offset)
	} else {
		this.command(gcode)
	}
}

func (this *zOffsetCalibrationPanel) CreateDecreaseZOffsetButton() {
	this.decreaseZOffsetButton = uiWidgets.CreateIncreaseZOffsetButton(false)
	_, err := this.decreaseZOffsetButton.Connect("clicked", this.decreaseZOffsetButtonClicked)
	if err != nil {
		panic(err)
	}
	this.Grid().Attach(this.decreaseZOffsetButton, 1, 0, 1, 1)
}

func (this *zOffsetCalibrationPanel) decreaseZOffsetButtonClicked() {
	if !this.zCalibrationMode {
		return
	}

	this.updateZOffset(this.zOffset - 0.02)
}

func (this *zOffsetCalibrationPanel) CreateIncreaseZOffsetButton() {
	this.increaseZOffsetButton = uiWidgets.CreateIncreaseZOffsetButton(true)
	_, err := this.increaseZOffsetButton.Connect("clicked", this.increaseZOffsetButtonClicked)
	if err != nil {
		panic(err)
	}
	this.Grid().Attach(this.increaseZOffsetButton, 2, 0, 1, 1)
}

func (this *zOffsetCalibrationPanel) increaseZOffsetButtonClicked() {
	if !this.zCalibrationMode {
		return
	}

	this.updateZOffset(this.zOffset + 0.02)
}


// Second row
func (this *zOffsetCalibrationPanel) CreateZOffsetLabel() *gtk.Label {
	label := utils.MustLabel("")
	label.SetVAlign(gtk.ALIGN_CENTER)
	label.SetHAlign(gtk.ALIGN_CENTER)
	label.SetVExpand(true)
	label.SetHExpand(true)
	label.SetLineWrap(true)

	return label
}



// Third row
func (this *zOffsetCalibrationPanel) CreateManualZCalibrationStepButton() {
	this.manualZCalibrationStepButton = uiWidgets.CreateManualZCalibrationStepButton()
	_, err := this.manualZCalibrationStepButton.Connect("clicked", this.manualZCalibrationStepButtonHandleClick)
	if err != nil {
		panic(err)
	}

	this.Grid().Attach(this.manualZCalibrationStepButton, 0, 2, 1, 1)
}

func (this *zOffsetCalibrationPanel) manualZCalibrationStepButtonHandleClick() {
	if this.manualZCalibrationStepButton.IsCalibrating() {
		// BUG: This does not work.  At least not on a Prusa i3.  Need to get this working with all printers.
		// NOTE: Running this also causes the machine to reboot.

		this.command("G28")				// G28 Return to Machine Zero Point
		this.command("T0")				// T0 Switch to first toolhead
		time.Sleep(time.Second * 1)
		this.command(fmt.Sprintf("G0 X%f Y%f F10000", this.cPoint.x, this.cPoint.y))
		this.command(fmt.Sprintf("G0 Z10 F2000"))
		this.command(fmt.Sprintf("G0 Z%f F400", this.cPoint.z))

		this.activeTool = 0
		this.updateZOffset(0)
	} else {
		this.zOffsetLabel.SetText("Press \"Z Offset\"\nbutton to start\nZ-Offset calibration")
	}
}

func (this *zOffsetCalibrationPanel) CreateAutoZCalibrationButton() gtk.IWidget {
	return utils.MustButtonImageStyle("Auto Z Calibration", "z-calibration.svg", "", func() {
		if this.zCalibrationMode {
			return
		}

		// BUG: This does not work.  At least not on a Prusa i3.  Need to get this working with all printers.
		// when RunZOffsetCalibrationRequest is called, it's returning a 404.
		cmd := &octoprint.RunZOffsetCalibrationRequest{}
		if err := cmd.Do(this.UI.Printer); err != nil {
			utils.LogError("z_offset_calibration.createAutoZCalibrationButton()", "Do(RunZOffsetCalibrationRequest)", err)
		}
	})
}

func (this *zOffsetCalibrationPanel) updateZOffset(value float64) {
	// BUG: This does not work.  At least not on a Prusa i3.  Need to get this working with all printers.

	this.zOffset = utils.ToFixed(value, 4)

	this.zOffsetLabel.SetText(fmt.Sprintf("Z-Offset: %.2f", this.zOffset))

	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string {
		fmt.Sprintf("SET_GCODE_OFFSET Z=%f", this.zOffset),
		"G0 Z0 F100",
	}
	if err := cmd.Do(this.UI.Printer); err != nil {
		utils.LogError("z_offset_calibration.updateZOffset()", "Do(CommandRequest)", err)
	}

	cmd2 := &octoprint.SetZOffsetRequest {
		Value: this.zOffset,
		Tool: this.activeTool,
	}
	if err := cmd2.Do(this.UI.Printer); err != nil {
		utils.LogError("z_offset_calibration.updateZOffset()", "Do(SetZOffsetRequest)", err)
	}
}


func (this *zOffsetCalibrationPanel) command(gcode string) error {
	cmd := &octoprint.CommandRequest{}
	cmd.Commands = []string{gcode}

	return cmd.Do(this.UI.Printer)
}
