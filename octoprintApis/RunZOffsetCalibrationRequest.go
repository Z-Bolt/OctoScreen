package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"

	"github.com/Z-Bolt/OctoScreen/logger"
)


type RunZOffsetCalibrationRequest struct {
	Command string `json:"command"`
}

func (this *RunZOffsetCalibrationRequest) Do(client *Client) error {
	this.Command = "run_zoffset_calibration"

	bytes := bytes.NewBuffer(nil)
	if err := json.NewEncoder(bytes).Encode(this); err != nil {
		logger.LogError("RunZOffsetCalibrationRequest.Do()", "json.NewEncoder(params).Encode(this)", err)
		return err
	}

	_, err := client.doJsonRequest("POST", PluginZBoltOctoScreenApiUri, bytes, ConnectionErrors)
	return err
}
