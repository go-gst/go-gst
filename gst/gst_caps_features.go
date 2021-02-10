package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Go casting of pre-baked caps features
var (
	CapsFeatureMemorySystemMemory string = C.GST_CAPS_FEATURE_MEMORY_SYSTEM_MEMORY
)

// CapsFeatures is a go representation of GstCapsFeatures.
type CapsFeatures struct {
	native *C.GstCapsFeatures
}

// NewCapsFeaturesEmpty returns a new empty CapsFeatures.
//
//   feats := gst.NewCapsFeaturesEmpty()
//   fmt.Println(feats.GetSize())
//   // 0
//
func NewCapsFeaturesEmpty() *CapsFeatures { return wrapCapsFeatures(C.gst_caps_features_new_empty()) }

// NewCapsFeaturesAny returns a new ANY CapsFeatures.
//
//   feats := gst.NewCapsFeaturesAny()
//   fmt.Println(feats.IsAny())
//   // true
//
func NewCapsFeaturesAny() *CapsFeatures { return wrapCapsFeatures(C.gst_caps_features_new_any()) }

// NewCapsFeaturesFromString creates new CapsFeatures from the given string.
func NewCapsFeaturesFromString(features string) *CapsFeatures {
	cStr := C.CString(features)
	defer C.free(unsafe.Pointer(cStr))
	capsFeatures := C.gst_caps_features_from_string(cStr)
	if capsFeatures == nil {
		return nil
	}
	return wrapCapsFeatures(capsFeatures)
}

// TypeCapsFeatures is the glib.Type for a CapsFeatures.
var TypeCapsFeatures = glib.Type(C.gst_caps_features_get_type())

var _ glib.ValueTransformer = &CapsFeatures{}

// ToGValue implements a glib.ValueTransformer
func (c *CapsFeatures) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(TypeCapsFeatures)
	if err != nil {
		return nil, err
	}
	C.gst_value_set_caps_features(
		(*C.GValue)(unsafe.Pointer(val.GValue)),
		c.Instance(),
	)
	return val, nil
}

// Instance returns the native underlying GstCapsFeatures instance.
func (c *CapsFeatures) Instance() *C.GstCapsFeatures {
	return C.toGstCapsFeatures(unsafe.Pointer(c.native))
}

// String implements a stringer on caps features.
//
//   feats := gst.NewCapsFeaturesAny()
//   fmt.Println(feats.String())
//   // ANY
//
func (c *CapsFeatures) String() string {
	return C.GoString(C.gst_caps_features_to_string(c.Instance()))
}

// Add adds the given feature to these.
//
//   feats := gst.NewCapsFeaturesEmpty()
//
//   fmt.Println(feats.GetSize())
//
//   feats.Add(gst.CapsFeatureMemorySystemMemory)
//
//   fmt.Println(feats.GetSize())
//   fmt.Println(feats.Contains(gst.CapsFeatureMemorySystemMemory))
//   fmt.Println(feats.GetNth(0))
//
//   // 0
//   // 1
//   // true
//   // memory:SystemMemory
//
func (c *CapsFeatures) Add(feature string) {
	cStr := C.CString(feature)
	defer C.free(unsafe.Pointer(cStr))
	C.gst_caps_features_add(
		c.Instance(),
		(*C.gchar)(unsafe.Pointer(cStr)),
	)
}

// Contains returns true if the given feature is included in these.
func (c *CapsFeatures) Contains(feature string) bool {
	cStr := C.CString(feature)
	defer C.free(unsafe.Pointer(cStr))
	return gobool(C.gst_caps_features_contains(
		c.Instance(),
		(*C.gchar)(unsafe.Pointer(cStr)),
	))
}

// Copy duplicates these features and all of it's values.
func (c *CapsFeatures) Copy() *CapsFeatures {
	return wrapCapsFeatures(C.gst_caps_features_copy(c.Instance()))
}

// Free frees the memory containing these features. Only call this if you
// do not intend to pass these features to other methods.
func (c *CapsFeatures) Free() { C.gst_caps_features_free(c.Instance()) }

// GetNth returns the feature at index.
func (c *CapsFeatures) GetNth(idx uint) string {
	feat := C.gst_caps_features_get_nth(c.Instance(), C.guint(idx))
	return C.GoString(feat)
}

// GetSize returns the number of features.
func (c *CapsFeatures) GetSize() uint {
	return uint(C.gst_caps_features_get_size(c.Instance()))
}

// IsAny returns true if these features match any.
func (c *CapsFeatures) IsAny() bool { return gobool(C.gst_caps_features_is_any(c.Instance())) }

// IsEqual returns true if the given CapsFeatures are equal to the provided ones.
// If the provided structure is nil, this function immediately returns false.
func (c *CapsFeatures) IsEqual(feats *CapsFeatures) bool {
	if feats == nil {
		return false
	}
	return gobool(C.gst_caps_features_is_equal(c.Instance(), feats.Instance()))
}

// Remove removes the given feature.
func (c *CapsFeatures) Remove(feature string) {
	cStr := C.CString(feature)
	defer C.free(unsafe.Pointer(cStr))
	C.gst_caps_features_remove(
		c.Instance(),
		(*C.gchar)(unsafe.Pointer(cStr)),
	)
}

// SetParentRefCount sets the parent_refcount field of CapsFeatures. This field is used to determine
// whether a caps features is mutable or not. This function should only be called by code implementing
// parent objects of CapsFeatures, as described in the MT Refcounting section of the design documents.
func (c *CapsFeatures) SetParentRefCount(refCount int) bool {
	gcount := C.gint(refCount)
	return gobool(C.gst_caps_features_set_parent_refcount(c.Instance(), &gcount))
}
