package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"log"
)


type RunZOffsetCalibrationRequest struct {
	Command string `json:"command"`
}

func (cmd *RunZOffsetCalibrationRequest) Do(c *Client) error {
	cmd.Command = "run_zoffset_calibration"

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		log.Println("zbolt.Do() 1 - NewEncoder() failed")
		return err
	}

	_, err := c.doJsonRequest("POST", PluginZBoltOctoScreenApiUri, b, ConnectionErrors)
	return err
}
