package gst

// CGO exports have to be defined in a separate file from where they are used or else
// there will be double linkage issues.

/*
#include <stdlib.h>
#include <gst/gst.h>
*/
import "C"

import (
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

//export goElementCallAsync
func goElementCallAsync(element *C.GstElement, userData C.gpointer) {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	defer gopointer.Unref(unsafe.Pointer(userData))
	f := iface.(func())
	f()
}

//export goPadStickyEventForEachFunc
func goPadStickyEventForEachFunc(gpad *C.GstPad, event **C.GstEvent, userData C.gpointer) C.gboolean {
	cbIface := gopointer.Restore(unsafe.Pointer(userData))
	cbFunc := cbIface.(StickyEventsForEachFunc)
	pad := wrapPad(toGObject(unsafe.Pointer(gpad)))
	ev := wrapEvent(*event)
	return gboolean(cbFunc(pad, ev))
}

//export goPadProbeFunc
func goPadProbeFunc(gstPad *C.GstPad, info *C.GstPadProbeInfo, userData C.gpointer) C.GstPadProbeReturn {
	cbIface := gopointer.Restore(unsafe.Pointer(userData))
	cbFunc := cbIface.(PadProbeCallback)
	pad := wrapPad(toGObject(unsafe.Pointer(gstPad)))
	return C.GstPadProbeReturn(cbFunc(pad, &PadProbeInfo{info}))
}

//export goPadForwardFunc
func goPadForwardFunc(gstPad *C.GstPad, userData C.gpointer) C.gboolean {
	cbIface := gopointer.Restore(unsafe.Pointer(userData))
	cbFunc := cbIface.(PadForwardFunc)
	pad := wrapPad(toGObject(unsafe.Pointer(gstPad)))
	return gboolean(cbFunc(pad))
}

//export goTagForEachFunc
func goTagForEachFunc(tagList *C.GstTagList, tag *C.gchar, userData C.gpointer) {
	cbIface := gopointer.Restore(unsafe.Pointer(userData))
	cbFunc := cbIface.(TagListForEachFunc)
	cbFunc(wrapTagList(tagList), Tag(C.GoString(tag)))
}

//export goBufferListForEachCb
func goBufferListForEachCb(buf **C.GstBuffer, idx C.guint, userData C.gpointer) C.gboolean {
	cbIface := gopointer.Restore(unsafe.Pointer(userData))
	cbFunc := cbIface.(func(*Buffer, uint) bool)
	return gboolean(cbFunc(wrapBuffer(*buf), uint(idx)))
}

//export goBufferMetaForEachCb
func goBufferMetaForEachCb(buf *C.GstBuffer, meta **C.GstMeta, userData C.gpointer) C.gboolean {
	cbIface := gopointer.Restore(unsafe.Pointer(userData))
	cbFunc := cbIface.(func(*Meta) bool)
	return gboolean(cbFunc(wrapMeta(*meta)))
}

//export structForEachCb
func structForEachCb(fieldID C.GQuark, val *C.GValue, chPtr C.gpointer) C.gboolean {
	ptr := gopointer.Restore(unsafe.Pointer(chPtr))
	resCh := ptr.(chan interface{})
	fieldName := C.GoString(C.g_quark_to_string(fieldID))

	var resValue interface{}

	gVal := glib.ValueFromNative(unsafe.Pointer(val))
	if resValue, _ = gVal.GoValue(); resValue == nil {
		// serialize the value if we can't do anything else with it
		serialized := C.gst_value_serialize(val)
		defer C.free(unsafe.Pointer(serialized))
		resValue = C.GoString(serialized)
	}

	resCh <- fieldName
	resCh <- resValue
	return gboolean(true)
}

//export goBusSyncHandler
func goBusSyncHandler(bus *C.GstBus, cMsg *C.GstMessage, userData C.gpointer) C.GstBusSyncReply {
	// wrap the message
	msg := wrapMessage(cMsg)

	// retrieve the ptr to the function
	ptr := unsafe.Pointer(userData)
	funcIface := gopointer.Restore(ptr)
	busFunc, ok := funcIface.(BusSyncHandler)

	if !ok {
		gopointer.Unref(ptr)
		return C.GstBusSyncReply(BusPass)
	}

	return C.GstBusSyncReply(busFunc(msg))
}

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

func getMetaInfoCbFuncs(meta *C.GstMeta) *MetaInfoCallbackFuncs {
	gapi := glib.Type(meta.info.api)
	gtype := glib.Type(meta.info._type)
	typeCbs := registeredMetas[gapi]
	if typeCbs == nil {
		return nil
	}
	return typeCbs[gtype.Name()]
}

//export goMetaFreeFunc
func goMetaFreeFunc(meta *C.GstMeta, buf *C.GstBuffer) {
	cbFuncs := getMetaInfoCbFuncs(meta)
	if cbFuncs != nil && cbFuncs.FreeFunc != nil {
		cbFuncs.FreeFunc(wrapBuffer(buf))
	}
}

//export goMetaInitFunc
func goMetaInitFunc(meta *C.GstMeta, params C.gpointer, buf *C.GstBuffer) C.gboolean {
	cbFuncs := getMetaInfoCbFuncs(meta)
	if cbFuncs != nil && cbFuncs.InitFunc != nil {
		paramsIface := gopointer.Restore(unsafe.Pointer(params))
		defer gopointer.Unref(unsafe.Pointer(params))
		return gboolean(cbFuncs.InitFunc(paramsIface, wrapBuffer(buf)))
	}
	return gboolean(true)
}

//export goMetaTransformFunc
func goMetaTransformFunc(transBuf *C.GstBuffer, meta *C.GstMeta, buffer *C.GstBuffer, mType C.GQuark, data C.gpointer) C.gboolean {
	cbFuncs := getMetaInfoCbFuncs(meta)
	if cbFuncs != nil && cbFuncs.TransformFunc != nil {
		transformData := (*C.GstMetaTransformCopy)(unsafe.Pointer(data))
		return gboolean(cbFuncs.TransformFunc(
			wrapBuffer(transBuf),
			wrapBuffer(buffer),
			quarkToString(mType),
			&MetaTransformCopy{
				Region: gobool(transformData.region),
				Offset: int64(transformData.offset),
				Size:   int64(transformData.size),
			},
		))
	}
	return gboolean(true)
}

//export goGDestroyNotifyFunc
func goGDestroyNotifyFunc(ptr C.gpointer) {
	funcIface := gopointer.Restore(unsafe.Pointer(ptr))
	defer gopointer.Unref(unsafe.Pointer(ptr))
	f := funcIface.(func())
	if f != nil {
		f()
	}
}

//export goGDestroyNotifyFuncNoRun
func goGDestroyNotifyFuncNoRun(ptr C.gpointer) {
	gopointer.Unref(unsafe.Pointer(ptr))
}

//export goCapsMapFunc
func goCapsMapFunc(features *C.GstCapsFeatures, structure *C.GstStructure, userData C.gpointer) C.gboolean {
	// retrieve the ptr to the function
	ptr := unsafe.Pointer(userData)
	funcIface := gopointer.Restore(ptr)
	mapFunc, ok := funcIface.(CapsMapFunc)

	if !ok {
		gopointer.Unref(ptr)
		return gboolean(false)
	}

	return gboolean(mapFunc(wrapCapsFeatures(features), wrapStructure(structure)))
}

//export goClockCb
func goClockCb(gclock *C.GstClock, clockTime C.GstClockTime, clockID C.GstClockID, userData C.gpointer) C.gboolean {
	// retrieve the ptr to the function
	ptr := unsafe.Pointer(userData)
	funcIface := gopointer.Restore(ptr)
	cb, ok := funcIface.(ClockCallback)

	if !ok {
		gopointer.Unref(ptr)
		return gboolean(false)
	}

	clock := wrapClock(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(gclock))})
	return gboolean(cb(clock, time.Duration(clockTime)))
}

//export goPluginInit
func goPluginInit(plugin *C.GstPlugin, userData C.gpointer) C.gboolean {
	ptr := unsafe.Pointer(userData)
	defer gopointer.Unref(ptr)
	funcIface := gopointer.Restore(ptr)
	cb, ok := funcIface.(PluginInitFunc)
	if !ok {
		return gboolean(false)
	}
	return gboolean(cb(wrapPlugin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(plugin))})))
}

//export goGlobalPluginInit
func goGlobalPluginInit(plugin *C.GstPlugin) C.gboolean {
	return gboolean(globalPluginInit(wrapPlugin(&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(plugin))})))
}
