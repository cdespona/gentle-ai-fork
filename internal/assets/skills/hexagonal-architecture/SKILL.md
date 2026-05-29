---
name: hexagonal-architecture
description: "Trigger: placing code across domain/application/infrastructure/interface layers or reviewing ports and adapters boundaries."
---

# Hexagonal Architecture

## Overview

Use this skill to keep changes aligned with Ports and Adapters. Preserve the dependency rule: domain code stays independent, application code orchestrates use cases, infrastructure implements ports, and interfaces translate external input/output.

For detailed layer contracts and review examples, read [references/ports-and-adapters.md](references/ports-and-adapters.md).

## Dependency Rule

- `interfaces -> application -> domain`
- `infrastructure -> domain`
- `domain` must not depend on application, infrastructure, interface, framework, persistence, or transport code
- `application` must depend on ports or abstractions, not concrete infrastructure adapters

## Placement Rules

- Put business concepts, invariants, entities, value objects, aggregates, domain services, and port interfaces in `domain`.
- Put use case orchestration, transaction boundaries, authorization coordination, and command/query handling in `application`.
- Put persistence, external clients, message brokers, framework configuration, adapter implementations, and mapping to external representations in `infrastructure`.
- Put REST, GraphQL, messaging handlers, CLI commands, UI controllers, request DTOs, response DTOs, and transport mapping in `interfaces`.

## Editing Workflow

Before editing:

1. Identify the change type: domain rule, use case orchestration, adapter detail, or transport mapping.
2. Place the code in the correct layer first.
3. Check whether a suitable port already exists before creating a new interface.
4. Choose the smallest change that preserves layer boundaries.

While editing:

- keep domain APIs stable unless the requested behavior requires changing them
- add or update mappers instead of leaking external DTOs into the domain
- introduce a port when a use case needs persistence or external data
- keep business rules out of controllers, repositories, clients, and adapter glue
- prefer narrow, behavior-specific ports over generic infrastructure-shaped interfaces

After editing:

1. Verify no forbidden imports crossed boundaries.
2. Verify business invariants remain in the domain.
3. Verify adapters implement ports, not the reverse.
4. Run relevant tests or build checks for touched modules.

## Common Smells

- Domain objects annotated with framework, ORM, DI, or transport annotations.
- Controllers or handlers making business decisions.
- Repositories returning persistence entities or query details through domain ports.
- Application services calling concrete adapters directly.
- Shared utilities that let layers bypass intended boundaries.
- DTOs from interfaces or infrastructure appearing in domain APIs.

## Review Checklist

- Is the domain framework-free?
- Are ports defined toward the domain/application need rather than the adapter technology?
- Do adapters implement ports?
- Do use cases orchestrate without taking over core business rules?
- Does the interface layer only handle transport concerns and mapping?
- Did the change preserve dependency direction?

## Output Style

When reviewing, lead with boundary violations and concrete file references. When designing or refactoring, explain where each new type belongs and why. Prefer minimal, surgical changes over broad architectural reshuffling.
