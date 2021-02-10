package gst

/*
#include "gst.go.h"

extern gboolean structForEachCb  (GQuark field_id, GValue * value, gpointer user_data);

gboolean structureForEach (GQuark field_id, GValue * value, gpointer user_data)
{
	return structForEachCb(field_id, value, user_data);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
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

// MarshalStructure will convert the given go struct into a GstStructure. Currently nested
// structs are not supported.
func MarshalStructure(data interface{}) *Structure {
	typeOf := reflect.TypeOf(data)
	valsOf := reflect.ValueOf(data)
	st := NewStructure(typeOf.Name())
	for i := 0; i < valsOf.NumField(); i++ {
		gval := valsOf.Field(i).Interface()
		fieldName := typeOf.Field(i).Name
		if err := st.SetValue(fieldName, gval); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set %v for %s", gval, fieldName)
		}
	}
	return st
}

// FromGstStructureUnsafe wraps the given unsafe.Pointer in a Structure. This is meant for internal usage
// and is exported for visibility to other packages.
func FromGstStructureUnsafe(st unsafe.Pointer) *Structure {
	return wrapStructure((*C.GstStructure)(st))
}

// UnmarshalInto will unmarshal this structure into the given pointer. The object
// reflected by the pointer must be non-nil.
func (s *Structure) UnmarshalInto(data interface{}) error {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Data is invalid (nil or non-pointer)")
	}

	val := reflect.ValueOf(data).Elem()
	nVal := rv.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		fieldName := val.Type().Field(i).Name
		val, err := s.GetValue(fieldName)
		if err == nil {
			nvField.Set(reflect.ValueOf(val))
		}
	}

	return nil
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
	C.gst_structure_set_value(s.Instance(), cKey, (*C.GValue)(unsafe.Pointer(gVal.GValue)))
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

// TypeStructure is the glib.Type for a Structure.
var TypeStructure = glib.Type(C.gst_structure_get_type())

var _ glib.ValueTransformer = &Structure{}

// ToGValue implements a glib.ValueTransformer
func (s *Structure) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeStructure)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_structure(
		(*C.GValue)(unsafe.Pointer(val.GValue)),
		s.Instance(),
	)
	return val, nil
}

func wrapStructure(st *C.GstStructure) *Structure {
	return &Structure{
		ptr:   unsafe.Pointer(st),
		gType: glib.Type(st._type),
	}
}
