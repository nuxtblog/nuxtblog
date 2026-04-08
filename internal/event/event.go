// Package event provides a lightweight in-process publish/subscribe bus.
//
// Register handlers with [On] (synchronous) or [OnAsync] (goroutine), then
// publish with [Emit].  Handlers are isolated from panics; async errors are
// only logged so they never affect the request path.
package event

import "context"

// Handler is called when a matching event is emitted.
//
// Returning a non-nil error from a synchronous handler propagates to the
// [Emit] caller.  Errors from asynchronous handlers are only logged.
type Handler func(ctx context.Context, e Event) error

// Event is the envelope delivered to every handler.
type Event struct {
	// Name identifies the event, e.g. "post.created".
	Name string
	// Payload carries the typed event data.
	// Use a type assertion to access the concrete type (see payload package).
	Payload any
}

// Default is the package-level bus used by [On], [OnAsync], and [Emit].
var Default = newBus()

// On registers h to run synchronously when an event named name is emitted.
// Sync handlers block [Emit] and can propagate errors to the caller.
func On(name string, h Handler) {
	Default.on(name, h, false)
}

// OnAsync registers h to run in a separate goroutine when name is emitted.
// Errors are logged but do not affect the emitting caller.
func OnAsync(name string, h Handler) {
	Default.on(name, h, true)
}

// Emit fires all handlers registered for name.
//
// Synchronous handlers are called in registration order; the first error is
// returned and remaining sync handlers are still executed.
// Asynchronous handlers are launched as goroutines and their errors are logged.
func Emit(ctx context.Context, name string, payload any) error {
	return Default.emit(ctx, Event{Name: name, Payload: payload})
}
