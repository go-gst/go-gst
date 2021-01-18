package gst

/*
#include "gst.go.h"

extern GstBusSyncReply   goBusSyncHandler (GstBus * bus, GstMessage * message, gpointer user_data);
extern gboolean          goBusFunc        (GstBus * bus, GstMessage * msg, gpointer user_data);

gboolean cgoBusFunc (GstBus * bus, GstMessage * msg, gpointer user_data)
{
	return goBusFunc(bus, msg, user_data);
}

GstBusSyncReply cgoBusSyncHandler (GstBus * bus, GstMessage * message, gpointer user_data)
{
	return goBusSyncHandler(bus, message, user_data);
}

*/
import "C"

import (
	"reflect"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// Bus is a Go wrapper around a GstBus. It provides convenience methods for
// popping messages from the queue.
type Bus struct {
	*Object
}

// NewBus returns a new Bus instance.
//
//   // Example of using the bus instance
//
//   package main
//
//   import (
//       "fmt"
//
//       "github.com/tinyzimmer/go-gst/gst"
//   )
//
//   func main() {
//       gst.Init(nil)
//
//       bus := gst.NewBus()
//       defer bus.Unref()
//
//       elem, err := gst.NewElement("fakesrc")
//       if err != nil {
//           panic(err)
//       }
//       defer elem.Unref()
//
//       bus.Post(gst.NewAsyncStartMessage(elem))
//
//       msg := bus.Pop()
//       defer msg.Unref()
//
//       fmt.Println(msg)
//   }
//
//   // > [fakesrc0] ASYNC-START - Async task started
//
func NewBus() *Bus {
	return FromGstBusUnsafeFull(unsafe.Pointer(C.gst_bus_new()))
}

// FromGstBusUnsafeNone wraps the given unsafe.Pointer in a bus. It takes a ref on the bus and sets
// a runtime finalizer on it.
func FromGstBusUnsafeNone(bus unsafe.Pointer) *Bus { return wrapBus(glib.TransferNone(bus)) }

// FromGstBusUnsafeFull wraps the given unsafe.Pointer in a bus. It does not increase the ref count
// and places a runtime finalizer on the instance.
func FromGstBusUnsafeFull(bus unsafe.Pointer) *Bus { return wrapBus(glib.TransferFull(bus)) }

// Instance returns the underlying GstBus instance.
func (b *Bus) Instance() *C.GstBus { return C.toGstBus(b.Unsafe()) }

// AddSignalWatch adds a bus signal watch to the default main context with the default priority (%G_PRIORITY_DEFAULT).
// It is also possible to use a non-default main context set up using g_main_context_push_thread_default (before one
// had to create a bus watch source and attach it to the desired main context 'manually').
//
// After calling this statement, the bus will emit the "message" signal for each message posted on the bus.
//
// This function may be called multiple times. To clean up, the caller is responsible for calling RemoveSignalWatch
// as many times as this function is called.
func (b *Bus) AddSignalWatch() { C.gst_bus_add_signal_watch(b.Instance()) }

// PopMessage attempts to pop a message from the bus. It returns nil if none are available.
// The message should be unreffed after usage.
//
// It is much safer and easier to use the AddWatch or other polling functions. Only use this method if you
// are unable to also run a MainLoop, or for convenience sake.
func (b *Bus) PopMessage(timeout int) *Message {
	return b.TimedPop(time.Duration(timeout) * time.Second)
}

// BlockPopMessage blocks until a message is available on the bus and then returns it.
// This function can return nil if the bus is closed. The message should be unreffed
// after usage.
//
// It is much safer and easier to use the AddWatch or other polling functions. Only use this method if you
// are unable to also run a MainLoop, or for convenience sake.
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

// AddWatch adds a watch to the default MainContext for messages emitted on this bus.
// This function is used to receive asynchronous messages in the main loop. There can
// only be a single bus watch per bus, you must remove it before you can set a new one.
// It is safe to unref the Bus after setting this watch, since the watch itself will take
// it's own reference to the Bus.
//
// The watch can be removed either by returning false from the function or by using RemoveWatch().
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

// CreateWatch creates a watch and returns the GSource to be added to a main loop.
// TODO: the return values from this function should be type casted and the MainLoop
// should offer methods for using the return of this function.
// func (b *Bus) CreateWatch() *C.GSource {
// 	return C.gst_bus_create_watch(b.Instance())
// }

// RemoveWatch will remove any watches installed on the bus. This can also be accomplished
// by returning false from a previously installed function.
//
// The function returns false if there was no watch on the bus.
func (b *Bus) RemoveWatch() bool {
	return gobool(C.gst_bus_remove_watch(b.Instance()))
}

// RemoveSignalWatch removes a signal watch previously added with AddSignalWatch.
func (b *Bus) RemoveSignalWatch() { C.gst_bus_remove_signal_watch(b.Instance()) }

// DisableSyncMessageEmission instructs GStreamer to stop emitting the "sync-message" signal for this bus.
// See EnableSyncMessageEmission for more information.
//
// In the event that multiple pieces of code have called EnableSyncMessageEmission, the sync-message emissions
// will only be stopped after all calls to EnableSyncMessageEmission were "cancelled" by calling this function.
// In this way the semantics are exactly the same as Ref that which calls enable should also call disable.
func (b *Bus) DisableSyncMessageEmission() { C.gst_bus_disable_sync_message_emission(b.Instance()) }

// EnableSyncMessageEmission instructs GStreamer to emit the "sync-message" signal after running the bus's sync handler.
// This function is here so that code can ensure that they can synchronously receive messages without having to affect
// what the bin's sync handler is.
//
// This function may be called multiple times. To clean up, the caller is responsible for calling DisableSyncMessageEmission
// as many times as this function is called.
//
// While this function looks similar to AddSignalWatch, it is not exactly the same -- this function enables *synchronous*
// emission of signals when messages arrive; AddSignalWatch adds an idle callback to pop messages off the bus asynchronously.
// The sync-message signal comes from the thread of whatever object posted the message; the "message" signal is marshalled
// to the main thread via the main loop.
func (b *Bus) EnableSyncMessageEmission() { C.gst_bus_enable_sync_message_emission(b.Instance()) }

// PollFd represents the possible values returned from a GetPollFd. On Windows, there will not be
// a Fd.
type PollFd struct {
	Fd              int
	Events, REvents uint
}

// GetPollFd gets the file descriptor from the bus which can be used to get notified about messages being available with
// functions like g_poll, and allows integration into other event loops based on file descriptors. Whenever a message is
// available, the POLLIN / G_IO_IN event is set.
//
// Warning: NEVER read or write anything to the returned fd but only use it for getting notifications via g_poll or similar
// and then use the normal GstBus API, e.g. PopMessage.
func (b *Bus) GetPollFd() *PollFd {
	var gpollFD C.GPollFD
	C.gst_bus_get_pollfd(b.Instance(), &gpollFD)
	pollFd := &PollFd{
		Events:  uint(gpollFD.events),
		REvents: uint(gpollFD.revents),
	}
	if fd := reflect.ValueOf(&gpollFD).Elem().FieldByName("fd"); fd.IsValid() {
		pollFd.Fd = int(fd.Interface().(C.gint))
	}
	return pollFd
}

// HavePending checks if there are pending messages on the bus that should be handled.
func (b *Bus) HavePending() bool {
	return gobool(C.gst_bus_have_pending(b.Instance()))
}

// Peek peeks the message on the top of the bus' queue. The message will remain on the bus'
// message queue. A reference is returned, and needs to be unreffed by the caller.
func (b *Bus) Peek() *Message {
	msg := C.gst_bus_peek(b.Instance())
	if msg == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(msg))
}

// Poll the bus for messages. Will block while waiting for messages to come. You can specify a maximum
// time to poll with the timeout parameter. If timeout is negative, this function will block indefinitely.
//
// All messages not in events will be popped off the bus and will be ignored. It is not possible to use message
// enums beyond MessageExtended in the events mask.
//
// Because poll is implemented using the "message" signal enabled by AddSignalWatch, calling Poll will cause the
// "message" signal to be emitted for every message that poll sees. Thus a "message" signal handler will see the
// same messages that this function sees -- neither will steal messages from the other.
//
// This function will run a main loop from the default main context when polling.
//
// You should never use this function, since it is pure evil. This is especially true for GUI applications based
// on Gtk+ or Qt, but also for any other non-trivial application that uses the GLib main loop. As this function
// runs a GLib main loop, any callback attached to the default GLib main context may be invoked. This could be
// timeouts, GUI events, I/O events etc.; even if Poll is called with a 0 timeout. Any of these callbacks
//  may do things you do not expect, e.g. destroy the main application window or some other resource; change other
// application state; display a dialog and run another main loop until the user clicks it away. In short, using this
// function may add a lot of complexity to your code through unexpected re-entrancy and unexpected changes to your
// application's state.
//
// For 0 timeouts use gst_bus_pop_filtered instead of this function; for other short timeouts use TimedPopFiltered;
// everything else is better handled by setting up an asynchronous bus watch and doing things from there.
func (b *Bus) Poll(msgTypes MessageType, timeout time.Duration) *Message {
	cTime := C.GstClockTime(timeout.Nanoseconds())
	mType := C.GstMessageType(msgTypes)
	msg := C.gst_bus_poll(b.Instance(), mType, cTime)
	if msg == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(msg))
}

// Pop pops a message from the bus, or returns nil if none are available.
func (b *Bus) Pop() *Message {
	msg := C.gst_bus_pop(b.Instance())
	if msg == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(msg))
}

// PopFiltered gets a message matching type from the bus. Will discard all messages on the bus that do not match type
// and that have been posted before the first message that does match type. If there is no message matching type on the
// bus, all messages will be discarded. It is not possible to use message enums beyond MessageExtended in the events mask.
func (b *Bus) PopFiltered(msgTypes MessageType) *Message {
	mType := C.GstMessageType(msgTypes)
	msg := C.gst_bus_pop_filtered(b.Instance(), mType)
	if msg == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(msg))
}

// Post a new message on the bus. The bus takes ownership of the message.
func (b *Bus) Post(msg *Message) bool {
	return gobool(C.gst_bus_post(b.Instance(), msg.Ref().Instance()))
}

// SetFlushing sets whether to flush out and unref any messages queued in the bus. Releases references to the message origin
// objects. Will flush future messages until SetFlushing sets flushing to FALSE.
func (b *Bus) SetFlushing(flushing bool) { C.gst_bus_set_flushing(b.Instance(), gboolean(flushing)) }

// BusSyncHandler will be invoked synchronously, when a new message has been injected into the bus. This function is mostly
// used internally. Only one sync handler can be attached to a given bus.
//
// If the handler returns BusDrop, it should unref the message, else the message should not be unreffed by the sync handler.
type BusSyncHandler func(msg *Message) BusSyncReply

// SetSyncHandler sets the synchronous handler on the bus. The function will be called every time a new message is posted on the bus.
// Note that the function will be called in the same thread context as the posting object. This function is usually only called by the
// creator of the bus. Applications should handle messages asynchronously using the watch and poll functions.
//
// Currently, destroyNotify funcs are not supported.
func (b *Bus) SetSyncHandler(f BusSyncHandler) {
	ptr := gopointer.Save(f)
	C.gst_bus_set_sync_handler(
		b.Instance(),
		C.GstBusSyncHandler(C.cgoBusSyncHandler),
		(C.gpointer)(unsafe.Pointer(ptr)),
		nil,
	)
}

// TimedPop gets a message from the bus, waiting up to the specified timeout. Unref returned messages after usage.
//
// If timeout is 0, this function behaves like Pop. If timeout is < 0, this function will block forever until a message was posted on the bus.
func (b *Bus) TimedPop(dur time.Duration) *Message {
	var cTime C.GstClockTime
	if dur == ClockTimeNone {
		cTime = C.GstClockTime(gstClockTimeNone)
	} else {
		cTime = C.GstClockTime(dur.Nanoseconds())
	}
	msg := C.gst_bus_timed_pop(b.Instance(), cTime)
	if msg == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(msg))
}

// TimedPopFiltered gets a message from the bus whose type matches the message type mask types, waiting up to the specified timeout
// (and discarding any messages that do not match the mask provided).
//
// If timeout is 0, this function behaves like PopFiltered. If timeout is < 0, this function will block forever until a matching message
// was posted on the bus.
func (b *Bus) TimedPopFiltered(dur time.Duration, msgTypes MessageType) *Message {
	var cTime C.GstClockTime
	if dur == ClockTimeNone {
		cTime = C.GstClockTime(gstClockTimeNone)
	} else {
		cTime = C.GstClockTime(dur.Nanoseconds())
	}
	msg := C.gst_bus_timed_pop_filtered(b.Instance(), cTime, C.GstMessageType(msgTypes))
	if msg == nil {
		return nil
	}
	return FromGstMessageUnsafeFull(unsafe.Pointer(msg))
}
