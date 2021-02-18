package uiWidgets

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type SelectToolStepButton struct {
	*StepButton
}

func CreateSelectExtruderStepButton(
	client							*octoprintApis.Client,
	includeBed						bool,
) *SelectToolStepButton {
	extruderCount := utils.GetExtruderCount(client)

	var steps []Step
	for i := 0; i < extruderCount; i++ {
		var step Step
		if i == 0 && extruderCount == 1 {
			step = Step {
				"Extruder",
				utils.GetExtruderFileName(1, extruderCount),
				nil,
				"tool0",
			}
		} else {
			step = Step {
				fmt.Sprintf("Extruder %d", i + 1),
				utils.GetExtruderFileName(i + 1, extruderCount),
				nil,
				fmt.Sprintf("tool%d", i),
			}
		}

		steps = append(steps, step)
	}

	if includeBed {
		steps = append(steps, Step{"Bed", "bed.svg", nil, "bed"})
	}

	base, err := CreateStepButton(
		1,
		steps...,
	)
	if err != nil {
		utils.LogError("PANIC!!! - CreateSelectExtruderStepButton()", "CreateStepButton()", err)
		panic(err)
	}

	instance := &SelectToolStepButton{
		StepButton: base,
	}

	return instance
}

func CreateSelectHotendStepButton(
	client							*octoprintApis.Client,
	includeBed						bool,
) *SelectToolStepButton {
	hotendCount := utils.GetHotendCount(client)

	var steps []Step
	for i := 0; i < hotendCount; i++ {
		var step Step
		if i == 0 && hotendCount == 1 {
			step = Step {
				"Hotend",
				utils.GetHotendFileName(1, hotendCount),
				nil,
				"tool0",
			}
		} else {
			step = Step {
				fmt.Sprintf("Hotend %d", i + 1),
				utils.GetHotendFileName(i + 1, hotendCount),
				nil,
				fmt.Sprintf("tool%d", i),
			}
		}

		steps = append(steps, step)
	}

	if includeBed {
		steps = append(steps, Step{"Bed", "bed.svg", nil, "bed"})
	}

	base, err := CreateStepButton(
		1,
		steps...,
	)
	if err != nil {
		utils.LogError("PANIC!!! - CreateSelectHotendStepButton()", "CreateStepButton()", err)
		panic(err)
	}

	instance := &SelectToolStepButton{
		StepButton: base,
	}

	return instance
}




func (this *SelectToolStepButton) Value() string  {
	return this.StepButton.Value().(string)
}

func (this *SelectToolStepButton) Index() int  {
	value := strings.Replace(this.Value(), "tool", "", -1)
	index, _ := strconv.Atoi(value)

	return index
}
