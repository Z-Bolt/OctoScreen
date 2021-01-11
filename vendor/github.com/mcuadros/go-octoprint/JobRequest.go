package octoprint

import (
	// "bytes"
	"encoding/json"
	// "io"
	// "log"
)

// JobRequest retrieve information about the current job (if there is one).
type JobRequest struct{}



// Job information response
// https://docs.octoprint.org/en/master/api/job.html#job-information-response

// JobResponse is the response from a job command.
type JobResponse struct {
	// Job contains information regarding the target of the current print job.
	Job JobInformation `json:"job"`

	// Progress contains information regarding the progress of the current job.
	Progress ProgressInformation `json:"progress"`

	//State StateInformation `json:"state"`
	State string `json:"state"`
}

// https://docs.octoprint.org/en/master/api/job.html
const JobTool = "/api/job"

// Do sends an API request and returns the API response.
func (cmd *JobRequest) Do(client *Client) (*JobResponse, error) {
	bytes, err := client.doJSONRequest("GET", JobTool, nil, nil)
	if err != nil {
		return nil, err
	}

	response := &JobResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
