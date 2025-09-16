package gst

import (
	"runtime"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gobject/v2"
)

// #cgo pkg-config: gstreamer-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/gst.h>
import "C"

// IDTakeValue wraps gst_structure_id_take_value
//
// The function takes the following parameters:
//
//   - field glib.Quark: a #GQuark representing a field
//   - value *gobject.Value: the new value of the field
//
// Sets the field with the given GQuark @field to @value.  If the field
// does not exist, it is created.  If the field exists, the previous
// value is replaced and freed.
func (structure *Structure) IDTakeValue(field glib.Quark, value any) {
	var carg0 *C.GstStructure // in, none, converted
	var carg1 C.GQuark        // in, none, casted, alias
	var carg2 *C.GValue       // in, full, converted

	carg0 = (*C.GstStructure)(UnsafeStructureToGlibNone(structure))
	carg1 = C.GQuark(field)
	v := gobject.NewValue(value)
	carg2 = (*C.GValue)(gobject.UnsafeValueToGlibFull(v))

	C.gst_structure_id_take_value(carg0, carg1, carg2)
	runtime.KeepAlive(structure)
	runtime.KeepAlive(field)
	runtime.KeepAlive(value)
	runtime.KeepAlive(v)
}

// TakeValue wraps gst_structure_take_value
//
// The function takes the following parameters:
//
//   - fieldname string: the name of the field to set
//   - value *gobject.Value: the new value of the field
//
// Sets the field with the given name @field to @value.  If the field
// does not exist, it is created.  If the field exists, the previous
// value is replaced and freed. The function will take ownership of @value.
func (structure *Structure) TakeValue(fieldname string, value any) {
	var carg0 *C.GstStructure // in, none, converted
	var carg1 *C.gchar        // in, none, string, casted *C.gchar
	var carg2 *C.GValue       // in, full, converted

	carg0 = (*C.GstStructure)(UnsafeStructureToGlibNone(structure))
	carg1 = (*C.gchar)(unsafe.Pointer(C.CString(fieldname)))
	defer C.free(unsafe.Pointer(carg1))
	v := gobject.NewValue(value)
	carg2 = (*C.GValue)(gobject.UnsafeValueToGlibFull(v))

	C.gst_structure_take_value(carg0, carg1, carg2)
	runtime.KeepAlive(structure)
	runtime.KeepAlive(fieldname)
	runtime.KeepAlive(value)
	runtime.KeepAlive(v)
}

// GetValue wraps gst_structure_get_value
//
// The function takes the following parameters:
//
//   - fieldname string: the name of the field to get
//
// The function returns the following values:
//
//   - goret any
//
// Get the value of the field with name @fieldname.
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
//
// The function takes the following parameters:
//
//   - fieldname string: the name of the field to set
//   - value any: the new value of the field
//
// Sets the field with the given name @field to @value.  If the field
// does not exist, it is created.  If the field exists, the previous
// value is replaced and freed.
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

// IDGetValue wraps gst_structure_id_get_value
//
// The function takes the following parameters:
//
//   - field glib.Quark: the #GQuark of the field to get
//
// The function returns the following values:
//
//   - goret any
//
// Get the value of the field with GQuark @field.
func (structure *Structure) IDGetValue(field glib.Quark) any {
	var carg0 *C.GstStructure // in, none, converted
	var carg1 C.GQuark        // in, none, casted, alias
	var cret *C.GValue        // return, none, converted

	carg0 = (*C.GstStructure)(UnsafeStructureToGlibNone(structure))
	carg1 = C.GQuark(field)

	cret = C.gst_structure_id_get_value(carg0, carg1)
	runtime.KeepAlive(structure)
	runtime.KeepAlive(field)

	var goret any

	goret = gobject.ValueFromNative(unsafe.Pointer(cret)).GoValue()

	return goret
}
