package utils

import (
	"log"
	"os"
)

// Required environment variables
const (
	EnvStylePath  = "OCTOSCREEN_STYLE_PATH"
	// EnvResolution = "OCTOSCREEN_RESOLUTION"
	EnvBaseURL    = "OCTOPRINT_HOST"
	EnvAPIKey     = "OCTOPRINT_APIKEY"
	// EnvConfigFile = "OCTOPRINT_CONFIG_FILE"
)

// Optional (but good to have) environment variables
const (
	EnvResolution = "OCTOSCREEN_RESOLUTION"
	EnvConfigFile = "OCTOPRINT_CONFIG_FILE"
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
		// OCTOPRINT_HOST must be in the form of "http://octopi.local" or "http://1.2.3.4"
		return EnvBaseURL
	}

	if( !environmentVariableIsSet(EnvAPIKey) ) {
		return EnvAPIKey
	}

	return "UNKNOWN"
}

func DumpEnvironmentVariables() {
	dumpEnvironmentVariable(EnvStylePath)
	dumpEnvironmentVariable(EnvResolution)
	dumpEnvironmentVariable(EnvBaseURL)
	dumpEnvironmentVariable(EnvAPIKey)
	dumpEnvironmentVariable(EnvConfigFile)
}

func dumpEnvironmentVariable(key string) {
	value := os.Getenv(key)
	if value == "" {
		value = ">>MISSING<<"
	}

	log.Println("key: " + key + ", value: " + value)
}
