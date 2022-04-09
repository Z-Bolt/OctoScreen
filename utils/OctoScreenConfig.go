package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Z-Bolt/OctoScreen/logger"
)


type octoScreenConfig struct {
	// Required configs
	OctoPrintConfig    *OctoPrintConfig
	CssStyleFilePath   string

	// Optional configs
	LogFilePath        string
	LogLevel           string
	Resolution         string
	Width              int
	Height             int
	DisplayCursor      bool
}

var octoScreenConfigInstance *octoScreenConfig

func GetOctoScreenConfigInstance() (*octoScreenConfig) {
	if octoScreenConfigInstance == nil {
		octoPrintConfig := ReadOctoPrintConfig()

		octoScreenConfigInstance = &octoScreenConfig{
			// Required configs
			OctoPrintConfig:  octoPrintConfig,
			CssStyleFilePath: "", // default to "" for now, but this must be set in the environment variables

			// Optional configs
			LogFilePath:      "",
			LogLevel:         "",
			Resolution:       "",
			Width:            -1,
			Height:           -1,
			DisplayCursor:    false,
		}

		octoScreenConfigInstance.overrideConfigsWithEnvironmentValues()
		octoScreenConfigInstance.updateValues()
	}

	return octoScreenConfigInstance
}

func (this *octoScreenConfig) overrideConfigsWithEnvironmentValues() {
	logger.TraceEnter("OctoScreenConfig.overrideConfigsWithEnvironmentValues()")

	this.OctoPrintConfig.OverrideConfigsWithEnvironmentValues()

	cssStyleFilePath := os.Getenv(EnvStylePath)
	if cssStyleFilePath != "" {
		this.CssStyleFilePath = cssStyleFilePath
	}

	logFilePath := os.Getenv(EnvLogFilePath)
	if logFilePath != "" {
		this.LogFilePath = logFilePath
	}

	logLevel := os.Getenv(EnvLogLevel)
	if logLevel != "" {
		this.LogLevel = logLevel
	}

	resolution := os.Getenv(EnvResolution)
	if resolution != "" {
		this.Resolution = resolution
	}

	// Set the width and height later, in updateValues()

	displayCursor := strings.ToLower(os.Getenv(EnvDisplayCursor))
	if displayCursor == "true" {
		this.DisplayCursor = true
	} else {
		this.DisplayCursor = false
	}

	logger.TraceLeave("OctoScreenConfig.overrideConfigsWithEnvironmentValues()")
}

func (this *octoScreenConfig) updateValues() {
	logger.TraceEnter("OctoScreenConfig.updateValues()")

	this.OctoPrintConfig.UpdateValues()

	this.setWidthAndHeight()

	logger.TraceLeave("OctoScreenConfig.updateValues()")
}

func (this *octoScreenConfig) setWidthAndHeight() {
	logger.TraceEnter("OctoScreenConfig.setWidthAndHeight()")

	if this.Resolution == "" {
		logger.Info("OctoScreenConfig.setWidthAndHeight() - Resolution is empty, defaulting to the default values defined in globalVars.go")
		this.Width = DefaultWindowWidth
		this.Height = DefaultWindowHeight
		logger.TraceLeave("OctoScreenConfig.setWidthAndHeight()")
		return
	} 

	var err error = nil
	width := -1
	height := -1
	parts := strings.SplitN(this.Resolution, "x", 2)
	if len(parts) != 2 {
		logger.Error("OctoScreenConfig.setWidthAndHeight() - SplitN() - len(parts) != 2")
		err = errors.New(fmt.Sprintf("%s is malformed\nvalue: %q", EnvResolution, this.Resolution))
	}

	if err == nil {
		width, err = strconv.Atoi(parts[0])
		if err != nil {
			logger.LogError("OctoScreenConfig.setWidthAndHeight()", "Atoi(parts[0])", err)
		} else if width < MinimumWindowWidth || width > MaximumWindowWidth {
			logger.Warn(fmt.Sprintf("window width setting (%d) is invalid (out of range)", width))
		}
	}

	if err == nil {
		height, err = strconv.Atoi(parts[1])
		if err != nil {
			logger.LogError("OctoScreenConfig.setWidthAndHeight()", "Atoi(parts[1])", err)
		} else if height < MinimumWindowHeight || height > MaximumWindowHeight {
			logger.Warn(fmt.Sprintf("window height setting (%d) is invalid (out of range)", height))
		}
	}

	if width != -1 {
		this.Width = width
	} else {
		logger.Warn(fmt.Sprintf("window width setting was not set, defaulting to the default value (%d)", DefaultWindowWidth))
		this.Width = DefaultWindowWidth
	}

	if height != -1 {
		this.Height = height
	} else {
		logger.Warn(fmt.Sprintf("window height setting was not set, defaulting to the default value (%d)", DefaultWindowHeight))
		this.Height = DefaultWindowHeight
	}

	logger.TraceLeave("OctoScreenConfig.setWidthAndHeight()")
}

func (this *octoScreenConfig) RequiredConfigsAreSet() bool {
	return this.MissingRequiredConfigName() == ""
}

func (this *octoScreenConfig) MissingRequiredConfigName() string {
	logger.TraceEnter("OctoScreenConfig.MissingRequiredConfigName()")

	missingOctoPrintConfigName := this.OctoPrintConfig.MissingRequiredConfigName()
	if missingOctoPrintConfigName != "" {
		return missingOctoPrintConfigName;
	}

	if this.CssStyleFilePath == "" {
		return "CssStyleFilePath"
	}

	logger.TraceLeave("OctoScreenConfig.MissingRequiredConfigName()")

	return ""
}

func (this *octoScreenConfig) DumpConfigs() {
	// Don't add TraceEnter/TraceLeave to this function.

	logger.Info("Dumping formatted configs...")

	// Required configs
	this.OctoPrintConfig.DumpConfigs()
	logger.Infof("%-16s: %q", "CssStyleFilePath", this.CssStyleFilePath)

	// Optional configs
	logger.Infof("%-16s: %q", "LogFilePath", this.LogFilePath)
	logger.Infof("%-16s: %q", "LogLevel", this.LogLevel)
	logger.Infof("%-16s: %q", "Resolution", this.Resolution)
	logger.Infof("%-16s: %d", "Width", this.Width)
	logger.Infof("%-16s: %d", "Height", this.Height)
	logger.Infof("%-16s: %t", "DisplayCursor", this.DisplayCursor)

	logger.Info("")
}
