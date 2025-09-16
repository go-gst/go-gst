package gst

import (
	"context"
	"iter"
	"sync/atomic"
)

type BusExtManual interface {
	// Messages adds a watch to the bus. This is a convenience function that
	// actually attaches a sync handler to the bus. This way you don't need to create a
	// main loop.
	//
	// Since this is a sync handler, make sure to handle the messages as fast as
	// possible. Otherwise your pipeline may block.
	Messages(context.Context) iter.Seq[*Message]
}

func (bus *BusInstance) Messages(ctx context.Context) iter.Seq[*Message] {
	messages := make(chan *Message, 20) // arbitrary cap to not block instantly

	c := atomic.Pointer[chan *Message]{}
	c.Store(&messages)

	bus.SetSyncHandler(func(bus Bus, message *Message) BusSyncReply {
		messages := c.Load()
		if messages == nil {
			return BusDrop
		}
		*messages <- message.Copy()
		return BusDrop
	})

	return func(yield func(*Message) bool) {

		defer func() {
			bus.SetSyncHandler(nil)
			c.Store(nil)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messages:
				if !yield(message) {
					return
				}
			}
		}
	}
}
