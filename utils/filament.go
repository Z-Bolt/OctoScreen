package utils

import (
	// "errors"

	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis"
	// "github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


func Extrude(
	client					*octoprintApis.Client,
	isForward				bool,
	extruderId				string,
	parentWindow			*gtk.Window,
	flowRatePercentage		int,
	length					int,
) {
	var action string
	if isForward {
		action = "extrude"
	} else {
		action = "retract"
	}

	if CheckIfHotendTemperatureIsTooLow(client, extruderId, action, parentWindow) {
		logger.Error("filament.Extrude() - temperature is too low")
		// No need to display an error - CheckIfHotendTemperatureIsTooLow() displays an error to the user
		// if the temperature is too low.
		return
	}

	if err := SelectTool(client, extruderId); err != nil {
		// TODO: display error?
		return
	}

	if err := SetFlowRate(client, flowRatePercentage); err != nil {
		// TODO: display error?
		return
	}

	if err := SendExtrudeRrequest(client, isForward, length); err != nil {
		// TODO: display error?
	}
}

func SelectTool(
	client					*octoprintApis.Client,
	extruderId				string,
) error {
	cmd := &octoprintApis.ToolSelectRequest{}
	cmd.Tool = extruderId

	logger.Infof("filament.SelectTool() - changing tool to %s", cmd.Tool)
	if err := cmd.Do(client); err != nil {
		logger.LogError("filament.SelectTool()", "Go(ToolSelectRequest)", err)
		return err
	}

	return nil
}

func SetFlowRate(
	client					*octoprintApis.Client,
	flowRatePercentage		int,
) error {
	cmd := &octoprintApis.ToolFlowRateRequest{}
	cmd.Factor = flowRatePercentage

	logger.Infof("filament.SetFlowRate() - changing flow rate to %d%%", cmd.Factor)
	if err := cmd.Do(client); err != nil {
		logger.LogError("filament.SetFlowRate()", "Go(ToolFlowRateRequest)", err)
		return err
	}

	return nil
}

func SendExtrudeRrequest(
	client					*octoprintApis.Client,
	isForward				bool,
	length					int,
) error {
	cmd := &octoprintApis.ToolExtrudeRequest{}
	if isForward {
		cmd.Amount = length
	} else {
		cmd.Amount = -length
	}

	logger.Infof("filament.SendExtrudeRrequest() - sending extrude request with length of: %d", cmd.Amount)
	if err := cmd.Do(client); err != nil {
		logger.LogError("filament.Extrude()", "Do(ToolExtrudeRequest)", err)
		return err
	}

	return nil
}
