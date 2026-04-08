package event

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

type entry struct {
	handler Handler
	async   bool
}

type bus struct {
	mu       sync.RWMutex
	handlers map[string][]entry
}

func newBus() *bus {
	return &bus{handlers: make(map[string][]entry)}
}

func (b *bus) on(name string, h Handler, async bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = append(b.handlers[name], entry{handler: h, async: async})
}

func (b *bus) emit(ctx context.Context, e Event) (firstErr error) {
	b.mu.RLock()
	entries := make([]entry, len(b.handlers[e.Name]))
	copy(entries, b.handlers[e.Name])
	b.mu.RUnlock()

	for _, en := range entries {
		if en.async {
			en := en // capture loop variable
			go func() {
				defer func() {
					if r := recover(); r != nil {
						g.Log().Warningf(ctx, "[event] async handler panic on %q: %v", e.Name, r)
					}
				}()
				if err := en.handler(ctx, e); err != nil {
					g.Log().Warningf(ctx, "[event] async handler error on %q: %v", e.Name, err)
				}
			}()
		} else {
			func() {
				defer func() {
					if r := recover(); r != nil {
						if firstErr == nil {
							firstErr = fmt.Errorf("handler panic on %q: %v", e.Name, r)
						}
					}
				}()
				if err := en.handler(ctx, e); err != nil && firstErr == nil {
					firstErr = err
				}
			}()
		}
	}
	return
}
