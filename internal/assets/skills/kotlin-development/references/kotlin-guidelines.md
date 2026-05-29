# Kotlin Guidelines Reference

## Core Philosophy

Write Kotlin, not Java in Kotlin syntax.

Embrace:

- properties instead of getters/setters
- null-safe operators
- data classes
- sealed classes and sealed interfaces
- extension functions
- default parameters and named arguments
- standard library collection functions
- coroutines and Flow where asynchronous behavior is needed

Avoid:

- unnecessary builders
- utility classes full of static-like functions
- callback-heavy APIs when suspend functions fit
- mutable state by default
- `!!` as routine null handling

## Data And Immutability

Use `data class` for DTOs, value objects, and data carriers.

```kotlin
data class User(
    val id: UserId,
    val name: String,
    val email: Email
)
```

Prefer `val` over `var` and immutable collections (`List`, `Set`, `Map`) by default. Use `copy()` for modified instances.

## Null Safety

Avoid `!!` unless there is a narrow, justified interop reason. Prefer:

- safe calls: `user?.name`
- Elvis: `name ?: "Unknown"`
- `let` for null-safe transformations
- `requireNotNull` or `checkNotNull` with messages when absence is invalid

## Scope Functions

Use scope functions deliberately:

- `let`: null-safe transformations or scoped value transformation
- `run`: object-scoped computation returning a result
- `with`: grouping calls on an existing object
- `apply`: object configuration returning the object
- `also`: side effects such as logging or tracing while returning the object

Avoid chaining many scope functions if it hides control flow.

## Extension Functions

Prefer extension functions over utility classes for focused operations on a type.

```kotlin
fun String.isValidEmail(): Boolean =
    contains("@") && contains(".")

fun List<Order>.totalRevenue(): Money =
    fold(Money.ZERO) { acc, order -> acc + order.total }
```

Keep extensions cohesive and close to the domain where possible.

## Sealed Types And `when`

Use sealed interfaces/classes for restricted variants. Combine them with `when` expressions for exhaustive handling.

```kotlin
sealed interface PaymentMethod {
    data class CreditCard(val number: String) : PaymentMethod
    data class BankTransfer(val iban: String) : PaymentMethod
    data object Cash : PaymentMethod
}
```

Avoid stringly typed variants when the compiler can enforce the allowed cases.

## Collections And Sequences

Use `map`, `filter`, `fold`, `groupBy`, `partition`, and related operations for transformations. Use `Sequence` for large collections or long chains where intermediate collections matter.

Use `forEach` only for side effects. Use loops when they are simpler, performance-sensitive, or involve early exits that collection chains obscure.

## Parameters And APIs

Use default parameters instead of overloads when Kotlin callers are primary. Use named arguments at call sites with multiple same-typed parameters or boolean flags.

For Java callers, consider `@JvmOverloads` only when the overloads materially improve interop.

## JVM Interop Annotations

- `@JvmInline value class`: type-safe primitive wrappers with low runtime overhead.
- `@JvmStatic`: companion object members that should appear static to Java.
- `@JvmOverloads`: generate overloads for Java callers.
- `@JvmField`: expose a field directly; use sparingly.
- `@JvmName`: resolve JVM signature conflicts or improve Java-facing names.
- `@Suppress`: use sparingly and explain why when the reason is not obvious.
- `@OptIn`: mark deliberate experimental API usage.
- `@Deprecated`: include a migration path with `replaceWith` where possible.

## Coroutines And Flow

Use structured concurrency:

- prefer `coroutineScope` or `supervisorScope`
- use `async`/`await` for parallel results
- use `launch` for scoped side effects
- avoid `GlobalScope`
- inject dispatchers when that improves testability
- use `withContext(Dispatchers.IO)` for blocking I/O and `Dispatchers.Default` for CPU work

Use `Flow` for cold streams, `StateFlow` for observable state, and `SharedFlow` for shared events. Handle errors with operators such as `catch` where appropriate.

Long-running loops should be cancellation-aware with `isActive` or `ensureActive()`.

## Naming

Follow Kotlin conventions:

- classes/interfaces: `UpperCamelCase`
- functions/properties: `lowerCamelCase`
- constants: `UPPER_SNAKE_CASE`
- packages: lowercase with dots
- backing properties: `_state` with public `state`

## Bug Patterns And Smells

- Resource management: use `use` for closeable resources.
- Equality: `==` is structural, `===` is referential.
- Nullable collections: distinguish `List<T>?` from `List<T?>`.
- Mutable state: prefer immutable values and copy-based updates.
- Callback nesting: prefer suspend functions or Flow when appropriate.
- Complex Elvis chains: extract named functions.
- Primitive obsession: use value classes for typed IDs and domain primitives.
- Builder pattern: prefer default parameters and named arguments unless Java interop or object construction complexity justifies a builder.

## Static Analysis Setup Notes

When asked to set up static analysis:

1. Prefer Detekt with project-specific rules.
2. Add SonarQube/SonarCloud for broader coverage when useful.
3. Use ktlint for formatting/style if the project wants automated style checks.
4. If configured tools fail, check project binding/tokens, CI execution, and IDE/plugin setup before switching tools.
