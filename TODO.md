# Tomorrow's Tasks

## 1. Add Form Validation
- Install: `bun add react-hook-form zod @hookform/resolvers` ✅
- Validate required fields (server type, version, region, name) ✅
- Show errors when user tries to submit (right now cant submit and see data for some reason) ✅

## 2. Wire Up Form Logic
- API endpoint to determine recommended package ✅
- Connect player count slider → auto-select RAM plan ✅
- Connect carousel selection to form state ✅
- Make "SELECT PLAN" button submit the form ✅

## 3. Connect to Backend (if time)
- Send form data to API ✅
- Update API create server endpoint ✅
- Handle success/error ✅
- Update Worker create server logic ✅

That's it. Keep it simple.

## Third Slice - UI Enhancements

## 1. Refactor create-server. Use custom hook. ✅
## 2. Fetch jar servers versions.
    - Create endpoint (/server/{server_type}/versions) ✅
    - Use version_provider.go ✅
    - Cache endpoint ✅
    - Fetch when select a server type.
    - Update options in the UI
## 3. Add routing to the UI
## 4. Update UI when server is created.
## 5. Create new view to show the new server info
## 6. Update architecture diagram

---

## Fourth Slice - Database & State Persistence

## 1. Design database schema (servers, server_events, users, server_plans)
## 2. Add PostgreSQL to docker-compose
## 3. Implement database migrations
## 4. Add database connection pooling
## 5. Create new State Manager service structure (`cmd/state-manager`)
## 6. Implement event consumer (consume from RedPanda)
## 7. Implement state persistence layer
## 8. Create server repository (CRUD operations)
## 9. Create server plans repository (CRUD operations)
## 10. Fetch the servers plans (RAM-based).
## 11. Add event history tracking
## 12. Create REST endpoints for server state queries
## 13. Add user-specific server listing
## 14. Implement server status endpoint
## 15. Add server ownership validation
## 16. Update architecture diagram (add State Manager + PostgreSQL)
## 17. Add State Manager to docker-compose
## 18. Document new event schema

---

## Fifth Slice - Lifecycle Management & Cleanup

## 1. Implement background job for server monitoring
## 2. Add 5-minute timeout policy for all servers (VPS resource constraints)
## 3. Publish shutdown events to RedPanda
## 4. Implement graceful server shutdown logic in Worker
## 5. Add server auto-start capability
## 6. Implement Docker state reconciliation job
## 7. Add orphaned container detection and cleanup
## 8. Create cleanup procedures for stopped servers
## 9. Add resource usage monitoring
## 10. Update README with lifecycle policies

