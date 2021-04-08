//package apis
package octoprintApis

import (
	"encoding/json"
	"fmt"
	// "strconv"
	// "strings"
	// "time"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// Retrieve a specific file’s or folder’s information
// GET /api/files/(string:location)/(path:filename)
// https://docs.octoprint.org/en/master/api/files.html#retrieve-a-specific-file-s-or-folder-s-information

// FileRequest retrieves the selected file’s or folder’s information.
type FileRequest struct {
	// Location of the file for which to retrieve the information/  Can be either
	// `local` or `sdcard`.
	Location dataModels.Location

	// Filename of the file for which to retrieve the information.
	Filename string

	// Recursive if set to true, return all files and folders recursively.
	// Otherwise only return items on same level.
	Recursive bool
}

// Do sends an API request and returns the API response
func (request *FileRequest) Do(c *Client) (*dataModels.FileResponse, error) {
	uri := fmt.Sprintf("%s/%s/%s?recursive=%t",
		FilesApiUri,
		request.Location,
		request.Filename,
		request.Recursive,
	)

	bytes, err := c.doJsonRequest("GET", uri, nil, FilesLocationGETErrors, true)
	if err != nil {
		return nil, err
	}

	response := &dataModels.FileResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
