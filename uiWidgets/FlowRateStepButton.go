package uiWidgets

import (
	"github.com/mcuadros/go-octoprint"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type FlowRateStepButton struct {
	*StepButton
	client			*octoprint.Client
}

func CreateFlowRateStepButton(
	client			*octoprint.Client,
) *FlowRateStepButton {
	base, err := CreateStepButton(
		1,
		Step{"Normal (100%)", "speed-normal.svg", nil, 100},
		Step{"Fast (125%)",   "speed-fast.svg",   nil, 125},
		Step{"Slow (75%)",    "speed-slow.svg",   nil,  75},
	)
	if err != nil {
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
	cmd := &octoprint.ToolFlowRateRequest{}
	cmd.Factor = this.Value()

	utils.Logger.Infof("FlowRateStepButton.SendChangeFlowRate() - changing flow rate to %d%%", cmd.Factor)
	if err := cmd.Do(this.client); err != nil {
		utils.LogError("FlowRateStepButton.SendChangeFlowRate()", "Go(ToolFlowRateRequest)", err)
		return err
	}

	return nil
}
