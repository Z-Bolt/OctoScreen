//package apis
package octoprintApis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)

// Set the current selections for the given tool to the specified spool
type FilamentManagerSetSelectionRequest struct {
	// Tool ID
	Tool int

	// Spool Id
	Spool int
}

func (request *FilamentManagerSetSelectionRequest) Do(c *Client) (*dataModels.FilamentManagerSelection, error) {
	buffer := bytes.NewBuffer(nil)
	uri := fmt.Sprintf(FilamentManagerSetSelectionUri, request.Tool)

	if err := request.encode(buffer); err != nil {
		logger.LogError("FilamentManagerSetSelectionRequest.Do()", "request.encode(buffer)", err)
		return nil, err
	}

	resp, err := c.doJsonRequest("PATCH", uri, buffer, nil, true)
	if err != nil {
		logger.Infof("doJsonRequest(PATCH)=>%v", string(resp[:]))
		logger.LogError("FilamentManagerSetSelectionRequest.Do()", "client.doJsonRequest(PATCH)", err)
		return nil, err
	}

	response := &dataModels.FilamentManagerSelection{}
	if err := json.Unmarshal(resp, response); err != nil {
		logger.LogError("FilamentManagerSetSelectionRequests.Do()", "json.Unmarshal()", err)
		return nil, err
	}

	return response, err
}

// The actual PATCH structure expected is a bit unusual, so just convert under
// the hood
type FilamentManagerSetSelectionRequestJson struct {
	Selection struct {
		Tool int `json:"tool"`
		Spool struct {
			Id int `json:"id"`
		} `json:"spool"`
	} `json:"selection"`
}

func (request *FilamentManagerSetSelectionRequest) encode(w io.Writer) error {
	req := FilamentManagerSetSelectionRequestJson{}

	req.Selection.Tool = request.Tool
	req.Selection.Spool.Id = request.Spool

	return json.NewEncoder(w).Encode(req)
}
