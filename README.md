# AuraX - Self-Healing IoT Fleet Management Platform

AuraX is an enterprise-grade, autonomous IoT device lifecycle management platform. It provides secure zero-touch provisioning, intelligent OTA updates with automatic rollback, and real-time fleet monitoring.

## 🚀 Features

### Zero-Touch Provisioning
- Factory-level device bootstrap with cryptographic challenge-response
- Automatic PKI certificate issuance (RSA 4096-bit CA)
- Secure device identity management

### Intelligent OTA Updates
- Canary deployments with progressive rollouts
- Real-time health monitoring via MQTT telemetry
- Automatic rollback on failure detection
- Multi-stage release management (canary → production)

### Self-Healing Architecture
- Health check-based automatic rollback
- Release health metrics tracking
- Device telemetry monitoring
- Distributed system resilience

### Security-First Design
- Zero-trust device authentication
- Short-lived PKI certificates
- End-to-end encrypted communication
- Challenge-response authentication

## 🏗️ Architecture

### Microservices

```
┌─────────────────┐     ┌──────────────┐     ┌─────────────────┐
│  Provisioning   │────▶│  PostgreSQL  │◀────│   API Server    │
│     Server      │     │   Database   │     │   (REST API)    │
│    (gRPC)       │     └──────────────┘     └─────────────────┘
└─────────────────┘              │                     │
                                 │                     │
                    ┌────────────▼──────────┐         │
                    │   OTA Orchestrator    │◀────────┘
                    │  (Health Monitoring)  │
                    └────────────┬──────────┘
                                 │
                         ┌───────▼────────┐
                         │  MQTT Broker   │
                         │  (Mosquitto)   │
                         └────────────────┘
                                 │
                         ┌───────▼────────┐
                         │  IoT Devices   │
                         └────────────────┘
```

### Components

**1. Provisioning Server (gRPC - Port 50051)**
- Handles device bootstrap requests
- Issues cryptographic challenges
- Provisions devices with certificates
- Manages device registration

**2. API Server (REST - Port 8080)**
- Device management endpoints
- Firmware upload and storage
- Release management
- Health monitoring dashboard

**3. OTA Orchestrator**
- Monitors active releases
- Sends update commands via MQTT
- Tracks deployment health metrics
- Triggers automatic rollbacks

**4. PostgreSQL Database**
- Device registry
- Firmware metadata
- Release tracking
- User management

**5. MQTT Broker (Mosquitto)**
- Device telemetry ingestion
- OTA update commands
- Real-time device communication

## 📦 Installation

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for local development)
- protoc compiler (for development)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/10xdev4u-alt/aura.git
cd aura

# Start all services
./deploy.sh
```

Services will be available at:
- **Provisioning Server (gRPC)**: `localhost:50051`
- **API Server (REST)**: `localhost:8080`
- **PostgreSQL**: `localhost:5432`
- **MQTT Broker**: `localhost:1883`

### Manual Build

```bash
# Build all binaries
make build-all

# Run individual services
make run          # Provisioning server
make run-api      # API server
make run-ota      # OTA orchestrator
```

## 🔌 API Reference

### Device Management

**List Devices**
```http
GET /api/v1/devices
```

**Get Device**
```http
GET /api/v1/devices/{id}
```

**Create Device**
```http
POST /api/v1/devices
Content-Type: application/json

{
  "bootstrap_token": "factory-token-123"
}
```

### Firmware Management

**Upload Firmware**
```http
POST /api/v1/firmware
Content-Type: multipart/form-data

version: "2.1.0"
description: "Security patch"
file: <binary>
```

**List Firmware**
```http
GET /api/v1/firmware
```

### Release Management

**Create Release**
```http
POST /api/v1/releases
Content-Type: application/json

{
  "firmware_id": "uuid",
  "target_fleet": "production",
  "health_policy": "auto-rollback"
}
```

**Update Release Status**
```http
PUT /api/v1/releases/{id}/status
Content-Type: application/json

{
  "status": "in_progress",
  "stage": "production"
}
```

## 🔐 Security

### PKI Architecture
- Self-signed root CA (RSA 4096-bit)
- Device certificates (RSA 2048-bit, 1-year validity)
- Automatic certificate rotation
- Certificate revocation support

### Authentication Flow
1. Device boots with factory bootstrap token
2. Server validates token and issues challenge
3. Device signs challenge with private key
4. Server verifies signature and provisions certificate
5. Device uses certificate for all future communication

## 📊 Monitoring

### Health Endpoints

```http
GET /health   # Service health check
GET /ready    # Readiness probe
```

### Telemetry Topics (MQTT)

```
aura/devices/{device_id}/telemetry        # Device status
aura/devices/{device_id}/update/status    # OTA progress
aura/devices/{device_id}/update/command   # Update commands
aura/devices/{device_id}/update/rollback  # Rollback commands
```

## 🛠️ Development

### Project Structure

```
aura/
├── cmd/
│   ├── auraserver/       # Provisioning server
│   ├── apiserver/        # REST API server
│   └── otaorchestrator/  # OTA orchestrator
├── pkg/
│   ├── api/              # API handlers & models
│   ├── config/           # Configuration management
│   ├── database/         # Database layer
│   ├── mqtt/             # MQTT client
│   ├── ota/              # OTA orchestrator logic
│   ├── pki/              # Certificate management
│   ├── provisioning/     # Provisioning service
│   └── storage/          # Firmware storage
├── gen/                  # Generated protobuf code
├── docs/                 # Documentation
├── Makefile             # Build automation
├── docker-compose.yml   # Container orchestration
└── deploy.sh            # Deployment script
```

### Building from Source

```bash
# Generate protobuf code
make proto

# Build specific service
make build          # Provisioning server
make build-api      # API server
make build-ota      # OTA orchestrator

# Run tests
make test

# Clean build artifacts
make clean
```

## 🐳 Docker Deployment

### Start Services
```bash
docker-compose up -d
```

### View Logs
```bash
docker-compose logs -f [service-name]
```

### Stop Services
```bash
docker-compose down
```

### Rebuild Images
```bash
docker-compose build --no-cache
```

## 📝 Configuration

Edit `config.yaml`:

```yaml
server:
  port: "50051"

database:
  host: "postgres"
  port: 5432
  user: "aura"
  password: "aura"
  dbname: "aura"
  sslmode: "disable"
```

Environment variables:
- `CONFIG_PATH` - Path to config file
- `GRPC_PORT` - Provisioning server port
- `API_PORT` - API server port
- `MQTT_BROKER` - MQTT broker hostname
- `STORAGE_PATH` - Firmware storage directory

## 🎯 Roadmap

- [x] Zero-touch provisioning
- [x] PKI certificate management
- [x] Firmware upload & storage
- [x] OTA orchestrator with health monitoring
- [x] MQTT integration
- [x] Docker deployment
- [ ] Web dashboard UI
- [ ] Kubernetes deployment
- [ ] Multi-region support
- [ ] Device groups & fleets
- [ ] Advanced analytics

## 📄 License

MIT License - See LICENSE file for details

## 🤝 Contributing

This is a private project developed for enterprise IoT fleet management.

## 📧 Contact

For inquiries: contact@aura-iot.example.com

---

**Built with Go, gRPC, PostgreSQL, MQTT, and Docker**