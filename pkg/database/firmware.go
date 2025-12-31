package database

import "fmt"

type Firmware struct {
	ID          string
	Version     string
	Description *string
	FilePath    string
	FileSize    int64
	Checksum    string
	CreatedAt   string
	UpdatedAt   string
}

type Release struct {
	ID           string
	FirmwareID   string
	Status       string
	Stage        string
	TargetFleet  *string
	HealthPolicy *string
	CreatedAt    string
	UpdatedAt    string
}

func (db *DB) CreateFirmware(version, description, filePath, checksum string, fileSize int64) (string, error) {
	var firmwareID string
	query := `INSERT INTO firmware (version, description, file_path, file_size, checksum) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(query, version, description, filePath, fileSize, checksum).Scan(&firmwareID)
	if err != nil {
		return "", fmt.Errorf("failed to create firmware: %w", err)
	}
	return firmwareID, nil
}

func (db *DB) GetFirmwareByID(firmwareID string) (*Firmware, error) {
	var firmware Firmware
	query := `SELECT id, version, description, file_path, file_size, checksum, created_at, updated_at 
	          FROM firmware WHERE id = $1`
	err := db.QueryRow(query, firmwareID).Scan(
		&firmware.ID, &firmware.Version, &firmware.Description, &firmware.FilePath,
		&firmware.FileSize, &firmware.Checksum, &firmware.CreatedAt, &firmware.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get firmware: %w", err)
	}
	return &firmware, nil
}

func (db *DB) ListFirmware() ([]Firmware, error) {
	query := `SELECT id, version, description, file_path, file_size, checksum, created_at, updated_at 
	          FROM firmware ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list firmware: %w", err)
	}
	defer rows.Close()

	var firmwares []Firmware
	for rows.Next() {
		var firmware Firmware
		err := rows.Scan(
			&firmware.ID, &firmware.Version, &firmware.Description, &firmware.FilePath,
			&firmware.FileSize, &firmware.Checksum, &firmware.CreatedAt, &firmware.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan firmware: %w", err)
		}
		firmwares = append(firmwares, firmware)
	}
	return firmwares, nil
}

func (db *DB) CreateRelease(firmwareID, targetFleet, healthPolicy string) (string, error) {
	var releaseID string
	query := `INSERT INTO releases (firmware_id, target_fleet, health_policy) 
	          VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, firmwareID, targetFleet, healthPolicy).Scan(&releaseID)
	if err != nil {
		return "", fmt.Errorf("failed to create release: %w", err)
	}
	return releaseID, nil
}

func (db *DB) GetReleaseByID(releaseID string) (*Release, error) {
	var release Release
	query := `SELECT id, firmware_id, status, stage, target_fleet, health_policy, created_at, updated_at 
	          FROM releases WHERE id = $1`
	err := db.QueryRow(query, releaseID).Scan(
		&release.ID, &release.FirmwareID, &release.Status, &release.Stage,
		&release.TargetFleet, &release.HealthPolicy, &release.CreatedAt, &release.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get release: %w", err)
	}
	return &release, nil
}

func (db *DB) UpdateReleaseStatus(releaseID, status, stage string) error {
	query := `UPDATE releases SET status = $1, stage = $2, updated_at = NOW() WHERE id = $3`
	_, err := db.Exec(query, status, stage, releaseID)
	if err != nil {
		return fmt.Errorf("failed to update release status: %w", err)
	}
	return nil
}

func (db *DB) ListReleases() ([]Release, error) {
	query := `SELECT id, firmware_id, status, stage, target_fleet, health_policy, created_at, updated_at 
	          FROM releases ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}
	defer rows.Close()

	var releases []Release
	for rows.Next() {
		var release Release
		err := rows.Scan(
			&release.ID, &release.FirmwareID, &release.Status, &release.Stage,
			&release.TargetFleet, &release.HealthPolicy, &release.CreatedAt, &release.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan release: %w", err)
		}
		releases = append(releases, release)
	}
	return releases, nil
}
