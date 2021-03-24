package dataModels


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
	TimeoutSdStatus float64 `json:"timeoutSdStatus"`

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

	// SupportResendsWithoutOk whether to support resends without follow-up ok
	// or not.
	SupportResendsWithoutOk string `json:"supportResendsWithoutOk"`

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
