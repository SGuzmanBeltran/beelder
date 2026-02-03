# Beelder

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Kafka](https://img.shields.io/badge/Kafka-RedPanda-FF6600?style=flat&logo=apache-kafka)](https://redpanda.com/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)

> Beelder is a learning project replicating a Minecraft server hosting platform. Built with Go (backend) and React (frontend), it showcases event-driven architecture with RedPanda, Docker containerization, PostgreSQL persistence, automated lifecycle management, and the Strategy pattern for multi-type server support.

---

## 🚀 Features

- **Dynamic Server Provisioning** - Create Minecraft servers on-demand
- **Multi-Type Support** - Paper, Forge, and Fabric servers via Strategy Pattern
- **Real-Time Updates** - SSE-based status notifications
- **Event-Driven Architecture** - Kafka/RedPanda for async processing
- **Docker Isolation** - Each server runs in its own container
- **Concurrent Builds** - Up to 3 servers building simultaneously
- **RAM-Based Plans** - Flexible memory allocation based on player needs
- **Database Persistence** - PostgreSQL for server state and user management
- **State Management** - Centralized event tracking and server lifecycle control
- **Automated Cleanup** - Auto-shutdown after 5 minutes and orphaned container removal
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

# Run Worker
cd core/cmd/worker
go run main.go

# Run Webapp
cd app
bun run dev
```

## Architecture
### First - Third slice
![Beelder Architecture](docs/images/architecture-1.svg)

### Fourth slice - Current
![Beelder Architecture](docs/diagrams/architecture-2.svg)
![Beelder database](docs/diagrams/database.svg)

**Creation server Flow:**
1. Client sends server creation request to API
2. API validates, creates the server registry in the database and publishes to Kafka
3. Worker consumes message and begins provisioning
4. Worker builds and deploys Docker container
5. Worker performs health checks
6. Worker publishes status updates to Kafka
    6.1 State manager consumes status updates and persists to database
7. Event Emitter broadcasts to connected clients
8. Client receives real-time status updates
    8.1. State manager reviews constantly if the container is being used.
    8.2. Shut down unused servers.


## ⚠️ Limitations
- **Auto-Shutdown Policy**: All servers automatically shut down after 5 minutes of uptime to optimize VPS resources.
- **Concurrent Builds**: Limited to 3 simultaneous server builds.
- **Supported Server Types**: Paper, Forge, Fabric (via Strategy Pattern).

## Roadmap
See [docs/roadmap.md](docs/roadmap.md) for the detailed development roadmap, including completed and upcoming slices.

For detailed component responsibilities, see [ARCHITECTURE.md](ARCHITECTURE.md).