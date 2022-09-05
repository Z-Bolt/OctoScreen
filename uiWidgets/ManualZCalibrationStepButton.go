package uiWidgets

type ManualZCalibrationStepButton struct {
	*StepButton
}

func CreateManualZCalibrationStepButton(
	colorVariation		int,
	clicked				func(),
) *ManualZCalibrationStepButton {
	base := CreateStepButton(
		colorVariation,
		clicked,
		Step{"Start Manual\nZ Calibration", "z-calibration.svg", nil, false},
		Step{"Stop Manual\nZ Calibration",  "z-calibration.svg", nil, true},
	)

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
