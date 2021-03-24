package dataModels


// CustomCommandsResponse is the response to a CustomCommandsRequest.
type CustomCommandsResponse struct {
	Controls []*ControlContainer `json:"controls"`
}
