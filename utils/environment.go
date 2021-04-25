package utils

import (
	"fmt"
	"os"
	//"strings"


	"github.com/Z-Bolt/OctoScreen/logger"
)


// OctoScreenVersion is set during compilation.
var OctoScreenVersion = "2.7.3"

const MISSING_ENV_TOKEN = ">>MISSING<<"
const INVALID_ENV_TOKEN = "!!!INVALID!!!"

// Required environment variables
const (
	EnvStylePath   = "OCTOSCREEN_STYLE_PATH"
	EnvBaseURL     = "OCTOPRINT_HOST"
	EnvAPIKey      = "OCTOPRINT_APIKEY"
)

// Optional (but good to have) environment variables
const (
	EnvLogLevel       = "OCTOSCREEN_LOG_LEVEL"
	EnvLogFilePath    = "OCTOSCREEN_LOG_FILE_PATH"
	EnvResolution     = "OCTOSCREEN_RESOLUTION"
	EnvConfigFile     = "OCTOPRINT_CONFIG_FILE"
	EnvDisplayCursor  = "DISPLAY_CURSOR"
)

func RequiredEnvironmentVariablesAreSet(apiKey string) bool {
	if( !environmentVariableIsSet(EnvStylePath) ) {
		return false
	}

	if( !environmentVariableIsSet(EnvBaseURL) ) {
		return false
	}

	// APIKey/OCTOPRINT_APIKEY can be set in either OctoScreen's config file,
	// or in OctoPrint's config file.  In main.init(), APIKey is initialized to whatever
	// it can find first.
	//
	// APIKey is global to the "main" namespace, but the "utils" namespace is a child,
	// and due to GoLang's rules, /main/utils doesn't have access to globals in /main,
	// so APIKey has to be passed into RequiredEnvironmentVariablesAreSet().
	//
	// if( !environmentVariableIsSet(EnvAPIKey) ) {
	// 	return false
	// }
	if apiKey == "" {
		return false
	}

	return true
}

func environmentVariableIsSet(environmentVariable string) bool {
	return os.Getenv(environmentVariable) != ""
}

func NameOfMissingRequiredEnvironmentVariable(apiKey string) string {
	if( !environmentVariableIsSet(EnvStylePath) ) {
		return EnvStylePath
	}

	if( !environmentVariableIsSet(EnvBaseURL) ) {
		return EnvBaseURL
	}

	// Similar comment as to the one that's in RequiredEnvironmentVariablesAreSet()...
	// Since the runtime value of APIKey is set in main.init(), and can be set by either
	// being defined in OctoScreen's config file or in OctoPrint's config file,
	// the value needs to be passed into NameOfMissingRequiredEnvironmentVariable().
	// if( !environmentVariableIsSet(EnvAPIKey) ) {
	// 	return EnvAPIKey
	// }
	if apiKey == "" {
		return EnvAPIKey
	}

	return "UNKNOWN"
}

func DumpSystemInformation() {
	logger.Info("System Information...")
	logger.Infof("OctoScreen version: %q", OctoScreenVersion)
	// More system stats to come...
	logger.Info("")
}

func DumpEnvironmentVariables() {
	logger.Info("Environment variables...")

	// Required environment variables
	logger.Info("Required environment variables:")
	dumpEnvironmentVariable(EnvBaseURL)

	// TODO: revisit this!
	// 1. remove OCTOPRINT_APIKEY from option settings
	// 2. make the octoprint config path required
	// 3. update code... use only one path to get the api key octoprint)
	// 4. update code... make octoprint config path required
	// 5. update code... read api key from octoprint config
	// 6. dump api key (obfuscated though)
	// 7. update docs
	// 8. make sure what's dumped to the log is correct, for both when present and when missing.
	dumpObfuscatedEnvironmentVariable(EnvAPIKey)

	dumpEnvironmentVariable(EnvStylePath)
	logger.Info("")

	// Optional environment variables
	logger.Info("Optional environment variables:")
	dumpEnvironmentVariable(EnvConfigFile)
	dumpEnvironmentVariable(EnvLogFilePath)
	dumpEnvironmentVariable(EnvLogLevel)

	dumpEnvironmentVariable(EnvResolution)
	// EnvResolution is optional.  If not set, the window size will
	// default to the values defined in globalVars.go.

	dumpEnvironmentVariable(EnvDisplayCursor)
	logger.Info("")
}

func dumpEnvironmentVariable(key string) {
	value := os.Getenv(key)
	if value == "" {
		value = MISSING_ENV_TOKEN
	}

	logger.Infof("key: %q, value: %q", key, value)
}

func dumpObfuscatedEnvironmentVariable(key string) {
	value := os.Getenv(key)
	if value == "" {
		value = MISSING_ENV_TOKEN
	} else {
		value = GetObfuscatedValue(value)
	}

	logger.Infof("key: %q, value: %q", key, value)
}

func GetObfuscatedValue(value string) string {
	length := len(value)

	obfuscatedValue := ""
	if length < 6 {
		obfuscatedValue = INVALID_ENV_TOKEN
	} else {
		if value == MISSING_ENV_TOKEN {
			return value
		} else {
			obfuscatedValue = fmt.Sprintf("%c%c%c---%c%c%c",
				value[0],
				value[1],
				value[2],
				value[length - 3],
				value[length - 2],
				value[length - 1],
			)
		}
	}

	return obfuscatedValue
}
