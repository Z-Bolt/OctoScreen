package octoprintApis

import (
	// "bytes"
	"encoding/json"
	// "io"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// JobRequest retrieve information about the current job (if there is one).
type JobRequest struct{}


// Do sends an API request and returns the API response.
func (cmd *JobRequest) Do(client *Client) (*dataModels.JobResponse, error) {
	bytes, err := client.doJsonRequest("GET", JobApiUri, nil, nil, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.JobResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
