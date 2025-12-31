package provisioning

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	pb "github.com/10xdev4u-alt/aura/gen/go/provisioning/v1"
	"github.com/10xdev4u-alt/aura/pkg/database"
	"github.com/10xdev4u-alt/aura/pkg/pki"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProvisioningService struct {
	pb.UnimplementedProvisioningServiceServer
	db         *database.DB
	pkiService *pki.PKIService
	challenges map[string]time.Time
}

func NewProvisioningService(db *database.DB, pkiService *pki.PKIService) *ProvisioningService {
	return &ProvisioningService{
		db:         db,
		pkiService: pkiService,
		challenges: make(map[string]time.Time),
	}
}

func (s *ProvisioningService) Bootstrap(ctx context.Context, req *pb.BootstrapRequest) (*pb.BootstrapResponse, error) {
	if req.BootstrapToken == "" {
		return nil, status.Error(codes.InvalidArgument, "bootstrap_token is required")
	}

	exists, err := s.db.BootstrapTokenExists(req.BootstrapToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to verify token")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "invalid bootstrap token")
	}

	challenge, err := generateChallenge()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate challenge")
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	s.challenges[challenge] = expiresAt

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

	expiresAt, exists := s.challenges[req.Challenge]
	if !exists {
		return nil, status.Error(codes.InvalidArgument, "invalid challenge")
	}
	if time.Now().After(expiresAt) {
		delete(s.challenges, req.Challenge)
		return nil, status.Error(codes.DeadlineExceeded, "challenge expired")
	}

	delete(s.challenges, req.Challenge)

	deviceID, err := s.db.CreateDevice()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create device")
	}

	clientCert, clientKey, err := s.pkiService.IssueCertificate(deviceID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to issue certificate")
	}

	err = s.db.MarkDeviceProvisioned(deviceID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update device status")
	}

	caCert := s.pkiService.GetCACertPEM()

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