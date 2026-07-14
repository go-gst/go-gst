package gst

import (
	"runtime"
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/userdata"
	"github.com/go-gst/go-glib/pkg/gobject/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
// extern gboolean _goglib_gst1_StructureForeachFunc(GQuark fieldid, GValue* value, gpointer user_data);
import "C"

// TakeValue takes ownership of the value, which could be useful but also break in very fun-to-debug ways.
// Maybe we'll add it when a user asks for it.
// // TakeValue wraps gst_structure_take_value
// func (structure *Structure) TakeValue(fieldname string, value any) {
// 	var carg0 *C.GstStructure // in, none, converted
// 	var carg1 *C.gchar        // in, none, string, casted *C.gchar
// 	var carg2 *C.GValue       // in, full, converted

// 	carg0 = (*C.GstStructure)(UnsafeStructureToGlibNone(structure))
// 	carg1 = (*C.gchar)(unsafe.Pointer(C.CString(fieldname)))
// 	defer C.free(unsafe.Pointer(carg1))
// 	v := gobject.NewValue(value)
// 	carg2 = (*C.GValue)(gobject.UnsafeValueToGlibFull(v))

// 	C.gst_structure_take_value(carg0, carg1, carg2)
// 	runtime.KeepAlive(structure)
// 	runtime.KeepAlive(fieldname)
// 	runtime.KeepAlive(value)
// 	runtime.KeepAlive(v)
// }

// GetValue wraps gst_structure_get_value
func (structure *Structure) GetValue(fieldname string) any {
	var carg0 *C.GstStructure // in, none, converted
	var carg1 *C.gchar        // in, none, string, casted *C.gchar
	var cret *C.GValue        // return, none, converted

	carg0 = (*C.GstStructure)(UnsafeStructureToGlibNone(structure))
	carg1 = (*C.gchar)(unsafe.Pointer(C.CString(fieldname)))
	defer C.free(unsafe.Pointer(carg1))

	cret = C.gst_structure_get_value(carg0, carg1)
	runtime.KeepAlive(structure)
	runtime.KeepAlive(fieldname)

	var goret any

	// FIXME: this must borrow, because we don't own the value
	goret = gobject.ValueFromNative(unsafe.Pointer(cret)).GoValue()

	return goret
}

// SetValue wraps gst_structure_set_value
func (structure *Structure) SetValue(fieldname string, value any) {
	var carg0 *C.GstStructure // in, none, converted
	var carg1 *C.gchar        // in, none, string, casted *C.gchar
	var carg2 *C.GValue       // in, none, converted

	carg0 = (*C.GstStructure)(UnsafeStructureToGlibNone(structure))
	carg1 = (*C.gchar)(unsafe.Pointer(C.CString(fieldname)))
	defer C.free(unsafe.Pointer(carg1))
	v := gobject.NewValue(value)
	carg2 = (*C.GValue)(gobject.UnsafeValueToGlibNone(v))

	C.gst_structure_set_value(carg0, carg1, carg2)
	runtime.KeepAlive(structure)
	runtime.KeepAlive(fieldname)
	runtime.KeepAlive(value)
}

// StructureForeachFunc wraps GstStructureForeachFunc
//
// We don't expose the quark to the user though, because that is not very go like.
//
// see also https://gstreamer.freedesktop.org/documentation/gstreamer/gststructure.html#GstStructureForeachFunc
type StructureForeachFunc func(field string, value any) (goret bool)

// ForEach wraps gst_caps_foreach
//
// see also https://gstreamer.freedesktop.org/documentation/gstreamer/gstcaps.html#gst_caps_foreach
func (structure *Structure) ForEach(fn StructureForeachFunc) bool {
	var carg0 *C.GstStructure           // in, none, converted
	var carg1 C.GstStructureForeachFunc // callback, scope: call, closure: carg2
	var carg2 C.gpointer                // implicit
	var cret C.gboolean                 // return

	carg0 = (*C.GstStructure)(UnsafeStructureToGlibNone(structure))
	carg1 = (*[0]byte)(C._goglib_gst1_StructureForeachFunc)
	carg2 = C.gpointer(userdata.Register(fn))
	defer userdata.Delete(unsafe.Pointer(carg2))

	cret = C.gst_structure_foreach(carg0, carg1, carg2)
	runtime.KeepAlive(structure)
	runtime.KeepAlive(fn)

	var goret bool

	if cret != 0 {
		goret = true
	}

	return goret
}
