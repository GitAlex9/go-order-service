package events

import "context"

type Handler interface {
	Handle(ctx context.Context, event any) error
}

type Named interface {
	EventName() string
}

type Dispatcher interface {
	Dispatch(ctx context.Context, events []any)
}
