package uiWidgets

import (
	"fmt"
	// "strconv"
	// "strings"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)


func CreateSelectHotendStepButton(
	client				*octoprintApis.Client,
	includeBed			bool,
	colorVariation		int,
	clicked				func(),
) *SelectToolStepButton {
	hotendCount := utils.GetHotendCount(client)

	var steps []Step
	for i := 0; i < hotendCount; i++ {
		label := ""
		if i == 0 && hotendCount == 1 {
			label = "Hotend"
		} else {
			label = fmt.Sprintf("Hotend %d", i + 1)
		}

		step := Step {
			label,
			utils.GetHotendFileName(i + 1, hotendCount),
			nil,
			fmt.Sprintf("tool%d", i),
		}

		steps = append(steps, step)
	}

	if includeBed {
		steps = append(steps, Step{"Bed", "bed.svg", nil, "bed"})
	}

	base, err := CreateStepButton(
		colorVariation,
		clicked,
		steps...,
	)
	if err != nil {
		logger.LogError("PANIC!!! - CreateSelectHotendStepButton()", "CreateStepButton()", err)
		panic(err)
	}

	instance := &SelectToolStepButton{
		StepButton: base,
	}

	return instance
}
