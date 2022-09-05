package uiWidgets

type AmountToExtrudeStepButton struct {
	*StepButton
}

func CreateAmountToExtrudeStepButton(
	colorVariation					int,
	clicked							func(),
) *AmountToExtrudeStepButton {
	base := CreateStepButton(
		colorVariation,
		clicked,
		Step{" 20mm", "move-step.svg", nil,  20},
		Step{" 50mm", "move-step.svg", nil,  50},
		Step{"100mm", "move-step.svg", nil, 100},
		Step{"  1mm", "move-step.svg", nil,   1},
		Step{"  5mm", "move-step.svg", nil,   5},
		Step{" 10mm", "move-step.svg", nil,  10},
	)

	instance := &AmountToExtrudeStepButton {
		StepButton: base,
	}

	return instance
}

func (this *AmountToExtrudeStepButton) Value() int {
	return this.StepButton.Value().(int)
}
