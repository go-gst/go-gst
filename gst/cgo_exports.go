package gst

// CGO exports have to be defined in a separate file from where they are used or else
// there will be double linkage issues.

/*
#include <stdlib.h>
#include <gst/gst.h>
*/
import "C"

import (
	"reflect"
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

//export goClassInit
func goClassInit(klass C.gpointer, klassData C.gpointer) {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	ptr := unsafe.Pointer(klassData)
	iface := gopointer.Restore(ptr)
	defer gopointer.Unref(ptr)

	data := iface.(*classData)
	registeredClasses[klass] = data.elem

	data.ext.InitClass(unsafe.Pointer(klass), data.elem)

	C.g_type_class_add_private(klass, C.gsize(unsafe.Sizeof(uintptr(0))))

	data.elem.ClassInit(wrapElementClass(klass))
}

//export goInstanceInit
func goInstanceInit(obj *C.GTypeInstance, klass C.gpointer) {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	elem := registeredClasses[klass].New()
	registeredClasses[klass] = elem

	ptr := gopointer.Save(elem)
	private := C.g_type_instance_get_private(obj, registeredTypes[reflect.TypeOf(registeredClasses[klass]).String()])
	C.memcpy(unsafe.Pointer(private), unsafe.Pointer(&ptr), C.gsize(unsafe.Sizeof(uintptr(0))))
}

//export goURIHdlrGetURIType
func goURIHdlrGetURIType(gtype C.GType) C.GstURIType {
	return C.GstURIType(globalURIHdlr.GetURIType())
}

//export goURIHdlrGetProtocols
func goURIHdlrGetProtocols(gtype C.GType) **C.gchar {
	protocols := globalURIHdlr.GetProtocols()
	size := C.size_t(unsafe.Sizeof((*C.gchar)(nil)))
	length := C.size_t(len(protocols))
	arr := (**C.gchar)(C.malloc(length * size))
	view := (*[1 << 30]*C.gchar)(unsafe.Pointer(arr))[0:len(protocols):len(protocols)]
	for i, proto := range protocols {
		view[i] = (*C.gchar)(C.CString(proto))
	}
	return arr
}

//export goURIHdlrGetURI
func goURIHdlrGetURI(hdlr *C.GstURIHandler) *C.gchar {
	iface := FromObjectUnsafePrivate(unsafe.Pointer(hdlr))
	return (*C.gchar)(unsafe.Pointer(C.CString(iface.(URIHandler).GetURI())))
}

//export goURIHdlrSetURI
func goURIHdlrSetURI(hdlr *C.GstURIHandler, uri *C.gchar, gerr **C.GError) C.gboolean {
	iface := FromObjectUnsafePrivate(unsafe.Pointer(hdlr))
	ok, err := iface.(URIHandler).SetURI(C.GoString(uri))
	if err != nil {
		C.g_set_error_literal(gerr, DomainLibrary.toQuark(), C.gint(LibraryErrorSettings), C.CString(err.Error()))
	}
	return gboolean(ok)
}

//export goObjectSetProperty
func goObjectSetProperty(obj *C.GObject, propID C.guint, val *C.GValue, param *C.GParamSpec) {
	iface := FromObjectUnsafePrivate(unsafe.Pointer(obj)).(interface {
		SetProperty(obj *Object, id uint, value *glib.Value)
	})
	iface.SetProperty(wrapObject(toGObject(unsafe.Pointer(obj))), uint(propID-1), glib.ValueFromNative(unsafe.Pointer(val)))
}

//export goObjectGetProperty
func goObjectGetProperty(obj *C.GObject, propID C.guint, value *C.GValue, param *C.GParamSpec) {
	iface := FromObjectUnsafePrivate(unsafe.Pointer(obj)).(interface {
		GetProperty(obj *Object, id uint) *glib.Value
	})
	val := iface.GetProperty(wrapObject(toGObject(unsafe.Pointer(obj))), uint(propID-1))
	if val == nil {
		return
	}
	C.g_value_copy((*C.GValue)(unsafe.Pointer(val.GValue)), value)
}

//export goObjectConstructed
func goObjectConstructed(obj *C.GObject) {
	iface := FromObjectUnsafePrivate(unsafe.Pointer(obj)).(interface {
		Constructed(*Object)
	})
	iface.Constructed(wrapObject(toGObject(unsafe.Pointer(obj))))
}

//export goObjectFinalize
func goObjectFinalize(obj *C.GObject, klass C.gpointer) {
	registerMutex.Lock()
	defer registerMutex.Unlock()
	delete(registeredClasses, klass)
	gopointer.Unref(privateFromObj(unsafe.Pointer(obj)))
}
