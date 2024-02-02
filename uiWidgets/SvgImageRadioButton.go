package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"

	// "github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type SvgImageRadioButton struct {
	gtk.Button

	image						*gtk.Image
	selectedImage				*gtk.Image
	caption						string

	labelCssClassName			string
	selectedLabelCssClassName	string

	Index						int
	IsSelected					bool
	clicked						func(*SvgImageRadioButton)
}

func CreateSvgImageRadioButton(
	image						*gtk.Image,
	selectedImage				*gtk.Image,
	caption						string,

	buttonCssClassName			string,
	labelCssClassName			string,
	selectedLabelCssClassName	string,

	index						int,
	isSelected					bool,
	clicked						func(*SvgImageRadioButton),
) *SvgImageRadioButton {
	var base *gtk.Button
	if isSelected {
		base = utils.MustButtonImageUsingImage(caption, selectedImage, nil)
	} else {
		base = utils.MustButtonImageUsingImage(caption, image, nil)
	}

	styleContext, _ := base.GetStyleContext()
	styleContext.AddClass(buttonCssClassName)
	styleContext.AddClass("text-shadow-none")

	instance := &SvgImageRadioButton {
		Button:						*base,

		image:						image,
		selectedImage:				selectedImage,
		caption:					caption,
		labelCssClassName:			labelCssClassName,
		selectedLabelCssClassName:	selectedLabelCssClassName,
		Index:						index,
		IsSelected:					isSelected,
		clicked:					clicked,
	}
	instance.updateStyles()
	instance.Connect("clicked", instance.handleClick)

	return instance
}

func (this *SvgImageRadioButton) Select() {
	this.IsSelected = true
	this.updateStyles()
}

func (this *SvgImageRadioButton) Unselect() {
	this.IsSelected = false
	this.updateStyles()
}

func (this *SvgImageRadioButton) handleClick() {
	if this.IsSelected {
		// At this time FilamentManager's API does not support "unsetting" a spool,
		// so guard against/don't support unchecking a spool.
		//
		// At this time, only FilamentManagerPanel uses SvgImageRadioButton, but if
		// this UI element is used elsewhere, this will need to be addressed.
		return
	}

	this.IsSelected = !this.IsSelected;
	this.updateStyles()
	if this.clicked != nil {
		this.clicked(this)
	}
}

func (this *SvgImageRadioButton) updateStyles() {
	styleContext, _ := this.GetStyleContext()
	if this.IsSelected {
		this.SetImage(this.selectedImage);
		if len(this.labelCssClassName) > 0 {
			styleContext.RemoveClass(this.labelCssClassName)
		}

		if len(this.selectedLabelCssClassName) > 0 {
			styleContext.AddClass(this.selectedLabelCssClassName)
		}
	} else {
		this.SetImage(this.image);
		if len(this.selectedLabelCssClassName) > 0 {
			styleContext.RemoveClass(this.selectedLabelCssClassName)
		}

		if len(this.labelCssClassName) > 0 {
			styleContext.AddClass(this.labelCssClassName)
		}
	}
}
