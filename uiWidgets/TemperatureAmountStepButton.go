package uiWidgets

import (
	"github.com/Z-Bolt/OctoScreen/logger"
)

type TemperatureAmountStepButton struct {
	*StepButton
}

func CreateTemperatureAmountStepButton() *TemperatureAmountStepButton {
	base, err := CreateStepButton(
		1,
		Step{"10°C", "move-step.svg", nil, 10.0},
		Step{"20°C", "move-step.svg", nil, 20.0},
		Step{"50°C", "move-step.svg", nil, 50.0},
		Step{" 1°C", "move-step.svg", nil,  1.0},
		Step{" 5°C", "move-step.svg", nil,  5.0},
	)
	if err != nil {
		logger.LogError("PANIC!!! - CreateTemperatureAmountStepButton()", "CreateStepButton()", err)
		panic(err)
	}

	instance := &TemperatureAmountStepButton{
		StepButton: base,
	}

	return instance
}

func (this *TemperatureAmountStepButton) Value() float64 {
	return this.StepButton.Value().(float64)
}
