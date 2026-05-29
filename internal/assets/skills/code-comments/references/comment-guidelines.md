# Comment Guidelines Reference

## Core Principle

Default to no comments. Prefer better names, smaller functions, clearer types, and simpler control flow over explanatory prose.

Comments are valuable when they preserve intent that code cannot express on its own. They are harmful when they repeat the code, drift out of date, or make unclear code look acceptable.

## Mandatory Comment Decision Flow

Before adding any comment:

1. Refactor first. Can a smaller function, extracted type, or simpler branch make intent obvious?
2. Rename second. Can variable, method, or type names communicate the idea?
3. Delete noise. If the comment explains what the code does, remove it.
4. Keep only why. Add a brief comment only if the rationale cannot be expressed in code.

If steps 1 through 3 solve clarity, do not add a comment.

## Good Comment Subjects

Use comments for durable, non-obvious context:

- Non-trivial algorithm choice and why it is appropriate.
- Business, legal, or regulatory constraints not inferable from code.
- Workarounds for external bugs, platform limitations, or library behavior.
- Intentional concurrency, performance, or security trade-offs.
- Public API behavioral contracts that callers need to rely on.

Good examples:

```python
# Required by PCI rule: mask PAN before persistence.
```

```typescript
// WORKAROUND(#4521): remove after upgrading library >= 2.4.0.
```

```go
// Binary search is safe here because events are sorted by timestamp upstream.
```

## Bad Comment Subjects

Remove or avoid comments that:

- restate obvious behavior
- narrate line-by-line actions
- duplicate names from nearby symbols
- describe behavior that changed
- explain unclear code that can reasonably be refactored
- add ceremony to private/internal code without a project requirement

Bad examples:

```python
# Increment count by one.
count += 1
```

```typescript
// User service class.
class UserService {}
```

```kotlin
// Get the active users.
val activeUsers = users.filter { it.active }
```

## Public API Documentation

Use the language's native API documentation style for public APIs when callers need contract clarity beyond the signature.

Common equivalents:

- Python: docstrings, commonly Google, NumPy, or reStructuredText style.
- JavaScript/TypeScript: JSDoc or TSDoc.
- Java/Kotlin: Javadoc or KDoc.
- C#: XML documentation comments.
- Rust: rustdoc comments.
- Go: doc comments beginning with the exported identifier name.

Document:

- purpose and behavioral contract
- preconditions
- postconditions
- important side effects
- thrown exceptions, returned errors, or failure modes when they are part of the contract

Do not add doc comments to private or internal code unless the project explicitly requires it.

Prefer:

```python
def reserve(product_id: ProductId, quantity: Quantity) -> Reservation:
    """Reserve stock for an order line.

    Raises:
        InsufficientStockError: Available stock is lower than quantity.
    """
```

Avoid:

```python
def name() -> str:
    """Gets the name."""
```

## TODO and FIXME Governance

Use TODO and FIXME comments sparingly. Include ownership and traceable context.

Allowed forms:

```python
# TODO(alex, 2026-06-30): replace polling with webhook after partner rollout.
```

```typescript
// FIXME(payments-1234): retry is not idempotent for partial captures.
```

Rules:

- No anonymous TODO or FIXME comments.
- No vague placeholders such as `TODO: cleanup`.
- Remove resolved items immediately.
- Convert long-lived TODOs into tracked tickets.

## Edit-Time Maintenance

Whenever touching a file:

1. Remove comments made redundant by the edit.
2. Update comments that no longer match behavior.
3. Replace explanatory comments with clearer code where feasible.
4. Keep rationale comments if they still explain a real, non-obvious constraint.

## Agent Behavior

Do not add comments as a substitute for unclear code. If clarity requires a large explanatory comment, prefer refactoring. If uncertain whether rationale is durable, improve naming and omit the comment. Ask the user only when the comment would encode an ambiguous business decision.

## Quick Review Checklist

- Does every comment explain why rather than what?
- Are stale comments removed?
- Are obvious comments removed?
- Are public API docs present where contract clarity is needed?
- Do TODO/FIXME comments include owner and traceable context?
- Was code readability improved before comments were added?
