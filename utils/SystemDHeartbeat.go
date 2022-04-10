package utils

import (
	// "os"
	// "strconv"
	"sync"
	// "time"

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
			systemDHeartbeatInstance.createBackgroundTask()
		})
	}

	return systemDHeartbeatInstance
}

func (this *systemDHeartbeat) createBackgroundTask() {
	logger.TraceEnter("SystemDHeartbeat.createBackgroundTask()")

	// Default timeout of 5 seconds
	duration := GetExperimentalFrequency(5, "EXPERIMENTAL_SYSTEMD_HEARTBEAT_UPDATE_FREQUENCY")
	systemDHeartbeatInstance.backgroundTask = CreateBackgroundTask(duration, func() {
		systemDHeartbeatInstance.sendHeartbeat()
	})

	logger.TraceLeave("SystemDHeartbeat.createBackgroundTask()")
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
