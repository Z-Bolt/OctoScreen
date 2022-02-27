package uiWidgets

import (
	"fmt"
	// "os"
	// "strings"

	"github.com/dustin/go-humanize"
	"github.com/gotk3/gotk3/gtk"
	"github.com/shirou/gopsutil/mem"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type SystemInformationInfoBox struct {
	*gtk.Box

	parentWindow					*gtk.Window
	logLevel						string
	uiScaleFactor					int
	memoryLabel						*gtk.Label
	// loadAverageLabel				*gtk.Label
	octoScreenResolutionLabel		*gtk.Label
	allocatedResolutionLabel		*gtk.Label
	currentResolutionLabel			*gtk.Label
	// uiScaleFactorLabel				*gtk.Label
}

func CreateSystemInformationInfoBox(
	parentWindow		*gtk.Window,
	uiScaleFactor		int,
) *SystemInformationInfoBox {
	base := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	base.SetVExpand(true)
	base.SetVAlign(gtk.ALIGN_CENTER)

	title := utils.MustLabel("<b>System Information</b>")
	title.SetMarginBottom(5)
	title.SetMarginTop(15)
	base.Add(title)

	instance := &SystemInformationInfoBox {
		Box:						base,
		parentWindow:				parentWindow,
		logLevel:					logger.LogLevel(),
		uiScaleFactor:				uiScaleFactor,
		memoryLabel:				createStyledLabel(),
		// loadAverageLabel:			createStyledLabel(),
		octoScreenResolutionLabel:	createStyledLabel(),
		allocatedResolutionLabel:	createStyledLabel(),
		currentResolutionLabel:		createStyledLabel(),
		// uiScaleFactorLabel:			createStyledLabel(),
	}

	instance.Add(instance.memoryLabel)
	// instance.Add(instance.loadAverageLabel)


	if instance.logLevel == "" {
		// If not set, default to warning level.
		instance.logLevel = "warn"
	}

	if instance.logLevel == "debug" {
		instance.Add(instance.octoScreenResolutionLabel)
		instance.Add(instance.allocatedResolutionLabel)
		instance.Add(instance.currentResolutionLabel)
		//instance.Add(instance.uiScaleFactorLabel)

		// Uncomment the following line to force the screen to expand.
		// instance.Add(utils.MustLabel("test"))
	} else {
		instance.Add(instance.currentResolutionLabel)
	}

	return instance
}

func createStyledLabel() *gtk.Label {
	label := utils.MustLabelWithCssClass("", "font-size-17")
	label.SetMarginBottom(2)

	return label
}

func (this *SystemInformationInfoBox) refreshMemoryLabel() {
	virtualMemoryStat, _ := mem.VirtualMemory()
	memoryString := fmt.Sprintf(
		"Memory: %s (free) / %s (total)",
		humanize.Bytes(virtualMemoryStat.Free),
		humanize.Bytes(virtualMemoryStat.Total),
	)

	this.memoryLabel.SetText(memoryString)
}

/*
func (this *SystemInformationInfoBox) refreshLoadAverageLabel() {
	avgStat, _ := load.Avg()
	loadAverageString := fmt.Sprintf(
		"Load average: %.2f, %.2f, %.2f",
		avgStat.Load1,
		avgStat.Load5,
		avgStat.Load15,
	)
	this.loadAverageLabel.SetText(loadAverageString)
}
*/

func (this *SystemInformationInfoBox) refreshOctoScreenResolutionLabel() {
	octoScreenConfig := utils.GetOctoScreenConfigInstance()
	resolution := octoScreenConfig.Resolution
	if resolution == "" {
		resolution = ">>MISSING<<"
	}
	octoScreenResolutionString := fmt.Sprintf(
		"OCTOSCREEN_RESOLUTION configuration setting: %s",
		resolution,
	)
	this.octoScreenResolutionLabel.SetText(octoScreenResolutionString)
}


func (this *SystemInformationInfoBox) refreshAllocatedResolutionLabel() {
	allocatedWidth := this.parentWindow.GetAllocatedWidth()
	allocatedHeight := this.parentWindow.GetAllocatedHeight()
	allocatedResolutionString := fmt.Sprintf(
		"Allocated window size: %dx%d",
		allocatedWidth,
		allocatedHeight,
	)
	this.allocatedResolutionLabel.SetText(allocatedResolutionString)
}


func (this *SystemInformationInfoBox) refreshCurrentResolutionLabel() {
	currentWidth, currentHeight := this.parentWindow.GetSize()
	// currentResolutionString := fmt.Sprintf(
	// 	"Current window size: %dx%d",
	// 	currentWidth,
	// 	currentHeight,
	// )

	currentResolutionString := fmt.Sprintf(
		"Current window size: %dx%d - UI scale factor: %d",
		currentWidth,
		currentHeight,
		this.uiScaleFactor,
	)

	this.currentResolutionLabel.SetText(currentResolutionString)
}

/*
func (this *SystemInformationInfoBox) refreshUiScaleFactorLabel() {
	uiScaleFactorString := fmt.Sprintf(
		"UI scale factor: %d",
		this.uiScaleFactor,
	)
	this.uiScaleFactorLabel.SetText(uiScaleFactorString)
}
*/

func (this *SystemInformationInfoBox) Refresh() {
	this.refreshMemoryLabel()
	//this.refreshLoadAverageLabel()

	if this.logLevel == "debug" {
		this.refreshOctoScreenResolutionLabel()
		this.refreshAllocatedResolutionLabel()
		this.refreshCurrentResolutionLabel()
		// this.refreshUiScaleFactorLabel()
	} else {
		this.refreshCurrentResolutionLabel()
	}
}
