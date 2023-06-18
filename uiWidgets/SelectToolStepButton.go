package uiWidgets

import (
	// "fmt"
	"strconv"
	"strings"

	// "github.com/Z-Bolt/OctoScreen/logger"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	// "github.com/Z-Bolt/OctoScreen/utils"
)

type SelectToolStepButton struct {
	*StepButton
}

func (this *SelectToolStepButton) Value() string  {
	return this.StepButton.Value().(string)
}

func (this *SelectToolStepButton) Index() int  {
	value := strings.Replace(this.Value(), "tool", "", -1)
	index, _ := strconv.Atoi(value)

	return index
}
