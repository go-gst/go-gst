// This is a GStreamer element implemented in Go that uses inbound data on a websocket
// connection as the source for the stream.
//
// In order to build the plugin for use by GStreamer, you can do the following:
//
//     $ go generate
//     $ go build -o libgstgofilesrc.so -buildmode c-shared .
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
	"fmt"
	"net/http"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

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
	gst.NewBoolParam(
		"retrieve-remote-addr",
		"Retrieve Remote Address",
		"Include the remote client's address in the buffer metadata",
		DefaultRetrieveRemoteAddr,
		gst.ParameterReadWrite,
	),
}

// Internals //

type state struct {
	started           bool
	server            *http.Server
	needInitialEvents bool
	needSegment       bool
}

type settings struct {
	address            string
	port               int
	retrieveRemoteAddr bool
}

func defaultSettings() *settings {
	return &settings{
		address:            DefaultAddress,
		port:               DefaultPort,
		retrieveRemoteAddr: DefaultRetrieveRemoteAddr,
	}
}

// Element implementation //

type websocketSrc struct {
	settings *settings
	state    *state
	srcpad   *gst.Pad
}

// // ObjectSubclass // //

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

// // Object // //
func (w *websocketSrc) SetProperty(self *gst.Object, id uint, value *glib.Value) {}

func (w *websocketSrc) GetProperty(self *gst.Object, id uint) *glib.Value { return nil }

func (w *websocketSrc) Constructed(self *gst.Object) {
	w.srcpad = gst.ToElement(self).GetStaticPad("src")

	w.srcpad.SetEventFunction(func(pad *gst.Pad, parent *gst.Object, event *gst.Event) bool {
		var ret bool

		self.Log(CAT, gst.LevelLog, fmt.Sprintf("Handling event: %s", event.Type()))

		switch event.Type() {
		case gst.EventTypeFlushStart:
			// TODO
		case gst.EventTypeFlushStop:
			// TODO
		case gst.EventTypeReconfigure:
			ret = true
		case gst.EventTypeLatency:
			ret = true
		default:
			ret = false
		}

		if ret {
			self.Log(CAT, gst.LevelLog, fmt.Sprintf("Handled event: %s", event.Type()))
		} else {
			self.Log(CAT, gst.LevelLog, fmt.Sprintf("Didn't handle event: %s", event.Type()))
		}

		return ret
	})

	w.srcpad.SetQueryFunction(func(pad *gst.Pad, parent *gst.Object, query *gst.Query) bool {
		var ret bool

		self.Log(CAT, gst.LevelLog, fmt.Sprintf("Handling query: %s", query.Type()))

		switch query.Type() {
		case gst.QueryLatency:
			query.SetLatency(true, 0, gst.ClockTimeNone)
			ret = true
		case gst.QueryScheduling:
			query.SetScheduling(gst.SchedulingFlagSequential, 1, -1, 0)
			query.AddSchedulingMode(gst.PadModePush)
			ret = true
		case gst.QueryCaps:
			query.SetCapsResult(query.ParseCaps())
			ret = true
		default:
			ret = false
		}

		return ret
	})
}
