package uiWidgets

import (
	"github.com/Z-Bolt/OctoScreen/logger"
)

type ManualZCalibrationStepButton struct {
	*StepButton
}

func CreateManualZCalibrationStepButton() *ManualZCalibrationStepButton {
	base, err := CreateStepButton(
		1,
		Step{"Start Manual\nZ Calibration", "z-calibration.svg", nil, false},
		Step{"Stop Manual\nZ Calibration",  "z-calibration.svg", nil, true},
	)
	if err != nil {
		logger.LogError("PANIC!!! - CreateManualZCalibrationStepButton()", "CreateStepButton()", err)
		panic(err)
	}

	instance := &ManualZCalibrationStepButton{
		StepButton: base,
	}

	return instance
}

// The value returned represents if it is running (true) or if idle (false).
func (this *ManualZCalibrationStepButton) Value() bool {
	return this.StepButton.Value().(bool)
}

func (this *ManualZCalibrationStepButton) IsCalibrating() bool {
	return this.Value()
}
