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
	"unsafe"
)

// NewElement is a generic wrapper around `gst_element_factory_make`.
func NewElement(name string) (*Element, error) {
	elemName := C.CString(name)
	defer C.free(unsafe.Pointer(elemName))
	elem := C.gst_element_factory_make((*C.gchar)(elemName), nil)
	if elem == nil {
		return nil, fmt.Errorf("Could not create element: %s", name)
	}
	return wrapElement(elem), nil
}

// NewElementMany is a convenience wrapper around building many GstElements in a
// single function call. It returns an error if the creation of any element fails. A
// map containing the ordinal of the argument to the Element created is returned.
func NewElementMany(elemNames ...string) (map[int]*Element, error) {
	elemMap := make(map[int]*Element)
	for idx, name := range elemNames {
		elem, err := NewElement(name)
		if err != nil {
			return nil, err
		}
		elemMap[idx] = elem
	}
	return elemMap, nil
}

// ElementFactory wraps the GstElementFactory
type ElementFactory struct{ *Object }

// Find returns the factory for the given plugin, or nil if it doesn't exist. Unref after usage.
func Find(name string) *ElementFactory {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	factory := C.gst_element_factory_find((*C.gchar)(cName))
	if factory == nil {
		return nil
	}
	return wrapElementFactory(factory)
}

// Instance returns the C GstFactory instance
func (e *ElementFactory) Instance() *C.GstElementFactory { return C.toGstElementFactory(e.unsafe()) }

// CanSinkAllCaps checks if the factory can sink all possible capabilities.
func (e *ElementFactory) CanSinkAllCaps(caps *C.GstCaps) bool {
	return gobool(C.gst_element_factory_can_sink_all_caps((*C.GstElementFactory)(e.Instance()), (*C.GstCaps)(caps)))
}

// CanSinkAnyCaps checks if the factory can sink any possible capability.
func (e *ElementFactory) CanSinkAnyCaps(caps *C.GstCaps) bool {
	return gobool(C.gst_element_factory_can_sink_any_caps((*C.GstElementFactory)(e.Instance()), (*C.GstCaps)(caps)))
}

// CanSourceAllCaps checks if the factory can src all possible capabilities.
func (e *ElementFactory) CanSourceAllCaps(caps *C.GstCaps) bool {
	return gobool(C.gst_element_factory_can_src_all_caps((*C.GstElementFactory)(e.Instance()), (*C.GstCaps)(caps)))
}

// CanSourceAnyCaps checks if the factory can src any possible capability.
func (e *ElementFactory) CanSourceAnyCaps(caps *C.GstCaps) bool {
	return gobool(C.gst_element_factory_can_src_any_caps((*C.GstElementFactory)(e.Instance()), (*C.GstCaps)(caps)))
}

// GetMetadata gets the metadata on this factory with key.
func (e *ElementFactory) GetMetadata(key string) string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	res := C.gst_element_factory_get_metadata((*C.GstElementFactory)(e.Instance()), (*C.gchar)(ckey))
	defer C.free(unsafe.Pointer(res))
	return C.GoString(res)
}

// GetMetadataKeys gets the available keys for the metadata on this factory.
func (e *ElementFactory) GetMetadataKeys() []string {
	keys := C.gst_element_factory_get_metadata_keys((*C.GstElementFactory)(e.Instance()))
	if keys == nil {
		return nil
	}
	defer C.g_strfreev(keys)
	size := C.sizeOfGCharArray(keys)
	return goStrings(size, keys)
}

func wrapElementFactory(factory *C.GstElementFactory) *ElementFactory {
	return &ElementFactory{wrapObject(C.toGstObject(unsafe.Pointer(factory)))}
}
