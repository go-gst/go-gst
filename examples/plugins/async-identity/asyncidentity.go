//lint:file-ignore U1000 Ignore all unused code, this is example code

// +plugin:Name=async-identity
// +plugin:Description=A go-gst example plugin with async state changes
// +plugin:Version=v0.0.1
// +plugin:License=gst.LicenseLGPL
// +plugin:Source=go-gst
// +plugin:Package=examples
// +plugin:Origin=https://github.com/go-gst/go-gst
// +plugin:ReleaseDate=2024-09-13
//
// +element:Name=asyncidentity
// +element:Rank=gst.RankNone
// +element:Impl=asyncidentity
// +element:Subclass=gst.ExtendsElement
//
//go:generate gst-plugin-gen
package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

var (
	_cat = gst.NewDebugCategory(
		"asyncidentity",
		gst.DebugColorNone,
		"asyncidentity element",
	)

	_srcPadTemplate = gst.NewPadTemplate("generic-src", gst.PadDirectionSource,
		gst.PadPresenceAlways, gst.NewAnyCaps())
	_sinkPadTemplate = gst.NewPadTemplate("generic-sink", gst.PadDirectionSink,
		gst.PadPresenceAlways, gst.NewAnyCaps())
	_properties = []*glib.ParamSpec{
		glib.NewUint64Param(
			"delay",
			"ns state change delay",
			"Duration in nanoseconds to wait until a state changes",
			_delayNsMin, _delayNsMax, _delayNsDefault,
			glib.ParameterReadWrite,
		),
	}
)

const (
	_propDelayNs = 0

	_delayNsMin     = uint64(0)
	_delayNsMax     = uint64(time.Second) * 10
	_delayNsDefault = uint64(time.Second)
)

func main() {}

type asyncidentity struct {
	// inner state
	sinkpad *gst.Pad
	srcpad  *gst.Pad

	asyncPending atomic.Bool

	// property storage
	delayNs atomic.Uint64
}

var _ glib.GoObjectSubclass = (*asyncidentity)(nil)

func (g *asyncidentity) New() glib.GoObjectSubclass { return &asyncidentity{} }

func (g *asyncidentity) ClassInit(klass *glib.ObjectClass) {
	class := gst.ToElementClass(klass)
	class.SetMetadata(
		"Async Identity Example",
		"General",
		"An async state changing identity like element",
		"Artem Martus <artemmartus2012@gmail.com>",
	)

	class.AddStaticPadTemplate(_srcPadTemplate)
	class.AddStaticPadTemplate(_sinkPadTemplate)

	class.InstallProperties(_properties)
}

var _ glib.GoObject = (*asyncidentity)(nil)

func (g *asyncidentity) SetProperty(obj *glib.Object, id uint, value *glib.Value) {
	self := gst.ToElement(obj)

	switch id {
	case _propDelayNs:
		newDelayErased, err := value.GoValue()
		if err != nil {
			self.Error("Failed unmarshalling the 'delay' property", err)
			return
		}
		newDelay, ok := newDelayErased.(uint64)
		if !ok {
			self.Error("Failed Go-casting the 'delay' interface{} into uint64",
				fmt.Errorf("interfaced value: %+v", newDelayErased))
			return
		}
		oldDelay := g.delayNs.Swap(newDelay)

		self.Log(_cat, gst.LevelInfo,
			fmt.Sprintf("Changed delay property %s => %s",
				time.Duration(oldDelay),
				time.Duration(newDelay),
			))
	default:
		self.Error("Tried to set unknown property",
			fmt.Errorf("prop id %d: %s", id, value.TypeName()))
	}
}

func (g *asyncidentity) GetProperty(obj *glib.Object, id uint) *glib.Value {
	var (
		out *glib.Value
		err error
	)

	switch id {
	case _propDelayNs:
		out, err = glib.GValue(g.delayNs.Load())
	default:
		err = fmt.Errorf("unknown property id: %d", id)
	}

	if err != nil {
		self := gst.ToElement(obj)
		self.Error("Get property error", err)
		out = nil
	}

	return out
}

func (g *asyncidentity) Constructed(self *glib.Object) {
	elem := gst.ToElement(self)
	srcPad := gst.NewPadFromTemplate(_srcPadTemplate, "src")
	sinkPad := gst.NewPadFromTemplate(_sinkPadTemplate, "sink")

	sinkPad.SetChainFunction(g.sink_chain_function)

	// Have to set proxy flags on a pads
	proxyFlags := gst.PadFlagProxyAllocation | gst.PadFlagProxyCaps | gst.PadFlagProxyScheduling
	sinkPad.SetFlags(proxyFlags)
	srcPad.SetFlags(proxyFlags)

	// Or setup query & event functions like so

	// sinkPad.SetQueryFunction(func(self *gst.Pad, parent *gst.Object, query *gst.Query) bool {
	// 	return srcPad.PeerQuery(query)
	// })
	// sinkPad.SetEventFunction(func(self *gst.Pad, parent *gst.Object, event *gst.Event) bool {
	// 	return srcPad.PushEvent(event)
	// })

	// srcPad.SetQueryFunction(func(self *gst.Pad, parent *gst.Object, query *gst.Query) bool {
	// 	return sinkPad.PeerQuery(query)
	// })
	// srcPad.SetEventFunction(func(self *gst.Pad, parent *gst.Object, event *gst.Event) bool {
	// 	return sinkPad.PushEvent(event)
	// })

	elem.AddPad(srcPad)
	elem.AddPad(sinkPad)

	g.srcpad = srcPad
	g.sinkpad = sinkPad

	g.delayNs.Store(_delayNsDefault)
}

func (g *asyncidentity) sink_chain_function(
	_self *gst.Pad,
	_parent *gst.Object,
	buffer *gst.Buffer,
) gst.FlowReturn {

	return g.srcpad.Push(buffer)
}

// var _ gst.ElementImpl = (*asyncidentity)(nil)

func (g *asyncidentity) ChangeState(el *gst.Element, transition gst.StateChange) gst.StateChangeReturn {
	if ret := el.ParentChangeState(transition); ret == gst.StateChangeFailure {
		return ret
	}

	switch transition {
	case gst.StateChangeNullToReady:
		// async will be ignored due to target state <= READY
	case gst.StateChangeReadyToPaused:
		// async will be ignored due to no_preroll
		return gst.StateChangeNoPreroll
	case gst.StateChangePausedToPlaying:
		fallthrough
	case gst.StateChangePlayingToPaused:
		g.asyncStateChange(el)
		return gst.StateChangeAsync
	case gst.StateChangePausedToReady:
		// async will be ignored due to target state <= READY
	case gst.StateChangeReadyToNull:
	}

	// check against forcing state change
	if g.asyncPending.Load() {
		return gst.StateChangeAsync
	}

	return gst.StateChangeSuccess
}

func (g *asyncidentity) asyncStateChange(el *gst.Element) {
	msg := gst.NewAsyncStartMessage(el)
	_ = el.PostMessage(msg)
	go func(el *gst.Element) {
		g.asyncPending.Store(true)
		delay := time.Duration(g.delayNs.Load())
		<-time.After(delay)
		msg := gst.NewAsyncDoneMessage(el, gst.ClockTimeNone)
		_ = el.PostMessage(msg)
		g.asyncPending.Store(false)
	}(el)
}
