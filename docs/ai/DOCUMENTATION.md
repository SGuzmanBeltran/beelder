# Beelder Documentation Plan

## Philosophy
Document **intentionally, not exhaustively**. Focus on architectural decisions and code that needs context. Avoid over-documenting obvious code or creating READMEs in every directory.

---

## 1. Inline Code Documentation

### ✅ Always Document
- **Exported types and interfaces** - purpose and usage contract
- **Public functions/methods** - what they do, params, return values, errors
- **Complex logic** - why, not just what
- **Non-obvious design decisions** - tradeoffs made

### Example Format
```go
// JarManager defines the contract for managing Minecraft server JAR files.
// Implementations can use different storage backends for caching downloaded JARs.
type JarManager interface {
    // GetJar retrieves a Minecraft server JAR file for the specified configuration.
    // It checks the cache first, downloads from mcutils.com if not cached, and returns the JAR bytes.
    // Returns an error if the download fails or the server type/version is invalid.
    GetJar(ctx context.Context, serverConfig *types.CreateServerConfig) ([]byte, error)
}
```

### ❌ Don't Over-Document
- Obvious getters/setters
- Self-explanatory code
- Implementation details that should be in ADRs

---

## 2. README Structure

### Top-Level README.md (Root)
**Purpose:** Project overview for first-time visitors

**Sections:**
- What is Beelder?
- Key Features
- Quick Start / How to Run
- Architecture Overview (diagram + brief explanation)
- Tech Stack
- Development Setup
- Contributing (if open source)

### core/README.md (Backend)
**Purpose:** Developer guide for backend contributors

**Sections:**
- Backend Architecture
- Component Interaction (Worker, Builder, JAR Manager, etc.)
- Development Workflow
- Testing Strategy
- Configuration Guide
- Common Patterns Used

### app/README.md (Frontend)
**Purpose:** Frontend-specific development guide

**Sections:**
- Tech Stack (React, Vite, etc.)
- Component Structure
- State Management
- Development Workflow
- Build & Deployment

### ❌ Don't Create
- `internal/README.md`
- `internal/worker/README.md`
- `internal/worker/builder/README.md`
- READMEs for every package/directory

**Rationale:** Code docs are sufficient for implementation details. READMEs should explain **systems**, not individual components.

---

## 3. Architecture Decision Records (ADRs)

### Location
`docs/decisions/`

### Format
```
docs/
  ├── decisions/
  │   ├── README.md              # Explains ADR process
  │   ├── 001-jar-manager-abstraction.md
  │   ├── 002-redpanda-over-rabbitmq.md
  │   ├── 003-docker-based-server-isolation.md
  │   └── template.md            # Template for new ADRs
  └── images/
      └── architecture.excalidraw
```

### ADR Template
```markdown
# ADR XXX: [Short Title]

## Status
[Proposed | Accepted | Deprecated | Superseded by ADR-XXX]

## Context
What is the issue or problem we're trying to solve?
What constraints are we working with?

## Decision
What are we doing about it?
Be specific and concrete.

## Consequences

### Positive
- What benefits does this bring?
- What problems does it solve?

### Negative
- What tradeoffs are we making?
- What new problems might this create?

### Future Considerations
- What might we need to revisit later?
- When should we reconsider this decision?
```

### When to Write an ADR
✅ **Yes:** Major architectural decisions
- Choosing message broker (Redpanda vs RabbitMQ)
- JAR manager abstraction design
- Docker-based server isolation strategy
- Database choice (if applicable)
- Authentication/authorization approach

❌ **No:** Minor implementation details
- Which HTTP router to use
- Logging library choice
- Code formatting preferences
- File organization tweaks

### Target: 3-5 High-Impact ADRs
Focus on decisions that:
1. Have significant impact on architecture
2. Demonstrate your thinking process
3. Show you understand tradeoffs
4. Would be interesting to discuss in interviews

---

## 4. Recommended ADRs for Beelder

### Priority ADRs to Write

1. **ADR 001: JAR Manager Abstraction**
   - Why download on-demand vs bundled JARs
   - Interface-based design for future S3 migration
   - Local caching strategy

2. **ADR 002: Docker-Based Server Isolation**
   - Why containers for each Minecraft server
   - Resource allocation strategy
   - Port management approach

3. **ADR 003: Redpanda for Message Queue**
   - Why Redpanda (or why you chose your message broker)
   - Async server creation flow
   - Worker concurrency model

4. **ADR 004: SSE for Real-Time Updates** (if applicable)
   - Why SSE over WebSockets
   - Event streaming architecture
   - Client-server communication pattern

5. **ADR 005: Monorepo Structure** (optional)
   - Why single repo for frontend + backend
   - Build and deployment strategy
   - Shared types/contracts

---

## 5. Portfolio Interview Talking Points

Having this documentation allows you to say:

> "I implemented ADRs to document major architectural decisions. For example, I designed the JAR manager with an interface to support different storage backends - currently using local filesystem caching to keep demo costs at zero, but architected to easily swap to S3 for production multi-worker deployments."

> "The codebase uses inline documentation for all exported APIs and includes a core README explaining how the worker, builder, and message queue components interact asynchronously."

**Shows:**
- ✅ Professional documentation practices
- ✅ Thoughtful architecture
- ✅ Understanding of tradeoffs
- ✅ Production-ready mindset
- ✅ Cost consciousness

---

## Notes
- Keep docs in sync with code - update ADRs if decisions change
- Don't document for documentation's sake
- Focus on **why** over **what**
- Use diagrams where they add value
- Remember: Portfolio reviewers have limited time - make it easy to understand your thinking