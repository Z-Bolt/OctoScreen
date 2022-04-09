package uiWidgets

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type CoolDownButton struct {
	*gtk.Button

	client						*octoprintApis.Client
	callback					func()
}

func CreateCoolDownButton(
	client						*octoprintApis.Client,
	callback					func(),
) *CoolDownButton {
	base := utils.MustButtonImage("All Off", "cool-down.svg", nil)

	instance := &CoolDownButton{
		Button:						base,
		client:						client,
		callback:					callback,
	}
	_, err := instance.Button.Connect("clicked", instance.handleClicked)
	if err != nil {
		logger.LogError("PANIC!!! - CreateCoolDownButton()", "instance.Button.Connect()", err)
		panic(err)
	}

	return instance
}

func (this *CoolDownButton) handleClicked() {
	TurnAllHeatersOff(this.client)

	if this.callback != nil {
		this.callback()
	}
}

func TurnAllHeatersOff(
	client						*octoprintApis.Client,
) {
	// Set the bed's temp.
	bedTargetRequest := &octoprintApis.BedTargetRequest {
		Target: 0.0,
	}
	err := bedTargetRequest.Do(client)
	if err != nil {
		logger.LogError("CoolDownButton.TurnAllHeatersOff()", "Do(BedTargetRequest)", err)
		return
	}

	// Set the temp of each hotend.
	hotendCount := utils.GetHotendCount(client)
	for i := 0; i < hotendCount; i++ {
		var toolTargetRequest = &octoprintApis.ToolTargetRequest {
			Targets: map[string]float64 {
				fmt.Sprintf("tool%d", i): 0.0,
			},
		}
		err = toolTargetRequest.Do(client)
		if err != nil {
			logger.LogError("CoolDownButton.TurnAllHeatersOff()", "Do(ToolTargetRequest)", err)
		}
	}
}
