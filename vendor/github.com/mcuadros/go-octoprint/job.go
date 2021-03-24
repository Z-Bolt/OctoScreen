package octoprint

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)


var JobToolErrors = StatusMapping{
	409: "Printer is not operational or the current print job state does not match the preconditions for the command.",
}




// JobRequest is now in JobRequest.go
// // JobRequest retrieve information about the current job (if there is one).
// type JobRequest struct{}

// // Do sends an API request and returns the API response.
// func (cmd *JobRequest) Do(c *Client) (*JobResponse, error) {
// 	bytes, err := c.doJSONRequest("GET", JobTool, nil, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	response := &JobResponse{}
// 	if err := json.Unmarshal(bytes, response); err != nil {
// 		return nil, err
// 	}

// 	return response, err
// }





// StartRequest starts the print of the currently selected file.
type StartRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *StartRequest) Do(c *Client) error {
	payload := map[string]string{"command": "start"}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", JobTool, buffer, JobToolErrors)
	return err
}

// CancelRequest cancels the current print job.
type CancelRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *CancelRequest) Do(c *Client) error {
	payload := map[string]string{"command": "cancel"}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", JobTool, buffer, JobToolErrors)
	return err
}

// RestartRequest restart the print of the currently selected file from the
// beginning. There must be an active print job for this to work and the print
// job must currently be paused
type RestartRequest struct{}

// Do sends an API request and returns an error if any.
func (cmd *RestartRequest) Do(c *Client) error {
	payload := map[string]string{"command": "restart"}

	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(payload); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", JobTool, buffer, JobToolErrors)
	return err
}

type PauseAction string

const (
	// Pause the current job if it’s printing, does nothing if it’s already paused.
	Pause PauseAction = "pause"

	// Resume the current job if it’s paused, does nothing if it’s printing.
	Resume PauseAction = "resume"

	// Toggle the pause state of the job, pausing it if it’s printing and
	// resuming it if it’s currently paused.
	Toggle PauseAction = "toggle"
)

// PauseRequest pauses/resumes/toggles the current print job.
type PauseRequest struct {
	// Action specifies which action to take.
	// In order to stay backwards compatible to earlier iterations of this API,
	// the default action to take if no action parameter is supplied is to
	// toggle the print job status.
	Action PauseAction `json:"action"`
}

// Do sends an API request and returns an error if any.
func (cmd *PauseRequest) Do(c *Client) error {
	buffer := bytes.NewBuffer(nil)
	if err := cmd.encode(buffer); err != nil {
		return err
	}

	_, err := c.doJSONRequest("POST", JobTool, buffer, JobToolErrors)
	return err
}






// Do sends an API request and returns an error if any.
func (cmd *PauseRequest) DoWithLogging(c *Client) error {
	log.Println("job.Do() - starting call")

	buffer := bytes.NewBuffer(nil)
	err := cmd.encode(buffer)
	if err != nil {
		log.Println("job.Do() - encode() failed")
		return err
	}

	// str := string([]byte(buffer.Bytes))
	str := string(buffer.Bytes())
	log.Printf("job.Do() - bytes are: '%s' abc", str)


	log.Println("job.Do() - about to call POST")
	body, err2 := c.doJSONRequestWithLogging("POST", JobTool, buffer, JobToolErrors)

	log.Println("job.Do() - finished calling POST")
	if err2 != nil {
		log.Println("job.Do() - call to doJSONRequest() failed")
		log.Printf("job.Do() - err2 is: %s", err2.Error())
	}

	if err2 == nil {
		log.Println("job.Do() - call to doJSONRequest() passed")
		str = string(body)
		log.Printf("job.Do() - body is: %s",  str)
	}

	return err
}




func (cmd *PauseRequest) encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(struct {
		Command string `json:"command"`
		PauseRequest
	}{
		Command:      "pause",
		PauseRequest: *cmd,
	})
}
