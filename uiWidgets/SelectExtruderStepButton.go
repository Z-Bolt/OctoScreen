package uiWidgets

import (
	"fmt"
	// "strconv"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

func CreateSelectExtruderStepButton(
	client *octoprintApis.Client,
	includeBed bool,
	colorVariation int,
	clicked func(),
) *SelectToolStepButton {
	extruderCount := utils.GetExtruderCount(client)

	var steps []Step
	for i := 0; i < extruderCount; i++ {
		label := ""
		if i == 0 && extruderCount == 1 {
			label = "Extruder"
		} else {
			label = fmt.Sprintf("Extruder %d", i+1)
		}

		step := Step{
			label,
			utils.GetExtruderFileName(i+1, extruderCount),
			nil,
			fmt.Sprintf("tool%d", i),
		}

		steps = append(steps, step)
	}

	if includeBed {
		steps = append(steps, Step{"Bed", "bed.svg", nil, "bed"})
	}

	base := CreateStepButton(
		colorVariation,
		clicked,
		steps...,
	)

	instance := &SelectToolStepButton{
		StepButton: base,
	}

	return instance
}
