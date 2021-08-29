package ui

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var filamentToolPanelInstance *filamentManagerToolPanel

type filamentManagerToolPanel struct {
	CommonPanel
	tool *dataModels.FilamentManagerSelection

	listBox					*gtk.Box
	toolName				*gtk.Label
	selected				*gtk.Label
	material				*gtk.Label
	vendor					*gtk.Label
	used					*gtk.Label
	cost					*gtk.Label

	spoolCount int
}

func FilamentManagerToolPanel(
	ui *UI,
	tool *dataModels.FilamentManagerSelection,
) *filamentManagerToolPanel {
	if filamentToolPanelInstance == nil {
		instance := &filamentManagerToolPanel {
			CommonPanel:		NewCommonPanel("FilamentToolPanel", ui),
		}
		instance.initialize()
		filamentToolPanelInstance = instance
	} else {
		filamentToolPanelInstance.tool = tool
		// be sure to refresh the labels/spools when reopening the panel
		filamentToolPanelInstance.update()
	}

	return filamentToolPanelInstance
}

func (this *filamentManagerToolPanel) initialize() {
	panel := utils.MustGrid()
	panel.Attach(this.createToolbar(), 0, 0, 1, 1)

	rows := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	rows.Add(this.createSpoolListWindow())

	actionBar := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBar.SetHAlign(gtk.ALIGN_END)

	backImage := utils.MustImageFromFileWithSize("back.svg",
		this.Scaled(40), this.Scaled(40))
	backButton := utils.MustButton(backImage, func() {
		this.UI.GoToPreviousPanel()
	})
	backButton.SetMarginEnd(10)
	actionBar.Add(backButton)

	rows.Add(actionBar)
	panel.Attach(rows, 1, 0, 5, 1)

	this.Grid().Attach(panel, 1, 1, 1, 1)

	this.update()
}

func (this *filamentManagerToolPanel) update() {
	logger.TraceEnter("filamentManagerToolPanel.update()")

	utils.EmptyTheContainer(&this.listBox.Container)
	this.spoolCount = 0

	this.updateLabels()
	this.updateSpools()
	this.listBox.ShowAll()

	logger.TraceLeave("filamentManagerToolPanel.update()")
}

func (this *filamentManagerToolPanel) updateLabels() {
	// This could be optimized to just request the current tool
	request := &octoprintApis.FilamentManagerSelectionsRequest {}
	response, _ := request.Do(this.UI.Client)

	for _, selection := range response.Selections {
		if this.tool == nil || selection.Tool == this.tool.Tool {
			this.tool = selection
			break
		}
	}

	this.toolName.SetText(fmt.Sprintf("Tool %d", this.tool.Tool))
	this.selected.SetText(fmt.Sprintf("Spool: %s", this.tool.Spool.Name))
	this.material.SetText(fmt.Sprintf("Material: %s", this.tool.Spool.Profile.Material))
	this.vendor.SetText(fmt.Sprintf("Brand: %s", this.tool.Spool.Profile.Vendor))
	this.used.SetText(fmt.Sprintf("Used: %.0fg/%.0fg",
		this.tool.Spool.Used, this.tool.Spool.Weight))
	// TODO check how to get currency units
	this.cost.SetText(fmt.Sprintf("Cost: %.0f", this.tool.Spool.Cost))

	// in case the spool change, updated the color
	color := getSpoolColor(&this.tool.Spool)
	if color != "" {
		style, _ := this.selected.GetStyleContext()
		style.AddClass(color)
	}
}

func (this *filamentManagerToolPanel) updateSpools() {
	request := &octoprintApis.FilamentManagerSpoolsRequest {}
	response, _ := request.Do(this.UI.Client)

	for _, spool := range response.Spools {
		this.addSpool(this.listBox, spool)
	}
}

func (this *filamentManagerToolPanel) createSpoolListWindow() gtk.IWidget {
	this.listBox = utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	this.listBox.SetVExpand(true)

	box := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)

	spoolListWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	spoolListWindow.SetProperty("overlay-scrolling", false)
	spoolListWindow.Add(this.listBox)
	box.Add(spoolListWindow)

	return box
}

func (this *filamentManagerToolPanel) createToolbar() gtk.IWidget {
	toolbar := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	toolbar.SetHExpand(true)

	toolbar.Add(this.createInfoBar())

	return toolbar
}

func (this *filamentManagerToolPanel) createInfoBar() gtk.IWidget {
	infoBar := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	infoBar.SetVExpand(true)
	infoBar.SetVAlign(gtk.ALIGN_START)
	infoBar.SetHAlign(gtk.ALIGN_START)

	infoBar.SetMarginEnd(25)
	infoBar.SetMarginStart(25)

	t1 := utils.MustLabel("<big><b>Tool Information</b></big>")
	t1.SetHAlign(gtk.ALIGN_START)
	t1.SetMarginTop(25)

	infoBar.Add(t1)

	this.toolName = utils.MustLabel("")
	this.toolName.SetHAlign(gtk.ALIGN_START)
	infoBar.Add(this.toolName)

	extruder := utils.MustImageFromFileWithSize("extruder.svg",
		this.Scaled(50), this.Scaled(50))
	extruder.SetHAlign(gtk.ALIGN_CENTER)
	infoBar.Add(extruder)

	t2 := utils.MustLabel("<b>Filament Information</b>")
	t2.SetHAlign(gtk.ALIGN_START)
	t2.SetMarginTop(25)
	infoBar.Add(t2)

	this.selected = utils.MustLabel("Name")
	this.selected.SetHAlign(gtk.ALIGN_START)
	this.selected.SetMarginTop(25)
	infoBar.Add(this.selected)

	this.material = utils.MustLabel("Material")
	this.material.SetHAlign(gtk.ALIGN_START)
	infoBar.Add(this.material)

	this.vendor = utils.MustLabel("Vendor")
	this.vendor.SetHAlign(gtk.ALIGN_START)
	infoBar.Add(this.vendor)

	this.used = utils.MustLabel("Used Filament:")
	this.used.SetHAlign(gtk.ALIGN_START)
	infoBar.Add(this.used)

	this.cost = utils.MustLabel("Cost:")
	this.cost.SetHAlign(gtk.ALIGN_START)
	infoBar.Add(this.cost)


	return infoBar
}

func (this *filamentManagerToolPanel) addSpoolInfo(spool *dataModels.FilamentManagerSpool) *gtk.Box {
	infoBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	infoBox.SetVExpand(true)
	infoBox.SetHExpand(true)
	infoBox.SetHAlign(gtk.ALIGN_START)
	infoBoxStyleContext, _ := infoBox.GetStyleContext()
	infoBoxStyleContext.AddClass("labels-box")

	name := utils.MustLabel(spool.Name)
	name.SetHAlign(gtk.ALIGN_START)
	name.SetMarkup(fmt.Sprintf("<big>%s</big>", utils.TruncateString(spool.Name, 24)))
	infoBox.Add(name)

	material := utils.MustLabel("")
	material.SetHAlign(gtk.ALIGN_START)
	material.SetMarkup(fmt.Sprintf("<small>Material: %s</small>",
		spool.Profile.Material))
	infoBox.Add(material)

	usage := utils.MustLabel("")
	usage.SetHAlign(gtk.ALIGN_START)
	usage.SetMarkup(fmt.Sprintf("<small>Used: %.0f/%.0fg</small>",
		spool.Used, spool.Weight))
	infoBox.Add(usage)

	brand := utils.MustLabel("")
	brand.SetHAlign(gtk.ALIGN_START)
	brand.SetMarkup(fmt.Sprintf("<small>Brand: %s</small>", spool.Profile.Vendor))
	infoBox.Add(brand)

	return infoBox
}

func (this *filamentManagerToolPanel) addSpool(box *gtk.Box, spool *dataModels.FilamentManagerSpool) {
	spoolRow := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 0)
	spoolRow.SetVExpand(true)
	spoolRow.SetHExpand(true)

	clicked := func() {
		// TODO: verify we aren't currently printing
		request := &octoprintApis.FilamentManagerSetSelectionRequest {}
		request.Tool = this.tool.Tool
		request.Spool = spool.Id

		_, err := request.Do(this.UI.Client)
		if err != nil {
			logger.LogError("filamentManagerToolPanel.addSpool()", "Do(FilamentManagerSetSelectionRequest)", err)
			utils.ErrorMessageDialogBox(this.UI.window, 
				fmt.Sprintf("Unable to change filament: %v", err))
		} else {
			// possible optimization, use the response from FilamentManagerSetSelectionRequest
			// (a FilamentManagerSelection) to do the update to avoid the extra round
			// trip
			this.update()
		}
	}

	image := utils.MustImageFromFileWithSize("filament-spool.svg", this.Scaled(25), this.Scaled(25))

	color := getSpoolColor(spool)
	if color != "" {
		style, _ := image.GetStyleContext()
		style.AddClass(color)
	}
	rowButton := utils.MustButton(image, clicked)
	spoolRow.Add(rowButton)

	if spool.Id == this.tool.Spool.Id {
		extruder := utils.MustImageFromFileWithSize("extruder.svg", this.Scaled(25), this.Scaled(25))
		spoolRow.Add(extruder)
	}

	/*
	I tried adding the whole row was a button, but it broke all of my
	formatting
	*/

	spoolRow.Add(this.addSpoolInfo(spool))

	if this.spoolCount % 2 != 0 {
		styleContext, _ := spoolRow.GetStyleContext()
		styleContext.AddClass("list-item-nth-child-even")
	}
	this.spoolCount++

	box.Add(spoolRow)
	this.listBox.Add(box)
}


func getSpoolColor(spool *dataModels.FilamentManagerSpool) string {
	var color string = ""

	name := strings.ToLower(spool.Name)

	if strings.Contains(name, "black") {
		// skip setting for black since black on black won't be seen
		color = ""
	} else if strings.Contains(name, "white") {
		color = "color-white"
	} else if strings.Contains(name, "red") {
		color = "color-red"
	} else if strings.Contains(name, "green") {
		color = "color-green"
	} else if strings.Contains(name, "blue") {
		color = "color-blue"
	} else if strings.Contains(name, "brown") {
		color = "color-brown"
	} else if strings.Contains(name, "yellow") {
		color = "color-yellow"
	} else if strings.Contains(name, "gold") {
		color = "color-gold"
	} else if strings.Contains(name, "pink") {
		color = "color-pink"
	} else if strings.Contains(name, "grey") {
		color = "color-grey"
	} else if strings.Contains(name, "lime") {
		color = "color-lime"
	} else if strings.Contains(name, "orange") {
		color = "color-orange"
	}

	return color
}
