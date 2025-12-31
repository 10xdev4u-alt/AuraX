package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/10xdev4u-alt/aura/pkg/config"
	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/10xdev4u-alt/aura/pkg/ota"
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

	orchestrator := ota.NewOrchestrator(db)

	go orchestrator.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down OTA Orchestrator...")
	orchestrator.Stop()
	log.Println("OTA Orchestrator stopped")
}
