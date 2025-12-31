# AuraX API Examples

This document provides practical examples for using the AuraX API.

## Prerequisites

```bash
# Set the API base URL
export API_BASE="http://localhost:8080"
```

## Device Management

### Create a New Device

```bash
curl -X POST $API_BASE/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "bootstrap_token": "factory-token-abc123"
  }'
```

Response:
```json
{
  "device": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "bootstrap_token": "factory-token-abc123",
    "created_at": "2025-12-31T10:30:00Z",
    "updated_at": "2025-12-31T10:30:00Z"
  }
}
```

### List All Devices

```bash
curl -X GET $API_BASE/api/v1/devices
```

### Get Specific Device

```bash
curl -X GET $API_BASE/api/v1/devices/550e8400-e29b-41d4-a716-446655440000
```

## Firmware Management

### Upload Firmware

```bash
curl -X POST $API_BASE/api/v1/firmware \
  -F "version=2.1.0" \
  -F "description=Security patch for CVE-2024-1234" \
  -F "file=@firmware-2.1.0.bin"
```

Response:
```json
{
  "firmware": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "version": "2.1.0",
    "description": "Security patch for CVE-2024-1234",
    "file_path": "/app/data/firmware/2.1.0.bin",
    "file_size": 2048576,
    "checksum": "sha256:abc123...",
    "created_at": "2025-12-31T10:35:00Z"
  }
}
```

### List All Firmware

```bash
curl -X GET $API_BASE/api/v1/firmware
```

### Get Specific Firmware

```bash
curl -X GET $API_BASE/api/v1/firmware/660e8400-e29b-41d4-a716-446655440000
```

## Release Management

### Create a Release

```bash
curl -X POST $API_BASE/api/v1/releases \
  -H "Content-Type: application/json" \
  -d '{
    "firmware_id": "660e8400-e29b-41d4-a716-446655440000",
    "target_fleet": "production",
    "health_policy": "auto-rollback-on-80-percent"
  }'
```

Response:
```json
{
  "release": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "firmware_id": "660e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "stage": "canary",
    "target_fleet": "production",
    "health_policy": "auto-rollback-on-80-percent",
    "created_at": "2025-12-31T10:40:00Z"
  }
}
```

### List All Releases

```bash
curl -X GET $API_BASE/api/v1/releases
```

### Get Specific Release

```bash
curl -X GET $API_BASE/api/v1/releases/770e8400-e29b-41d4-a716-446655440000
```

### Update Release Status

```bash
curl -X PUT $API_BASE/api/v1/releases/770e8400-e29b-41d4-a716-446655440000/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "stage": "production"
  }'
```

## Health Checks

### Service Health

```bash
curl -X GET $API_BASE/health
```

Response:
```json
{
  "status": "healthy",
  "service": "aura-api-server"
}
```

### Readiness Check

```bash
curl -X GET $API_BASE/ready
```

## Device Provisioning (gRPC)

### Using grpcurl

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Bootstrap request
grpcurl -plaintext -d '{
  "bootstrap_token": "factory-token-abc123"
}' localhost:50051 aura.provisioning.v1.ProvisioningService/Bootstrap

# Provision request (after bootstrap)
grpcurl -plaintext -d '{
  "challenge": "base64-encoded-challenge",
  "signed_challenge": "base64-encoded-signature"
}' localhost:50051 aura.provisioning.v1.ProvisioningService/Provision
```

## MQTT Communication

### Subscribe to Device Telemetry

```bash
# Install mosquitto clients
# sudo apt-get install mosquitto-clients  # Debian/Ubuntu
# brew install mosquitto                  # macOS

# Subscribe to all device telemetry
mosquitto_sub -h localhost -t "aura/devices/+/telemetry" -v

# Subscribe to specific device
mosquitto_sub -h localhost -t "aura/devices/550e8400-e29b-41d4-a716-446655440000/telemetry" -v
```

### Publish Test Telemetry

```bash
mosquitto_pub -h localhost \
  -t "aura/devices/550e8400-e29b-41d4-a716-446655440000/telemetry" \
  -m '{
    "device_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": 1704020400,
    "battery_level": 87.5,
    "temperature": 23.4,
    "uptime": 86400,
    "firmware_version": "2.0.0",
    "status": "healthy"
  }'
```

### Monitor Update Status

```bash
# Subscribe to update status from all devices
mosquitto_sub -h localhost -t "aura/devices/+/update/status" -v
```

### Simulate Update Status

```bash
mosquitto_pub -h localhost \
  -t "aura/devices/550e8400-e29b-41d4-a716-446655440000/update/status" \
  -m '{
    "device_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "downloading",
    "progress": 45,
    "error": ""
  }'
```

## Complete Workflow Example

### 1. Create Device

```bash
DEVICE_ID=$(curl -s -X POST $API_BASE/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{"bootstrap_token": "test-token-123"}' | jq -r '.device.id')

echo "Created device: $DEVICE_ID"
```

### 2. Upload Firmware

```bash
FIRMWARE_ID=$(curl -s -X POST $API_BASE/api/v1/firmware \
  -F "version=3.0.0" \
  -F "description=Major update" \
  -F "file=@firmware.bin" | jq -r '.firmware.id')

echo "Uploaded firmware: $FIRMWARE_ID"
```

### 3. Create Release

```bash
RELEASE_ID=$(curl -s -X POST $API_BASE/api/v1/releases \
  -H "Content-Type: application/json" \
  -d "{
    \"firmware_id\": \"$FIRMWARE_ID\",
    \"target_fleet\": \"production\",
    \"health_policy\": \"auto-rollback\"
  }" | jq -r '.release.id')

echo "Created release: $RELEASE_ID"
```

### 4. Monitor Release

```bash
# Check release status
curl -s $API_BASE/api/v1/releases/$RELEASE_ID | jq

# Monitor MQTT for device responses
mosquitto_sub -h localhost -t "aura/devices/+/update/status" -v
```

## Troubleshooting

### Check Service Logs

```bash
# API Server
docker-compose logs -f apiserver

# Provisioning Server
docker-compose logs -f auraserver

# OTA Orchestrator
docker-compose logs -f otaorchestrator

# MQTT Broker
docker-compose logs -f mosquitto
```

### Verify Database Connection

```bash
# Connect to PostgreSQL
docker exec -it aura-postgres psql -U aura -d aura

# List devices
SELECT id, bootstrap_token, provisioned_at FROM devices;

# List firmware
SELECT id, version, file_size FROM firmware;

# List releases
SELECT id, firmware_id, status, stage FROM releases;
```

### Test MQTT Connectivity

```bash
# Check MQTT broker is running
docker-compose ps mosquitto

# Test MQTT connection
mosquitto_pub -h localhost -t test -m "hello"
mosquitto_sub -h localhost -t test -C 1
```
