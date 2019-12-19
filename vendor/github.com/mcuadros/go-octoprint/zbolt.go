package octoprint

import (
	"bytes"
	"encoding/json"
)

const URIZBoltRequest = "/api/plugin/zbolt"
const URIZBoltOctoScreenRequest = "/api/plugin/zbolt_octoscreen"

type RunZOffsetCalibrationRequest struct {
	Command string `json:"command"`
}

func (cmd *RunZOffsetCalibrationRequest) Do(c *Client) error {
	cmd.Command = "run_zoffset_calibration"

	b := bytes.NewBuffer(nil)
	if err := json.NewEncoder(b).Encode(cmd); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", URIZBoltRequest, b, ConnectionErrors)
	return err
}

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

	b, err := c.doJSONRequest("POST", URIZBoltOctoScreenRequest, params, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	r := &GetNotificationResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

type GetSettingsRequest struct {
	Command string `json:"command"`
}
type GetSettingsResponse struct {
	// Job contains information regarding the target of the current print job.
	FilamentInLength  float64    `json:"filament_in_length"`
	FilamentOutLength float64    `json:"filament_out_length"`
	ToolChanger       bool       `json:"toolchanger"`
	ZAxisInverted     bool       `json:"z_axis_inverted"`
	MenuStructure     []MenuItem `json:"menu_structure"`
	GCodes            struct {
		AutoBedLevel string `json:"auto_bed_level"`
	} `json:"gcodes"`
}

func (cmd *GetSettingsRequest) Do(c *Client) (*GetSettingsResponse, error) {
	cmd.Command = "get_settings"

	params := bytes.NewBuffer(nil)
	if err := json.NewEncoder(params).Encode(cmd); err != nil {
		return nil, err
	}

	b, err := c.doJSONRequest("POST", URIZBoltOctoScreenRequest, params, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	r := &GetSettingsResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	return r, err
}

type MenuItem struct {
	Name  string     `json:"name"`
	Icon  string     `json:"icon"`
	Panel string     `json:"panel"`
	Items []MenuItem `json:"items"`
}
