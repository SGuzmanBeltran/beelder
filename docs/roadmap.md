# Beelder Development Roadmap

This document outlines the planned development slices for the Beelder project, a Minecraft server hosting platform built with Go and React. Each slice represents a milestone with specific features and tasks.

## ✅ Completed Slices

### First Slice - Core Server Creation
- REST API for server creation
- RedPanda event streaming
- Worker service with Docker provisioning
- Multi-type server support (Paper, Forge, Fabric) via Strategy Pattern
- Concurrent build system (max 3 simultaneous)
- SSE real-time status updates
- Configuration validation
- Health checking system

### Second Slice - Frontend Foundation
- React + TypeScript frontend
- Server creation form with validation (react-hook-form + zod)
- RAM plan selection carousel
- Dynamic RAM recommendation based on player count
- API integration for server creation
- Asset download and versioning system

### Third Slice - UI Enhancements
- Custom React hooks for server management
- Dynamic RAM plan and version fetching
- Client-side routing with React Router
- Real-time UI updates via SSE
- Server details view
- Ordered server versions
- Updated architecture diagram

## 🚧 Upcoming Slices

### Fourth Slice - Database & State Persistence
- Design database schema (servers, server_events, users, server_plans)
- Add PostgreSQL to docker-compose
- Implement database migrations
- Add database connection pooling
- Create new State Manager service structure (`cmd/state-manager`)
- Implement event consumer (consume from RedPanda)
- Implement state persistence layer
- Create server repository (CRUD operations)
- Create server plans repository (CRUD operations)
- Fetch the servers plans (RAM-based)
- Add event history tracking
- Create REST endpoints for server state queries
- Add user-specific server listing
- Implement server status endpoint
- Add server ownership validation
- Update architecture diagram (add State Manager + PostgreSQL)
- Add State Manager to docker-compose
- Document new event schema

### Fifth Slice - Lifecycle Management & Cleanup
- Implement background job for server monitoring
- Add 5-minute timeout policy for all servers (VPS resource constraints)
- Publish shutdown events to RedPanda
- Implement graceful server shutdown logic in Worker
- Add server auto-start capability
- Implement Docker state reconciliation job
- Add orphaned container detection and cleanup
- Create cleanup procedures for stopped servers
- Add resource usage monitoring
- Update README with lifecycle policies

## 📋 Notes
- Each slice is designed to be independent and incrementally build upon the previous one.
- The project follows an event-driven architecture with microservices (API, Worker, State Manager).
- Resource constraints (VPS limitations) drive design decisions like auto-shutdown policies.
- For more details, see [TODO.md](../TODO.md) and the main [README.md](../README.md).