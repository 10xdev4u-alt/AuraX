# 4. API Documentation

This document provides detailed documentation for all of Aura's public and internal APIs.

## 4.1. Provisioning API (gRPC)

The gRPC API for device provisioning is the primary interface for devices to connect and get credentials.

*   **Source of Truth:** The Protobuf definition is the source of truth for this API's contract.
*   **Location:** The definition can be found in the `/api/proto/v1/provisioning.proto` file within the main repository.
*   **Documentation:** Detailed, per-method documentation will be generated from the comments within the `.proto` file.

## 4.2. Operator API (REST/GraphQL)

This API is used by the web dashboard and CLI tools to manage the Aura platform.

*   **Status:** Not yet implemented.
*   **Specification:** An OpenAPI (Swagger) or GraphQL schema will be published here once the API is designed.

---
*This document will be automatically updated as the APIs are developed.*
