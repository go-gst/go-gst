// This is a GStreamer element implemented in Go that uses inbound data on a websocket
// connection as the source for the stream.
//
// In order to build the plugin for use by GStreamer, you can do the following:
//
//     $ go generate
//     $ go build -o libgstwebsocketsrc.so -buildmode c-shared .
//
//
//go:generate gst-plugin-gen
//
// +plugin:Name=websocketsrc
// +plugin:Description=GStreamer Websocket Source
// +plugin:Version=v0.0.1
// +plugin:License=gst.LicenseLGPL
// +plugin:Source=go-gst
// +plugin:Package=examples
// +plugin:Origin=https://github.com/tinyzimmer/go-gst
// +plugin:ReleaseDate=2021-01-10
//
// +element:Name=websocketsrc
// +element:Rank=gst.RankNone
// +element:Impl=websocketSrc
// +element:Subclass=gst.ExtendsElement
package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"golang.org/x/net/websocket"
)

// MaxPayloadSize to accept over websocket connections. Also the size of buffers.
const MaxPayloadSize = 1024

// Defaults //
var (
	DefaultAddress            string = "0.0.0.0"
	DefaultPort               int    = 5000
	DefaultRetrieveRemoteAddr bool   = true
)

func main() {}

// CAT is the log category for the websocketsrc.
var CAT = gst.NewDebugCategory(
	"websocketsrc",
	gst.DebugColorNone,
	"WebsocketSrc Element",
)

var properties = []*gst.ParamSpec{
	gst.NewStringParam(
		"address",
		"Server Address",
		"The address to bind the server to",
		&DefaultAddress,
		gst.ParameterReadWrite,
	),
	gst.NewIntParam(
		"port",
		"Server Port",
		"The port to bind the server to",
		1024, 65535,
		DefaultPort,
		gst.ParameterReadWrite,
	),
	// not implemented yet
	gst.NewBoolParam(
		"retrieve-remote-addr",
		"Retrieve Remote Address",
		"Include the remote client's address in the buffer metadata",
		DefaultRetrieveRemoteAddr,
		gst.ParameterReadWrite,
	),
}

// Internals //

// A private settings struct to hold the values of the above parameters
type settings struct {
	address            string
	port               int
	retrieveRemoteAddr bool
}

// Helper function to retrieve a settings object set to the default values.
func defaultSettings() *settings {
	return &settings{
		address:            DefaultAddress,
		port:               DefaultPort,
		retrieveRemoteAddr: DefaultRetrieveRemoteAddr,
	}
}

// The internal state object
type state struct {
	serverStarted, channelsStarted, sentInitialEvents, sentSegment bool
	server                                                         *http.Server
	srcpad                                                         *gst.Pad
	bufferpool                                                     *gst.BufferPool
	bufferchan                                                     chan []byte
	stopchan                                                       chan struct{}

	mux     sync.Mutex
	connmux sync.Mutex
}

// Base struct definition for the websocket src
type websocketSrc struct {
	settings *settings
	state    *state
}

// prepare verifies the src pad has been added to the element, and then sets up server
// handlers and a buffer pool
func (w *websocketSrc) prepare(elem *gst.Element) error {
	w.state.mux.Lock()
	defer w.state.mux.Unlock()

	// Make sure we have a srcpad
	if w.state.srcpad == nil {
		w.setupSrcPad(elem)
	}

	elem.Log(CAT, gst.LevelDebug, "Creating channels for goroutines")

	// Setup a channel for handling buffers
	w.state.bufferchan = make(chan []byte)
	w.state.stopchan = make(chan struct{})

	elem.Log(CAT, gst.LevelDebug, "Setting up the HTTP server")

	// Setup the HTTP server instance
	w.state.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", w.settings.address, w.settings.port),
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		Handler: &websocket.Server{
			// Don't check the Origin header
			Handshake: func(*websocket.Config, *http.Request) error { return nil },
			Handler: func(conn *websocket.Conn) {
				elem.Log(CAT, gst.LevelInfo, fmt.Sprintf("Received new connection from: %s", conn.Request().RemoteAddr))

				// Only allow a stream from one client at a time
				w.state.connmux.Lock()
				defer w.state.connmux.Unlock()

				conn.PayloadType = websocket.BinaryFrame
				conn.MaxPayloadBytes = MaxPayloadSize

				for {
					// Read the PayloadSize into a bytes slice
					buf := make([]byte, conn.MaxPayloadBytes)
					size, err := conn.Read(buf)
					if err != nil {
						elem.ErrorMessage(gst.DomainStream, gst.StreamErrorFailed, "Error reading bytes from client", err.Error())
						return
					}

					// The goroutine listening for buffers will use the size to determine offsets,
					// So trim the zeroes if we receive a buffer less than the requested size.
					if size < conn.MaxPayloadBytes {
						trimmed := make([]byte, size)
						copy(trimmed, buf)
						buf = trimmed
					}

					// Queue the buffer for processing
					elem.Log(CAT, gst.LevelLog, fmt.Sprintf("Queueing %d bytes for processing", len(buf)))
					w.state.bufferchan <- buf
				}
			},
		},
	}

	elem.Log(CAT, gst.LevelDebug, "Configuring a buffer pool")

	// Configure a buffer pool
	w.state.bufferpool = gst.NewBufferPool()
	cfg := w.state.bufferpool.GetConfig()
	cfg.SetParams(nil, MaxPayloadSize, 0, 0)
	w.state.bufferpool.SetConfig(cfg)
	w.state.bufferpool.SetActive(true)

	return nil
}

// This runs in a goroutine and checks for pause events or new buffers to push onto the pad.
func (w *websocketSrc) watchChannels(elem *gst.Element) {
	for {
		select {

		case data, more := <-w.state.bufferchan:
			if !more {
				elem.Log(CAT, gst.LevelInfo, "Buffer channel has closed, stopping processing")
				return
			}
			elem.Log(CAT, gst.LevelDebug, "Retrieving buffer from the pool")

			buf, ret := w.state.bufferpool.AcquireBuffer(nil)
			if ret != gst.FlowOK {
				elem.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed,
					fmt.Sprintf("Could not allocate buffer for data: %s", ret), "")
				return
			}

			elem.Log(CAT, gst.LevelDebug, "Writing data to buffer")
			buf.Map(gst.MapWrite).WriteData(data)
			buf.Unmap()
			buf.SetSize(int64(len(data)))

			elem.Log(CAT, gst.LevelDebug, "Pushing buffer onto src pad")
			w.pushPrelude(elem)
			if ret := w.state.srcpad.Push(buf); ret == gst.FlowError {
				elem.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed,
					fmt.Sprintf("Failed to push buffer to srcpad: %s", ret), "")
				return
			}

		case <-w.state.stopchan:
			elem.Log(CAT, gst.LevelInfo, "Received signal on stopchan to halt buffer processing")
			return

		}
	}
}

// start will start the websocket server and the buffer processing goroutines.
func (w *websocketSrc) start(elem *gst.Element) {
	w.state.mux.Lock()
	defer w.state.mux.Unlock()

	if !w.state.serverStarted {
		elem.Log(CAT, gst.LevelInfo, "Starting the HTTP server")
		go w.startServer(elem)
		w.state.serverStarted = true
	}
	if !w.state.channelsStarted {
		elem.Log(CAT, gst.LevelInfo, "Starting channel goroutine")
		go w.watchChannels(elem)
		w.state.channelsStarted = true
	}
	elem.Log(CAT, gst.LevelInfo, "WebsocketSrc has started")
}

// starts the server, is called as a goroutine.
func (w *websocketSrc) startServer(elem *gst.Element) {
	if err := w.state.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			elem.Log(CAT, gst.LevelInfo, "Server exited cleanly")
			return
		}
		elem.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed, "Failed to start websocket server", err.Error())
	}
}

// Checks if initial stream events were sent and pushes them onto the pad if needed.
func (w *websocketSrc) pushPrelude(elem *gst.Element) {
	w.state.mux.Lock()
	defer w.state.mux.Unlock()

	if !w.state.sentInitialEvents {
		elem.Log(CAT, gst.LevelDebug, "Sending stream start event")

		streamid := "blahblahblah"
		ev := gst.NewStreamStartEvent(streamid)
		if res := w.state.srcpad.PushEvent(ev); !res {
			elem.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, "Failed to notify elements of stream start", "")
			return
		}
		w.state.sentInitialEvents = true
	}
	if !w.state.sentSegment {
		elem.Log(CAT, gst.LevelDebug, "Sending new segment event")

		ev := gst.NewSegmentEvent(gst.NewFormattedSegment(gst.FormatTime))
		if res := w.state.srcpad.PushEvent(ev); !res {
			elem.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, "Failed to notify elements of new segment", "")
			return
		}
		w.state.sentSegment = true
	}
}

// Stops the goroutines and the websocket server
func (w *websocketSrc) stop(elem *gst.Element) {
	w.state.mux.Lock()
	defer w.state.mux.Unlock()

	if w.state.channelsStarted {
		elem.Log(CAT, gst.LevelInfo, "Sending stop signal to go routines")
		w.state.stopchan <- struct{}{}
		w.state.channelsStarted = false
	}

	if w.state.serverStarted {
		elem.Log(CAT, gst.LevelInfo, "Shutting down HTTP server")
		w.state.server.Shutdown(context.Background())
		w.state.serverStarted = false
	}
}

// Just stops the buffer processing routine, but leaves the server running
func (w *websocketSrc) pause(elem *gst.Element) {
	w.state.mux.Lock()
	defer w.state.mux.Unlock()
	elem.Log(CAT, gst.LevelDebug, "Sending stop signal to go routines")
	w.state.stopchan <- struct{}{}
	w.state.channelsStarted = false
}

// Tears down all resources for the element.
func (w *websocketSrc) unprepare(elem *gst.Element) {
	w.state.mux.Lock()
	defer w.state.mux.Unlock()

	elem.Log(CAT, gst.LevelDebug, "Freeing pads and buffers")

	w.state.bufferpool.SetActive(false)
	w.state.bufferpool.Unref()

	elem.Log(CAT, gst.LevelDebug, "Closing channels and clearing state")

	close(w.state.bufferchan)
	close(w.state.stopchan)
	w.state = &state{}
}

// Sets up a src pad for an element and adds the necessary callbacks.
func (w *websocketSrc) setupSrcPad(elem *gst.Element) {
	// Configure the src pad
	elem.Log(CAT, gst.LevelDebug, "Configuring the src pad")

	w.state.srcpad = gst.NewPadFromTemplate(elem.GetPadTemplates()[0], "src")
	elem.AddPad(w.state.srcpad)

	// Set a function for handling events
	w.state.srcpad.SetEventFunction(func(pad *gst.Pad, parent *gst.Object, event *gst.Event) bool {
		var ret bool

		pad.Log(CAT, gst.LevelLog, fmt.Sprintf("Handling event: %s", event.Type()))

		switch event.Type() {
		case gst.EventTypeReconfigure:
			ret = true
		case gst.EventTypeLatency:
			ret = true
		default:
			ret = false
		}

		if ret {
			pad.Log(CAT, gst.LevelDebug, fmt.Sprintf("Handled event: %s", event.Type()))
		} else {
			pad.Log(CAT, gst.LevelLog, fmt.Sprintf("Didn't handle event: %s", event.Type()))
		}

		return ret
	})

	// Set a query handler for the src pad
	w.state.srcpad.SetQueryFunction(func(pad *gst.Pad, parent *gst.Object, query *gst.Query) bool {
		var ret bool

		pad.Log(CAT, gst.LevelLog, fmt.Sprintf("Handling query: %s", query.Type()))

		switch query.Type() {
		case gst.QueryLatency:
			query.SetLatency(true, 0, gst.ClockTimeNone)
			ret = true
		case gst.QueryScheduling:
			query.SetScheduling(gst.SchedulingFlagSequential, 1, -1, 0)
			query.AddSchedulingMode(gst.PadModePush)
			ret = true
		case gst.QueryCaps:
			query.SetCapsResult(gst.NewAnyCaps())
			ret = true
		default:
			ret = false
		}

		if ret {
			pad.Log(CAT, gst.LevelDebug, fmt.Sprintf("Handled query: %s", query.Type()))
		} else {
			pad.Log(CAT, gst.LevelLog, fmt.Sprintf("Didn't handle query: %s", query.Type()))
		}

		return ret
	})
}

// * ObjectSubclass * //

func (w *websocketSrc) New() gst.GoElement {
	return &websocketSrc{
		settings: defaultSettings(),
		state:    &state{},
	}
}

func (w *websocketSrc) TypeInit(instance *gst.TypeInstance) {}

func (w *websocketSrc) ClassInit(klass *gst.ElementClass) {
	klass.SetMetadata(
		"Websocket Src",
		"Src/Websocket",
		"Write stream from a connection over a websocket server",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
	klass.AddPadTemplate(gst.NewPadTemplate(
		"src",
		gst.PadDirectionSource,
		gst.PadPresenceAlways,
		gst.NewAnyCaps(),
	))
	klass.InstallProperties(properties)
}

// * Object * //
func (w *websocketSrc) SetProperty(self *gst.Object, id uint, value *glib.Value) {
	prop := properties[id]

	switch prop.Name() {
	case "address":
		val, err := value.GetString()
		if err != nil {
			gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
				"Could not get string from GValue",
				err.Error(),
			)
			return
		}
		w.settings.address = val
	case "port":
		val, err := value.GoValue()
		if err != nil {
			gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
				"Could not get go value from GValue",
				err.Error(),
			)
			return
		}
		intval, ok := val.(int)
		if !ok {
			gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
				fmt.Sprintf("Could not coerce govalue %v to integer", val),
				err.Error(),
			)
			return
		}
		w.settings.port = intval
	case "retrieve-remote-addr":
		val, err := value.GoValue()
		if err != nil {
			gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
				"Could not get go value from GValue",
				err.Error(),
			)
			return
		}
		boolval, ok := val.(bool)
		if !ok {
			gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
				fmt.Sprintf("Could not coerce govalue %v to bool", val),
				err.Error(),
			)
			return
		}
		w.settings.retrieveRemoteAddr = boolval
	default:
		gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorSettings,
			fmt.Sprintf("Cannot set invalid property %s", prop.Name()), "")

	}
}

func (w *websocketSrc) GetProperty(self *gst.Object, id uint) *glib.Value {
	prop := properties[id]

	var localVal interface{}

	switch prop.Name() {
	case "address":
		localVal = w.settings.address
	case "port":
		localVal = w.settings.port
	case "retrieve-remote-addr":
		localVal = w.settings.retrieveRemoteAddr
	default:
		gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorSettings,
			fmt.Sprintf("Cannot get invalid property %s", prop.Name()), "")
		return nil
	}

	val, err := glib.GValue(localVal)
	if err != nil {
		gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
			fmt.Sprintf("Could not convert %v to GValue", localVal),
			err.Error(),
		)
	}

	return val
}

func (w *websocketSrc) Constructed(self *gst.Object) {
	elem := gst.ToElement(self)
	w.setupSrcPad(elem)
}

// * Element * //

func (w *websocketSrc) ChangeState(self *gst.Element, transition gst.StateChange) (ret gst.StateChangeReturn) {
	self.Log(CAT, gst.LevelTrace, fmt.Sprintf("Changing state: %s", transition))

	ret = gst.StateChangeSuccess

	switch transition {
	case gst.StateChangeNullToReady:
		if err := w.prepare(self); err != nil {
			self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed, err.Error(), "")
			return gst.StateChangeFailure
		}
	case gst.StateChangePlayingToPaused:
		w.pause(self)
	case gst.StateChangeReadyToNull:
		w.unprepare(self)
	}

	// Apply the transition to the parent element
	if ret = self.ParentChangeState(transition); ret == gst.StateChangeFailure {
		return
	}

	switch transition {
	case gst.StateChangeReadyToPaused:
		ret = gst.StateChangeNoPreroll
	case gst.StateChangePausedToPlaying:
		w.start(self)
	case gst.StateChangePlayingToPaused:
		ret = gst.StateChangeNoPreroll
	case gst.StateChangePausedToReady:
		w.stop(self)
	}

	return
}
