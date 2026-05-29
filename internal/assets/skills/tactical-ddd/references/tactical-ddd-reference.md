# Tactical DDD Reference

## Building Blocks

### Aggregate Roots

Aggregate roots are top-level entities that own a consistency boundary. They contain child entities and value objects, expose behavior-oriented methods, and enforce all rules required to keep the aggregate valid.

Good aggregate roots:
- expose behavior such as `addItem`, `confirm`, `withdraw`, or `applyDiscount`
- validate invariants before changing state
- hide internal mutable collections
- keep persistent identity stable
- prevent invalid construction through factories or controlled constructors

Avoid anemic aggregates that only expose fields while services mutate them from the outside.

### Entities

Entities have identity that persists over time. Equality should be based on the ID, not on every property. Entities may change, but changes should occur through aggregate methods that preserve the aggregate invariants.

Kotlin sketch:

```kotlin
data class OrderItem(
    val id: OrderItemId,
    val productId: ProductId,
    val quantity: Int
) {
    override fun equals(other: Any?): Boolean =
        this === other || (other is OrderItem && id == other.id)

    override fun hashCode(): Int = id.hashCode()
}
```

### Value Objects

Value objects have no identity. Equality is structural, they should be immutable, and they should reject invalid values when created.

Kotlin sketch:

```kotlin
data class Money(
    val amount: BigDecimal,
    val currency: Currency
) {
    init {
        require(amount >= BigDecimal.ZERO) { "Amount cannot be negative" }
    }

    fun add(other: Money): Money {
        require(currency == other.currency) { "Cannot add different currencies" }
        return Money(amount + other.amount, currency)
    }
}
```

## Tactical Practices

### Use Domain-Specific Types

Prefer explicit types such as `UserId`, `OrderId`, `Email`, `Money`, and `CustomerId` over raw primitives. This prevents accidental mixing and gives validation a natural home.

```kotlin
data class Email(val value: String) {
    init {
        require(value.contains("@")) { "Invalid email format" }
    }
}

data class User(val id: UserId, val email: Email)
```

### Prevent Infrastructure Leakage

The domain layer should not import framework, ORM, web, messaging, or DI concerns. Avoid annotations such as `@Entity`, `@Table`, `@Service`, `@Repository`, or `@RestController` in domain objects. Put persistence entities and mappers in infrastructure.

Dependency direction should be:

```text
interfaces -> application services -> domain <- infrastructure adapters
```

Infrastructure implements domain ports; the domain does not know the infrastructure exists.

### Put Business Logic in the Domain

Business decisions belong in aggregates, value objects, or domain services. Application services should orchestrate calls and handle boundaries, not decide validity.

```kotlin
data class ShoppingCart(
    val id: CartId,
    private val items: List<CartItem>,
    val customerId: CustomerId
) {
    fun applyDiscount(discount: Discount): ShoppingCart {
        require(discount.isValid()) { "Discount has expired" }
        val newTotal = calculateTotal().minus(discount.amount)
        require(newTotal >= Money.ZERO) { "Discount exceeds cart total" }
        return copy(items = items.map { it.applyDiscount(discount) })
    }

    private fun calculateTotal(): Money =
        items.fold(Money.ZERO) { acc, item -> acc.add(item.subtotal) }
}
```

### Enforce Invariants at Creation and Transition

Invariants must always hold. Validate during construction and before each state change.

```kotlin
data class Order private constructor(
    val id: OrderId,
    val items: List<OrderItem>,
    val status: OrderStatus
) {
    companion object {
        fun create(): Order =
            Order(OrderId.generate(), emptyList(), OrderStatus.DRAFT)
    }

    fun addItem(item: OrderItem): Order {
        require(status == OrderStatus.DRAFT) { "Cannot add items to a confirmed order" }
        require(items.none { it.productId == item.productId }) { "Item already in order" }
        return copy(items = items + item)
    }

    fun confirm(): Order {
        require(items.isNotEmpty()) { "Cannot confirm an empty order" }
        require(status == OrderStatus.DRAFT) { "Only draft orders can be confirmed" }
        return copy(status = OrderStatus.CONFIRMED)
    }
}
```

Avoid public constructors that allow invalid states, mutable aggregate internals, or partial validation that only occurs at some transitions.

### Use Explicit Variants

Use sealed types, enums, or dedicated classes for real domain variants. Avoid strings that require callers to remember valid values.

```kotlin
sealed interface PaymentMethod {
    data class CreditCard(val cardNumber: String, val expiryDate: String) : PaymentMethod
    data class BankTransfer(val accountNumber: String) : PaymentMethod
    data class PayPal(val email: Email) : PaymentMethod
}
```

## Domain Services

Use domain services sparingly, when behavior truly spans multiple aggregates or does not naturally belong to one aggregate.

Good domain services:
- are stateless
- are named after behavior, such as `MoneyTransferService`
- operate on aggregate roots and domain types
- avoid duplicating logic already present in aggregates

Avoid services such as `OrderTotalCalculationService` when the behavior clearly belongs on `Order`.

For Kotlin, prefer extension functions for aggregate collection operations:

```kotlin
fun List<Order>.totalRevenue(): Money =
    fold(Money.ZERO) { acc, order -> acc.add(order.calculateTotal()) }

fun List<Order>.confirmedOrders(): List<Order> =
    filter { it.status == OrderStatus.CONFIRMED }
```

For Java, use focused domain services for collection operations when needed:

```java
public class OrderCollectionService {
    public Money calculateTotalRevenue(List<Order> orders) {
        return orders.stream()
            .map(Order::calculateTotal)
            .reduce(Money.ZERO, Money::add);
    }
}
```

## Application and Infrastructure Services

Application services coordinate a use case. They may handle transactions, authorization, logging, request-to-domain mapping, persistence, and integration with external systems.

They should:
- load aggregates through repository ports
- call domain methods to make business decisions
- persist updated aggregates
- translate results for the interface layer

They should not:
- mutate persistence rows directly to bypass aggregate behavior
- implement business rules that belong in the domain
- expose database details to callers
- hide domain errors behind vague failures

## Repository Pattern

Define repository interfaces in the domain or application boundary as ports. Implement them in infrastructure.

```kotlin
interface OrderRepository {
    fun save(order: Order): Order
    fun findById(id: OrderId): Order?
    fun findByCustomer(customerId: CustomerId): List<Order>
}
```

Infrastructure adapters can use ORM entities internally, but they should map persistence models to domain models before returning them.

## Common Smells

- Primitive obsession: important concepts are represented by raw strings or numbers.
- Anemic model: aggregates are passive data while services contain all decisions.
- Framework pollution: domain objects carry persistence or web annotations.
- Invalid construction: public constructors allow impossible business states.
- Leaky repositories: callers receive database rows, ORM entities, or query builders instead of aggregates.
- Vague services: names like `Manager`, `Processor`, `Helper`, or `DataService`.
- Aggregate boundary violations: unrelated entities are updated together without a clear consistency rule.
- Mutable internals: callers can change aggregate collections without invoking domain behavior.
