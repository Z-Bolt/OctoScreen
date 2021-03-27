package dataModels

// Job information response
// https://docs.octoprint.org/en/master/api/job.html#job-information-response

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
