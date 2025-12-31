package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/10xdev4u-alt/aura/gen/go/provisioning/v1"
	"github.com/10xdev4u-alt/aura/pkg/config"
	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/10xdev4u-alt/aura/pkg/pki"
	"github.com/10xdev4u-alt/aura/pkg/provisioning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort = "50051"
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
		log.Printf("Warning: Database connection failed: %v", err)
		log.Println("Continuing without database (limited functionality)")
		db = nil
	} else {
		defer db.Close()
		if err := db.InitSchema(); err != nil {
			log.Fatalf("Failed to initialize database schema: %v", err)
		}
	}

	pkiService, err := pki.NewPKIService()
	if err != nil {
		log.Fatalf("Failed to initialize PKI service: %v", err)
	}
	log.Println("PKI service initialized with new CA")

	port := cfg.Server.Port
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	provisioningService := provisioning.NewProvisioningService(db, pkiService)
	pb.RegisterProvisioningServiceServer(grpcServer, provisioningService)

	reflection.Register(grpcServer)

	log.Printf("Aura Provisioning Server listening on port %s", port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}