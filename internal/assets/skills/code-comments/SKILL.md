---
name: code-comments
description: "Trigger: writing, editing, or reviewing code comments, docstrings, API docs, TODOs, or stale/noisy comments."
---

# Code Comments

## Overview

Use this skill to keep code self-explanatory and comments rare, durable, and useful. Prefer clearer structure, naming, and smaller functions before adding explanatory text.

For examples and detailed rules, read [references/comment-guidelines.md](references/comment-guidelines.md).

## Decision Flow

Before adding or keeping a comment, apply this order:

1. Refactor first: make structure reveal intent.
2. Rename second: use names that communicate the idea.
3. Delete noise: remove comments that explain what the code already says.
4. Keep only durable rationale: add a brief comment only when the reason cannot be expressed in code.

If refactoring, naming, or deleting solves the clarity issue, do not add a comment.

## Allowed Comments

Add comments only when they explain non-obvious rationale, such as:

- a non-trivial algorithm choice
- a business, legal, or regulatory constraint not inferable from code
- a workaround for an external bug or limitation
- an intentional concurrency, performance, or security trade-off
- a public API contract where callers need behavioral guarantees

## Forbidden Comments

Do not add or keep comments that:

- restate obvious behavior
- narrate line-by-line actions
- duplicate type, variable, or method names
- describe stale behavior after code changed
- compensate for unclear code that could reasonably be refactored
- use anonymous or vague TODO/FIXME placeholders

## Public API Documentation

Use the language's native API documentation style for public APIs when the contract is not obvious from the signature: docstrings in Python, JSDoc/TSDoc in JavaScript/TypeScript, KDoc/Javadoc on the JVM, XML documentation in C#, rustdoc in Rust, Go doc comments in Go, and equivalent conventions elsewhere.

Document purpose, preconditions, postconditions, important side effects, and exceptions or errors that are part of the contract.

Avoid doc comments on internal or private code unless project policy explicitly requires them.

## TODO and FIXME Governance

Use action comments sparingly and include ownership plus traceable context:

- `TODO(owner, yyyy-mm-dd): short action`
- `FIXME(owner, ticket): short issue`

Remove resolved items immediately. Convert long-lived TODOs into tracked tickets.

## Maintenance Checklist

When touching a file:

- remove redundant comments
- update comments that no longer match behavior
- replace explanatory comments with clearer code where feasible
- preserve domain rationale comments that still add non-obvious context

## Output Style

When reviewing, lead with stale, misleading, or noisy comments and include file references. When editing, keep comments short, specific, and tied to durable constraints. Ask for confirmation only when a proposed comment encodes an ambiguous business decision.
