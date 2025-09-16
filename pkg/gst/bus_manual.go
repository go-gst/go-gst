package gst

import "iter"

type BusExtManual interface {
	// Messages adds a watch to the bus. This is a convenience function that
	// actually attaches a sync handler to the bus. This way you don't need to create a
	// main loop.
	//
	// Since this is a sync handler, make sure to handle the messages as fast as
	// possible. Otherwise your pipeline may block.
	Messages() iter.Seq[*Message]
}

func (bus *BusInstance) Messages() iter.Seq[*Message] {
	messages := make(chan *Message, 20) // arbitrary cap to not block instantly

	bus.SetSyncHandler(func(bus Bus, message *Message) BusSyncReply {
		messages <- message.Copy()
		return BusDrop
	})

	return func(yield func(*Message) bool) {
		for message := range messages {
			if !yield(message) {
				bus.SetSyncHandler(nil)
				close(messages)
				return
			}
		}
	}
}
