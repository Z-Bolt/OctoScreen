package uiWidgets

type AmountToMoveStepButton struct {
	*StepButton
}

func CreateAmountToMoveStepButton(
	colorVariation		int,
	clicked				func(),
) *AmountToMoveStepButton {
	base := CreateStepButton(
		colorVariation,
		clicked,
		Step{"10mm",   "move-step.svg", nil, 10.00},
		Step{"20mm",   "move-step.svg", nil, 20.00},
		Step{"50mm",   "move-step.svg", nil, 50.00},
		Step{"0.02mm", "move-step.svg", nil,  0.02},
		Step{"0.1mm",  "move-step.svg", nil,  0.10},
		Step{" 1mm",   "move-step.svg", nil,  1.00},
	)

	instance := &AmountToMoveStepButton{
		StepButton: base,
	}

	return instance
}

func (this *AmountToMoveStepButton) Value() float64  {
	return this.StepButton.Value().(float64)
}
