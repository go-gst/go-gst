package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

func getParent(parent *C.GstObject) *Object {
	if parent == nil {
		return nil
	}
	return wrapObject(toGObject(unsafe.Pointer(parent)))
}

//export goGstPadFuncDestroyNotify
func goGstPadFuncDestroyNotify(notifyInfo *C.PadDestroyNotifyInfo) {
	padPtr := unsafe.Pointer(notifyInfo.pad_ptr)
	funcMapPtr := unsafe.Pointer(notifyInfo.func_map_ptr)

	defer gopointer.Unref(padPtr)
	defer gopointer.Unref(funcMapPtr)

	pad := gopointer.Restore(padPtr).(unsafe.Pointer)
	funcMap := gopointer.Restore(funcMapPtr).(PadFuncMap)

	funcMap.RemoveFuncForPad(pad)
}

//export goGstPadActivateFunction
func goGstPadActivateFunction(pad *C.GstPad, parent *C.GstObject) C.gboolean {
	f := padActivateFuncs.FuncForPad(unsafe.Pointer(pad)).(PadActivateFunc)
	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
	))
}

//export goGstPadActivateModeFunction
func goGstPadActivateModeFunction(pad *C.GstPad, parent *C.GstObject, mode C.GstPadMode, active C.gboolean) C.gboolean {
	f := padActivateModeFuncs.FuncForPad(unsafe.Pointer(pad)).(PadActivateModeFunc)
	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		PadMode(mode),
		gobool(active),
	))
}

//export goGstPadChainFunction
func goGstPadChainFunction(pad *C.GstPad, parent *C.GstObject, buffer *C.GstBuffer) C.GstFlowReturn {
	f := padChainFuncs.FuncForPad(unsafe.Pointer(pad)).(PadChainFunc)
	buf := FromGstBufferUnsafeFull(unsafe.Pointer(buffer))
	defer buf.Unref()
	return C.GstFlowReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		buf,
	))
}

//export goGstPadChainListFunction
func goGstPadChainListFunction(pad *C.GstPad, parent *C.GstObject, list *C.GstBufferList) C.GstFlowReturn {
	f := padChainListFuncs.FuncForPad(unsafe.Pointer(pad)).(PadChainListFunc)
	buflist := FromGstBufferListUnsafeFull(unsafe.Pointer(list))
	defer buflist.Unref()
	return C.GstFlowReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		buflist,
	))
}

//export goGstPadEventFullFunction
func goGstPadEventFullFunction(pad *C.GstPad, parent *C.GstObject, event *C.GstEvent) C.GstFlowReturn {
	f := padEventFullFuncs.FuncForPad(unsafe.Pointer(pad)).(PadEventFullFunc)
	ev := FromGstEventUnsafeFull(unsafe.Pointer(event))
	defer ev.Unref()
	return C.GstFlowReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		ev,
	))
}

//export goGstPadEventFunction
func goGstPadEventFunction(pad *C.GstPad, parent *C.GstObject, event *C.GstEvent) C.gboolean {
	f := padEventFuncs.FuncForPad(unsafe.Pointer(pad)).(PadEventFunc)
	ev := FromGstEventUnsafeFull(unsafe.Pointer(event))
	defer ev.Unref()
	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		ev,
	))
}

//export goGstPadGetRangeFunction
func goGstPadGetRangeFunction(pad *C.GstPad, parent *C.GstObject, offset C.guint64, length C.guint, buffer **C.GstBuffer) C.GstFlowReturn {
	f := padGetRangeFuncs.FuncForPad(unsafe.Pointer(pad)).(PadGetRangeFunc)
	ret, buf := f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		uint64(offset),
		uint(length),
	)
	if ret == FlowOK {
		C.memcpy(unsafe.Pointer(*buffer), unsafe.Pointer(buf.Instance()), C.sizeof_GstBuffer)
	}
	return C.GstFlowReturn(ret)
}

//export goGstPadIterIntLinkFunction
func goGstPadIterIntLinkFunction(pad *C.GstPad, parent *C.GstObject) *C.GstIterator {
	f := padIterIntLinkFuncs.FuncForPad(unsafe.Pointer(pad)).(PadIterIntLinkFunc)
	pads := f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
	)
	if len(pads) == 0 {
		return nil
	}
	if len(pads) == 1 {
		val, err := glib.ValueInit(glib.Type(C.gst_pad_get_type()))
		if err != nil {
			return nil
		}
		val.SetInstance(uintptr(unsafe.Pointer(pads[0].Instance())))
		return C.gst_iterator_new_single(C.gst_pad_get_type(), (*C.GValue)(unsafe.Pointer(val.GValue)))
	}
	return nil
}

//export goGstPadLinkFunction
func goGstPadLinkFunction(pad *C.GstPad, parent *C.GstObject, peer *C.GstPad) C.GstPadLinkReturn {
	f := padLinkFuncs.FuncForPad(unsafe.Pointer(pad)).(PadLinkFunc)
	return C.GstPadLinkReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		wrapPad(toGObject(unsafe.Pointer(peer))),
	))
}

//export goGstPadQueryFunction
func goGstPadQueryFunction(pad *C.GstPad, parent *C.GstObject, query *C.GstQuery) C.gboolean {
	f := padQueryFuncs.FuncForPad(unsafe.Pointer(pad)).(PadQueryFunc)
	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		wrapQuery(query),
	))
}

//export goGstPadUnlinkFunction
func goGstPadUnlinkFunction(pad *C.GstPad, parent *C.GstObject) {
	f := padUnlinkFuncs.FuncForPad(unsafe.Pointer(pad)).(PadUnlinkFunc)
	f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
	)
}
