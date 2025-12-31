#!/bin/bash

set -e

echo "🚀 Starting AuraX deployment..."

echo "📦 Building Docker images..."
docker-compose build

echo "🔧 Starting services..."
docker-compose up -d

echo "⏳ Waiting for services to be healthy..."
sleep 10

echo "✅ Checking service status..."
docker-compose ps

echo ""
echo "🎉 AuraX is now running!"
echo ""
echo "📡 Services:"
echo "  - Provisioning Server (gRPC): localhost:50051"
echo "  - API Server (REST):          localhost:8080"
echo "  - PostgreSQL:                 localhost:5432"
echo "  - MQTT Broker:                localhost:1883"
echo ""
echo "🔍 View logs:"
echo "  docker-compose logs -f [service-name]"
echo ""
echo "🛑 Stop services:"
echo "  docker-compose down"
