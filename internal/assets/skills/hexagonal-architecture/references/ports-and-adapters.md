# Ports and Adapters Reference

## Core Objective

Hexagonal Architecture protects the domain from external concerns. The center contains business meaning; the edges translate between that center and the outside world.

The dependency rule:

```text
interfaces -> application -> domain
infrastructure -> domain
domain -X-> application/interfaces/infrastructure
application -X-> concrete infrastructure adapters
```

## Layer Contracts

### Domain

The domain contains the business model and the abstractions it needs to express business behavior.

Belongs here:

- aggregates
- entities
- value objects
- domain services
- business invariants
- domain events
- ports for persistence or external systems, when the port is expressed in domain language

Allowed dependencies:

- language standard library
- other domain types

Forbidden:

- framework annotations
- ORM or database APIs
- HTTP clients or messaging clients
- interface-layer DTOs
- infrastructure persistence models

### Application

The application layer implements use cases and orchestrates domain operations. It coordinates flow, but should not own core business rules.

Belongs here:

- use case services or handlers
- commands and queries
- transaction boundaries
- authorization coordination
- calls to domain ports
- application-specific result mapping

Allowed dependencies:

- domain model and ports
- transaction, clock, identity, or security abstractions

Forbidden:

- SQL details
- concrete HTTP clients
- framework-specific persistence entities
- direct dependency on infrastructure adapter classes

### Infrastructure

Infrastructure implements ports and owns technical integration details.

Belongs here:

- repository implementations
- ORM entities
- database queries
- external service clients
- message broker adapters
- filesystem, cache, email, search, or payment adapters
- mappings between domain and external representations

Allowed dependencies:

- domain models and ports
- frameworks and libraries

Forbidden:

- business rules that belong in the domain
- use case decisions that belong in application

### Interfaces

Interfaces are entry points into the system. They translate external protocols into application calls and translate application results back to the outside world.

Belongs here:

- REST controllers
- GraphQL resolvers
- message handlers
- CLI commands
- UI controllers
- request and response DTOs
- transport validation and mapping

Allowed dependencies:

- application use cases
- transport-specific frameworks

Forbidden:

- direct repository access
- domain mutation logic
- business decisions
- persistence or external client details

## Good Placement Examples

- Validate whether an order can be confirmed: `domain`
- Execute the place order workflow: `application`
- Persist an order through JPA, SQL, MongoDB, or a file: `infrastructure`
- Call a remote payment provider: `infrastructure`
- Define the payment gateway capability needed by the domain/use case: `domain` or `application` port, depending on the project convention
- Map an HTTP request to a command: `interfaces`
- Convert a persistence entity into an aggregate: `infrastructure`

## Port Design

Ports should be shaped by the domain or use case need, not by the adapter technology.

Prefer:

```kotlin
interface OrderRepository {
    fun save(order: Order): Order
    fun findById(id: OrderId): Order?
}
```

Avoid:

```kotlin
interface OrderJpaRepositoryLikePort {
    fun findByStatusSql(status: String): List<OrderEntity>
}
```

Good ports use domain names, domain IDs, domain results, and behavior-oriented operations. They do not expose persistence entities, SQL, HTTP payloads, framework types, or adapter exceptions.

## Mapper Guidance

Use mappers at boundaries:

- Interface mapper: request DTO to command/query; application result to response DTO.
- Infrastructure mapper: persistence/external representation to domain model; domain model to persistence/external representation.

Avoid mapping inside the domain. Avoid passing DTOs through multiple layers to save typing; that shortcut leaks policy and makes later changes harder.

## Anti-Patterns

- Domain model annotated as a database entity.
- Application service imports a concrete repository adapter.
- Controller calls a repository and mutates an aggregate directly.
- Repository returns ORM entities through a domain port.
- External API response object appears in domain logic.
- Shared utility package becomes a hidden dependency bridge between layers.
- Business rules are duplicated in controllers, repositories, and domain objects.

## Minimal Review Checklist

- Domain is framework-free.
- Ports are defined in the inner layers and implemented by adapters.
- Use cases orchestrate but do not contain core business rules.
- Interface layer only handles transport concerns.
- Infrastructure owns technical details and maps them to domain concepts.
- No dependency rule violations were introduced.
- Tests cover the layer where behavior changed.
