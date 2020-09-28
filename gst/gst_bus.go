package gst

// #include "gst.go.h"
import "C"

import (
	"sync"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

// Bus is a Go wrapper around a GstBus. It provides convenience methods for
// popping messages from the queue.
type Bus struct {
	*Object

	msgChannels []chan *Message
	mux         sync.Mutex
}

// Instance returns the underlying GstBus instance.
func (b *Bus) Instance() *C.GstBus { return C.toGstBus(b.Unsafe()) }

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

// PopMessage attempts to pop a message from the bus. It returns nil if none are available.
// The message should be unreffed after usage.
func (b *Bus) PopMessage(timeout int) *Message {
	if b.Instance() == nil {
		return nil
	}
	dur := time.Duration(timeout) * time.Second
	cTimeout := C.GstClockTime(dur.Nanoseconds())
	msg := C.gst_bus_timed_pop(b.Instance(), cTimeout)
	if msg == nil {
		return nil
	}
	return wrapMessage(msg)
}

// BlockPopMessage blocks until a message is available on the bus and then returns it.
// This function can return nil if the bus is closed. The message should be unreffed
// after usage.
func (b *Bus) BlockPopMessage() *Message {
	for {
		if b.Instance() == nil {
			return nil
		}
		msg := b.PopMessage(1)
		if msg == nil {
			continue
		}
		return msg
	}
}

// BusWatchFunc is a go representation of a GstBusFunc. It takes a message as a single argument
// and returns a bool value for whether to continue processing messages or not. There is no need to unref
// the message unless addtional references are placed on it during processing.
type BusWatchFunc func(msg *Message) bool

//export goBusFunc
func goBusFunc(bus *C.GstBus, cMsg *C.GstMessage, userData C.gpointer) C.gboolean {
	// wrap the message
	msg := wrapMessage(cMsg)

	// retrieve the ptr to the function
	ptr := unsafe.Pointer(userData)
	funcIface := gopointer.Restore(ptr)
	busFunc, ok := funcIface.(BusWatchFunc)
	if !ok {
		gopointer.Unref(ptr)
		return gboolean(false)
	}

	// run the call back
	if cont := busFunc(msg); !cont {
		gopointer.Unref(ptr)
		return gboolean(false)
	}

	return gboolean(true)
}

// AddWatch adds a watch to the default MainContext for messages emitted on this bus.
// This function is used to receive asynchronous messages in the main loop. There can
// only be a single bus watch per bus, you must remove it before you can set a new one.
// It is safe to unref the Bus after setting this watch, since the watch itself will take
// it's own reference to the Bus.
//
// The watch can be removed either by returning false from the function or using RemoveWatch.
// A MainLoop must be running for bus watches to work.
//
// The return value reflects whether the watch was successfully added. False is returned if there
// is already a function registered.
func (b *Bus) AddWatch(busFunc BusWatchFunc) bool {
	fPtr := gopointer.Save(busFunc)
	return gobool(
		C.int(C.gst_bus_add_watch(
			b.Instance(),
			C.GstBusFunc(C.cgoBusFunc),
			(C.gpointer)(unsafe.Pointer(fPtr)),
		)),
	)
}

// RemoveWatch will remove any watches installed on the bus. This can also be accomplished
// by returning false from a previously installed function.
//
// The function returns false if there was no watch on the bus.
func (b *Bus) RemoveWatch() bool {
	return gobool(C.gst_bus_remove_watch(b.Instance()))
}
