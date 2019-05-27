package octoprint

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

var Version = "0.1"

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
	FilePosition uint64 `json:"filepos"`
}

// ProgressInformation contains information regarding the progress of the
// current print job.
type ProgressInformation struct {
	// Completion percentage of completion of the current print job.
	Completion float64 `json:"completion"`
	// FilePosition current position in the file being printed, in bytes
	// from the beginning.
	FilePosition uint64 `json:"filepos"`
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
	Files    []*FileInformation
	Children []*FileInformation
	// Free is the amount of disk space in bytes available in the local disk
	// space (refers to OctoPrint’s `uploads` folder). Only returned if file
	// list was requested for origin `local` or all origins.
	Free uint64
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
	Size uint64 `json:"size"`
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
		Length uint32 `json:"length"`
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

// SystemCommandsResponse is the response to a SystemCommandsRequest.
type SystemCommandsResponse struct {
	Core   []*CommandDefinition `json:"core"`
	Custom []*CommandDefinition `json:"custom"`
}

// CommandSource is the source of the command definition.
type CommandSource string

const (
	// Core for system actions defined by OctoPrint itself.
	Core CommandSource = "core"
	// Custom for custom system commands defined by the user through `config.yaml`.
	Custom CommandSource = "custom"
)

// CommandDefinition describe a system command.
type CommandDefinition struct {
	// Name of the command to display in the System menu.
	Name string `json:"name"`
	// Command is the full command line to execute for the command.
	Command string `json:"command"`
	// Action is an identifier to refer to the command programmatically. The
	// special `action` string divider signifies a `divider` in the menu.
	Action string `json:"action"`
	// Confirm if present and set, this text will be displayed to the user in a
	// confirmation dialog they have to acknowledge in order to really execute
	// the command.
	RawConfirm json.RawMessage `json:"confirm"`
	Confirm    string          `json:"-"`
	// Async whether to execute the command asynchronously or wait for its
	// result before responding to the HTTP execution request.
	Async bool `json:"async"`
	// Ignore whether to ignore the return code of the command’s execution.
	Ignore bool `json:"ignore"`
	// Source of the command definition.
	Source CommandSource `json:"source"`
	// Resource is the URL of the command to use for executing it.
	Resource string `json:"resource"`
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

// CustomCommandsResponse is the response to a CustomCommandsRequest.
type CustomCommandsResponse struct {
	Controls []*ControlContainer `json:"controls"`
}

// ControlContainer describes a control container.
type ControlContainer struct {
	// Name to display above the container, basically a section header.
	Name string `json:"name"`
	// Children a list of children controls or containers contained within this
	// container.
	Children []*ControlDefinition `json:"children"`
	// Layout  to use for laying out the contained children, either from top to
	// bottom (`vertical`) or from left to right (`horizontal``). Defaults to a
	// vertical layout.
	Layout string `json:"layout"`
}

// ControlDefinition describe a system control.
type ControlDefinition struct {
	// Name of the control, will be displayed either on the button if it’s a
	// control sending a command or as a label for controls which only display
	// output.
	Name string `json:"name"`
	// Command a single GCODE command to send to the printer. Will be rendered
	// as a button which sends the command to the printer upon click. The button
	// text will be the value of the `name` attribute. Mutually exclusive with
	// `commands` and `script`. The rendered button be disabled if the printer
	// is currently offline or printing or alternatively if the requirements
	// defined via the `enabled` attribute are not met.
	Command string `json:"command"`
	// Command a list of GCODE commands to send to the printer. Will be rendered
	// as a button which sends the commands to the printer upon click. The
	// button text will be the value of the `name` attribute. Mutually exclusive
	// with `command` and `script`. The rendered button will be disabled if the
	// printer is currently offline or printing or alternatively if the
	// requirements defined via the `enabled` attribute are not met.
	Commands []string `json:"commands"`
	// Script is the name of a full blown GCODE script to send to the printer.
	// Will be rendered as a button which sends the script to the printer upon
	// click. The button text will be the value of the name attribute. Mutually
	// exclusive with `command` and `commands`. The rendered button will be
	// disabled if the printer is currently offline or printing or alternatively
	// if the requirements defined via the `enabled`` attribute are not met.
	//
	// Values of input parameters will be available in the template context
	// under the `parameter` variable (e.g. an input parameter speed will be
	// available in the script template as parameter.speed). On top of that all
	// other variables defined in the GCODE template context will be available.
	Script string `json:"script"`
	// JavaScript snippet to be executed when the button rendered for `command`
	// or `commands` is clicked. This allows to override the direct sending of
	// the command or commands to the printer with more sophisticated behaviour.
	// The JavaScript snippet is `eval`’d and processed in a context where the
	// control it is part of is provided as local variable `data` and the
	// `ControlViewModel` is available as self.
	JavasScript string `json:"javascript"`
	// Enabled a JavaScript snippet to be executed when the button rendered for
	// `command` or `commands` is clicked. This allows to override the direct
	// sending of the command or commands to the printer with more sophisticated
	// behaviour. The JavaScript snippet is `eval`’d and processed in a context
	// where the control it is part of is provided as local variable `data` and
	// the `ControlViewModel` is available as `self`.
	Enabled bool `json:"enabled"`
	// Input a list of definitions of input parameters for a command or
	// commands, to be rendered as additional input fields.
	Input *ControlInput `json:"input"`
	// Regex a regular expression to match against lines received from the
	// printer to retrieve information from it (e.g. specific output). Together
	// with template this allows rendition of received data from the printer
	// within the UI.
	Regex string `json:"regex"`
	// Template to use for rendering the match of `regex`. May contain
	// placeholders in Python Format String Syntax[1] for either named groups
	// within the regex (e.g. `Temperature: {temperature}` for a regex
	// `T:\s*(?P<temperature>\d+(\.\d*)`) or positional groups within the regex
	// (e.g. `Position: X={0}, Y={1}, Z={2}, E={3}` for a regex
	// `X:([0-9.]+) Y:([0-9.]+) Z:([0-9.]+) E:([0-9.]+)`).
	// https://docs.python.org/2/library/string.html#format-string-syntax
	Template string `json:"template"`
	// Confirm a text to display to the user to confirm his button press. Can
	// be used with sensitive custom controls like changing EEPROM values in
	// order to prevent accidental clicks.
	Confirm string `json:"confirm"`
}

// ControlInput a list of definitions of input parameters for a command or
// commands, to be rendered as additional input fields.
type ControlInput struct {
	// Name to display for the input field.
	Name string `json:"name"`
	// Parameter name for the input field, used as a placeholder in
	// command/commands.
	Parameter string `json:"parameter"`
	// Default value for the input field.
	Default interface{} `json:"default"`
	// Slider if defined instead of an input field a slider control will be
	// rendered.
	Slider struct {
		// Minimum value of the slider, defaults to 0.
		Min int `json:"min"`
		// Maximum value of the slider, defaults to 255.
		Maximum int `json:"max"`
		// Step size per slider “tick”, defaults to 1.
		Step int `json:"step"`
	} `json:"slider"`
}

// Settings are the current configuration of OctoPrint.
type Settings struct {
	// API REST API settings.
	API *APIConfig `json:"api"`
	// Features settings to enable or disable OctoPrint features.
	Feature *FeaturesConfig `json:"feature"`
	//Folder settings to set custom paths for folders used by OctoPrint.
	Folder *FolderConfig `json:"folder"`
	// Serial settings to configure the serial connection to the printer.
	Serial *SerialConfig `json:"serial"`
	// Server settings to configure the server.
	Server *ServerConfig `json:"server"`
	// Temperature profiles which will be displayed in the temperature tab.
	Temperature *TemperatureConfig `json:"temperature"`
	// TerminalFilters to display in the terminal tab for filtering certain
	// lines from the display terminal log.
	TerminalFilters []*TerminalFilter `json:"terminalFilters"`
	// Webcam settings to configure webcam support.
	Webcam *WebcamConfig `json:"json"`

	// Un-handled values
	Appearance interface{} `json:"appearance"`
	Plugins    interface{} `json:"plugins"`
	Printer    interface{} `json:"printer"`
}

// APIConfig REST API settings.
type APIConfig struct {
	// Enabled whether to enable the API.
	Enabled bool `json:"enabled"`
	// Key current API key needed for accessing the API
	Key string `json:"key"`
}

// FeaturesConfig settings to enable or disable OctoPrint features.
type FeaturesConfig struct {
	// SizeThreshold maximum size a GCODE file may have to automatically be
	// loaded into the viewer, defaults to 20MB. Maps to
	// gcodeViewer.sizeThreshold in config.yaml.
	SizeThreshold uint64
	// MobileSizeThreshold maximum size a GCODE file may have on mobile devices
	// to automatically be loaded into the viewer, defaults to 2MB. Maps to
	// gcodeViewer.mobileSizeThreshold in config.yaml.
	MobileSizeThreshold uint64 `json:"mobileSizeThreshold"`
	// TemperatureGraph whether to enable the temperature graph in the UI or not.
	TemperatureGraph bool `json:"temperatureGraph"`
	// WaitForStart specifies whether OctoPrint should wait for the start
	// response from the printer before trying to send commands during connect.
	WaitForStart bool `json:"waitForStart"`
	// AlwaysSendChecksum specifies whether OctoPrint should send linenumber +
	// checksum with every printer command. Needed for successful communication
	// with Repetier firmware.
	AlwaysSendChecksum bool `json:"alwaysSendChecksum"`
	NeverSendChecksum  bool `json:"neverSendChecksum"`
	// SDSupport specifies whether support for SD printing and file management
	// should be enabled
	SDSupport bool `json:"sdSupport"`
	// SDAlwaysAvailable whether to always assume that an SD card is present in
	// the printer. Needed by some firmwares which don't report the SD card
	// status properly.
	SDAlwaysAvailable bool `json:"sdAlwaysAvailable"`
	// SDReleativePath Specifies whether firmware expects relative paths for
	// selecting SD files.
	SDRelativePath bool `json:"sdRelativePath"`
	// SwallowOkAfterResend whether to ignore the first ok after a resend
	// response. Needed for successful communication with Repetier firmware.
	SwallowOkAfterResend bool `json:"swallowOkAfterResend"`
	// RepetierTargetTemp whether the printer sends repetier style target
	// temperatures in the format `TargetExtr0:<temperature>` instead of
	// attaching that information to the regular M105 responses.
	RepetierTargetTemp bool `json:"repetierTargetTemp"`
	// ExternalHeatupDetection whether to enable external heatup detection (to
	// detect heatup triggered e.g. through the printer's LCD panel or while
	// printing from SD) or not. Causes issues with Repetier's "first ok then
	// response" approach to communication, so disable for printers running
	// Repetier firmware.
	ExternalHeatupDetection bool `json:"externalHeatupDetection"`
	// KeyboardControl whether to enable the keyboard control feature in the
	// control tab.
	KeyboardControl bool `json:"keyboardControl"`
	// PollWatched whether to actively poll the watched folder (true) or to rely
	// on the OS's file system notifications instead (false).
	PollWatched bool `json:"pollWatched"`
	// IgnoreIdenticalResends whether to ignore identical resends from the
	// printer (true, repetier) or not (false).
	IgnoreIdenticalResends bool `json:"ignoreIdenticalResends"`
	// ModelSizeDetection whether to enable model size detection and warning
	// (true) or not (false)
	ModelSizeDetection bool `json:"modelSizeDetection"`
	// FirmwareDetection whether to attempt to auto detect the firmware of the
	// printer and adjust settings  accordingly (true) or not and rely on manual
	// configuration (false)
	FirmwareDetection bool `json:"firmwareDetection"`
	// PrintCancelConfirmation whether to show a confirmation on print
	// cancelling (true) or not (false).
	PrintCancelConfirmation bool `json:"printCancelConfirmation"`
	// BlockWhileDwelling whether to block all sending to the printer while a G4
	// (dwell) command is active (true, repetier) or not (false).
	BlockWhileDwelling bool `json:"blockWhileDwelling"`
}

// FolderConfig settings to set custom paths for folders used by OctoPrint.
type FolderConfig struct {
	// Uploads absolute path where to store gcode uploads. Defaults to the
	// uploads folder in the OctoPrint settings folder.
	Uploads string `json:"uploads"`
	// Timelapse absolute path where to store finished timelapse recordings.
	// Defaults to the timelapse folder in the OctoPrint settings dir.
	Timelapse string `json:"timelapse"`
	// TimelapseTmp absolute path where to store temporary timelapse files.
	// Defaults to the timelapse/tmp folder in the OctoPrint settings dir Maps
	// to folder.timelapse_tmp in config.yaml.
	TimelapseTmp string `json:"timelapseTmp"`
	// Logs absolute path where to store log files. Defaults to the logs folder
	// in the OctoPrint settings dir
	Logs string `json:"logs"`
	// Watched absolute path to a folder being watched for new files which then
	// get automatically added to OctoPrint (and deleted from that folder).
	// Can e.g. be used to define a folder which can then be mounted from remote
	// machines and used as local folder for quickly adding downloaded and/or
	// sliced objects to print in the future.
	Watched string `json:"watched"`
}

// SerialConfig settings to configure the serial connection to the printer.
type SerialConfig struct {
	// Port is the default serial port, defaults to unset (= AUTO)
	Port string `json:"port"`
	// Baudrate is the default baudrate, defaults to unset (= AUTO)
	Baudrate int `json:"baudrate"`
	// 	Available serial ports
	PortOptions []string `json:"portOptions"`
	//	Available serial baudrates
	BaudrateOptions []int `json:"baudrateOptions"`
	// Autoconnect whether to automatically connect to the printer on server
	// startup (if available).
	Autoconnect bool `json:"autoconnect"`
	// TimeoutConnection for waiting to establish a connection with the selected
	// port, in seconds. Defaults to 2 sec. Maps to serial.timeout.connection in
	// config.yaml
	TimeoutConnection float64 `json:"timeoutConnection"`
	// TimeoutDetection for waiting for a response from the currently tested
	// port during autodetect, in seconds. Defaults to 0.5 sec Maps to
	// serial.timeout.detection in config.yaml
	TimeoutDetection float64 `json:"timeoutDetection"`
	// TimeoutCommunication during serial communication, in seconds. Defaults to
	// 30 sec. Maps to serial.timeout.communication in config.yaml
	TimeoutCommunication float64 `json:"timeoutCommunication"`
	// TimeoutTemperature after which to query temperature when no target is
	// set. Maps to serial.timeout.temperature in config.yaml
	TimeoutTemperature float64 `json:"timeoutTemperature"`
	// TimeoutTemperatureTargetSet after which to query temperature when a
	// target is set. Maps to serial.timeout.temperatureTargetSet in config.yaml
	TimeoutTemperatureTargetSet float64 `json:"timeoutTemperatureTargetSet"`
	// TimeoutSDStatus after which to query the SD status while SD printing.
	// Maps to serial.timeout.sdStatus in config.yaml
	TimeoutSDStatus float64 `json:"timeoutSdStatus"`
	// Log whether to log whole communication to serial.log (warning: might
	// decrease performance)
	Log bool `json:"log"`
	// AdditionalPorts use this to define additional patterns to consider for
	// serial port listing. Must be a valid "glob" pattern (see
	// http://docs.python.org/2/library/glob.html). Defaults to not set.
	AdditionalPorts []string `json:"additionalPorts"`
	// AdditionalBaudrates use this to define additional baud rates to offer for
	// connecting to serial ports. Must be a valid integer. Defaults to not set
	AdditionalBaudrates []int `json:"additionalBaudrates"`
	// LongRunningCommands which are known to take a long time to be
	// acknowledged by the firmware. E.g. homing, dwelling, auto leveling etc.
	LongRunningCommands []string `json:"longRunningCommands"`
	// ChecksumRequiringCommands which need to always be send with a checksum.
	// Defaults to only M110
	ChecksumRequiringCommands []string `json:"checksumRequiringCommands"`
	// HelloCommand to send in order to initiate a handshake with the printer.
	// Defaults to "M110 N0" which simply resets the line numbers in the
	// firmware and which should be acknowledged with a simple "ok".
	HelloCommand string `json:"helloCommand"`
	// IgnoreErrorsFromFirmware whether to completely ignore errors from the
	// firmware or not
	IgnoreErrorsFromFirmware bool `json:"ignoreErrorsFromFirmware"`
	// DisconnectOnErrors whether to disconnect on errors or not.
	DisconnectOnErrors bool `json:"disconnectOnErrors"`
	// TriggerOkForM29 whether to "manually" trigger an ok for M29 (a lot of
	// versions of this command are buggy and the responds skips on the ok)
	TriggerOkForM29 bool `json:"triggerOkForM29"`
	// SupportResendsWIthoutOk whether to support resends without follow-up ok
	// or not.
	SupportResendsWIthoutOk string `json:"supportResendsWIthoutOk"`
	// Maps to serial.maxCommunicationTimeouts.idle in config.yaml
	MaxTimeoutsIdle float64 `json:"maxTimeoutsIdle"`
	// MaxTimeoutsPrinting maximum number of consecutive communication timeouts
	// after which the printer will be considered dead and OctoPrint disconnects
	// with an error. Maps to serial.maxCommunicationTimeouts.printing in
	// config.yaml
	MaxTimeoutsPrinting float64 `json:"maxTimeoutsPrinting"`
	// MaxTimeoutsPrinting maximum number of consecutive communication timeouts
	// after which the printer will be considered dead and OctoPrint disconnects
	// with an error. Maps to serial.maxCommunicationTimeouts.log in config.yaml
	MaxTimeoutsLong float64 `json:"maxTimeoutsLong"`
}

// ServerConfig settings to configure the server.
type ServerConfig struct {
	// Commands to restart/shutdown octoprint or the system it's running on.
	Commands struct {
		// ServerRestartCommand to restart OctoPrint, defaults to being unset
		ServerRestartCommand string `json:"serverRestartCommand"`
		//SystemRestartCommand  to restart the system OctoPrint is running on,
		// defaults to being unset
		SystemRestartCommand string `json:"systemRestartCommand"`
		// SystemShutdownCommand Command to shut down the system OctoPrint is
		// running on, defaults to being unset
		SystemShutdownCommand string `json:"systemShutdownCommand"`
	} `json:"commands"`
	// Diskspace settings of when to display what disk space warning
	Diskspace struct {
		// Warning threshold (bytes) after which to consider disk space becoming
		// sparse, defaults to 500MB.
		Warning uint64 `json:"warning"`
		// Critical threshold (bytes) after which to consider disk space becoming
		// critical, defaults to 200MB.
		Critical uint64 `json:"critical"`
	} `json:"diskspace"`
	// OnlineCheck configuration of the regular online connectivity check.
	OnlineCheck struct {
		// Enabled whether the online check is enabled, defaults to false due to
		// valid privacy concerns.
		Enabled bool `json:"enabled"`
		// Interval in which to check for online connectivity (in seconds)
		Interval int `json:"interval"`
		// Host DNS host against which to check (default: 8.8.8.8 aka Google's DNS)
		Host string `json:"host"`
		// DNS port against which to check (default: 53 - the default DNS port)
		Port int `json:"port"`
	} `json:"onlineCheck"`
	// PluginBlacklist configuration of the plugin blacklist
	PluginBlacklist struct {
		// Enabled whether use of the blacklist is enabled, defaults to false
		Enabled bool `json:"enabled"`
		/// URL from which to fetch the blacklist
		URL string `json:"url"`
		// TTL is time to live of the cached blacklist, in secs (default: 15mins)
		TTL int `json:"ttl"`
	} `json:"pluginBlacklist"`
}

// TemperatureConfig temperature profiles which will be displayed in the
// temperature tab.
type TemperatureConfig struct {
	// Graph cutoff in minutes.
	Cutoff int `json:"cutoff"`
	// Profiles  which will be displayed in the temperature tab.
	Profiles []*TemperatureProfile `json:"profiles"`
	// SendAutomatically enable this to have temperature fine adjustments you
	// do via the + or - button be sent to the printer automatically.
	SendAutomatically bool `json:"sendAutomatically"`
	// SendAutomaticallyAfter OctoPrint will use this delay to limit the number
	// of sent temperature commands should you perform multiple fine adjustments
	// in a short time.
	SendAutomaticallyAfter float64 `json:"sendAutomaticallyAfter"`
}

// TerminalFilter to display in the terminal tab for filtering certain lines
// from the display terminal log.
type TerminalFilter struct {
	Name  string `json:"name"`
	RegEx string `json:"regex"`
}

// WebcamConfig settings to configure webcam support.
type WebcamConfig struct {
	// StreamUrl use this option to enable display of a webcam stream in the
	// UI, e.g. via MJPG-Streamer. Webcam support will be disabled if not
	// set. Maps to webcam.stream in config.yaml.
	StreamURL string `json:"streamUrl"`
	// SnapshotURL use this option to enable timelapse support via snapshot,
	// e.g. via MJPG-Streamer. Timelapse support will be disabled if not set.
	// Maps to webcam.snapshot in config.yaml.
	SnapshotURL string `json:"snapshotUrl"`
	// FFmpegPath path to ffmpeg binary to use for creating timelapse
	// recordings. Timelapse support will be disabled if not set. Maps to
	// webcam.ffmpeg in config.yaml.
	FFmpegPath string `json:"ffmpegPath"`
	// Bitrate to use for rendering the timelapse video. This gets directly
	// passed to ffmpeg.
	Bitrate int `json:"bitrate"`
	// FFmpegThreads number of how many threads to instruct ffmpeg to use for
	// encoding. Defaults to 1. Should be left at 1 for RPi1.
	FFmpegThreads int `json:"ffmpegThreads"`
	// Watermark whether to include a "created with OctoPrint" watermark in the
	// generated timelapse movies.
	Watermark string `json:"watermark"`
	// FlipH whether to flip the webcam horizontally.
	FlipH bool `json:"flipH"`
	// FlipV whether to flip the webcam vertically.
	FlipV bool `json:"flipV"`
	// Rotate90 whether to rotate the webcam 90° counter clockwise.
	Rotate90 bool `json:"rotate90"`
}

// TemperatureProfile describes the temperature profile preset for a given
// material.
type TemperatureProfile struct {
	Name     string  `json:"name"`
	Bed      float64 `json:"bed"`
	Extruder float64 `json:"extruder"`
}
