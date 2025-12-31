package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceTelemetry struct {
	DeviceID    string  `json:"device_id"`
	Timestamp   int64   `json:"timestamp"`
	BatteryLevel float64 `json:"battery_level,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	Uptime      int64   `json:"uptime,omitempty"`
	FirmwareVersion string `json:"firmware_version,omitempty"`
	Status      string  `json:"status"`
}

type UpdateCommand struct {
	DeviceID    string `json:"device_id"`
	FirmwareURL string `json:"firmware_url"`
	Version     string `json:"version"`
	Checksum    string `json:"checksum"`
}

type UpdateStatus struct {
	DeviceID string `json:"device_id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Error    string `json:"error,omitempty"`
}

type TelemetryHandler func(telemetry *DeviceTelemetry)
type UpdateStatusHandler func(status *UpdateStatus)

func (c *Client) SubscribeToTelemetry(handler TelemetryHandler) error {
	topic := "aura/devices/+/telemetry"
	return c.Subscribe(topic, func(client mqtt.Client, msg mqtt.Message) {
		var telemetry DeviceTelemetry
		if err := json.Unmarshal(msg.Payload(), &telemetry); err != nil {
			log.Printf("Error parsing telemetry: %v", err)
			return
		}
		handler(&telemetry)
	})
}

func (c *Client) SubscribeToUpdateStatus(handler UpdateStatusHandler) error {
	topic := "aura/devices/+/update/status"
	return c.Subscribe(topic, func(client mqtt.Client, msg mqtt.Message) {
		var status UpdateStatus
		if err := json.Unmarshal(msg.Payload(), &status); err != nil {
			log.Printf("Error parsing update status: %v", err)
			return
		}
		handler(&status)
	})
}

func (c *Client) PublishUpdateCommand(deviceID string, cmd *UpdateCommand) error {
	topic := "aura/devices/" + deviceID + "/update/command"
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return c.Publish(topic, payload)
}

func (c *Client) PublishRollbackCommand(deviceID string) error {
	topic := "aura/devices/" + deviceID + "/update/rollback"
	payload := []byte(`{"action":"rollback"}`)
	return c.Publish(topic, payload)
}
