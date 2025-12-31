package handlers

import (
	"net/http"

	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	db *database.DB
}

func NewDeviceHandler(db *database.DB) *DeviceHandler {
	return &DeviceHandler{db: db}
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	devices, err := h.db.ListDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
		"total":   len(devices),
	})
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
	deviceID := c.Param("id")

	device, err := h.db.GetDeviceByID(deviceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"device": device})
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req struct {
		BootstrapToken string `json:"bootstrap_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deviceID, err := h.db.CreateDeviceWithToken(req.BootstrapToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
		return
	}

	device, err := h.db.GetDeviceByID(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"device": device})
}
