package octoprint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
		log.Println("zbolt.Do() 1 - NewEncoder() failed")
		return err
	}

	_, err := c.doJSONRequest("POST", URIZBoltRequest, b, ConnectionErrors)
	return err
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
		log.Println("zbolt.Do() 3 - Encode() failed")
		return nil, err
	}

	// b, err := c.doJSONRequest("POST", URIZBoltRequest, params, ConnectionErrors)
	b, err := c.doJSONRequest("GET", URIZBoltRequest, params, ConnectionErrors)
	if err != nil {
		log.Println("zbolt.Do() 3 - doJSONRequest() failed")
		return nil, err
	}

	r := &GetZOffsetResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		log.Println("zbolt.Do() 3 - Unmarshal() failed")
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
	target := fmt.Sprintf("%s?command=get_notification", URIZBoltOctoScreenRequest)
	bytes, err := c.doJSONRequest("GET", target, nil, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	response := &GetNotificationResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}


type GetSettingsRequest struct {
	Command string `json:"command"`
}

type MenuItem struct {
	Name  string     `json:"name"`
	Icon  string     `json:"icon"`
	Panel string     `json:"panel"`
	Items []MenuItem `json:"items"`
}

type GetSettingsResponse struct {
	// Job contains information regarding the target of the current print job.
	FilamentInLength  float64    `json:"filament_in_length"`
	FilamentOutLength float64    `json:"filament_out_length"`
	ToolChanger       bool       `json:"toolchanger"`
	XAxisInverted     bool       `json:"x_axis_inverted"`
	YAxisInverted     bool       `json:"y_axis_inverted"`
	ZAxisInverted     bool       `json:"z_axis_inverted"`
	MenuStructure     []MenuItem `json:"menu_structure"`
	GCodes            struct {
		AutoBedLevel string `json:"auto_bed_level"`
	} `json:"gcodes"`
}

func (cmd *GetSettingsRequest) Do(c *Client) (*GetSettingsResponse, error) {
	target := fmt.Sprintf("%s?command=get_settings", URIZBoltOctoScreenRequest)
	bytes, err := c.doJSONRequest("GET", target, nil, ConnectionErrors)
	if err != nil {
		return nil, err
	}

	response := &GetSettingsResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
