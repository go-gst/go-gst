package gst

import (
	"context"
	"iter"
	"runtime"

	"github.com/go-gst/go-glib/pkg/core/userdata"
	"github.com/go-gst/go-gst/pkg/gst/internal/channel"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
// extern GstBusSyncReply _gogst_gst1_BusSyncHandler(GstBus*, GstMessage*, gpointer);
// extern void destroyUserdata(gpointer);
import "C"

type BusExtManual interface {
	// Messages adds a watch to the bus. This is a convenience function that
	// actually attaches a sync handler to the bus. This way you don't need to create a
	// main loop.
	//
	// Since this is a sync handler, make sure to handle the messages as fast as
	// possible. Otherwise your pipeline may block.
	Messages(context.Context) iter.Seq[*Message]

	// SetSyncHandler wraps gst_bus_set_sync_handler
	//
	// The function takes the following parameters:
	//
	// 	- fn BusSyncHandler (nullable): The handler function to install
	//
	// Sets the synchronous handler on the bus. The function will be called
	// every time a new message is posted on the bus. Note that the function
	// will be called in the same thread context as the posting object. This
	// function is usually only called by the creator of the bus. Applications
	// should handle messages asynchronously using the gst_bus watch and poll
	// functions.
	//
	// Before 1.16.3 it was not possible to replace an existing handler and
	// clearing an existing handler with %NULL was not thread-safe.
	SetSyncHandler(BusSyncHandler)
}

func (bus *BusInstance) Messages(ctx context.Context) iter.Seq[*Message] {
	ctx, cancel := context.WithCancel(ctx)

	messages := channel.NewGrowable[*Message]()

	go func() {
		<-ctx.Done()
		messages.Close()
	}()

	bus.SetSyncHandler(func(bus Bus, message *Message) BusSyncReply {
		messages.Send(message)
		return BusDrop
	})

	return func(yield func(*Message) bool) {

		for {
			message, ok := messages.Receive()
			if !ok {
				break
			}

			if !yield(message) {
				break
			}
		}

		bus.SetSyncHandler(nil)
		cancel()
	}
}

// SetSyncHandler wraps gst_bus_set_sync_handler
//
// The function takes the following parameters:
//
//   - fn BusSyncHandler (nullable): The handler function to install
//
// Sets the synchronous handler on the bus. The function will be called
// every time a new message is posted on the bus. Note that the function
// will be called in the same thread context as the posting object. This
// function is usually only called by the creator of the bus. Applications
// should handle messages asynchronously using the gst_bus watch and poll
// functions.
//
// Before 1.16.3 it was not possible to replace an existing handler and
// clearing an existing handler with %NULL was not thread-safe.
func (bus *BusInstance) SetSyncHandler(fn BusSyncHandler) {
	var carg0 *C.GstBus           // in, none, converted
	var carg1 C.GstBusSyncHandler // callback, scope: notified, closure: carg2, destroy: carg3, nullable
	var carg2 C.gpointer          // implicit
	var carg3 C.GDestroyNotify    // implicit

	carg0 = (*C.GstBus)(UnsafeBusToGlibNone(bus))
	if fn != nil {
		carg1 = (*[0]byte)(C._gogst_gst1_BusSyncHandler)
		carg2 = C.gpointer(userdata.Register(fn))
		carg3 = (C.GDestroyNotify)((*[0]byte)(C.destroyUserdata))
	}

	C.gst_bus_set_sync_handler(carg0, carg1, carg2, carg3)
	runtime.KeepAlive(bus)
	runtime.KeepAlive(fn)
}

// BusSyncHandler wraps GstBusSyncHandler
//
// The function takes the following parameters:
//
//   - bus Bus: the #GstBus that sent the message
//   - message *Message: the #GstMessage
//
// The function returns the following values:
//
//   - goret BusSyncReply
//
// Handler will be invoked synchronously, when a new message has been injected
// into the bus. This function is mostly used internally. Only one sync handler
// can be attached to a given bus.
type BusSyncHandler func(bus Bus, message *Message) (goret BusSyncReply)
