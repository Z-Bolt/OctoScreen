package utils

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"github.com/Z-Bolt/OctoScreen/logger"
)

var systemDHeartbeatOnce sync.Once

type systemDHeartbeat struct {
	backgroundTask *BackgroundTask
}

var systemDHeartbeatInstance *systemDHeartbeat

func GetSystemDHeartbeatInstance() (*systemDHeartbeat) {
	if systemDHeartbeatInstance == nil {
		_, err := daemon.SdNotify(false, daemon.SdNotifyReady)
		if err != nil {
			logger.Errorf("SystemDHeartbeat.GetSystemDHeartbeatInstance() - SdNotify() returned an error: %q", err)
		}

		systemDHeartbeatOnce.Do(func() {
			systemDHeartbeatInstance = &systemDHeartbeat{}

			// Default timeout of 5 seconds
			duration := time.Second * 5

			// Experimental, set the timeout based on config setting, but only if the config is pressent.
			updateFrequency := os.Getenv("EXPERIMENTAL_SYSTEMD_HEARTBEAT_UPDATE_FREQUENCY")
			if updateFrequency != "" {
				logger.Infof("SystemDHeartbeat.GetSystemDHeartbeatInstance() - EXPERIMENTAL_SYSTEMD_HEARTBEAT_UPDATE_FREQUENCY is present, frequency is %s", updateFrequency)
				val, err := strconv.Atoi(updateFrequency)
				if err == nil {
					duration = time.Second * time.Duration(val)
				} else {
					logger.LogError("SystemDHeartbeat.GetSystemDHeartbeatInstance()", "strconv.Atoi()", err)
				}
			}

			systemDHeartbeatInstance.backgroundTask = CreateBackgroundTask(duration, func() {
				systemDHeartbeatInstance.sendHeartbeat()
			})
		})
	}

	return systemDHeartbeatInstance
}

func (this *systemDHeartbeat) Start() {
	this.backgroundTask.Start()
}

func (this *systemDHeartbeat) Stop() {
	this.backgroundTask.Close()
}

func (this *systemDHeartbeat) sendHeartbeat() {
	_, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog)
	if err != nil {
		logger.Errorf("SystemDHeartbeat.sendHeartbeat() - SdNotify() returned an error: %q", err)
	}
}

