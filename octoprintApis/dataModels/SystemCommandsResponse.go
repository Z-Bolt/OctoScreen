package dataModels


// SystemCommandsResponse is the response to a SystemCommandsRequest.
type SystemCommandsResponse struct {
	Core   []*CommandDefinition `json:"core"`
	Custom []*CommandDefinition `json:"custom"`
}
