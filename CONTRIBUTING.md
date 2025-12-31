# Contributing to AuraX

Thank you for your interest in contributing to AuraX!

## Development Setup

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 16
- protoc compiler
- Make

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/10xdev4u-alt/aura.git
cd aura
```

2. **Install dependencies**
```bash
go mod download
```

3. **Generate protobuf code**
```bash
make proto
```

4. **Build all services**
```bash
make build-all
```

5. **Start infrastructure**
```bash
docker-compose up -d postgres mosquitto
```

6. **Run services locally**
```bash
# Terminal 1: Provisioning Server
make run

# Terminal 2: API Server
make run-api

# Terminal 3: OTA Orchestrator
make run-ota
```

## Code Style

### Go Conventions
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused
- Handle errors explicitly

### Commit Messages
Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(api): Add device filtering by status

Implement query parameter to filter devices by provisioned status.
Add database index for better query performance.

feat(ota): Implement progressive rollout stages

Add canary, staging, and production deployment stages.
Include configurable device percentages per stage.
```

## Testing

### Run Tests
```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/provisioning/...
```

### Manual Testing
```bash
# Run test script
./test.sh
```

## Pull Request Process

1. **Create a feature branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes**
- Write clean, documented code
- Add tests for new functionality
- Update documentation if needed

3. **Test your changes**
```bash
./test.sh
make build-all
```

4. **Commit your changes**
```bash
git add .
git commit -m "feat: your feature description"
```

5. **Push to your branch**
```bash
git push origin feature/your-feature-name
```

6. **Open a Pull Request**
- Provide clear description
- Reference related issues
- Include test results

## Project Structure

```
aura/
├── cmd/                    # Main applications
│   ├── auraserver/        # Provisioning server
│   ├── apiserver/         # REST API server
│   └── otaorchestrator/   # OTA orchestrator
├── pkg/                   # Shared packages
│   ├── api/              # API handlers & models
│   ├── config/           # Configuration
│   ├── database/         # Database layer
│   ├── mqtt/             # MQTT client
│   ├── ota/              # OTA logic
│   ├── pki/              # PKI/certificates
│   ├── provisioning/     # Provisioning service
│   └── storage/          # Storage abstraction
├── docs/                 # Documentation
├── gen/                  # Generated code
└── Makefile             # Build automation
```

## Adding New Features

### New API Endpoint
1. Define handler in `pkg/api/handlers/`
2. Add route in `cmd/apiserver/main.go`
3. Update `docs/05-api-examples.md`
4. Add database methods if needed

### New Database Table
1. Update schema in `pkg/database/database.go`
2. Add CRUD methods in separate file
3. Add migrations (if applicable)
4. Update documentation

### New MQTT Topic
1. Define message structure in `pkg/mqtt/messages.go`
2. Add publish/subscribe methods
3. Update orchestrator logic
4. Document in API examples

## Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Commit messages follow format
- [ ] No hardcoded credentials
- [ ] Error handling is appropriate
- [ ] Logging is adequate
- [ ] No breaking changes (or documented)

## Getting Help

- Read the documentation in `docs/`
- Check existing issues
- Ask questions in discussions
- Review API examples

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
