---
name: java-development
description: "Trigger: writing, reviewing, or refactoring Java code, language idioms, static analysis, or Java code smells."
---

# Java Development

## Overview

Use this skill to produce idiomatic, maintainable Java. Favor clarity, immutability, expressive types, small methods, and modern language features when they improve readability.

For detailed rules and examples, read [references/java-guidelines.md](references/java-guidelines.md).

## Core Practices

- Use records for DTOs and immutable data carriers.
- Prefer immutable objects: `final` fields, unmodifiable collections, `List.of`, `Map.of`, and `Stream.toList`.
- Use `var` only when the type is obvious from the right-hand side.
- Use pattern matching for `instanceof` and `switch` when it simplifies branching.
- Use streams, lambdas, and method references for clear collection transformations.
- Prefer `Optional<T>` for possibly absent return values; do not use it as a blanket replacement for all nullable inputs or fields.
- Use `Objects.requireNonNull` and `Objects.equals` where they make null handling explicit.
- Keep names expressive and conventional: `UpperCamelCase` types, `lowerCamelCase` methods/variables, `UPPER_SNAKE_CASE` constants, lowercase packages.

## Static Analysis

When project setup is in scope, ask whether the team wants static analysis. Prefer SonarQube/SonarCloud plus IDE SonarLint when available. If Sonar is unavailable or declined, fall back to project-appropriate tools such as SpotBugs, PMD, or Checkstyle.

When fixing issues, use analyzer output as concrete evidence and reference rule keys when they are available. Do not block ordinary implementation work just because static analysis is not configured.

## Review Checklist

- Could a data-only class be a record?
- Are objects and collections immutable where practical?
- Is `var` used only where it improves readability?
- Are object equality checks using `.equals` or `Objects.equals` rather than `==`?
- Are resources closed with try-with-resources?
- Are nulls avoided or handled deliberately?
- Are methods small enough to read without heavy nesting?
- Are long parameter lists grouped into a value object or builder when appropriate?
- Are repeated literals named as constants or enums?
- Is dead code removed?

## Build And Verification

After code changes, run the project's normal verification path when feasible:

- Maven: `mvn clean install`
- Gradle: `./gradlew build`
- Include tests and configured static analysis in the verification path when available.

## Output Style

When reviewing, lead with concrete bugs, smells, and maintainability risks. When editing, prefer small improvements that fit the local code style over broad rewrites.
