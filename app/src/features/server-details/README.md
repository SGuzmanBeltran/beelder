# Hexagonal Architecture - Server Details Feature

## 🏗️ Architecture Overview

This feature follows **Hexagonal Architecture (Ports & Adapters)** pattern to achieve:
- ✅ Framework independence
- ✅ Testability
- ✅ Flexibility
- ✅ Maintainability

## 📁 Directory Structure

```
features/server-details/
├── domain/                 # CORE - Business Logic (Framework Independent)
│   ├── models/            # Domain entities
│   ├── ports/             # Interfaces (contracts)
│   └── services/          # Pure business logic
│
├── infrastructure/        # ADAPTERS - External Dependencies
│   ├── api/              # HTTP implementation
│   └── store/            # State management
│
├── application/           # USE CASES - Orchestration Layer
│   └── hooks/            # React hooks that connect UI to domain
│
└── presentation/          # UI LAYER - React Components
    └── components/       # UI components
```

## 🔄 Data Flow

```
UI Component
    ↓
Application Hook (use-server-details)
    ↓
Domain Service (ServerService)
    ↓
Port Interface (ServerRepository)
    ↓
Adapter (ServerApiAdapter)
    ↓
External API
```

## 📝 Layer Responsibilities

### **Domain Layer** (Core)
- **Models**: Pure TypeScript interfaces/types
- **Ports**: Interface definitions (contracts)
- **Services**: Business logic (NO framework dependencies)

**Rules:**
- ❌ No React imports
- ❌ No Axios/HTTP libraries
- ❌ No Zustand/Redux
- ✅ Pure TypeScript only
- ✅ 100% testable without mocks

### **Infrastructure Layer** (Adapters)
- **API Adapter**: Implements ports using HTTP/Axios
- **Store Adapter**: Implements ports using Zustand

**Rules:**
- ✅ Implements port interfaces
- ✅ Handles external communication
- ✅ Maps external data to domain models

### **Application Layer** (Use Cases)
- **Hooks**: Orchestrate domain + infrastructure
- **DTOs**: Data transfer objects

**Rules:**
- ✅ Can use React hooks
- ✅ Connects UI to domain
- ✅ Handles side effects

### **Presentation Layer** (UI)
- **Components**: React components
- **Pages**: Page-level components

**Rules:**
- ✅ Only UI concerns
- ✅ Uses application hooks
- ❌ No direct API calls
- ❌ No business logic

## 🧪 Testing Benefits

### Domain Layer Tests (Pure)
```typescript
describe('ServerService', () => {
  it('should prevent starting an already running server', async () => {
    const mockRepo = createMockRepository();
    const service = new ServerService(mockRepo);
    
    await expect(service.startServer('running-server'))
      .rejects.toThrow('Server is already running');
  });
});
```

### Infrastructure Tests (Integration)
```typescript
describe('ServerApiAdapter', () => {
  it('should fetch server by id', async () => {
    // Test with real API or mock HTTP
  });
});
```

### UI Tests (Component)
```typescript
describe('ServerDetails', () => {
  it('should display server name', () => {
    // Test with mocked hook
  });
});
```

## 🔄 Swapping Implementations

### Example: Replace Axios with Fetch
```typescript
// Just create a new adapter!
export class ServerFetchAdapter implements ServerRepository {
  async getById(id: string): Promise<Server> {
    const response = await fetch(`${API_URL}/server/${id}`);
    return response.json();
  }
  // ... implement other methods
}

// Update dependency injection
const repository = new ServerFetchAdapter(); // ← Changed here only!
const service = new ServerService(repository);
```

### Example: Add Caching Layer
```typescript
export class ServerCachedAdapter implements ServerRepository {
  constructor(
    private readonly apiAdapter: ServerApiAdapter,
    private readonly cache: Cache
  ) {}

  async getById(id: string): Promise<Server> {
    const cached = this.cache.get(id);
    if (cached) return cached;
    
    const server = await this.apiAdapter.getById(id);
    this.cache.set(id, server);
    return server;
  }
}
```

## 💡 Key Benefits

1. **Business Logic is Portable**
   - Can reuse in mobile app, CLI, desktop app
   - No React dependency in core logic

2. **Easy to Test**
   - Test business rules without React
   - Mock adapters easily
   - Fast unit tests

3. **Flexible**
   - Swap HTTP library
   - Change state management
   - Add caching/logging without touching core

4. **Maintainable**
   - Clear separation of concerns
   - Easy to find code
   - Changes are localized

## 🚀 Next Steps

To fully adopt this architecture:

1. **Move existing hooks** to application layer
2. **Extract API calls** to adapters
3. **Create domain services** for business rules
4. **Add dependency injection** container
5. **Write tests** for each layer

## 📚 Resources

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Ports & Adapters Pattern](https://herbertograca.com/2017/09/14/ports-adapters-architecture/)
