package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type IncreaseZOffsetButton struct {
	*gtk.Button
	isIncrease					bool
}

func CreateIncreaseZOffsetButton(
	isIncrease					bool,
) *IncreaseZOffsetButton {
	var base *gtk.Button
	if isIncrease {
		base = utils.MustButtonImageStyle("Increase Offset", "z-offset-increase.svg", "", nil)
	} else {
		base = utils.MustButtonImageStyle("Decrease Offset", "z-offset-decrease.svg", "", nil)
	}

	instance := &IncreaseZOffsetButton{
		Button:						base,
		isIncrease:					isIncrease,
	}

	return instance
}

// NOTE: leave the heavy lifting to the parent panel.  If this button is use in multiple places,
// move the logic into here, and also add the dependency controls, and also any internal values
// that are needed.
