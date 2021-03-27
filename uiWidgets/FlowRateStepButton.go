package uiWidgets

import (
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type FlowRateStepButton struct {
	*StepButton
	client			*octoprintApis.Client
}

func CreateFlowRateStepButton(
	client			*octoprintApis.Client,
) *FlowRateStepButton {
	base, err := CreateStepButton(
		1,
		Step{"Normal (100%)", "speed-normal.svg", nil, 100},
		Step{"Fast (125%)",   "speed-fast.svg",   nil, 125},
		Step{"Slow (75%)",    "speed-slow.svg",   nil,  75},
	)
	if err != nil {
		logger.LogError("PANIC!!! - CreateFlowRateStepButton()", "CreateStepButton()", err)
		panic(err)
	}

	instance := &FlowRateStepButton{
		StepButton:		base,
		client:			client,
	}

	return instance
}

func (this *FlowRateStepButton) Value() int {
	return this.StepButton.Value().(int)
}

func (this *FlowRateStepButton) SendChangeFlowRate() error {
	return utils.SetFlowRate(this.client, this.Value())
}
