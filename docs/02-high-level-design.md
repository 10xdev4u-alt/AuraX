# 2. High-Level Design (HLD)

## 2.1. System Architecture Overview

Aura is designed as a distributed system of microservices, primarily written in Go. This architecture ensures scalability, separation of concerns, and resilience. The major components communicate via gRPC for internal service calls and expose APIs for device and user interactions.

The system is composed of four primary services:

1.  **Aura Provisioning Service (gRPC):** The public-facing endpoint for new devices. It handles the initial bootstrap, claim, and certificate issuance process.
2.  **Aura API Server (REST/GraphQL):** The primary interface for the web dashboard and operator CLI. Used for defining device profiles, uploading firmware, and managing fleets.
3.  **Aura OTA Orchestrator (Internal Service):** The "brain" of the update system. It manages the state of rollouts, monitors release health, and triggers promotions or rollbacks.
4.  **Aura PKI Service (Internal gRPC):** A dedicated service that manages the Certificate Authority (CA), and handles the signing and revocation of device certificates.

## 2.2. Data Flow and Storage

*   **PostgreSQL Database:** The primary database for storing relational data such as device metadata, user accounts, fleet configurations, and audit logs.
*   **MQTT Broker (e.g., VerneMQ, EMQX):** The central message bus for all device telemetry data. Aura services subscribe to specific topics to receive data for OTA health monitoring.
*   **Object Storage (e.g., Minio, S3):** Used for storing firmware binaries.

## 2.3. Core Process Flows

### Provisioning Flow

1.  Device boots in factory state, connects to Wi-Fi, and calls the **Provisioning Service**.
2.  Operator "claims" the device via the **API Server**.
3.  The **API Server** instructs the **Provisioning Service** to proceed.
4.  **Provisioning Service** requests a new certificate from the **PKI Service**.
5.  **PKI Service** signs a certificate and returns it.
6.  **Provisioning Service** sends the new certificate and MQTT endpoint details to the device.
7.  Device connects to the MQTT broker and begins sending telemetry.

### OTA Update Flow

1.  Operator uploads firmware and defines a rollout policy via the **API Server**.
2.  The **API Server** stores the firmware in Object Storage and creates a new "release" record in the database.
3.  The **OTA Orchestrator** detects the new release.
4.  The **Orchestrator** sends an update command via MQTT to the first canary group of devices.
5.  Devices download the firmware, verify its signature, and reboot.
6.  The **Orchestrator** monitors the telemetry from the updated devices.
7.  If health metrics are stable, it promotes the release to the next stage. If not, it sends a rollback command via MQTT.
