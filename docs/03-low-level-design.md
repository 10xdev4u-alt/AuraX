# 3. Low-Level Design (LLD)

## 3.1. Provisioning Service (gRPC)

This service is the public-facing entry point for devices. It must be highly available and secure.

### 3.1.1. Protobuf Definition (`provisioning.proto`)

We will define our service using Protocol Buffers. This gives us a strongly-typed, language-agnostic API contract.

```protobuf
syntax = "proto3";

package aura.provisioning.v1;

import "google/protobuf/timestamp.proto";

// The primary service for device provisioning
service ProvisioningService {
  // Handles the initial handshake from a device in bootstrap mode.
  // The device sends a temporary bootstrap token, and the server returns a
  // unique challenge to be signed.
  rpc Bootstrap(BootstrapRequest) returns (BootstrapResponse);

  // The device sends the signed challenge. If valid, the server returns
  // the permanent client certificate, private key, and MQTT server details.
  rpc Provision(ProvisionRequest) returns (ProvisionResponse);
}

message BootstrapRequest {
  // A single-use token hard-coded at the factory.
  string bootstrap_token = 1;
}

message BootstrapResponse {
  // A unique, server-generated challenge (e.g., a JWT or random string)
  // that the device must sign with its factory private key.
  string challenge = 1;
  google.protobuf.Timestamp expires_at = 2;
}

message ProvisionRequest {
  // The original challenge received from the server.
  string challenge = 1;
  // The challenge, signed by the device's private key.
  bytes signed_challenge = 2;
}

message ProvisionResponse {
  // The unique ID for this device in the Aura system.
  string device_id = 1;
  // The PEM-encoded client certificate for this device.
  string client_certificate = 2;
  // The PEM-encoded private key for this device.
  string client_key = 3;
  // The PEM-encoded CA certificate to validate the server.
  string ca_certificate = 4;
  // The hostname of the MQTT broker to connect to.
  string mqtt_host = 5;
  // The port of the MQTT broker.
  int32 mqtt_port = 6;
}
```

### 3.1.2. Database Schema

The provisioning service will interact with a `devices` table.

*   **devices**
    *   `id` (UUID, Primary Key)
    *   `bootstrap_token` (TEXT, UNIQUE, Indexed) - The factory token. Null after provisioning.
    *   `claimed_by_user_id` (UUID, Foreign Key to `users.id`)
    *   `claimed_at` (TIMESTAMPTZ)
    *   `provisioned_at` (TIMESTAMPTZ)
    *   `certificate_serial` (TEXT) - Serial number of the current active certificate.
    *   `created_at` (TIMESTAMPTZ)
    *   `updated_at` (TIMESTAMPTZ)

---
*Further details on other services will be added as they are developed.*
