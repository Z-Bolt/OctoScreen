package octoprintApis

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"log"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


type ZOffsetRequest struct {
	Command string `json:"command"`
	Tool    int    `json:"tool"`
}

func (cmd *ZOffsetRequest) Do(c *Client) (*dataModels.ZOffsetResponse, error) {
	cmd.Command = "get_z_offset"

	params := bytes.NewBuffer(nil)
	if err := json.NewEncoder(params).Encode(cmd); err != nil {
		log.Println("zbolt.Do() 3 - Encode() failed")
		return nil, err
	}

	// b, err := c.doJsonRequest("POST", URIZBoltRequest, params, ConnectionErrors)
	b, err := c.doJsonRequest("GET", PluginZBoltApiUri, params, ConnectionErrors)
	if err != nil {
		log.Println("zbolt.Do() 3 - doJsonRequest() failed")
		return nil, err
	}

	r := &dataModels.ZOffsetResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		log.Println("zbolt.Do() 3 - Unmarshal() failed")
		return nil, err
	}

	return r, err
}




// SetZOffsetRequest - retrieves the current configuration of OctoPrint.
type SetZOffsetRequest struct {
	Command string  `json:"command"`
	Tool    int     `json:"tool"`
	Value   float64 `json:"value"`
}

func (cmd *SetZOffsetRequest) Do(c *Client) error {
	cmd.Command = "set_z_offset"

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		log.Println("zbolt.Do() 2 - Encode() failed")
		return err
	}

	_, err := c.doJsonRequest("POST", PluginZBoltApiUri, b, ConnectionErrors)
	return err
}
