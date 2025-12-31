package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/10xdev4u-alt/aura/gen/go/provisioning/v1"
	"github.com/10xdev4u-alt/aura/pkg/provisioning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort = "50051"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	provisioningService := provisioning.NewProvisioningService()
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
