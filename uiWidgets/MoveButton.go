package uiWidgets

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type MoveButton struct {
	*gtk.Button

	client					*octoprintApis.Client
	amountToMoveStepButton	*AmountToMoveStepButton
	axis					octoprintApis.Axis
	direction				float64
}

func CreateMoveButton(
	client					*octoprintApis.Client,
	amountToMoveStepButton	*AmountToMoveStepButton,
	label					string,
	image					string,
	axis					octoprintApis.Axis,
	direction				float64,
) *MoveButton {
	// A little bit of a "chicken or the egg" situation here.  Create the
	// instance first so there is a reference to the callback...
	instance := &MoveButton{
		Button:					nil,
		client:					client,
		amountToMoveStepButton:	amountToMoveStepButton,
		axis:					axis,
		direction:				direction,
	}
	base := MustPressedButton(label, image, instance.handlePressed, 200)
	// ... and then set the button
	instance.Button = base

	return instance
}

func (this *MoveButton) handlePressed() {
	distance := this.amountToMoveStepButton.Value() * this.direction
	cmd := &octoprintApis.PrintHeadJogRequest{}
	switch this.axis {
		case octoprintApis.XAxis:
			cmd.X = distance

		case octoprintApis.YAxis:
			cmd.Y = distance

		case octoprintApis.ZAxis:
			cmd.Z = distance
	}

	if err := cmd.Do(this.client); err != nil {
		utils.LogError("MoveButton.handlePressed()", "Do(PrintHeadJogRequest)", err)
		return
	}
}
