// Package notify provides an extensible multi-channel notification dispatcher.
// Channels (InApp, Email, SMS, …) register themselves via init() and are
// fan-out called by Dispatch on every notification event.
package notify

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Message carries all the data needed to render a notification in any channel.
type Message struct {
	UserID      int64
	Type        string // follow | like | comment | reply | mention | system
	SubType     string
	ActorID     *int64
	ActorName   string
	ActorAvatar string
	ObjectType  string
	ObjectID    *int64
	ObjectTitle string
	ObjectLink  string
	Content     string
}

// Channel is implemented by every notification delivery backend.
type Channel interface {
	Name() string
	// Enabled returns true when the channel is configured and should be used.
	Enabled(ctx context.Context) bool
	Send(ctx context.Context, msg Message) error
}

var channels []Channel

// Register adds a channel to the dispatcher. Call from channel init() functions.
func Register(ch Channel) {
	channels = append(channels, ch)
}

// Dispatch fans the message out to all enabled channels.
// Call it in a goroutine from business logic; it runs with a 30-second timeout.
func Dispatch(ctx context.Context, msg Message) {
	tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for _, ch := range channels {
		if ch.Enabled(tctx) {
			if err := ch.Send(tctx, msg); err != nil {
				g.Log().Warningf(ctx, "[notify] channel %s error: %v", ch.Name(), err)
			}
		}
	}
}
