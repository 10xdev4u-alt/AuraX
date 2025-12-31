package provisioning

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	pb "github.com/10xdev4u-alt/aura/gen/go/provisioning/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProvisioningService struct {
	pb.UnimplementedProvisioningServiceServer
}

func NewProvisioningService() *ProvisioningService {
	return &ProvisioningService{}
}

func (s *ProvisioningService) Bootstrap(ctx context.Context, req *pb.BootstrapRequest) (*pb.BootstrapResponse, error) {
	if req.BootstrapToken == "" {
		return nil, status.Error(codes.InvalidArgument, "bootstrap_token is required")
	}

	challenge, err := generateChallenge()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate challenge")
	}

	expiresAt := time.Now().Add(5 * time.Minute)

	return &pb.BootstrapResponse{
		Challenge: challenge,
		ExpiresAt: timestamppb.New(expiresAt),
	}, nil
}

func (s *ProvisioningService) Provision(ctx context.Context, req *pb.ProvisionRequest) (*pb.ProvisionResponse, error) {
	if req.Challenge == "" {
		return nil, status.Error(codes.InvalidArgument, "challenge is required")
	}
	if len(req.SignedChallenge) == 0 {
		return nil, status.Error(codes.InvalidArgument, "signed_challenge is required")
	}

	deviceID := generateDeviceID()

	clientCert := "-----BEGIN CERTIFICATE-----\nMIIC...(placeholder)\n-----END CERTIFICATE-----"
	clientKey := "-----BEGIN PRIVATE KEY-----\nMIIE...(placeholder)\n-----END PRIVATE KEY-----"
	caCert := "-----BEGIN CERTIFICATE-----\nMIIC...(placeholder CA)\n-----END CERTIFICATE-----"

	return &pb.ProvisionResponse{
		DeviceId:          deviceID,
		ClientCertificate: clientCert,
		ClientKey:         clientKey,
		CaCertificate:     caCert,
		MqttHost:          "mqtt.aura.example.com",
		MqttPort:          8883,
	}, nil
}

func generateChallenge() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateDeviceID() string {
	return fmt.Sprintf("device-%d", time.Now().Unix())
}
