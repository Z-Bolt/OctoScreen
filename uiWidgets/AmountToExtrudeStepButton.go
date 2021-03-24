package uiWidgets

import (
	"github.com/Z-Bolt/OctoScreen/logger"
)

type AmountToExtrudeStepButton struct {
	*StepButton
}

func CreateAmountToExtrudeStepButton() *AmountToExtrudeStepButton {
	base, err := CreateStepButton(
		2,
		Step{" 20mm", "move-step.svg", nil,  20},
		Step{" 50mm", "move-step.svg", nil,  50},
		Step{"100mm", "move-step.svg", nil, 100},
		Step{"  1mm", "move-step.svg", nil,   1},
		Step{"  5mm", "move-step.svg", nil,   5},
		Step{" 10mm", "move-step.svg", nil,  10},
	)
	if err != nil {
		logger.LogError("PANIC!!! - CreateAmountToExtrudeStepButton()", "CreateStepButton()", err)
		panic(err)
	}

	instance := &AmountToExtrudeStepButton{
		StepButton: base,
	}

	return instance
}

func (this *AmountToExtrudeStepButton) Value() int {
	return this.StepButton.Value().(int)
}
