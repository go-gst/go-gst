package base

/*
#include "gst.go.h"
*/
import "C"
import (
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"

	"github.com/tinyzimmer/go-gst/gst"
)

//export goGstCollectPadsBufferFunc
func goGstCollectPadsBufferFunc(pads *C.GstCollectPads, data *C.GstCollectData, buf *C.GstBuffer, userData C.gpointer) C.GstFlowReturn {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	collectPads := iface.(*CollectPads)
	f := collectPads.funcMap.bufferFunc

	var wrappedBuf *gst.Buffer
	var wrappedData *CollectData
	if buf != nil {
		wrappedBuf = gst.FromGstBufferUnsafeNone(unsafe.Pointer(buf))
		defer wrappedBuf.Unref()
	}
	if data != nil {
		wrappedData = wrapCollectData(data)
	}

	return C.GstFlowReturn(f(collectPads, wrappedData, wrappedBuf))
}

//export goGstCollectPadsClipFunc
func goGstCollectPadsClipFunc(pads *C.GstCollectPads, data *C.GstCollectData, inbuf *C.GstBuffer, outbuf **C.GstBuffer, userData C.gpointer) C.GstFlowReturn {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	collectPads := iface.(*CollectPads)
	f := collectPads.funcMap.clipFunc

	buf := gst.FromGstBufferUnsafeNone(unsafe.Pointer(inbuf))
	defer buf.Unref()

	ret, gooutbuf := f(collectPads, wrapCollectData(data), buf)
	if gooutbuf != nil {
		C.memcpy(unsafe.Pointer(*outbuf), unsafe.Pointer(gooutbuf.Instance()), C.sizeof_GstBuffer)
	}

	return C.GstFlowReturn(ret)
}

//export goGstCollectPadsCompareFunc
func goGstCollectPadsCompareFunc(pads *C.GstCollectPads, data1 *C.GstCollectData, ts1 C.GstClockTime, data2 *C.GstCollectData, ts2 C.GstClockTime, userData C.gpointer) C.gint {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	collectPads := iface.(*CollectPads)
	f := collectPads.funcMap.compareFunc

	return C.gint(f(collectPads, wrapCollectData(data1), time.Duration(ts1), wrapCollectData(data2), time.Duration(ts2)))
}

//export goGstCollectPadsEventFunc
func goGstCollectPadsEventFunc(pads *C.GstCollectPads, data *C.GstCollectData, event *C.GstEvent, userData C.gpointer) C.gboolean {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	collectPads := iface.(*CollectPads)
	f := collectPads.funcMap.eventFunc

	return gboolean(f(collectPads, wrapCollectData(data), gst.FromGstEventUnsafeNone(unsafe.Pointer(event))))
}

//export goGstCollectPadsFlushFunc
func goGstCollectPadsFlushFunc(pads *C.GstCollectPads, userData C.gpointer) {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	collectPads := iface.(*CollectPads)
	f := collectPads.funcMap.flushFunc

	f(collectPads)
}

//export goGstCollectPadsFunc
func goGstCollectPadsFunc(pads *C.GstCollectPads, userData C.gpointer) C.GstFlowReturn {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	collectPads := iface.(*CollectPads)
	f := collectPads.funcMap.funcFunc

	return C.GstFlowReturn(f(collectPads))
}

//export goGstCollectPadsQueryFunc
func goGstCollectPadsQueryFunc(pads *C.GstCollectPads, data *C.GstCollectData, query *C.GstQuery, userData C.gpointer) C.gboolean {
	iface := gopointer.Restore(unsafe.Pointer(userData))
	collectPads := iface.(*CollectPads)
	f := collectPads.funcMap.queryFunc

	return gboolean(f(collectPads, wrapCollectData(data), gst.FromGstQueryUnsafeNone(unsafe.Pointer(query))))
}
