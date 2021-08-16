package ui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

var filamentManagerPanelInstance *filamentManagerPanel

// This page acts as the dispatch to change the filament on a specific tool.
// Many users won't see this page at all

type filamentManagerPanel struct {
	CommonPanel

	listBox				*gtk.Box

	// this is mostly to make alternating colors easier
	toolCount int
}

func FilamentManagerPanel(
	ui				*UI,
) *filamentManagerPanel {
	if filamentManagerPanelInstance == nil {
		instance := &filamentManagerPanel {
			CommonPanel: NewCommonPanel("FilamentManagerPanel", ui),
		}
		instance.initialize()
		filamentManagerPanelInstance = instance
	}

	return filamentManagerPanelInstance
}

func (this *filamentManagerPanel) initialize() {
	this.listBox = utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	this.listBox.SetVExpand(true)

	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetProperty("overlay-scrolling", false)
	scroll.Add(this.listBox)

	box := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(scroll)
	box.Add(this.createActionFooter())
	this.Grid().Add(box)

	this.doLoad()
}

func (this *filamentManagerPanel) createActionFooter() *gtk.Box {
	actionBar := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBar.SetHAlign(gtk.ALIGN_END)
	actionBar.SetHExpand(true)
	actionBar.SetMarginTop(5)
	actionBar.SetMarginBottom(5)
	actionBar.SetMarginEnd(5)
	
	backImage := utils.MustImageFromFileWithSize("back.svg", this.Scaled(40), this.Scaled(40))
	backButton := utils.MustButton(backImage, func() {
		this.UI.GoToPreviousPanel()
	})

	actionBar.Add(backButton)

	return actionBar
}

func (this *filamentManagerPanel) loadSelections() *dataModels.FilamentManagerSelections {
	request := &octoprintApis.FilamentManagerSelectionsRequest {}
	response, _ := request.Do(this.UI.Client)

	return response
}

func (this *filamentManagerPanel) doLoad() {
	utils.EmptyTheContainer(&this.listBox.Container)
	this.toolCount = 0

	selections := this.loadSelections()

	if len(selections.Selections) == 0 {
		listBoxRow, _ := gtk.ListBoxRowNew()

		nameLabel := utils.MustLabel("Tool")
		nameLabel.SetMarkup("Filament Manager is not configured, or no spools are loaded")
		nameLabel.SetHExpand(true)
		nameLabel.SetHAlign(gtk.ALIGN_START)

		listBoxRow.Add(nameLabel)

		this.listBox.Add(listBoxRow)
	} else {
		for _, selection := range selections.Selections {
			this.addToolhead(selection)
		}
	}
}

func (this *filamentManagerPanel) addToolhead(selection *dataModels.FilamentManagerSelection) {
	listBoxRow, _ := gtk.ListBoxRowNew()

	image := utils.MustImageFromFileWithSize("extruder.svg", this.Scaled(45), this.Scaled(45))

	topBox := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	topBox.Add(image)

	labelsBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 5)
	labelsBox.SetVExpand(false)
	labelsBox.SetVAlign(gtk.ALIGN_CENTER)
	labelsBox.SetHAlign(gtk.ALIGN_START)
	labelsBoxStyleContext, _ := labelsBox.GetStyleContext()
	labelsBoxStyleContext.AddClass("labels-box")

	nameLabel := utils.MustLabel("Tool")
	nameLabel.SetMarkup(fmt.Sprintf("<big>Tool %d: %s</big>",
		selection.Tool, selection.Spool.Name))
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk.ALIGN_START)
	labelsBox.Add(nameLabel)

	spoolLabel := utils.MustLabel("")
	spoolLabel.SetHAlign(gtk.ALIGN_START)
	spoolLabel.SetMarkup(fmt.Sprintf("<small>Spool: %s by %s | Material: %s</small>",
		selection.Spool.Name, selection.Spool.Profile.Vendor, selection.Spool.Profile.Material))
	labelsBox.Add(spoolLabel)

	usageLabel := utils.MustLabel("")
	usageLabel.SetHAlign(gtk.ALIGN_START)
	usageLabel.SetMarkup(fmt.Sprintf("<small>Used: %.0f/%.0fg</small>",
		selection.Spool.Used, selection.Spool.Weight))
	labelsBox.Add(usageLabel)


	topBox.Add(labelsBox)

	button, _ := gtk.ButtonNew()
	button.Connect("clicked", func() {
		this.UI.GoToPanel(FilamentManagerToolPanel(this.UI, selection))
	})

	button.Add(topBox)

	listBoxRow.Add(button)
	this.listBox.Add(listBoxRow)
}