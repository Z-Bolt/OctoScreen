package octoprintApis

import (
	"fmt"

	"github.com/Z-Bolt/OctoScreen/octoprintApis/dataModels"
)


// DeleteFileRequest delete the selected path on the selected location.
type DeleteFileRequest struct {
	// Location is the target location on which to delete the file, either
	// `local` (for OctoPrint’s uploads folder) or \sdcard\ for the printer’s
	// SD card (if available)
	Location dataModels.Location

	// Path of the file to delete
	Path string
}

// Do sends an API request and returns error if any.
func (this *DeleteFileRequest) Do(client *Client) error {
	uri := fmt.Sprintf("%s/%s/%s", FilesApiUri, this.Location, this.Path)
	if _, err := client.doJsonRequest("DELETE", uri, nil, FilesLocationDeleteErrors, true); err != nil {
		return err
	}

	return nil
}
