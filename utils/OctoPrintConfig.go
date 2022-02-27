package utils

import (
	"io/ioutil"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Z-Bolt/OctoScreen/logger"

	"gopkg.in/yaml.v1"
)


var (
	configLocation = ".octoprint/config.yaml"
	homeOctoPi     = "/home/pi/"
)

type OctoPrintConfig struct {
	// API Settings.
	API struct {
		// Key is the current API key needed for accessing the API.
		Key string
	}

	// Server settings.
	Server struct {
		// Hosts defines the host to which to bind the server.
		Host string

		// Port defines the port to which to bind the server.
		Port int
	}
}

func ReadOctoPrintConfig() *OctoPrintConfig {
	logger.TraceEnter("OctoPrintConfig.ReadOctoPrintConfig()")

	configFilePath := os.Getenv(EnvConfigFilePath)
	if configFilePath == "" {
		configFilePath = findOctoPrintConfigFilePath()
	}

	if configFilePath == "" {
		panic("OctoPrintConfig.ReadOctoPrintConfig() - configFilePath is empty")
	}

	logger.Infof("Path to OctoPrint's config file: %q", configFilePath)

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("OctoPrintConfig.ReadOctoPrintConfig() - ReadFile() returned an error: %q", err))
	}

	cfg := &OctoPrintConfig{}
	err = yaml.Unmarshal([]byte(data), cfg)
	if err != nil {
		panic(fmt.Sprintf("OctoPrintConfig.ReadOctoPrintConfig() - error decoding YAML config file %q: %s", configFilePath, err))
	}

	logger.Infof("OctoPrintConfig.ReadOctoPrintConfig() - server host is: %q", cfg.Server.Host)
	logger.Infof("OctoPrintConfig.ReadOctoPrintConfig() - server port is: %d", cfg.Server.Port)

	logger.TraceLeave("OctoPrintConfig.ReadOctoPrintConfig()")
	return cfg
}

func findOctoPrintConfigFilePath() string {
	logger.TraceEnter("OctoPrintConfig.FindOctoPrintConfigFilePath()")

	filePath := filepath.Join(homeOctoPi, configLocation)
	if _, err := os.Stat(filePath); err == nil {
		logger.Info("OctoPrintConfig.FindOctoPrintConfigFilePath() - doFindOctoPrintConfigFilePath() found a file")
		logger.TraceLeave("OctoPrintConfig.FindOctoPrintConfigFilePath(), returning the file")
		return filePath
	}

	usr, err := user.Current()
	if err != nil {
		logger.LogError("OctoPrintConfig.FindOctoPrintConfigFilePath()", "Current()", err)
		logger.TraceLeave("OctoPrintConfig.FindOctoPrintConfigFilePath(), returning an empty string")
		return ""
	}

	octoPrintConfigFilePath := filepath.Join(usr.HomeDir, configLocation)

	logger.TraceLeave("main.FindOctoPrintConfigFilePath(), returning octoPrintConfigFilePath")
	return octoPrintConfigFilePath
}

func (this *OctoPrintConfig) OverrideConfigsWithEnvironmentValues() {
	logger.TraceEnter("OctoPrintConfig.OverrideConfigsWithEnvironmentValues()")

	apiKey := os.Getenv(EnvOctoPrintApiKey)
	if apiKey != "" {
		this.API.Key = apiKey
	}

	host := os.Getenv(EnvOctoPrintHost)
	if host != "" {
		this.Server.Host = host
	}

	// The port is not set via an environment variable.
	// Might want to add one if there is interest.

	logger.TraceLeave("OctoPrintConfig.OverrideConfigsWithEnvironmentValues()")
}

func (this *OctoPrintConfig) UpdateValues() {
	logger.TraceEnter("OctoPrintConfig.UpdateValues()")

	if this.Server.Host == "" {
		logger.Infof("Server host is empty, defaulting to the default value (%s)", DefaultServerHost)
		this.Server.Host = DefaultServerHost
	}

	if this.Server.Port == 0 || this.Server.Port == -1 {
		logger.Infof("Server port is 0, defaulting to the default value (%d)", NoServerPort)
		this.Server.Port = NoServerPort
	}

	url := strings.ToLower(this.Server.Host)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		logger.Warn("WARNING!  OCTOPRINT_HOST requires the transport protocol ('http://' or 'https://') but is missing.  'http://' is being added to Server.Host.");
		this.Server.Host = fmt.Sprintf("http://%s", this.Server.Host)
	}

	if strings.Count(this.Server.Host, ":") >= 2 {
		// The Host has a port specified.
		// One ":" is for the leading "http://"
		// And the second ":" is for the trailing port.

		if this.Server.Port != NoServerPort {
			logger.Warn("WARNING!  Server.Host includes a port value, but Server.Port has also been defined") 
			logger.Warn("WARNING!  Ignoring Server.Port and just using Server.Host") 
		}
	} else {
		// The Host doesn't specify a port.

		if this.Server.Port != NoServerPort {
			// If the user specified a port to use, append it to Server.Host.
			this.Server.Host = fmt.Sprintf("%s:%d", this.Server.Host, this.Server.Port)
		}
	}
	
	logger.TraceLeave("OctoPrintConfig.UpdateValues()")
}

func (this *OctoPrintConfig) MissingRequiredConfigName() string {
	logger.TraceEnter("OctoPrintConfig.MissingRequiredConfigName()")

	if this.API.Key == "" {
		return "API.Key"
	}

	if this.Server.Host == "" || this.Server.Host == "http://" {
		return "Server.Host"
	}

	logger.TraceLeave("OctoPrintConfig.MissingRequiredConfigName()")

	return ""
}

func (this *OctoPrintConfig) DumpConfigs() {
	// Don't add TraceEnter/TraceLeave to this function.

	logger.Infof("%-16s: %q", "API.Key", GetObfuscatedValue(this.API.Key))
	logger.Infof("%-16s: %q", "Server.Host", this.Server.Host)
	logger.Infof("%-16s: %d", "Server.Port", this.Server.Port)
}