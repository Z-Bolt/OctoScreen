package octoprintApis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// Set the current selections for the given tool to the specified spool.
type FilamentManagerSetSelectionRequest struct {
	// Tool ID
	Tool int

	// Spool Id
	Spool int
}

func (this *FilamentManagerSetSelectionRequest) Do(client *Client) (*dataModels.FilamentManagerSelectionsResponse, error) {
	logger.TraceEnter("FilamentManagerSetSelectionRequest.Do()")

	buffer := bytes.NewBuffer(nil)
	uri := fmt.Sprintf(FilamentManagerSetSelectionUri, this.Tool)

	if err := this.encode(buffer); err != nil {
		logger.LogError("FilamentManagerSetSelectionRequest.Do()", "request.encode(buffer)", err)
		logger.TraceLeave("FilamentManagerSetSelectionRequest.Do()")
		return nil, err
	}

	bytes, err := client.doJsonRequest("PATCH", uri, buffer, nil, true)
	if err != nil {
		logger.Infof("doJsonRequest(PATCH)=>%v", string(bytes[:]))
		logger.LogError("FilamentManagerSetSelectionRequest.Do()", "client.doJsonRequest(PATCH)", err)
		logger.TraceLeave("FilamentManagerSetSelectionRequest.Do()")
		return nil, err
	}

	response := &dataModels.FilamentManagerSelectionsResponse {}
	if err := json.Unmarshal(bytes, response); err != nil {
		logger.LogError("FilamentManagerSetSelectionRequests.Do()", "json.Unmarshal()", err)
		logger.TraceLeave("FilamentManagerSetSelectionRequest.Do()")
		return nil, err
	}

	logger.TraceLeave("FilamentManagerSetSelectionRequest.Do()")
	return response, err
}

// The actual PATCH structure expected is a bit unusual, so just convert under the hood.
type FilamentManagerSetSelectionRequestJson struct {
	Selection struct {
		Tool int `json:"tool"`
		Spool struct {
			Id int `json:"id"`
		} `json:"spool"`
	} `json:"selection"`
}

func (this *FilamentManagerSetSelectionRequest) encode(ioWriter io.Writer) error {
	request := FilamentManagerSetSelectionRequestJson{}

	request.Selection.Tool = this.Tool
	request.Selection.Spool.Id = this.Spool

	return json.NewEncoder(ioWriter).Encode(request)
}
