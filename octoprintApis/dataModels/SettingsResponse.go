package dataModels


// Settings are the current configuration of OctoPrint.
type SettingsResponse struct {
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
