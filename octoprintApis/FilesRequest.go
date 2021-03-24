package octoprintApis

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "io"
	// "mime/multipart"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// FilesRequest retrieve information regarding all files currently available and
// regarding the disk space still available locally in the system.
type FilesRequest struct {
	// Location is the target location .
	Location dataModels.Location

	// Recursive if set to true, return all files and folders recursively.
	// Otherwise only return items on same level.
	Recursive bool
}

// Do sends an API request and returns the API response.
func (cmd *FilesRequest) Do(c *Client) (*dataModels.FilesResponse, error) {
	uri := fmt.Sprintf("%s?recursive=%t", FilesApiUri, cmd.Recursive)
	if cmd.Location != "" {
		uri = fmt.Sprintf("%s/%s?recursive=%t", FilesApiUri, cmd.Location, cmd.Recursive)
	}

	b, err := c.doJsonRequest("GET", uri, nil, FilesLocationGETErrors)
	if err != nil {
		return nil, err
	}

	r := &dataModels.FilesResponse{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	if len(r.Children) > 0 {
		r.Files = r.Children
	}

	return r, err
}
