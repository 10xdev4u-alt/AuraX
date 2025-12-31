package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type DB struct {
	*sql.DB
}

func NewDatabase(cfg Config) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return &DB{db}, nil
}

func (db *DB) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS devices (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		bootstrap_token TEXT UNIQUE,
		claimed_by_user_id UUID,
		claimed_at TIMESTAMPTZ,
		provisioned_at TIMESTAMPTZ,
		certificate_serial TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_devices_bootstrap_token ON devices(bootstrap_token);
	CREATE INDEX IF NOT EXISTS idx_devices_claimed_by_user ON devices(claimed_by_user_id);

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email TEXT UNIQUE NOT NULL,
		name TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS firmware (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		version TEXT UNIQUE NOT NULL,
		description TEXT,
		file_path TEXT NOT NULL,
		file_size BIGINT NOT NULL,
		checksum TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS releases (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		firmware_id UUID NOT NULL REFERENCES firmware(id),
		status TEXT NOT NULL DEFAULT 'pending',
		stage TEXT NOT NULL DEFAULT 'canary',
		target_fleet TEXT,
		health_policy TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_releases_firmware ON releases(firmware_id);
	CREATE INDEX IF NOT EXISTS idx_releases_status ON releases(status);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database schema initialized")
	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) BootstrapTokenExists(token string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM devices WHERE bootstrap_token = $1 AND provisioned_at IS NULL)`
	err := db.QueryRow(query, token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check token: %w", err)
	}
	return exists, nil
}

func (db *DB) CreateDevice() (string, error) {
	var deviceID string
	query := `INSERT INTO devices (bootstrap_token) VALUES (NULL) RETURNING id`
	err := db.QueryRow(query).Scan(&deviceID)
	if err != nil {
		return "", fmt.Errorf("failed to create device: %w", err)
	}
	return deviceID, nil
}

func (db *DB) MarkDeviceProvisioned(deviceID string) error {
	query := `UPDATE devices SET provisioned_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := db.Exec(query, deviceID)
	if err != nil {
		return fmt.Errorf("failed to mark device provisioned: %w", err)
	}
	return nil
}

func (db *DB) GetDeviceByID(deviceID string) (*Device, error) {
	var device Device
	query := `SELECT id, bootstrap_token, claimed_by_user_id, claimed_at, provisioned_at, 
	          certificate_serial, created_at, updated_at FROM devices WHERE id = $1`
	err := db.QueryRow(query, deviceID).Scan(
		&device.ID, &device.BootstrapToken, &device.ClaimedByUserID,
		&device.ClaimedAt, &device.ProvisionedAt, &device.CertificateSerial,
		&device.CreatedAt, &device.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("device not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return &device, nil
}

type Device struct {
	ID                string
	BootstrapToken    *string
	ClaimedByUserID   *string
	ClaimedAt         *string
	ProvisionedAt     *string
	CertificateSerial *string
	CreatedAt         string
	UpdatedAt         string
}

func (db *DB) ListDevices() ([]Device, error) {
	query := `SELECT id, bootstrap_token, claimed_by_user_id, claimed_at, provisioned_at, 
	          certificate_serial, created_at, updated_at FROM devices ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		err := rows.Scan(
			&device.ID, &device.BootstrapToken, &device.ClaimedByUserID,
			&device.ClaimedAt, &device.ProvisionedAt, &device.CertificateSerial,
			&device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating devices: %w", err)
	}

	return devices, nil
}

func (db *DB) CreateDeviceWithToken(token string) (string, error) {
	var deviceID string
	query := `INSERT INTO devices (bootstrap_token) VALUES ($1) RETURNING id`
	err := db.QueryRow(query, token).Scan(&deviceID)
	if err != nil {
		return "", fmt.Errorf("failed to create device with token: %w", err)
	}
	return deviceID, nil
}