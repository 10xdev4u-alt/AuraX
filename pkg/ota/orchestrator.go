package ota

import (
	"log"
	"time"

	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/10xdev4u-alt/aura/pkg/mqtt"
)

type Orchestrator struct {
	db            *database.DB
	mqttClient    *mqtt.Client
	pollInterval  time.Duration
	stopChan      chan struct{}
	healthMetrics map[string]*ReleaseHealth
}

type ReleaseHealth struct {
	ReleaseID      string
	SuccessCount   int
	FailureCount   int
	TotalDevices   int
	CurrentStage   string
	LastUpdateTime time.Time
}

func NewOrchestrator(db *database.DB, mqttClient *mqtt.Client) *Orchestrator {
	return &Orchestrator{
		db:            db,
		mqttClient:    mqttClient,
		pollInterval:  30 * time.Second,
		stopChan:      make(chan struct{}),
		healthMetrics: make(map[string]*ReleaseHealth),
	}
}

func (o *Orchestrator) Start() {
	log.Println("OTA Orchestrator started")

	o.mqttClient.SubscribeToTelemetry(o.handleTelemetry)
	o.mqttClient.SubscribeToUpdateStatus(o.handleUpdateStatus)

	ticker := time.NewTicker(o.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.processReleases()
		case <-o.stopChan:
			log.Println("OTA Orchestrator stopped")
			return
		}
	}
}

func (o *Orchestrator) handleTelemetry(telemetry *mqtt.DeviceTelemetry) {
	log.Printf("Received telemetry from device %s: status=%s, version=%s",
		telemetry.DeviceID, telemetry.Status, telemetry.FirmwareVersion)
}

func (o *Orchestrator) handleUpdateStatus(status *mqtt.UpdateStatus) {
	log.Printf("Update status from device %s: %s (progress: %d%%)",
		status.DeviceID, status.Status, status.Progress)

	if status.Status == "completed" {
		log.Printf("Device %s successfully updated", status.DeviceID)
	} else if status.Status == "failed" {
		log.Printf("Device %s update failed: %s", status.DeviceID, status.Error)
	}
}

func (o *Orchestrator) Stop() {
	close(o.stopChan)
}

func (o *Orchestrator) processReleases() {
	releases, err := o.db.ListReleases()
	if err != nil {
		log.Printf("Error fetching releases: %v", err)
		return
	}

	for _, release := range releases {
		if release.Status == "pending" {
			o.startRelease(release.ID)
		} else if release.Status == "in_progress" {
			o.monitorRelease(release.ID)
		}
	}
}

func (o *Orchestrator) startRelease(releaseID string) {
	log.Printf("Starting release: %s", releaseID)

	err := o.db.UpdateReleaseStatus(releaseID, "in_progress", "canary")
	if err != nil {
		log.Printf("Error starting release %s: %v", releaseID, err)
		return
	}

	o.healthMetrics[releaseID] = &ReleaseHealth{
		ReleaseID:      releaseID,
		CurrentStage:   "canary",
		LastUpdateTime: time.Now(),
	}

	release, err := o.db.GetReleaseByID(releaseID)
	if err != nil {
		log.Printf("Error getting release %s: %v", releaseID, err)
		return
	}

	firmware, err := o.db.GetFirmwareByID(release.FirmwareID)
	if err != nil {
		log.Printf("Error getting firmware for release %s: %v", releaseID, err)
		return
	}

	devices, err := o.db.ListDevices()
	if err != nil {
		log.Printf("Error getting devices for release %s: %v", releaseID, err)
		return
	}

	canaryCount := 5
	if len(devices) < canaryCount {
		canaryCount = len(devices)
	}

	for i := 0; i < canaryCount; i++ {
		device := devices[i]
		cmd := &mqtt.UpdateCommand{
			DeviceID:    device.ID,
			FirmwareURL: "https://firmware.aura.example.com/" + firmware.ID,
			Version:     firmware.Version,
			Checksum:    firmware.Checksum,
		}
		if err := o.mqttClient.PublishUpdateCommand(device.ID, cmd); err != nil {
			log.Printf("Error sending update command to device %s: %v", device.ID, err)
		}
	}

	log.Printf("Release %s moved to canary stage, sent to %d devices", releaseID, canaryCount)
}

func (o *Orchestrator) monitorRelease(releaseID string) {
	health, exists := o.healthMetrics[releaseID]
	if !exists {
		health = &ReleaseHealth{
			ReleaseID:      releaseID,
			CurrentStage:   "canary",
			LastUpdateTime: time.Now(),
		}
		o.healthMetrics[releaseID] = health
	}

	health.SuccessCount++
	health.TotalDevices++

	successRate := float64(health.SuccessCount) / float64(health.TotalDevices)

	if successRate < 0.8 && health.TotalDevices > 5 {
		log.Printf("Release %s failed health check (success rate: %.2f), rolling back", releaseID, successRate)
		o.rollbackRelease(releaseID)
		return
	}

	if health.TotalDevices > 10 && health.CurrentStage == "canary" {
		log.Printf("Release %s passed canary stage, promoting to production", releaseID)
		o.promoteRelease(releaseID, "production")
	}

	if health.TotalDevices > 50 && health.CurrentStage == "production" {
		log.Printf("Release %s completed successfully", releaseID)
		o.completeRelease(releaseID)
	}
}

func (o *Orchestrator) promoteRelease(releaseID, stage string) {
	err := o.db.UpdateReleaseStatus(releaseID, "in_progress", stage)
	if err != nil {
		log.Printf("Error promoting release %s: %v", releaseID, err)
		return
	}

	if health, exists := o.healthMetrics[releaseID]; exists {
		health.CurrentStage = stage
	}

	log.Printf("Release %s promoted to %s stage", releaseID, stage)
}

func (o *Orchestrator) rollbackRelease(releaseID string) {
	err := o.db.UpdateReleaseStatus(releaseID, "rolled_back", "rollback")
	if err != nil {
		log.Printf("Error rolling back release %s: %v", releaseID, err)
		return
	}

	devices, err := o.db.ListDevices()
	if err != nil {
		log.Printf("Error getting devices for rollback: %v", err)
		return
	}

	for _, device := range devices {
		if err := o.mqttClient.PublishRollbackCommand(device.ID); err != nil {
			log.Printf("Error sending rollback command to device %s: %v", device.ID, err)
		}
	}

	delete(o.healthMetrics, releaseID)
	log.Printf("Release %s rolled back, rollback commands sent to all devices", releaseID)
}

func (o *Orchestrator) completeRelease(releaseID string) {
	err := o.db.UpdateReleaseStatus(releaseID, "completed", "completed")
	if err != nil {
		log.Printf("Error completing release %s: %v", releaseID, err)
		return
	}

	delete(o.healthMetrics, releaseID)
	log.Printf("Release %s completed", releaseID)
}

func (o *Orchestrator) GetReleaseHealth(releaseID string) *ReleaseHealth {
	return o.healthMetrics[releaseID]
}