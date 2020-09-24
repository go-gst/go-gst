package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"

import "sync"

// Bus is a Go wrapper around a GstBus. It provides convenience methods for
// popping messages from the queue.
type Bus struct {
	*Object

	msgChannels []chan *Message
	mux         sync.Mutex
}

// Instance returns the underlying GstBus instance.
func (b *Bus) Instance() *C.GstBus { return C.toGstBus(b.unsafe()) }

func (b *Bus) deliverMessages() {
	for {
		msg := b.BlockPopMessage()
		if msg == nil {
			return
		}
		b.mux.Lock()
		for _, ch := range b.msgChannels {
			ch <- msg.Ref()
		}
		b.mux.Unlock()
		msg.Unref()
	}
}

// MessageChan returns a new channel to listen for messages asynchronously. Messages
// should be unreffed after each usage. Messages are delivered to channels in the
// order in which this function was called.
//
// While a message is being delivered to created channels, there is a lock on creating
// new ones.
func (b *Bus) MessageChan() chan *Message {
	b.mux.Lock()
	defer b.mux.Unlock()
	ch := make(chan *Message)
	b.msgChannels = append(b.msgChannels, ch)
	if len(b.msgChannels) == 1 {
		go b.deliverMessages()
	}
	return ch
}

// BlockPopMessage blocks until a message is available on the bus and then returns it.
// This function can return nil if the bus is closed. The message should be unreffed
// after usage.
func (b *Bus) BlockPopMessage() *Message {
	// I think this is ok since no other main loop is running
	msg := C.gst_bus_poll(
		(*C.GstBus)(b.Instance()),
		C.GST_MESSAGE_ANY,
		C.GST_CLOCK_TIME_NONE,
	)
	if msg == nil {
		return nil
	}
	return wrapMessage(msg)
}

func wrapBus(bus *C.GstBus) *Bus {
	return &Bus{
		Object:      wrapObject(&bus.object),
		msgChannels: make([]chan *Message, 0),
	}
}
