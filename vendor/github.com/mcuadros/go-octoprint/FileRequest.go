//package apis
package octoprint

import (
	"fmt"
	"encoding/json"

	// "log"
	// "strconv"
	// "strings"
	// "time"

	// "../octoprint"
	//"github.com/mcuadros/go-octoprint"
)

// Retrieve a specific file’s or folder’s information
// GET /api/files/(string:location)/(path:filename)
// https://docs.octoprint.org/en/master/api/files.html#retrieve-a-specific-file-s-or-folder-s-information


// FileRequest retrieves the selected file’s or folder’s information.
type FileRequest struct {
	// Location of the file for which to retrieve the information/  Can be either
	// `local` or `sdcard`.
	Location Location

	// Filename of the file for which to retrieve the information.
	Filename string

	// Recursive if set to true, return all files and folders recursively.
	// Otherwise only return items on same level.
	Recursive bool
}

// FileResponse contains information regarding a file.
// https://docs.octoprint.org/en/master/api/datamodel.html#file-information
type FileResponse struct {
	// Name is name of the file without path. E.g. “file.gco” for a file
	// “file.gco” located anywhere in the file system.
	Name string `json:"name"`

	// The name of the file without the path.
	Display string `json:"display"`

	// Path is the path to the file within the location. E.g.
	//“folder/subfolder/file.gco” for a file “file.gco” located within “folder”
	// and “subfolder” relative to the root of the location.
	Path string `json:"path"`

	// Type of file. model or machinecode.  Or folder if it’s a folder, in
	// which case the children node will be populated.
	Type string `json:"type"`

	// TypePath path to type of file in extension tree. E.g. `["model", "stl"]`
	// for .stl files, or `["machinecode", "gcode"]` for .gcode files.
	// `["folder"]` for folders.
	TypePath []string `json:"typePath"`




	// Additional properties depend on type. For a type value of folder, see Folders. For any other value see Files.

	// * Folders
	//     --children
	//     --size



	// * Files
	// Hash is the MD5 hash of the file.  Only available for `local` files.
	Hash string `json:"hash"`

	// Size of the file in bytes.  Only available for `local` files or `sdcard`
	// files if the printer supports file sizes for sd card files.
	Size uint64 `json:"size"`

	// Date when this file was uploaded.  Only available for `local` files.
	Date JSONTime `json:"date"`

	// Origin of the file, `local` when stored in OctoPrint’s `uploads` folder,
	// `sdcard` when stored on the printer’s SD card (if available).
	Origin string `json:"origin"`

	// Refs references relevant to this file, left out in abridged version.
	Refs Reference `json:"refs"`

	// GCodeAnalysis information from the analysis of the GCODE file, if
	// available. Left out in abridged version.
	GCodeAnalysis GCodeAnalysisInformation `json:"gcodeAnalysis"`




	// * Additional properties not listed in the SDK...

	// Print information from the print stats of a file.
	Print PrintStats `json:"print"`


	// Relative path to the preview thumbnail image (if it exists)
	// The PrusaSlicer Thumbnails plug-in is required or this.
	Thumbnail string `json:"thumbnail"`
}

// IsFolder it returns true if the file is a folder.
func (response *FileResponse) IsFolder() bool {
	if len(response.TypePath) == 1 && response.TypePath[0] == "folder" {
		return true
	}

	return false
}

const URIFiles = "/api/files"

// Do sends an API request and returns the API response
func (request *FileRequest) Do(c *Client) (*FileResponse, error) {
	uri := fmt.Sprintf("%s/%s/%s?recursive=%t",
		URIFiles,
		request.Location,
		request.Filename,
		request.Recursive,
	)

	bytes, err := c.doJSONRequest("GET", uri, nil, FilesLocationGETErrors)
	if err != nil {
		return nil, err
	}

	response := &FileResponse{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, err
}
