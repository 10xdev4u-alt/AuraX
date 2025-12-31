package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) SaveFirmware(firmwareID string, data io.Reader) (string, int64, error) {
	filePath := filepath.Join(s.basePath, firmwareID+".bin")

	file, err := os.Create(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create firmware file: %w", err)
	}
	defer file.Close()

	size, err := io.Copy(file, data)
	if err != nil {
		return "", 0, fmt.Errorf("failed to write firmware: %w", err)
	}

	return filePath, size, nil
}

func (s *LocalStorage) GetFirmware(firmwareID string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.basePath, firmwareID+".bin")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open firmware file: %w", err)
	}
	return file, nil
}

func (s *LocalStorage) DeleteFirmware(firmwareID string) error {
	filePath := filepath.Join(s.basePath, firmwareID+".bin")
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete firmware: %w", err)
	}
	return nil
}

func (s *LocalStorage) FirmwareExists(firmwareID string) bool {
	filePath := filepath.Join(s.basePath, firmwareID+".bin")
	_, err := os.Stat(filePath)
	return err == nil
}
