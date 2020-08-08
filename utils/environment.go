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

// RequiredEnvironmentVariablesAreSet - (Captain Obvious says...) verifies that all required environment variables are set.
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

// NameOfMissingRequiredEnvironmentVariable - (Captain Obvious says...) returns the name of the missing required environment variable.
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

// DumpEnvironmentVariables - (Captain Obvious says...) dumps (logs) all environment variables.
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
