package octoprint

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Axis string

const (
	XAxis Axis = "x"
	YAxis Axis = "y"
	ZAxis Axis = "z"
)

// FullStateResponse contains informantion about the current state of the printer.
type FullStateResponse struct {
	//Temperature is the printer’s temperature state data.
	Temperature TemperatureState `json:"temperature"`
	// SD is the printer’s sd state data.
	SD SDState `json:"sd"`
	// State is the printer’s general state.
	State PrinterState `json:"state"`
}

// JobResponse is the response from a job command.
type JobResponse struct {
	// Job contains information regarding the target of the current print job.
	Job JobInformation `json:"job"`
	// Progress contains information regarding the progress of the current job.
	Progress ProgressInformation `json:"progress"`
}

// JobInformation contains information regarding the target of the current job.
type JobInformation struct {
	// File is the file that is the target of the current print job.
	File FileInformation `json:"file"`
	// EstimatedPrintTime is the estimated print time for the file, in seconds.
	EstimatedPrintTime float64 `json:"estimatedPrintTime"`
	// LastPrintTime is the print time of the last print of the file, in seconds.
	LastPrintTime float64 `json:"lastPrintTime"`
	// Filament contains Information regarding the estimated filament
	// usage of the print job.
	Filament struct {
		// Length of filament used, in mm
		Length float64 `json:"length"`
		// Volume of filament used, in cm³
		Volume float64 `json:"volume"`
	} `json:"filament"`
	FilePosition int `json:"filepos"`
}

// ProgressInformation contains information regarding the progress of the
// current print job.
type ProgressInformation struct {
	// Completion percentage of completion of the current print job.
	Completion float64 `json:"completion"`
	// FilePosition current position in the file being printed, in bytes
	// from the beginning.
	FilePosition int `json:"filepos"`
	// PrintTime is time already spent printing, in seconds
	PrintTime float64 `json:"printTime"`
	// PrintTimeLeft is estimate of time left to print, in seconds
	PrintTimeLeft float64 `json:"printTimeLeft"`
}

// TemperatureState is the printer’s temperature state data.
type TemperatureState temperatureState
type temperatureState struct {
	// Current temperature stats.
	Current map[string]TemperatureData `json:"current"`
	// Temperature history.
	History []*HistoricTemperatureData `json:"history"`
}

func (r *TemperatureState) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	history := raw["history"]
	delete(raw, "history")
	b, _ = json.Marshal(map[string]interface{}{
		"current": raw,
		"history": history,
	})

	i := &temperatureState{}
	if err := json.Unmarshal(b, i); err != nil {
		return err
	}

	*r = TemperatureState(*i)
	return nil
}

// TemperatureData is temperature stats for a tool.
type TemperatureData struct {
	// Actual current temperature.
	Actual float64 `json:"actual"`
	// Target temperature, may be nil if no target temperature is set.
	Target float64 `json:"target"`
	// Offset currently configured temperature offset to apply, will be left
	// out for historic temperature information.
	Offset float64 `json:"offset"`
}

// PrinterState current state of the printer.
type PrinterState struct {
	Text  string `json:"text"`
	Flags struct {
		Operations    bool `json:"operational"`
		Paused        bool `json:"paused"`
		Printing      bool `json:"printing"`
		SDReady       bool `json:"sdReady"`
		Error         bool `json:"error"`
		Ready         bool `json:"ready"`
		ClosedOnError bool `json:"closedOrError"`
	} `json:"flags"`
}

// SDState is the state of the sd reader.
type SDState struct {
	Ready bool `json:"ready"`
}

// HistoricTemperatureData is temperature historic stats for a tool.
type HistoricTemperatureData historicTemperatureData
type historicTemperatureData struct {
	// Time of this data point.
	Time JSONTime `json:"time"`
	// Tools is temperature stats a set of tools.
	Tools map[string]TemperatureData `json:"tools"`
}

func (h *HistoricTemperatureData) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	ts := raw["time"]
	delete(raw, "time")
	b, _ = json.Marshal(map[string]interface{}{
		"time":  ts,
		"tools": raw,
	})

	i := &historicTemperatureData{}
	if err := json.Unmarshal(b, i); err != nil {
		return err
	}

	*h = HistoricTemperatureData(*i)
	return nil
}

// VersionResponse is the response from a job command.
type VersionResponse struct {
	// API is the API version.
	API string `json:"api"`
	// Server is the server version.
	Server string `json:"server"`
}

type ConnectionState string

const (
	Operational ConnectionState = "Operational"
)

// The states are  based on:
// https://github.com/foosel/OctoPrint/blob/77753ca02602d3a798d6b0a22535e6fd69ff448a/src/octoprint/util/comm.py#L549
// no comments :(

func (s ConnectionState) IsOperational() bool {
	return strings.HasPrefix(string(s), "Operational")
}

func (s ConnectionState) IsPrinting() bool {
	return strings.HasPrefix(string(s), "Printing") ||
		strings.HasPrefix(string(s), "Sending") ||
		strings.HasPrefix(string(s), "Paused") ||
		strings.HasPrefix(string(s), "Transfering") ||
		strings.HasPrefix(string(s), "Paused")
}

func (s ConnectionState) IsOffline() bool {
	return strings.HasPrefix(string(s), "Offline") ||
		strings.HasPrefix(string(s), "Closed")
}

func (s ConnectionState) IsError() bool {
	return strings.HasPrefix(string(s), "Error") ||
		strings.HasPrefix(string(s), "Unknown")
}

func (s ConnectionState) IsConnecting() bool {
	return strings.HasPrefix(string(s), "Opening") ||
		strings.HasPrefix(string(s), "Detecting") ||
		strings.HasPrefix(string(s), "Connecting") ||
		strings.HasPrefix(string(s), "Detecting")
}

// ConnectionResponse is the response from a connection command.
type ConnectionResponse struct {
	Current struct {
		// State current state of the connection.
		State ConnectionState `json:"state"`
		// Port to connect to.
		Port string `json:"port"`
		// BaudRate speed of the connection.
		BaudRate int `json:"baudrate"`
		// PrinterProfile profile to use for connection.
		PrinterProfile string `json:"printerProfile"`
	}
	Options struct {
		// Ports list of available ports.
		Ports []string `json:"ports"`
		// BaudRates list of available speeds.
		BaudRates []int `json:"baudrates"`
		// PrinterProfile list of available profiles.
		PrinterProfiles []*Profile `json:"printerProfiles"`
		// PortPreference default port.
		PortPreference string `json:"portPreference"`
		// BaudRatePreference default speed.
		BaudRatePreference int `json:"baudratePreference"`
		// PrinterProfilePreference default profile.
		PrinterProfilePreference string `json:"printerProfilePreference"`
		// Autoconnect whether to automatically connect to the printer on
		// OctoPrint’s startup in the future.
		Autoconnect bool `json:"autoconnect"`
	}
}

// Profile describe a printer profile.
type Profile struct {
	// ID is the identifier of the profile.
	ID string `json:"id"`
	// Name is the display name of the profile.
	Name string `json:"name"`
}

// FilesResponse is the response to a FilesRequest.
type FilesResponse struct {
	// Files is the list of requested files. Might be an empty list if no files
	// are available
	Files []*FileInformation
	// Free is the amount of disk space in bytes available in the local disk
	// space (refers to OctoPrint’s `uploads` folder). Only returned if file
	// list was requested for origin `local` or all origins.
	Free int
}

// FileInformation contains information regarding a file.
type FileInformation struct {
	// Name is name of the file without path. E.g. “file.gco” for a file
	// “file.gco” located anywhere in the file system.
	Name string `json:"name"`
	// Path is the path to the file within the location. E.g.
	//“folder/subfolder/file.gco” for a file “file.gco” located within “folder”
	// and “subfolder” relative to the root of the location.
	Path string `json:"path"`
	// Type of file. model or machinecode. Or folder if it’s a folder, in
	// which case the children node will be populated.
	Type string `json:"type"`
	// TypePath path to type of file in extension tree. E.g. `["model", "stl"]`
	// for .stl files, or `["machinecode", "gcode"]` for .gcode files.
	// `["folder"]` for folders.
	TypePath []string `json:"typePath"`
	// Hash is the MD5 hash of the file. Only available for `local` files.
	Hash string `json:"hash"`
	// Size of the file in bytes. Only available for `local` files or `sdcard`
	// files if the printer supports file sizes for sd card files.
	Size int `json:"size"`
	// Date when this file was uploaded. Only available for `local` files.
	Date JSONTime `json:"date"`
	// Origin of the file, `local` when stored in OctoPrint’s `uploads` folder,
	// `sdcard` when stored on the printer’s SD card (if available)
	Origin string `json:"origin"`
	// Refs references relevant to this file, left out in abridged versio
	Refs Reference `json:"refs"`
	// GCodeAnalysis information from the analysis of the GCODE file, if
	// available. Left out in abridged version.
	GCodeAnalysis GCodeAnalysisInformation `json:"gcodeAnalysis"`
	// Print information from the print stats of a file.
	Print PrintStats `json:"print"`
}

// IsFolder it returns true if the file is a folder.
func (f *FileInformation) IsFolder() bool {
	if len(f.TypePath) == 1 && f.TypePath[0] == "folder" {
		return true
	}

	return false
}

// Reference of a file.
type Reference struct {
	// Resource that represents the file or folder (e.g. for issuing commands
	// to or for deleting)
	Resource string `json:"resource"`
	// Download URL for the file. Never present for folders.
	Download string `json:"download"`
	// Model from which this file was generated (e.g. an STL, currently not
	// used). Never present for folders.
	Model string `json:"model"`
}

// GCodeAnalysisInformation Information from the analysis of the GCODE file.
type GCodeAnalysisInformation struct {
	// EstimatedPrintTime is the estimated print time of the file, in seconds.
	EstimatedPrintTime float64 `json:"estimatedPrintTime"`
	// Filament estimated usage of filament
	Filament struct {
		// Length estimated of filament used, in mm
		Length int `json:"length"`
		// Volume estimated of filament used, in cm³
		Volume float64 `json:"volume"`
	} `json:"filament"`
}

// PrintStats information from the print stats of a file.
type PrintStats struct {
	// Failure number of failed prints.
	Failure int `json:"failure"`
	// Success number of success prints.
	Success int `json:"success"`
	// Last print information.
	Last struct {
		// Date of the last print.
		Date JSONTime `json:"date"`
		// Success or not.
		Success bool `json:"success"`
	} `json:"last"`
}

// UploadFileResponse is the response to a UploadFileRequest.
type UploadFileResponse struct {
	// Abridged information regarding the file that was just uploaded. If only
	// uploaded to local this will only contain the local property. If uploaded
	// to SD card, this will contain both local and sdcard properties. Only
	// contained if a file was uploaded, not present if only a new folder was
	// created.
	File struct {
		// Local is the information regarding the file that was just uploaded
		// to the local storage.
		Local *FileInformation `json:"local"`
		// SDCard is the information regarding the file that was just uploaded
		// to the printer’s SD card.
		SDCard *FileInformation `json:"sdcard"`
	} `json:"files"`
	// Done whether any file processing after upload has already finished or
	// not, e.g. due to first needing to perform a slicing step. Clients may
	// use this information to direct progress displays related to the upload.
	Done bool `json:"done"`
}

type JSONTime struct{ time.Time }

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t.Time).Unix(), 10)), nil
}

func (t *JSONTime) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)
	if r == "null" {
		return nil
	}

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}

	t.Time = time.Unix(q, 0)
	return
}
