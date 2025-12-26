# Project Aura: Progress Tracker

This document tracks the completed tasks and the planned future work for the Aura project.

## ✅ Done

- **[X] Project Initialization**
  - Initialized Git repository.
  - Configured user (`10xdev4u-alt`) for commits.

- **[X] Boilerplate Files**
  - Created and committed `.gitignore` for a Go project.
  - Created and committed `LICENSE` file (MIT).
  - Created and committed initial `README.md`.

- **[X] Go Module Setup**
  - Initialized Go module: `github.com/10xdev4u-alt/aura`.
  - Committed `go.mod` and `go.sum` files.
  - Added `google.golang.org/grpc` and `google.golang.org/protobuf` dependencies.

- **[X] Documentation Scaffolding**
  - Created `docs` directory.
  - Wrote and committed `01-introduction.md`.
  - Wrote and committed `02-high-level-design.md`.
  - Wrote and committed `03-low-level-design.md` (with initial gRPC spec).
  - Wrote and committed `04-api-documentation.md` placeholder.

- **[X] API Definition & Tooling**
  - Created and committed `pkg/api/v1/provisioning.proto` with the full API contract.
  - Installed Go plugins for gRPC code generation (`protoc-gen-go`, `protoc-gen-go-grpc`).

## ⏳ To Do

- **[ ] 🔴 Blocked: Generate Go Code from Protobuf**
  - **Action:** Run the `protoc` compiler to generate Go code from `provisioning.proto`.
  - **Blocker:** The `protoc` executable is not installed on the system. **User action required.**

- **[ ] Implement gRPC Server**
  - Create the main application skeleton in `cmd/auraserver/main.go`.
  - Implement the gRPC server logic.
  - Create a skeleton implementation of the `ProvisioningService` interface.

- **[ ] Database Integration**
  - Set up database connection logic.
  - Implement the database schema defined in the LLD.

- **[ ] Implement Core Provisioning Logic**
  - Flesh out the `Bootstrap` and `Provision` RPC methods.
  - Integrate with the PKI service (to be built).

- **[ ] Build API Server & OTA Orchestrator**
  - Begin development on the other core microservices as defined in the HLD.
