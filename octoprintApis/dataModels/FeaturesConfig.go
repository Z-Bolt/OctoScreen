package dataModels


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
	SdSupport bool `json:"sdSupport"`

	// SDAlwaysAvailable whether to always assume that an SD card is present in
	// the printer. Needed by some firmwares which don't report the SD card
	// status properly.
	SdAlwaysAvailable bool `json:"sdAlwaysAvailable"`

	// SDReleativePath Specifies whether firmware expects relative paths for
	// selecting SD files.
	SdRelativePath bool `json:"sdRelativePath"`

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
