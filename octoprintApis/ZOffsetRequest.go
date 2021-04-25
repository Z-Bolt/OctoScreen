package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


type ZOffsetRequest struct {
	Command string `json:"command"`
	Tool    int    `json:"tool"`
}

func (this *ZOffsetRequest) Do(client *Client) (*dataModels.ZOffsetResponse, error) {
	this.Command = "get_z_offset"

	params := bytes.NewBuffer(nil)
	if err := json.NewEncoder(params).Encode(this); err != nil {
		logger.LogError("ZOffsetRequest.Do()", "json.NewEncoder(params).Encode(this)", err)
		return nil, err
	}

	// bytes, err := client.doJsonRequest("POST", UriZBoltRequest, params, ConnectionErrors)
	bytes, err := client.doJsonRequest("GET", PluginZBoltApiUri, params, ConnectionErrors, true)
	if err != nil {
		logger.LogError("ZOffsetRequest.Do()", "client.doJsonRequest()", err)
		return nil, err
	}

	response := &dataModels.ZOffsetResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("ZOffsetRequest.Do()", "json.Unmarshal()", err)
		return nil, err
	}

	return response, err
}




// SetZOffsetRequest - retrieves the current configuration of OctoPrint.
type SetZOffsetRequest struct {
	Command string  `json:"command"`
	Tool    int     `json:"tool"`
	Value   float64 `json:"value"`
}

func (this *SetZOffsetRequest) Do(client *Client) error {
	this.Command = "set_z_offset"

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(this); err != nil {
		logger.LogError("SetZOffsetRequest.Do()", "json.NewEncoder(params).Encode(this)", err)
		return err
	}

	_, err := client.doJsonRequest("POST", PluginZBoltApiUri, buffer, ConnectionErrors, true)
	return err
}
