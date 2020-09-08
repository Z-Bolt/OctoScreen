package uiWidgets

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type SelectToolStepButton struct {
	*StepButton
}

func CreateSelectToolStepButton(
	client							*octoprint.Client,
	useToolheadForDisplay			bool,
	includeBed						bool,
) *SelectToolStepButton {
	strLabel:= ""
	if useToolheadForDisplay {
		strLabel = "Toolhead"
	} else {
		strLabel = "Hotend"
	}

	toolheadCount := utils.GetToolheadCount(client)

	var steps []Step
	for i := 0; i < toolheadCount; i++ {
		var step Step
		if i == 0 && toolheadCount == 1 {
			step = Step {
				strLabel,
				"toolhead.svg",
				nil,
				"tool0",
			}
		} else {
			step = Step {
				fmt.Sprintf("%s %d", strLabel, i + 1),
				fmt.Sprintf("toolhead-%d.svg", i),
				nil,
				fmt.Sprintf("tool%d", i),
			}
		}

		steps = append(steps, step)
	}

	if includeBed {
		steps = append(steps, Step{"Bed", "bed.svg", nil, "bed"})
	}

	base, err := CreateStepButton(steps...)
	if err != nil {
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
