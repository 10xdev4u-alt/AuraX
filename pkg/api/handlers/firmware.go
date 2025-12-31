package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/10xdev4u-alt/aura/pkg/storage"
	"github.com/gin-gonic/gin"
)

type FirmwareHandler struct {
	db      *database.DB
	storage *storage.LocalStorage
}

func NewFirmwareHandler(db *database.DB, storage *storage.LocalStorage) *FirmwareHandler {
	return &FirmwareHandler{
		db:      db,
		storage: storage,
	}
}

func (h *FirmwareHandler) UploadFirmware(c *gin.Context) {
	version := c.PostForm("version")
	description := c.PostForm("description")

	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version is required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	hash := sha256.New()
	tempFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer tempFile.Close()

	fileContent := make([]byte, file.Size)
	_, err = tempFile.Read(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file content"})
		return
	}
	hash.Write(fileContent)
	checksum := hex.EncodeToString(hash.Sum(nil))

	firmwareID := ""
	filePath, fileSize, err := h.storage.SaveFirmware(version, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save firmware"})
		return
	}

	firmwareID, err = h.db.CreateFirmware(version, description, filePath, checksum, fileSize)
	if err != nil {
		h.storage.DeleteFirmware(version)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create firmware record"})
		return
	}

	firmware, err := h.db.GetFirmwareByID(firmwareID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve firmware"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"firmware": firmware})
}

func (h *FirmwareHandler) ListFirmware(c *gin.Context) {
	firmwares, err := h.db.ListFirmware()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list firmware"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"firmwares": firmwares,
		"total":     len(firmwares),
	})
}

func (h *FirmwareHandler) GetFirmware(c *gin.Context) {
	firmwareID := c.Param("id")

	firmware, err := h.db.GetFirmwareByID(firmwareID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "firmware not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"firmware": firmware})
}
