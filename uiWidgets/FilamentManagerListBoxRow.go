package uiWidgets

import (
	// "fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type FilamentManagerListBoxRow struct {
	*ListBoxRow

	SvgImageRadioButton			*SvgImageRadioButton
	FilamentManagerSpool		*dataModels.FilamentManagerSpool
	rowIndex					int
}

func CreateFilamentManagerListBoxRow(
	filamentManagerSpool		*dataModels.FilamentManagerSpool,
	rowIndex					int,
	spoolIsSelected				bool,
	clicked 					func(*SvgImageRadioButton),
) *FilamentManagerListBoxRow {
	const ROW_PADDING = 0
	base := CreateListBoxRow(rowIndex, ROW_PADDING)

	instance := &FilamentManagerListBoxRow {
		ListBoxRow:				base,
		SvgImageRadioButton:	nil,
		FilamentManagerSpool:	filamentManagerSpool,
		rowIndex:				rowIndex,
	}

	instance.SvgImageRadioButton = instance.createSvgImageRadioButton(spoolIsSelected, clicked)
	instance.Add(instance.SvgImageRadioButton)

	spoolInfoBox := instance.createSpoolInfoBox()
	instance.Add(spoolInfoBox)

	return instance
}

func (this *FilamentManagerListBoxRow) createSvgImageRadioButton(
	isSelected					bool,
	clicked 					func(*SvgImageRadioButton),
) *SvgImageRadioButton {
	spoolColor := this.getSpoolColor()
	filamentSpoolImage, _ := utils.CreateFilamentSpoolImage(spoolColor)
	filamentSpoolWithCheckMarkImage, _ := utils.CreateFilamentSpoolWithCheckMarkImage(spoolColor)

	buttonCssClassName := ""
	labelCssClassName := ""
	selectedLabelCssClassName := "white-foreground"
	if this.rowIndex % 2 == 0 {
		buttonCssClassName = "list-item-nth-child-odd-background-color"
		labelCssClassName = "list-item-nth-child-odd-foreground-color"
	} else {
		buttonCssClassName = "list-item-nth-child-even-background-color"
		labelCssClassName = "list-item-nth-child-even-foreground-color"
	}

	svgImageRadioButton := CreateSvgImageRadioButton(
		filamentSpoolImage,
		filamentSpoolWithCheckMarkImage,
		"Selected",

		buttonCssClassName,
		labelCssClassName,
		selectedLabelCssClassName,

		this.rowIndex,
		isSelected,
		clicked,
	)

	return svgImageRadioButton
}

func (this *FilamentManagerListBoxRow) getSpoolColor() string {
	name := strings.ToLower(this.FilamentManagerSpool.Name)
	if strings.Contains(name, "red") {
		return "red"
	}

	if strings.Contains(name, "green") {
		return "green"
	}

	if strings.Contains(name, "blue") {
		return "blue"
	}

	if strings.Contains(name, "cyan") {
		// return "cyan"
		return "#64cacc" // brightness: 80%
	}

	if strings.Contains(name, "magenta") {
		return "magenta"
	}

	if strings.Contains(name, "yellow") {
		// return "yellow"
		return "#cccc43" // brightness: 80%
	}

	if strings.Contains(name, "black") {
		return "black"
	}

	if strings.Contains(name, "white") {
		return "white"
	}

	if strings.Contains(name, "gray") || strings.Contains(name, "grey") {
		return "gray"
	}

	if strings.Contains(name, "orange") {
		return "orange"
	}

	if strings.Contains(name, "pink") {
		return "pink"
	}

	if strings.Contains(name, "purple") {
		return "purple"
	}

	if strings.Contains(name, "brown") {
		// return "brown"
		return "#874d19" // brightness: 53%
	}

	if strings.Contains(name, "gold") {
		// return "gold"
		return "#ccb13b" // brightness: 80%
	}

	if strings.Contains(name, "silver") {
		return "silver"
	}
	
	// Default to passing the background color of the row.
	if this.rowIndex % 2 != 0 {
		// .list-item-nth-child-even-background-color
		return "#34383C"
	} else {
		// .list-item-nth-child-odd-background-color
		return "#13181C"
	}
}

func (this *FilamentManagerListBoxRow) createSpoolInfoBox() *gtk.Box {
	spoolInfoBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	spoolInfoBox.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.SetMarginTop(5)
	spoolInfoBox.SetMarginBottom(5)

	nameLabel := utils.MustLabel("<big>%s</big>", utils.TruncateString(this.FilamentManagerSpool.Name, 25))
	nameLabel.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.Add(nameLabel)

	materialLabel := utils.MustLabel("<small>Material: %s</small>", this.FilamentManagerSpool.Profile.Material)
	materialLabel.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.Add(materialLabel)
	
	vendorLabel := utils.MustLabel("<small>Vendor: %s</small>", this.FilamentManagerSpool.Profile.Vendor)
	vendorLabel.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.Add(vendorLabel)
	
	weightLabel := utils.MustLabel("<small>Weight: %.0fg</small>", this.FilamentManagerSpool.Weight)
	weightLabel.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.Add(weightLabel)
	
	remainingLabel := utils.MustLabel("<small>Remaining: %.0fg</small>", this.FilamentManagerSpool.Weight - this.FilamentManagerSpool.Used)
	remainingLabel.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.Add(remainingLabel)

	var percentageUsedLabel *gtk.Label
	if this.FilamentManagerSpool.Weight != 0 {
		percentageUsed := (this.FilamentManagerSpool.Used / this.FilamentManagerSpool.Weight) * 100.0
		percentageUsedLabel = utils.MustLabel("<small>Used: %.0f%%</small>", percentageUsed)
	} else {
		percentageUsedLabel = utils.MustLabel("<small>Used: unknown</small>")
	}

	percentageUsedLabel.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.Add(percentageUsedLabel)

	return spoolInfoBox;
}

/*
TODO:


* when a spool is checked, call the API to update Filament Manager


* make sure this works (updating the check list) when there is only 1 extruder
	...that the dashes aren't present
	...that the "circle-1" doesn't appear
	...or, should the button not appear at all?
		...what do the other panels do?
			...do they display the extruder button when there is only 1 extruder?
* make sure this works when there are multiple extruders, but a single toolhead
* make sure that when there is only 1 extruder




* change the entire row to be a button?


* the back button isn't under the scroll bar
	maybe make the panel more like the files panel

* add a refresh button?
*/
