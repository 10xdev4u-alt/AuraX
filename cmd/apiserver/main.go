package main

import (
	"log"
	"os"

	"github.com/10xdev4u-alt/aura/pkg/api/handlers"
	"github.com/10xdev4u-alt/aura/pkg/api/middleware"
	"github.com/10xdev4u-alt/aura/pkg/config"
	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/10xdev4u-alt/aura/pkg/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.DefaultConfig()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		loadedCfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Printf("Warning: Failed to load config from %s, using defaults: %v", configPath, err)
		} else {
			cfg = loadedCfg
		}
	}

	dbCfg := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewDatabase(dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./data/firmware"
	}
	localStorage, err := storage.NewLocalStorage(storagePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	deviceHandler := handlers.NewDeviceHandler(db)
	firmwareHandler := handlers.NewFirmwareHandler(db, localStorage)
	releaseHandler := handlers.NewReleaseHandler(db)

	v1 := router.Group("/api/v1")
	{
		devices := v1.Group("/devices")
		{
			devices.GET("", deviceHandler.ListDevices)
			devices.GET("/:id", deviceHandler.GetDevice)
			devices.POST("", deviceHandler.CreateDevice)
		}

		firmware := v1.Group("/firmware")
		{
			firmware.GET("", firmwareHandler.ListFirmware)
			firmware.GET("/:id", firmwareHandler.GetFirmware)
			firmware.POST("", firmwareHandler.UploadFirmware)
		}

		releases := v1.Group("/releases")
		{
			releases.GET("", releaseHandler.ListReleases)
			releases.GET("/:id", releaseHandler.GetRelease)
			releases.POST("", releaseHandler.CreateRelease)
			releases.PUT("/:id/status", releaseHandler.UpdateReleaseStatus)
		}
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Aura API Server listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}