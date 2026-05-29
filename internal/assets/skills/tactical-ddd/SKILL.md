---
name: tactical-ddd
description: "Trigger: designing or reviewing domain models, aggregates, value objects, domain services, or tactical DDD boundaries."
---

# Tactical DDD

## Overview

Use this skill to make domain code express business behavior directly. Keep the domain layer free of framework concerns, place invariants close to the state they protect, and use services only when behavior genuinely spans objects or layers.

For detailed patterns and examples, read [references/tactical-ddd-reference.md](references/tactical-ddd-reference.md).

## Working Loop

1. Identify the domain concept being modeled and name the bounded context if it matters.
2. Classify each object as an aggregate root, entity, value object, domain service, repository port, application service, or infrastructure adapter.
3. Move business rules into aggregates, value objects, or domain services; leave orchestration, transactions, authorization, and mapping at the application/infrastructure boundary.
4. Enforce invariants during creation and every state transition.
5. Replace primitives with meaningful domain types where they carry business semantics.
6. Check dependency direction: infrastructure may depend on domain, but domain must not depend on frameworks or persistence details.
7. Review tests against domain behavior, especially invalid states and disallowed transitions.

## Modeling Rules

- Aggregate roots own consistency boundaries, child entities, value objects, and all rules needed to keep the aggregate valid.
- Entities have stable identity; equality is by identity, not by matching attributes.
- Value objects have no identity; make them immutable and validate them at construction.
- Factory methods or controlled constructors should prevent invalid initial state.
- State-changing methods should validate preconditions before returning or applying a new state.
- Sealed types or explicit variants are preferable to stringly typed domain categories.
- Repositories are domain ports that expose aggregate operations, not database details.
- Domain services are stateless and behavior-named; use them only for logic that does not belong naturally to one aggregate.
- Application or infrastructure services orchestrate use cases, transactions, security, mapping, and persistence; they should not implement business rules.

## Language Guidance

- In Kotlin, prefer immutable data, copy-based transitions, sealed interfaces/classes for variants, and extension functions for collection operations over aggregate lists.
- In Java, use focused domain services when collection-level aggregate behavior cannot be expressed as extension functions.
- In both languages, prefer domain-specific IDs and value types such as `OrderId`, `Email`, `Money`, or `CustomerId` over raw `String`, `Long`, or `BigDecimal` when semantics matter.

## Review Checklist

- Is the domain model behavior-rich, or are services manipulating passive data?
- Can invalid objects be constructed directly?
- Can callers mutate aggregate internals and bypass invariants?
- Do domain classes import persistence, web, DI, logging, or framework annotations?
- Are repositories returning aggregates and domain types rather than persistence rows, query builders, or ORM entities?
- Are application services coordinating domain calls instead of deciding business validity themselves?
- Are domain service names behavior-specific, such as `MoneyTransferService`, instead of vague names like `Manager`, `Processor`, or `Helper`?
- Are errors and results expressed in domain language rather than hidden behind generic exceptions?

## Output Style

When reviewing code, lead with concrete findings and file references. When designing, propose a small domain model first, then explain aggregate boundaries, invariants, and service placement. Prefer short, tactical guidance over abstract DDD theory.
