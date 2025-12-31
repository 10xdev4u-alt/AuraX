package models

import "time"

type Device struct {
	ID                string     `json:"id"`
	BootstrapToken    *string    `json:"bootstrap_token,omitempty"`
	ClaimedByUserID   *string    `json:"claimed_by_user_id,omitempty"`
	ClaimedAt         *time.Time `json:"claimed_at,omitempty"`
	ProvisionedAt     *time.Time `json:"provisioned_at,omitempty"`
	CertificateSerial *string    `json:"certificate_serial,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateDeviceRequest struct {
	BootstrapToken string `json:"bootstrap_token" binding:"required"`
}

type CreateDeviceResponse struct {
	Device Device `json:"device"`
}

type ListDevicesResponse struct {
	Devices []Device `json:"devices"`
	Total   int      `json:"total"`
}

type Firmware struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	FileURL     string    `json:"file_url"`
	FileSize    int64     `json:"file_size"`
	Checksum    string    `json:"checksum"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateFirmwareRequest struct {
	Version     string `json:"version" binding:"required"`
	Description string `json:"description"`
}

type Release struct {
	ID           string    `json:"id"`
	FirmwareID   string    `json:"firmware_id"`
	Status       string    `json:"status"`
	Stage        string    `json:"stage"`
	TargetFleet  string    `json:"target_fleet"`
	HealthPolicy string    `json:"health_policy"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
