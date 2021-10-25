package base

/*
#include "gst.go.h"

#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
)

// TypeFindHelper tries to find what type of data is flowing from the given source GstPad.
// Returns nil if no Caps matches the data stream. Unref after usage.
func TypeFindHelper(pad *gst.Pad, size uint64) *gst.Caps {
	caps := C.gst_type_find_helper((*C.GstPad)(unsafe.Pointer(pad.Instance())), C.guint64(size))
	if caps == nil {
		return nil
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// TypeFindHelperForBuffer tries to find what type of data is contained in the given GstBuffer,
// the assumption being that the buffer represents the beginning of the stream or file.
//
// All available typefinders will be called on the data in order of rank. If a typefinding function
// returns a probability of gst.TypeFindMaximum, typefinding is stopped immediately and the found
// caps will be returned right away. Otherwise, all available typefind functions will the tried, and
// the caps with the highest probability will be returned, or nil if the content of the buffer could
// not be identified.
//
// Object can either be nil or the object doing the typefinding (used for logging). Caps should be unrefed
// after usage.
func TypeFindHelperForBuffer(obj *gst.Object, buffer *gst.Buffer) (*gst.Caps, gst.TypeFindProbability) {
	var prob C.GstTypeFindProbability
	var cobj *C.GstObject
	if obj != nil {
		cobj = (*C.GstObject)(obj.Unsafe())
	}
	caps := C.gst_type_find_helper_for_buffer(cobj, (*C.GstBuffer)(unsafe.Pointer(buffer.Instance())), &prob)
	if caps == nil {
		return nil, gst.TypeFindProbability(prob)
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps)), gst.TypeFindProbability(prob)
}

// TypeFindHelperForBufferWithExtension ries to find what type of data is contained in the given GstBuffer,
// the assumption being that the buffer represents the beginning of the stream or file.
//
// All available typefinders will be called on the data in order of rank. If a typefinding function returns
// a probability of gst.TypeFindMaximum, typefinding is stopped immediately and the found caps will be returned
// right away. Otherwise, all available typefind functions will the tried, and the caps with the highest
// probability will be returned, or nil if the content of the buffer could not be identified.
//
// When extension is not empty, this function will first try the typefind functions for the given extension,
// which might speed up the typefinding in many cases.
//
// Unref caps after usage.
func TypeFindHelperForBufferWithExtension(obj *gst.Object, buffer *gst.Buffer, extension string) (*gst.Caps, gst.TypeFindProbability) {
	var prob C.GstTypeFindProbability
	var cobj *C.GstObject
	var cext *C.gchar
	if obj != nil {
		cobj = (*C.GstObject)(obj.Unsafe())
	}
	if extension != "" {
		cstr := C.CString(extension)
		defer C.free(unsafe.Pointer(cstr))
		cext = (*C.gchar)(unsafe.Pointer(cstr))
	}
	caps := C.gst_type_find_helper_for_buffer_with_extension(cobj, (*C.GstBuffer)(unsafe.Pointer(buffer.Instance())), cext, &prob)
	if caps == nil {
		return nil, gst.TypeFindProbability(prob)
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps)), gst.TypeFindProbability(prob)
}

// TypeFindHelperForData tries to find what type of data is contained in the given data,
// the assumption being that the buffer represents the beginning of the stream or file.
//
// All available typefinders will be called on the data in order of rank. If a typefinding function
// returns a probability of gst.TypeFindMaximum, typefinding is stopped immediately and the found
// caps will be returned right away. Otherwise, all available typefind functions will the tried, and
// the caps with the highest probability will be returned, or nil if the content of the buffer could
// not be identified.
//
// Object can either be nil or the object doing the typefinding (used for logging). Caps should be unrefed
// after usage.
func TypeFindHelperForData(obj *gst.Object, data []byte) (*gst.Caps, gst.TypeFindProbability) {
	var prob C.GstTypeFindProbability
	var cobj *C.GstObject
	if obj != nil {
		cobj = (*C.GstObject)(obj.Unsafe())
	}
	caps := C.gst_type_find_helper_for_data(cobj, (*C.guint8)(unsafe.Pointer(&data[0])), C.gsize(len(data)), &prob)
	if caps == nil {
		return nil, gst.TypeFindProbability(prob)
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps)), gst.TypeFindProbability(prob)
}

// TypeFindHelperForDataWithExtension ries to find what type of data is contained in the given data,
// the assumption being that the buffer represents the beginning of the stream or file.
//
// All available typefinders will be called on the data in order of rank. If a typefinding function returns
// a probability of gst.TypeFindMaximum, typefinding is stopped immediately and the found caps will be returned
// right away. Otherwise, all available typefind functions will the tried, and the caps with the highest
// probability will be returned, or nil if the content of the buffer could not be identified.
//
// When extension is not empty, this function will first try the typefind functions for the given extension,
// which might speed up the typefinding in many cases.
//
// Object can either be nil or the object doing the typefinding (used for logging). Unref caps after usage.
func TypeFindHelperForDataWithExtension(obj *gst.Object, data []byte, extension string) (*gst.Caps, gst.TypeFindProbability) {
	var prob C.GstTypeFindProbability
	var cobj *C.GstObject
	var cext *C.gchar
	if obj != nil {
		cobj = (*C.GstObject)(obj.Unsafe())
	}
	if extension != "" {
		cstr := C.CString(extension)
		defer C.free(unsafe.Pointer(cstr))
		cext = (*C.gchar)(unsafe.Pointer(cstr))
	}
	caps := C.gst_type_find_helper_for_data_with_extension(cobj, (*C.guint8)(unsafe.Pointer(&data[0])), C.gsize(len(data)), cext, &prob)
	if caps == nil {
		return nil, gst.TypeFindProbability(prob)
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps)), gst.TypeFindProbability(prob)
}

// TypeFindHelperForExtension tries to find the best GstCaps associated with extension.
//
// All available typefinders will be checked against the extension in order of rank. The caps of the first typefinder
// that can handle extension will be returned.
//
// Object can either be nil or the object doing the typefinding (used for logging). Unref caps after usage.
func TypeFindHelperForExtension(obj *gst.Object, extension string) *gst.Caps {
	var cobj *C.GstObject
	if obj != nil {
		cobj = (*C.GstObject)(obj.Unsafe())
	}
	cext := C.CString(extension)
	defer C.free(unsafe.Pointer(cext))
	caps := C.gst_type_find_helper_for_extension(cobj, (*C.gchar)(unsafe.Pointer(cext)))
	if caps == nil {
		return nil
	}
	return gst.FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}
