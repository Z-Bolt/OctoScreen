package octoprint

import (
	"bytes"
	"encoding/json"
)

const URIZBoltRequest = "/api/plugin/zbolt"

// SettingsRequest retrieves the current configuration of OctoPrint.
type SetZOffsetRequest struct {
	Command string  `json:"command"`
	Tool    int     `json:"tool"`
	Value   float64 `json:"value"`
}

func (cmd *SetZOffsetRequest) Do(c *Client) error {
	cmd.Command = "set_z_offset"

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIZBoltRequest, b, ConnectionErrors)
	return err
}

type GetZOffsetRequest struct {
	Command string `json:"command"`
	Tool    int    `json:"tool"`
}

type GetZOffsetResponse struct {
	// Job contains information regarding the target of the current print job.
	Offset float64 `json:"offset"`
}

func (cmd *GetZOffsetRequest) Do(c *Client) (*GetZOffsetResponse, error) {
	cmd.Command = "get_z_offset"

	params := bytes.NewBuffer(nil)
	if err := json.NewEncoder(params).Encode(cmd); err != nil {
		return nil, err
	}

	b, err := c.doJSONRequest("POST", URIZBoltRequest, params, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	r := &GetZOffsetResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

type GetNotificationRequest struct {
	Command string `json:"command"`
}
type GetNotificationResponse struct {
	// Job contains information regarding the target of the current print job.
	Message string `json:"message"`
}

func (cmd *GetNotificationRequest) Do(c *Client) (*GetNotificationResponse, error) {
	cmd.Command = "get_notification"

	params := bytes.NewBuffer(nil)
	if err := json.NewEncoder(params).Encode(cmd); err != nil {
		return nil, err
	}

	b, err := c.doJSONRequest("POST", URIZBoltRequest, params, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	r := &GetNotificationResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}
