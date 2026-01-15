# Beelder

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Kafka](https://img.shields.io/badge/Kafka-RedPanda-FF6600?style=flat&logo=apache-kafka)](https://redpanda.com/)

> A learning project that replicates a Minecraft server platform using Go, demonstrating event-driven architecture, containerization, and the Strategy pattern.

---

## 🚀 Features

- **Dynamic Server Provisioning** - Create Minecraft servers on-demand
- **Multi-Type Support** - Paper, Forge, and Fabric servers via Strategy Pattern
- **Real-Time Updates** - SSE-based status notifications
- **Event-Driven Architecture** - Kafka/RedPanda for async processing
- **Docker Isolation** - Each server runs in its own container
- **Concurrent Builds** - Up to 3 servers building simultaneously
- **RAM-Based Plans** - Flexible memory allocation based on player needs

---

## 📋 Quick Start

### Prerequisites

```bash
# Required
- Go 1.21+
- Docker
- RedPanda (Kafka-compatible)
```

### Installation

```bash
# Clone the repository
git clone https://github.com/SGuzmanBeltran/beelder.git
cd beelder

# Start dependencies (Kafka/RedPanda)
docker-compose up -d

# Run API server
cd core/cmd/api
go run main.go

# Run Worker (in another terminal)
cd core/cmd/worker
go run main.go
```

---

## 📦 Development Slices

### ✅ First Slice - Core Server Creation (COMPLETED)
Basic server provisioning with event-driven architecture.

**Implemented Features:**
- ✅ REST API for server creation
- ✅ RedPanda event streaming
- ✅ Worker service with Docker provisioning
- ✅ Multi-type server support (Paper, Forge, Fabric) via Strategy Pattern
- ✅ Concurrent build system (max 3 simultaneous)
- ✅ SSE real-time status updates
- ✅ Configuration validation
- ✅ Health checking system

**Configuration Options:**
- Server name/MOTD
- Server type (Paper, Forge, Fabric)
- Player count (affects memory allocation)
- RAM plan (memory-based allocation)
- Difficulty (Peaceful, Easy, Normal, Hard)
- Online mode (official Minecraft accounts only)

**Progress States:**
- Created → Building → Running → Stopped → Error

---

### ✅ Second Slice - Frontend Foundation (COMPLETED)
User interface for server creation and configuration.

**Implemented Features:**
- ✅ React + TypeScript frontend
- ✅ Server creation form with validation (react-hook-form + zod)
- ✅ RAM plan selection carousel
- ✅ Dynamic RAM recommendation based on player count
- ✅ API integration for server creation
- ✅ Asset download and versioning system

---

### 🚧 Third Slice - UI Enhancements (IN PROGRESS)
Improved user experience and server management interface.

**Planned Features:**
- Custom React hooks for server management
- Dynamic RAM plan and version fetching
- Client-side routing
- Real-time UI updates via SSE
- Server details

---

### 📋 Fourth Slice - Database & State Persistence (PLANNED)
Centralized state with database persistence and query capabilities.

**Planned Features:**
- PostgreSQL database integration
- State Manager Service (new microservice)
- Event history tracking
- User ownership and multi-tenancy
- REST API for state queries
- Server status persistence

**Architecture Changes:**
- New State Manager microservice
- PostgreSQL for persistent storage
- Enhanced event schema for state updates

---

### 📋 Fifth Slice - Lifecycle Management & Cleanup (PLANNED)
Automated server lifecycle and resource management.

**Planned Features:**
- Auto-shutdown after 5 minutes of uptime (VPS resource constraints)
- Docker state reconciliation
- Orphaned container cleanup
- Server auto-restart capability
- Resource usage monitoring

**Resource Policy:**
> All servers auto-shutdown after 5 minutes to optimize VPS resource usage

## Architecture

![Beelder Architecture](docs/images/architecture-1.png)

**Flow:**
1. Client sends server creation request to API
2. API validates and publishes to Kafka
3. Worker consumes message and begins provisioning
4. Worker builds and deploys Docker container
5. Worker performs health checks
6. Worker publishes status updates to Kafka
7. Event Emitter broadcasts to connected clients
8. Client receives real-time status updates

For detailed component responsibilities, see [ARCHITECTURE.md](ARCHITECTURE.md).