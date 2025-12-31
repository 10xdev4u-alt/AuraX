#!/bin/bash

set -e

echo "🧪 AuraX Test Suite"
echo ""

echo "1️⃣  Testing Go builds..."
make clean
make build-all

if [ ! -f "bin/auraserver" ] || [ ! -f "bin/apiserver" ] || [ ! -f "bin/otaorchestrator" ]; then
    echo "❌ Build failed - binaries not found"
    exit 1
fi
echo "✅ All binaries built successfully"

echo ""
echo "2️⃣  Testing Docker builds..."
docker-compose build --quiet

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed"
    exit 1
fi
echo "✅ Docker images built successfully"

echo ""
echo "3️⃣  Testing Docker Compose configuration..."
docker-compose config > /dev/null

if [ $? -ne 0 ]; then
    echo "❌ Docker Compose configuration invalid"
    exit 1
fi
echo "✅ Docker Compose configuration valid"

echo ""
echo "4️⃣  Checking Go module dependencies..."
go mod verify

if [ $? -ne 0 ]; then
    echo "❌ Go module verification failed"
    exit 1
fi
echo "✅ Go modules verified"

echo ""
echo "5️⃣  Running go vet..."
go vet ./...

if [ $? -ne 0 ]; then
    echo "❌ Go vet found issues"
    exit 1
fi
echo "✅ Go vet passed"

echo ""
echo "6️⃣  Checking code formatting..."
gofmt -l . > /tmp/gofmt-output.txt

if [ -s /tmp/gofmt-output.txt ]; then
    echo "❌ Code formatting issues found:"
    cat /tmp/gofmt-output.txt
    exit 1
fi
echo "✅ Code properly formatted"

echo ""
echo "🎉 All tests passed!"
echo ""
echo "📊 Summary:"
echo "  ✅ Binaries: 3/3"
echo "  ✅ Docker images: 3/3"
echo "  ✅ Configuration: Valid"
echo "  ✅ Dependencies: Verified"
echo "  ✅ Code quality: Passed"
