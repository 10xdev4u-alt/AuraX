package handlers

import (
	"net/http"

	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/gin-gonic/gin"
)

type ReleaseHandler struct {
	db *database.DB
}

func NewReleaseHandler(db *database.DB) *ReleaseHandler {
	return &ReleaseHandler{db: db}
}

func (h *ReleaseHandler) CreateRelease(c *gin.Context) {
	var req struct {
		FirmwareID   string `json:"firmware_id" binding:"required"`
		TargetFleet  string `json:"target_fleet"`
		HealthPolicy string `json:"health_policy"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	releaseID, err := h.db.CreateRelease(req.FirmwareID, req.TargetFleet, req.HealthPolicy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create release"})
		return
	}

	release, err := h.db.GetReleaseByID(releaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve release"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"release": release})
}

func (h *ReleaseHandler) ListReleases(c *gin.Context) {
	releases, err := h.db.ListReleases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list releases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"releases": releases,
		"total":    len(releases),
	})
}

func (h *ReleaseHandler) GetRelease(c *gin.Context) {
	releaseID := c.Param("id")

	release, err := h.db.GetReleaseByID(releaseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "release not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"release": release})
}

func (h *ReleaseHandler) UpdateReleaseStatus(c *gin.Context) {
	releaseID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
		Stage  string `json:"stage" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.UpdateReleaseStatus(releaseID, req.Status, req.Stage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update release"})
		return
	}

	release, err := h.db.GetReleaseByID(releaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve release"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"release": release})
}
