package octoprint

import (
	"encoding/json"
	"time"
)

type Axis string

const (
	XAxis Axis = "x"
	YAxis Axis = "y"
	ZAxis Axis = "z"
)

// JobInformation contains information regarding the target of the current job.
type JobInformation struct {
	// File is the file that is the target of the current print job.
	File FileInformation `json:"file"`
	// EstimatedPrintTime is the estimated print time for the file, in seconds.
	EstimatedPrintTime int `json:"estimatedPrintTime"`
	// LastPrintTime is the print time of the last print of the file, in seconds.
	LastPrintTime int `json:"lastPrintTime"`
	// Filament contains Information regarding the estimated filament
	// usage of the print job.
	Filament struct {
		// Length of filament used, in mm
		Length int `json:"length"`
		// Volume of filament used, in cm³
		Volume float64 `json:"volume"`
	} `json:"filament"`
	FilePosition int `json:"filepos"`
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
	TypePath string `json:"typePath"`
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
	PrintTime int `json:"printTime"`
	// PrintTimeLeft is estimate of time left to print, in seconds
	PrintTimeLeft int `json:"printTimeLeft"`
}

// CurrentState of a tool.
type CurrentState struct {
	State
	Offset float64 `json:"offset"`
}

// State of a tool.
type State struct {
	Actual float64 `json:"actual"`
	Target float64 `json:"target"`
}

// PrinterState current state of the printer.
type PrinterState struct {
	Text  string `json:"text"`
	Flags struct {
		Operations    bool `json:"operational"`
		Puased        bool `json:"paused"`
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

// History of a tool.
type History history
type history struct {
	Time  time.Time        `json:"time"`
	Tools map[string]State `json:"tools"`
}

func (h *History) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	ts := raw["time"]
	delete(raw, "time")
	b, _ = json.Marshal(map[string]interface{}{
		"time":  time.Unix(int64(ts.(float64)), 0),
		"tools": raw,
	})

	i := &history{}
	if err := json.Unmarshal(b, i); err != nil {
		return err
	}

	*h = History(*i)
	return nil
}
