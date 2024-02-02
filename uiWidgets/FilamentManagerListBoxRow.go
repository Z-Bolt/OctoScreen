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
	ClickableListBoxRow

	ScreenWidth							int
	ContentsBox							*gtk.Box
	FilamentSpoolImage					*gtk.Image
	FilamentSpoolWithCheckMarkImage		*gtk.Image
	SvgImageRadioButton					*SvgImageRadioButton
	FilamentManagerSpool				*dataModels.FilamentManagerSpool
}

func CreateFilamentManagerListBoxRow(
	screenWidth				int,
	extruderCount			int,
	filamentManagerSpool	*dataModels.FilamentManagerSpool,
	spoolIsSelected			bool,
	rowIndex				int,
	rowClickHandler			func (button *gtk.Button, index int),
) *FilamentManagerListBoxRow {
	const ROW_PADDING = 0
	base := CreateClickableListBoxRow(rowIndex, ROW_PADDING, rowClickHandler)

	instance := &FilamentManagerListBoxRow {
		ClickableListBoxRow:				*base,
		ScreenWidth:						screenWidth,
		ContentsBox:						nil,
		FilamentSpoolImage:					nil,
		FilamentSpoolWithCheckMarkImage:	nil,
		SvgImageRadioButton:				nil,
		FilamentManagerSpool:				filamentManagerSpool,
	}

	instance.ContentsBox = createContentsBox(ROW_PADDING)

	instance.SvgImageRadioButton = instance.createSvgImageRadioButton(spoolIsSelected, nil)
	instance.ContentsBox.Add(instance.SvgImageRadioButton)

	spoolInfoBox := instance.createSpoolInfoBox(extruderCount)
	instance.ContentsBox.Add(spoolInfoBox)

	instance.Add(instance.ContentsBox)

	return instance
}

func (this *FilamentManagerListBoxRow) createSvgImageRadioButton(
	isSelected					bool,
	clicked						func (*SvgImageRadioButton),
) *SvgImageRadioButton {
	spoolColor := this.getSpoolColor()
	filamentSpoolImage, _ := utils.CreateFilamentSpoolImage(spoolColor)
	filamentSpoolWithCheckMarkImage, _ := utils.CreateFilamentSpoolWithCheckMarkImage(spoolColor)

	buttonCssClassName := ""
	labelCssClassName := ""
	selectedLabelCssClassName := "white-foreground"
	if this.RowIndex % 2 == 0 {
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

		this.RowIndex,
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
	if this.RowIndex % 2 != 0 {
		// .list-item-nth-child-even-background-color
		return "#34383C"
	} else {
		// .list-item-nth-child-odd-background-color
		return "#13181C"
	}
}

func (this *FilamentManagerListBoxRow) createSpoolImage() *gtk.Image {
	spoolColor := this.getSpoolColor()
	filamentSpoolImage, _ := utils.CreateFilamentSpoolImage(spoolColor)

	return filamentSpoolImage
}

func (this *FilamentManagerListBoxRow) createSpoolInfoBox(extruderCount int) *gtk.Box {
	spoolInfoBox := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	spoolInfoBox.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.SetMarginStart(10)
	spoolInfoBox.SetMarginTop(5)
	spoolInfoBox.SetMarginBottom(5)

	// Warning: this has to potential to expand the screen if the name is too long.
	// If bug reports come in, the name will need to be truncated more.
	maxNameLength := 0
	if extruderCount >= 2 {
		// If there are multiple extruders, the extruder step button will be visible, and the
		// width available to display the name will be less.

		if (this.ScreenWidth >= 760) {
			maxNameLength = 30
		} else if (this.ScreenWidth >= 600) {
			maxNameLength = 22
		} else {
			maxNameLength = 15
		}
	} else {
		if (this.ScreenWidth >= 760) {
			maxNameLength = 45
		} else if (this.ScreenWidth >= 600) {
			maxNameLength = 32
		} else {
			maxNameLength = 20
		}
	}

	truncatedName := utils.TruncateString(this.FilamentManagerSpool.Name, maxNameLength)
	nameLabel := utils.MustLabel("<big>%s</big>", truncatedName)
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

	costLabel := utils.MustLabel("<small>Cost: %.0f</small>", this.FilamentManagerSpool.Cost)
	costLabel.SetHAlign(gtk.ALIGN_START)
	spoolInfoBox.Add(costLabel)

	return spoolInfoBox;
}
