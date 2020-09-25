package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	gopointer "github.com/mattn/go-pointer"
)

// Structure is a go implementation of a C GstStructure.
type Structure struct {
	ptr   unsafe.Pointer
	gType glib.Type
}

// NewStructure returns a new empty structure with the given name.
func NewStructure(name string) *Structure {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	structure := C.gst_structure_new_empty(cName)
	return wrapStructure(structure)
}

// NewStructureFromString builds a new GstStructure from the given string.
func NewStructureFromString(stStr string) *Structure {
	cStr := C.CString(stStr)
	defer C.free(unsafe.Pointer(cStr))
	structure := C.gst_structure_from_string(cStr, nil)
	if structure == nil {
		return nil
	}
	return wrapStructure(structure)
}

// StructureFromGValue extracts the GstStructure from a glib.Value.
func StructureFromGValue(gval *glib.Value) *Structure {
	st := C.gst_value_get_structure((*C.GValue)(gval.Native()))
	return wrapStructure(st)
}

// Instance returns the native GstStructure instance.
func (s *Structure) Instance() *C.GstStructure { return C.toGstStructure(s.ptr) }

// Free frees the memory for the underlying GstStructure.
func (s *Structure) Free() { C.gst_structure_free(s.Instance()) }

// String implement a stringer on a GstStructure.
func (s *Structure) String() string {
	str := C.gst_structure_to_string(s.Instance())
	defer C.g_free((C.gpointer)(str))
	return C.GoString(str)
}

// Name returns the name of this structure.
func (s *Structure) Name() string {
	return C.GoString(C.gst_structure_get_name(s.Instance()))
}

// Size returns the number of fields inside this structure.
func (s *Structure) Size() int {
	return int(C.gst_structure_n_fields(s.Instance()))
}

// SetValue sets the data at key to the given value.
func (s *Structure) SetValue(key string, value interface{}) error {
	gVal, err := glib.GValue(value)
	if err != nil {
		return err
	}
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	C.gst_structure_set_value(s.Instance(), cKey, (*C.GValue)(gVal.GetPointer()))
	return nil
}

// GetValue retrieves the value at key.
func (s *Structure) GetValue(key string) (interface{}, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	gVal := C.gst_structure_get_value(s.Instance(), cKey)
	if gVal == nil {
		return nil, fmt.Errorf("No value exists at %s", key)
	}
	return glib.ValueFromNative(unsafe.Pointer(gVal)).GoValue()
}

// RemoveValue removes the value at the given key. If the key does not exist,
// the structure is unchanged.
func (s *Structure) RemoveValue(key string) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	C.gst_structure_remove_field(s.Instance(), cKey)
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

// Values returns a map of all the values inside this structure. If values cannot be
// converted to an equivalent go type, they are serialized to a string.
func (s *Structure) Values() map[string]interface{} {
	out := make(map[string]interface{})
	resCh := make(chan interface{})
	chPtr := gopointer.Save(resCh)
	defer gopointer.Unref(chPtr)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < s.Size(); i++ {
			nameIface := <-resCh
			valIface := <-resCh
			fieldName := nameIface.(string)
			out[fieldName] = valIface
		}
	}()

	C.gst_structure_foreach(s.Instance(), C.GstStructureForeachFunc(C.structureForEach), (C.gpointer)(unsafe.Pointer(chPtr)))
	wg.Wait()

	return out
}

func wrapStructure(st *C.GstStructure) *Structure {
	return &Structure{
		ptr:   unsafe.Pointer(st),
		gType: glib.Type(st._type),
	}
}
