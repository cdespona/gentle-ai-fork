# Java Guidelines Reference

## Modern Java Features

### Records

Use records for classes whose primary purpose is carrying immutable data, such as DTOs, query results, configuration values, or simple value objects.

Prefer:

```java
public record UserDto(UserId id, String name, Email email) {}
```

Avoid hand-written data carriers with boilerplate getters, equality, and constructors unless behavior or framework constraints require a class.

### Pattern Matching

Use pattern matching for `instanceof` and `switch` when it removes casts and clarifies branching.

```java
if (value instanceof Order order) {
    return order.total();
}
```

### Type Inference

Use `var` for locals only when the type is obvious:

```java
var users = userRepository.findActiveUsers();
```

Avoid `var` when the expression hides the type or reduces readability:

```java
var result = service.execute(input);
```

## Immutability

Favor immutable objects:

- make fields `final`
- avoid setters unless mutation is part of the model
- prefer immutable collections
- return new instances instead of mutating state where practical

For fixed collections, use `List.of`, `Set.of`, or `Map.of`. For stream results, prefer `Stream.toList()` when an unmodifiable list is acceptable.

## Streams And Lambdas

Use streams when they make collection transformations clearer:

```java
var activeUsers = users.stream()
    .filter(User::isActive)
    .map(User::summary)
    .toList();
```

Avoid streams for complex control flow, heavy side effects, or logic that is clearer as a named method.

## Null Handling

Avoid returning `null` for absent values; prefer `Optional<T>` for return types where absence is expected.

Use:

- `Objects.requireNonNull` for mandatory dependencies and parameters
- `Objects.equals` for null-safe equality
- explicit validation for invalid input

Avoid:

- returning `null` from repository or lookup methods when `Optional<T>` would clarify absence
- accepting `Optional<T>` as a parameter unless the local project already uses that convention
- storing `Optional<T>` in fields or DTOs without a strong reason

## Naming

Follow conventional Java style:

- classes and interfaces: `UpperCamelCase`
- methods and variables: `lowerCamelCase`
- constants: `UPPER_SNAKE_CASE`
- packages: lowercase

Use nouns for types and verbs for behavior. Avoid abbreviations, Hungarian notation, and vague suffixes such as `Manager`, `Helper`, or `Util` unless the concept is genuinely generic.

## Common Bug Patterns

- Resource leaks: close files, sockets, and streams with try-with-resources.
- Equality mistakes: use `.equals` or `Objects.equals` for object equality.
- Redundant casts: fix generics or type flow instead of casting around the issue.
- Always-true or always-false conditions: treat them as likely bugs or dead code.
- Null-sensitive comparisons: prefer constants on the left or `Objects.equals` where helpful.

## Code Smells

- Long parameter list: group related values into a value object or builder.
- Large method: extract behavior into named methods.
- High cognitive complexity: reduce nesting with guard clauses, polymorphism, or strategy objects.
- Duplicated literals: promote meaningful literals to constants or enums.
- Dead code: remove unused assignments, branches, and variables.
- Magic numbers or strings: name them when the literal carries business meaning.

## Static Analysis Setup Notes

When asked to set up static analysis:

1. Prefer SonarQube or SonarCloud with SonarLint in the IDE and scanner execution in CI.
2. Create a project key and store scanner tokens in CI secrets.
3. If Sonar is declined or blocked, use SpotBugs, PMD, or Checkstyle as project-appropriate fallbacks.
4. If analysis is configured but failing, check binding/token, CI scanner execution, and IDE/plugin setup before switching tools.
