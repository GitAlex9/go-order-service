package events

import (
	"context"
	"errors"
	"testing"

	domainevents "github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

type testEvent struct {
	name string
}

func (e testEvent) EventName() string { return e.name }

type recordingHandler struct {
	received []any
	err      error
}

func (h *recordingHandler) Handle(ctx context.Context, event any) error {
	h.received = append(h.received, event)
	return h.err
}

type notNamedEvent struct{}

func TestDispatcher_Dispatch_CallsRegisteredHandler(t *testing.T) {
	registry := NewRegistry()
	handler := &recordingHandler{}
	registry.Register("test.event", handler)

	dispatcher := NewDispatcher(registry, logger.New("test"))

	evt := testEvent{name: "test.event"}
	dispatcher.Dispatch(context.Background(), []any{evt})

	if got, want := len(handler.received), 1; got != want {
		t.Fatalf("handler received %d events, want %d", got, want)
	}
	if got, want := handler.received[0], any(evt); got != want {
		t.Errorf("handler received = %v, want %v", got, want)
	}
}

func TestDispatcher_Dispatch_MultipleHandlersForSameEvent(t *testing.T) {
	registry := NewRegistry()
	handlerA := &recordingHandler{}
	handlerB := &recordingHandler{}
	registry.Register("test.event", handlerA)
	registry.Register("test.event", handlerB)

	dispatcher := NewDispatcher(registry, logger.New("test"))
	dispatcher.Dispatch(context.Background(), []any{testEvent{name: "test.event"}})

	if got := len(handlerA.received); got != 1 {
		t.Errorf("handlerA received %d events, want 1", got)
	}
	if got := len(handlerB.received); got != 1 {
		t.Errorf("handlerB received %d events, want 1", got)
	}
}

func TestDispatcher_Dispatch_NoHandlerRegistered_DoesNotPanic(t *testing.T) {
	registry := NewRegistry()
	dispatcher := NewDispatcher(registry, logger.New("test"))

	dispatcher.Dispatch(context.Background(), []any{testEvent{name: "unregistered.event"}})
	// Se chegou até aqui sem panic, o teste passa.
}

func TestDispatcher_Dispatch_EventNotNamed_IsSkipped(t *testing.T) {
	registry := NewRegistry()
	dispatcher := NewDispatcher(registry, logger.New("test"))

	dispatcher.Dispatch(context.Background(), []any{notNamedEvent{}})
}

func TestDispatcher_Dispatch_HandlerError_DoesNotStopOtherHandlers(t *testing.T) {
	registry := NewRegistry()

	failingHandler := &recordingHandler{err: errors.New("handler failed")}
	succeedingHandler := &recordingHandler{}

	registry.Register("test.event", failingHandler)
	registry.Register("test.event", succeedingHandler)

	dispatcher := NewDispatcher(registry, logger.New("test"))
	dispatcher.Dispatch(context.Background(), []any{testEvent{name: "test.event"}})

	if got := len(succeedingHandler.received); got != 1 {
		t.Errorf("succeedingHandler received %d events, want 1 (não deveria ser impedido pela falha do outro handler)", got)
	}
}

func TestDispatcher_Dispatch_MultipleEventsInOneCall(t *testing.T) {
	registry := NewRegistry()
	handler := &recordingHandler{}
	registry.Register("test.event", handler)

	dispatcher := NewDispatcher(registry, logger.New("test"))

	dispatcher.Dispatch(context.Background(), []any{
		testEvent{name: "test.event"},
		testEvent{name: "test.event"},
		testEvent{name: "other.event"},
	})

	if got, want := len(handler.received), 2; got != want {
		t.Errorf("handler received %d events, want %d", got, want)
	}
}

var _ domainevents.Dispatcher = (*dispatcher)(nil)
