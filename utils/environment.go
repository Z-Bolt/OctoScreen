package utils

import (
	"os"
	//"strings"
)

// Required environment variables
const (
	EnvStylePath   = "OCTOSCREEN_STYLE_PATH"
	EnvBaseURL     = "OCTOPRINT_HOST"
	EnvAPIKey      = "OCTOPRINT_APIKEY"
)

// Optional (but good to have) environment variables
const (
	EnvLogLevel    = "OCTOSCREEN_LOG_LEVEL"
	EnvLogFilePath = "OCTOSCREEN_LOG_FILE_PATH"
	EnvResolution  = "OCTOSCREEN_RESOLUTION"
	EnvConfigFile  = "OCTOPRINT_CONFIG_FILE"
)

func RequiredEnvironmentVariablesAreSet() bool {
	if( !environmentVariableIsSet(EnvStylePath) ) {
		return false
	}

	if( !environmentVariableIsSet(EnvBaseURL) ) {
		return false
	}

	if( !environmentVariableIsSet(EnvAPIKey) ) {
		return false
	}

	return true
}

func environmentVariableIsSet(environmentVariable string) bool {
	return os.Getenv(environmentVariable) != ""
}

func NameOfMissingRequiredEnvironmentVariable() string {
	if( !environmentVariableIsSet(EnvStylePath) ) {
		return EnvStylePath
	}

	if( !environmentVariableIsSet(EnvBaseURL) ) {
		return EnvBaseURL
	}

	if( !environmentVariableIsSet(EnvAPIKey) ) {
		return EnvAPIKey
	}

	return "UNKNOWN"
}

func DumpEnvironmentVariables() {
	Logger.Info("")
	Logger.Info("")
	Logger.Info("Environment variables...")

	// Required environment variables
	Logger.Infof("Required environment variables:")
	dumpEnvironmentVariable(EnvBaseURL)
	dumpEnvironmentVariable(EnvAPIKey)
	dumpEnvironmentVariable(EnvStylePath)

	// Optional environment variables
	Logger.Info("")
	Logger.Infof("Optional environment variables:")
	dumpEnvironmentVariable(EnvConfigFile)
	dumpEnvironmentVariable(EnvLogFilePath)
	dumpEnvironmentVariable(EnvLogLevel)
	dumpEnvironmentVariable(EnvResolution)
	// EnvResolution is optional.  If not set, the window size will
	// default to the values defined in globalVars.go.

	Logger.Info("")
	Logger.Info("")
}

func dumpEnvironmentVariable(key string) {
	value := os.Getenv(key)
	if value == "" {
		value = ">>MISSING<<"
	}

	Logger.Infof("key: %q, value: %q", key, value)
}
