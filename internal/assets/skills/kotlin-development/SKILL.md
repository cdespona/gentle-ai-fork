---
name: kotlin-development
description: "Trigger: writing, reviewing, or refactoring Kotlin code, null-safety, sealed types, coroutines, or Kotlin code smells."
---

# Kotlin Development

## Overview

Use this skill to write Kotlin as Kotlin, not Java translated into Kotlin syntax. Prefer immutability, expressions, null-safety, sealed models, extension functions, standard library collection operations, and structured concurrency.

For detailed rules and examples, read [references/kotlin-guidelines.md](references/kotlin-guidelines.md).

## Core Practices

- Prefer `val` and immutable collections.
- Use `data class` for data carriers and value-like models.
- Use `@JvmInline value class` for type-safe primitive wrappers such as IDs and typed strings when appropriate.
- Avoid `!!`; use safe calls, Elvis, `let`, `requireNotNull`, or explicit validation.
- Use `when` as an expression, especially with sealed classes/interfaces for exhaustive handling.
- Prefer extension functions and top-level functions over utility objects.
- Use default parameters and named arguments instead of overload-heavy or builder-style APIs.
- Use collection operations for transformations; use loops when they are clearer or more efficient.
- Use coroutines with structured concurrency; avoid `GlobalScope`.

## Static Analysis

When project setup is in scope, ask whether the team wants static analysis. Prefer Detekt for Kotlin-specific rules, optionally paired with SonarQube/SonarCloud. Use ktlint for formatting/style when appropriate.

When fixing issues, use Detekt/Sonar/IDE output as concrete evidence and reference rule keys when available. Do not block ordinary implementation work just because analysis is not configured.

## Review Checklist

- Is this Kotlin idiomatic, or Java style in Kotlin syntax?
- Could mutable state become `val`, immutable collections, or copy-based updates?
- Can null handling avoid `!!`?
- Would a sealed type make variants explicit and exhaustive?
- Would an extension function make type-specific behavior clearer?
- Are scope functions used for their intended purpose rather than cleverness?
- Are coroutines scoped and cancellable?
- Are Java interop annotations used only when needed?
- Are long parameter lists handled with defaults, named arguments, or data classes?

## Build And Verification

After code changes, run the project's normal verification path when feasible:

- Gradle: `./gradlew build`
- Maven: `mvn clean install`
- Detekt: `./gradlew detekt` when configured
- ktlint: `./gradlew ktlintCheck` when configured

## Output Style

When reviewing, lead with correctness, null-safety, coroutine, and maintainability risks. When editing, prefer small idiomatic changes that match the existing project style.
