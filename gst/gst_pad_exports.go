package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
	gopointer "github.com/go-gst/go-pointer"
)

func getParent(parent *C.GstObject) *Object {
	if parent == nil {
		return nil
	}
	return wrapObject(toGObject(unsafe.Pointer(parent)))
}

//export goGstPadFuncDestroyNotify
func goGstPadFuncDestroyNotify(notifyInfo *C.PadDestroyNotifyInfo) {
	funcMapPtr := unsafe.Pointer(notifyInfo.func_map_ptr)

	defer gopointer.Unref(funcMapPtr)
	funcMap := gopointer.Restore(funcMapPtr).(padFuncMapLike)

	funcMap.removeFuncForPad(notifyInfo.pad_ptr)
}

//export goGstPadActivateFunction
func goGstPadActivateFunction(pad *C.GstPad, parent *C.GstObject) C.gboolean {
	f := padActivateFuncs.funcForPad(pad)
	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
	))
}

//export goGstPadActivateModeFunction
func goGstPadActivateModeFunction(pad *C.GstPad, parent *C.GstObject, mode C.GstPadMode, active C.gboolean) C.gboolean {
	f := padActivateModeFuncs.funcForPad(pad)
	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		PadMode(mode),
		gobool(active),
	))
}

//export goGstPadChainFunction
func goGstPadChainFunction(pad *C.GstPad, parent *C.GstObject, buffer *C.GstBuffer) C.GstFlowReturn {
	f := padChainFuncs.funcForPad(pad)

	// do not work with a finalizer here, because they are too unreliable for such short lived objects
	buf := ToGstBuffer(unsafe.Pointer(buffer))
	// defer buf.Unref() // FIXME: the buffer leaks in case we don't pass to ChainDefault() or similar, but we cannot take a reference on the buffer because that breakes writability

	return C.GstFlowReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		buf,
	))
}

//export goGstPadChainListFunction
func goGstPadChainListFunction(pad *C.GstPad, parent *C.GstObject, list *C.GstBufferList) C.GstFlowReturn {
	f := padChainListFuncs.funcForPad(pad)

	// do not work with a finalizer here, because they are too unreliable for such short lived objects
	buflist := ToGstBufferList(unsafe.Pointer(list))
	// defer buflist.Unref() // FIXME: the buffer leaks in case we don't pass to ChainDefault() or similar, but we cannot take a reference on the buffer because that breakes writability

	return C.GstFlowReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		buflist,
	))
}

//export goGstPadEventFullFunction
func goGstPadEventFullFunction(pad *C.GstPad, parent *C.GstObject, event *C.GstEvent) C.GstFlowReturn {
	f := padEventFullFuncs.funcForPad(pad)

	// do not work with a finalizer here, because they are too unreliable for such short lived objects
	ev := ToGstEvent(unsafe.Pointer(event))
	// defer ev.Unref() // FIXME: the buffer leaks in case we don't pass to ChainDefault() or similar, but we cannot take a reference on the buffer because that breakes writability

	return C.GstFlowReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		ev,
	))
}

//export goGstPadEventFunction
func goGstPadEventFunction(pad *C.GstPad, parent *C.GstObject, event *C.GstEvent) C.gboolean {
	f := padEventFuncs.funcForPad(pad)

	// do not work with a finalizer here, because they are too unreliable for such short lived objects
	ev := ToGstEvent(unsafe.Pointer(event))
	// defer ev.Unref() // FIXME: the event leaks when not passed to EventDefault or beeing manually Unreffed.

	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		ev,
	))
}

//export goGstPadGetRangeFunction
func goGstPadGetRangeFunction(pad *C.GstPad, parent *C.GstObject, offset C.guint64, length C.guint, buffer **C.GstBuffer) C.GstFlowReturn {
	f := padGetRangeFuncs.funcForPad(pad)
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
	f := padIterIntLinkFuncs.funcForPad(pad)
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
		val.SetInstance(unsafe.Pointer(pads[0].Instance()))
		return C.gst_iterator_new_single(C.gst_pad_get_type(), (*C.GValue)(unsafe.Pointer(val.GValue)))
	}
	return nil
}

//export goGstPadLinkFunction
func goGstPadLinkFunction(pad *C.GstPad, parent *C.GstObject, peer *C.GstPad) C.GstPadLinkReturn {
	f := padLinkFuncs.funcForPad(pad)
	return C.GstPadLinkReturn(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		wrapPad(toGObject(unsafe.Pointer(peer))),
	))
}

//export goGstPadQueryFunction
func goGstPadQueryFunction(pad *C.GstPad, parent *C.GstObject, query *C.GstQuery) C.gboolean {
	f := padQueryFuncs.funcForPad(pad)
	return gboolean(f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
		wrapQuery(query),
	))
}

//export goGstPadUnlinkFunction
func goGstPadUnlinkFunction(pad *C.GstPad, parent *C.GstObject) {
	f := padUnlinkFuncs.funcForPad(pad)
	f(
		wrapPad(toGObject(unsafe.Pointer(pad))),
		getParent(parent),
	)
}
