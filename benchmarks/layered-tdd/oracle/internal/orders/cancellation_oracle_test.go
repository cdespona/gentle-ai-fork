package orders

import (
	"context"
	"errors"
	"testing"
)

type countingRepository struct {
	order     Order
	getErr    error
	saveCount int
}

func (r *countingRepository) Get(context.Context, string) (Order, error) {
	return r.order, r.getErr
}

func (r *countingRepository) Save(_ context.Context, order Order) error {
	r.order = order
	r.saveCount++
	return nil
}

type countingNotifier struct {
	count int
}

func (n *countingNotifier) OrderCancelled(context.Context, string) error {
	n.count++
	return nil
}

func TestBenchmarkOracleDomainCancellation(t *testing.T) {
	pending := Order{ID: "order-1", Status: StatusPending}
	cancelled, changed, err := pending.Cancel()
	if err != nil || !changed || cancelled.Status != StatusCancelled {
		t.Fatalf("first cancellation: order=%#v changed=%v err=%v", cancelled, changed, err)
	}

	again, changed, err := cancelled.Cancel()
	if err != nil || changed || again != cancelled {
		t.Fatalf("repeated cancellation: order=%#v changed=%v err=%v", again, changed, err)
	}

	_, changed, err = (Order{ID: "order-2", Status: StatusFulfilled}).Cancel()
	if !errors.Is(err, ErrCannotCancel) || changed {
		t.Fatalf("fulfilled cancellation: changed=%v err=%v", changed, err)
	}
}

func TestBenchmarkOracleServiceCancellationIsIdempotent(t *testing.T) {
	repository := &countingRepository{order: Order{ID: "order-1", Status: StatusPending}}
	notifier := &countingNotifier{}
	service := NewService(repository, notifier)

	first, err := service.Cancel(context.Background(), "order-1")
	if err != nil || first.Status != StatusCancelled {
		t.Fatalf("first cancellation: order=%#v err=%v", first, err)
	}
	second, err := service.Cancel(context.Background(), "order-1")
	if err != nil || second.Status != StatusCancelled {
		t.Fatalf("second cancellation: order=%#v err=%v", second, err)
	}
	if repository.saveCount != 1 {
		t.Fatalf("save count: got %d, want 1", repository.saveCount)
	}
	if notifier.count != 1 {
		t.Fatalf("notification count: got %d, want 1", notifier.count)
	}
}

