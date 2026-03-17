# Shipment Tracking Microservice

A gRPC microservice for managing shipments and tracking their lifecycle through valid status transitions. Implements hexagonal architecture to enforce clean separation between domain logic, business use cases, and infrastructure.

## Quick Start

### Prerequisites

- Go 1.25.4+
- PostgreSQL 12+ (or Docker)
- golang-migrate CLI (or Docker)

### Setup Option 1: Using Docker Compose (Recommended)

Everything (PostgreSQL, migrations, app) runs in containers:

```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f

# Test the service (should be running on localhost:8080)
grpcurl -plaintext localhost:8080 list

# Stop all services
docker-compose down

# Clean up (including data)
docker-compose down -v
```

Services started:
- **postgres**: PostgreSQL on `localhost:5432`
- **migrate**: Applies database migrations automatically
- **app**: gRPC server on `localhost:8080`

### Setup Option 2: Local Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Start PostgreSQL with Docker:**
   ```bash
   docker run --name shipment-postgres \
     -e POSTGRES_USER=user \
     -e POSTGRES_PASSWORD=pass \
     -e POSTGRES_DB=shipment \
     -p 5432:5432 \
     -d postgres:latest
   ```

3. **Run database migrations:**
   ```bash
   migrate -path internal/migrations -database "postgres://user:pass@localhost:5432/shipment?sslmode=disable" up
   ```

4. **Start the gRPC server:**
   ```bash
   go run cmd/main.go
   ```

   Server listens on `localhost:8080`

## Running Tests

### All tests:
```bash
go test ./...
```

### Specific package:
```bash
go test ./internal/domain
go test ./internal/application
```

### With coverage:
```bash
go test -cover ./...
```

## Docker

### Understanding docker-compose.yml

The compose file orchestrates three services:

1. **postgres** - PostgreSQL database
   - Starts first
   - Health check ensures readiness before other services
   - Data persisted in named volume `postgres_data`

2. **migrate** - Database migrations
   - Waits for postgres to be healthy
   - Applies all migrations from `internal/migrations/`
   - Exits after successful migration

3. **app** - gRPC application
   - Builds from Dockerfile
   - Waits for migrations to complete
   - Runs on port 8080
   - Uses postgres service name in connection string

### Common Docker Commands

```bash
# Start all services in background
docker-compose up -d

# Start with logs visible
docker-compose up

# View logs
docker-compose logs app
docker-compose logs postgres
docker-compose logs migrate

# Restart specific service
docker-compose restart app

# Run command in app container
docker-compose exec app sh

# Stop services (keeps data)
docker-compose stop

# Stop and remove services
docker-compose down

# Remove everything including volumes
docker-compose down -v

# View running containers
docker-compose ps

# Rebuild image
docker-compose build
docker-compose up -d
```

### Rebuilding the App Image

After code changes:
```bash
docker-compose build app
docker-compose up -d app
```

### Running Tests in Docker

Build and run tests in container:
```bash
docker build -t shipment-test .
docker run --rm shipment-test go test ./...
```

## Testing the Service

### Using grpcurl

```bash
# Create shipment
grpcurl -plaintext -d '{
  "id": "ship-001",
  "reference": 123,
  "origin": "New York",
  "destination": "Los Angeles",
  "cost": 1500.50,
  "driver_revenue": 300.25
}' localhost:8080 shipment.ShipmentService/CreateShipment

# Get shipment
grpcurl -plaintext -d '{"id": "ship-001"}' localhost:8080 shipment.ShipmentService/GetShipment

# Add status event
grpcurl -plaintext -d '{
  "id": "ship-001",
  "status": 1
}' localhost:8080 shipment.ShipmentService/AddEvent

# Get history
grpcurl -plaintext -d '{"id": "ship-001"}' localhost:8080 shipment.ShipmentService/GetShipmentHistory
```

### Using Postman

1. New → gRPC Request
2. Server URL: `localhost:8080`
3. Import proto file: `api/shipment.proto`
4. Select method and fill request JSON
5. Click "Invoke"

## Project Structure

```
shipmentVektor/
├── api/
│   └── shipment.proto                # gRPC service contract
│
├── cmd/
│   └── main.go                       # Application entry point
│
├── internal/
│   ├── adapters/
│   │   ├── grpc/                     # inbound adapter
│   │   │   └── handler.go
│   │   └── postgres/                 # outbound adapter
│   │       └── repo.go
│   │
│   ├── application/
│   │   ├── service.go                # use case orchestration
│   │   └── service_test.go
│   │
│   ├── domain/
│   │   ├── shipment.go               # shipment aggregate
│   │   ├── event.go                  # domain events
│   │   ├── status.go                 # status definitions
│   │   └── shipment_test.go
│   │
│   ├── ports/
│   │   ├── inbound/
│   │   │   └── usecase.go
│   │   └── outbound/
│   │       └── repo.go
│   │
│   ├── migrations/
│   │   ├── 000001_init.up.sql
│   │   └── 000001_init.down.sql
│   │
│   └── pkg/
│       ├── config.go
│       └── db.go
│
├── README.md
├── go.mod
├── go.sum
└── .env
```

## Architecture

### Hexagonal Architecture (Ports & Adapters)

The service uses hexagonal architecture to isolate domain logic from infrastructure:

```
gRPC Request
    ↓
[gRPC Handler] (Inbound Adapter)
    ↓
[Port: ShipmentUseCase]
    ↓
[ShipmentService] (Application/Use Case)
    ↓
[Domain Layer] 
  - Shipment aggregate
  - Status validation
  - Event logic
    ↓
[Port: ShipmentRepository]
    ↓
[PostgreSQL Adapter] (Outbound Adapter)
    ↓
Database
```

### Core Design Principles

**1. Isolated Domain Logic**

The domain layer (`internal/domain/`) contains pure business logic with no external dependencies:
- Shipment aggregate enforces valid state transitions
- Status enum defines allowed statuses
- Event structure models state changes

Domain logic can be tested without gRPC, database, or any framework.

**2. Dependency Inversion**

Service layer depends on interfaces (ports), not concrete implementations:
- `ShipmentUseCase` interface defined in `ports/inbound/`
- `ShipmentRepository` interface defined in `ports/outbound/`

Adapters (gRPC handler, PostgreSQL repository) implement these interfaces. This enables:
- Independent testing via mock implementations
- Swapping infrastructure without changing domain logic

**3. Stateful Aggregate Pattern**

`Shipment` is implemented as an aggregate:
- Encapsulates internal state (status, events)
- Enforces invariants through `CanUpdate()` and `ApplyEvent()`
- Never leaves invalid state

State transitions follow a single valid path:
```
PENDING → SHIPPED → ON_THE_WAY → DELIVERED
```

**4. Event Sourcing (Partial)**

Shipment lifecycle is modeled as a sequence of events:
- Initial state: PENDING event created with shipment
- Each status change: recorded as new event
- Current state: derived from latest event
- Full history: all events persisted in `shipment_events` table

## Business Logic & Assumptions

### Shipment Lifecycle

A shipment follows this state machine:

```
PENDING (initial)
  ↓
SHIPPED (picked up)
  ↓
ON_THE_WAY (in transit)
  ↓
DELIVERED (final)
```

Transitions are **strictly validated**:
- Can only move forward in this sequence
- Cannot skip stages
- Cannot move backward
- Cannot transition to invalid states

### Key Assumptions

**1. Immutable Events**
Once recorded, events cannot be modified or deleted. This provides an audit trail.

**2. Single Origin of Truth**
Current shipment status always reflects the latest valid event. If the database is consulted, the status is recalculated from events.

**3. No Concurrent Updates**
The service does not implement optimistic/pessimistic locking. Concurrent updates to the same shipment may cause issues. In production, would need:
- Version fields with optimistic locking
- Or database-level constraints
- Or event versioning

**4. Synchronous Processing**
All operations complete synchronously. No async event processing or message queues.

**5. Single Service Instance**
No distributed locking for shipments. Works correctly with single instance.

### Enforced Business Rules

```go
// Cannot transition to invalid states
if !shipment.CanUpdate(newStatus) {
    return error("invalid status transition")
}

// Status only advances in one direction
StatusPending → StatusShipped (only valid next state)
StatusShipped → StatusOnTheWay (only valid next state)
StatusOnTheWay → StatusDelivered (only valid next state)
```

## Test Coverage

Tests focus on business behavior and domain logic:

### Domain Tests (`internal/domain/shipment_test.go`)
- Valid status transitions
- Invalid status transitions  
- Sequential state progression
- Initial state correctness
- Transition validation logic

### Application Tests (`internal/application/service_test.go`)
- Shipment creation
- Shipment retrieval
- Event addition and status updates
- History retrieval
- Error handling (missing shipments)

All tests use in-memory mock repository, testing application and domain logic independently from the database.

## Database Schema

Two tables maintain shipment data and event history:

**shipments** - current state
```sql
CREATE TABLE shipments (
    id TEXT PRIMARY KEY,
    reference INT NOT NULL,
    origin TEXT NOT NULL,
    destination TEXT NOT NULL,
    status TEXT NOT NULL,
    cost NUMERIC NOT NULL,
    driver_revenue NUMERIC NOT NULL
);
```

**shipment_events** - immutable event history
```sql
CREATE TABLE shipment_events (
    id SERIAL PRIMARY KEY,
    shipment_id TEXT NOT NULL REFERENCES shipments(id),
    status TEXT NOT NULL,
    timestamp BIGINT NOT NULL
);
```

### Migrations

Migrations versioned using golang-migrate:

```bash
# Apply
migrate -path internal/migrations -database "postgres://..." up

# Rollback
migrate -path internal/migrations -database "postgres://..." down
```

## Configuration

Environment variables in `.env`:

```env
DB_CONN="postgres://user:pass@localhost:5432/shipment?sslmode=disable"
PORT=8080
```

## Design Decisions Explained

### 1. Why Aggregate Pattern for Shipment?

**Decision:** Model Shipment as an aggregate with internal consistency rules.

**Rationale:** 
- Shipments are complex entities with multiple related pieces (events, status, history)
- Invalid states must be prevented at the object level
- Aggregates provide clear boundaries and encapsulation
- Status validation happens in one place, not scattered across layers

**Alternative Considered:** Anemic data model with validation in service layer. Rejected because business logic should belong in domain.

### 2. Why Event Sourcing for History?

**Decision:** Store all status changes as immutable events rather than just current state.

**Rationale:**
- Provides complete audit trail of shipment lifecycle
- Enables history query without separate logging system
- Can reconstruct past states if needed
- Clear separation: events are immutable facts, shipment state is current projection

**Alternative Considered:** Only store current status. Rejected because history requirement needs full event chain.

### 3. Why Hexagonal Architecture?

**Decision:** Separate domain, application, and infrastructure layers with interface boundaries.

**Rationale:**
- Domain logic can be tested without external dependencies
- Easy to swap database implementations (PostgreSQL → MySQL → etc.)
- Easy to test with mock repositories
- Clear responsibility boundaries
- Infrastructure changes don't affect business logic

**How It Works Here:**
- Domain knows nothing about gRPC or SQL
- Application (use cases) knows about domain logic but not frameworks
- Adapters know about domain and their specific technology (gRPC protocol, SQL database)
- Ports (interfaces) define contracts between layers

### 4. Why Interface-Based Repositories?

**Decision:** Use `ShipmentRepository` interface; implement with PostgreSQL adapter.

**Rationale:**
- Tests use MockRepository for isolated testing
- Can swap database without changing application code
- Dependency inversion: service depends on interface, not implementation
- Clear contract of what data operations are needed

**Benefit for This Project:**
- Service tests in `service_test.go` use `MockRepository`
- Domain tests don't need repository at all
- Real PostgreSQL implementation isolated in `adapters/postgres/`

### 5. Why Status as Constants?

**Decision:** Define statuses as Go constants, not database enum or arbitrary strings.

**Rationale:**
- Type-safe: can't accidentally use invalid status
- IDE autocomplete works
- Compiler catches typos
- Status transitions validated at compile time in `CanUpdate()`

### 6. Why Raw SQL Instead of ORM?

**Decision:** Use raw SQL in repository layer.

**Rationale:**
- Explicit control over queries
- No ORM overhead or magic
- SQL is clear and auditable
- Easy to add indices and optimize

**Trade-off:** More boilerplate scanning rows, but code clarity wins for this project size.

## Extensions & Maintenance

### Adding a New Status

To add a new status (e.g., `RETURNED`):

1. Add to `domain/status.go`:
   ```go
   StatusReturned Status = "returned"
   ```

2. Update `domain/shipment.go` in `CanUpdate()`:
   ```go
   StatusDelivered: {StatusReturned}, // Can return after delivery
   ```

3. Update proto file if needed

4. Tests automatically catch missing cases

### Adding a New Operation

New operations follow this pattern:

1. Define RPC in `api/shipment.proto`
2. Implement in `application/service.go` 
3. Add handler in `adapters/grpc/handler.go`
4. Add corresponding repository method if needed
5. Add tests

### Replacing PostgreSQL

To use MySQL instead:

1. Update `internal/pkg/db.go` to use MySQL driver
2. Implement new `adapters/mysql/repo.go` with `ShipmentRepository` interface
3. Update migrations to MySQL syntax
4. Only adapter changes needed; domain and application layers untouched

## Key Files

- `api/shipment.proto` - gRPC contract definition
- `internal/domain/shipment.go` - Core business logic, state machine
- `internal/application/service.go` - Use case orchestration
- `internal/adapters/grpc/handler.go` - gRPC request/response handling
- `internal/adapters/postgres/repo.go` - Data persistence
- `internal/ports/` - Interface contracts
- `cmd/main.go` - Service initialization and startup



