package gst

// #include "gst.go.h"
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// Caps is a go wrapper around GstCaps.
type Caps struct {
	native *C.GstCaps
}

// NewAnyCaps creates a new caps that indicate compatibility with any format.
func NewAnyCaps() *Caps { return wrapCaps(C.gst_caps_new_any()) }

// NewEmptyCaps creates a new empty caps object. This is essentially the opposite of
// NewAnyCamps.
func NewEmptyCaps() *Caps { return wrapCaps(C.gst_caps_new_empty()) }

// NewEmptySimpleCaps returns a new empty caps object with the given media format.
func NewEmptySimpleCaps(mediaFormat string) *Caps {
	cFormat := C.CString(mediaFormat)
	defer C.free(unsafe.Pointer(cFormat))
	caps := C.gst_caps_new_empty_simple(cFormat)
	return wrapCaps(caps)
}

// NewFullCaps creates a new caps from the given structures.
func NewFullCaps(structures ...*Structure) *Caps {
	caps := NewEmptyCaps()
	for _, st := range structures {
		caps.AppendStructure(st)
	}
	return caps
}

// NewSimpleCaps creates new caps with the given media format and key value pairs.
// The key of each pair must be a string, followed by any field that can be converted
// to a GType.
func NewSimpleCaps(mediaFormat string, fieldVals ...interface{}) (*Caps, error) {
	if len(fieldVals)%2 != 0 {
		return nil, errors.New("Received odd number of key/value pairs")
	}
	caps := NewEmptySimpleCaps(mediaFormat)
	strParts := make([]string, 0)
	for i := 0; i < len(fieldVals); i = i + 2 {
		fieldKey, ok := fieldVals[i].(string)
		if !ok {
			return nil, errors.New("One or more field keys are not a valid string")
		}
		strParts = append(strParts, fmt.Sprintf("%s=%v", fieldKey, fieldVals[i+1]))
	}
	structure := NewStructureFromString(strings.Join(strParts, ", "))
	if structure == nil {
		return nil, errors.New("Could not build structure from the provided arguments")
	}
	caps.AppendStructure(structure)
	return caps, nil
}

// NewCapsFromString creates a new Caps object from the given string.
func NewCapsFromString(capsStr string) *Caps {
	cStr := C.CString(capsStr)
	defer C.free(unsafe.Pointer(cStr))
	caps := C.gst_caps_from_string(cStr)
	if caps == nil {
		return nil
	}
	return wrapCaps(caps)
}

// NewRawCaps returns new GstCaps with the given format, sample-rate, and channels.
func NewRawCaps(format string, rate, channels int) *Caps {
	capsStr := fmt.Sprintf("audio/x-raw, format=%s, rate=%d, channels=%d", format, rate, channels)
	return NewCapsFromString(capsStr)
}

func (c *Caps) unsafe() unsafe.Pointer { return unsafe.Pointer(c.native) }

// Ref increases the ref count on these caps by one.
func (c *Caps) Ref() { C.gst_caps_ref(c.Instance()) }

// Unref decreases the ref count on these caps by one.
func (c *Caps) Unref() { C.gst_caps_unref(c.Instance()) }

// Instance returns the native GstCaps instance
func (c *Caps) Instance() *C.GstCaps { return C.toGstCaps(c.unsafe()) }

// IsAny returns true if these caps match any media format.
func (c *Caps) IsAny() bool { return gobool(C.gst_caps_is_any(c.Instance())) }

// IsEmpty returns true if these caps are empty.
func (c *Caps) IsEmpty() bool { return gobool(C.gst_caps_is_empty(c.Instance())) }

// AppendStructure appends the given structure to this caps instance.
func (c *Caps) AppendStructure(st *Structure) {
	C.gst_caps_append_structure(c.Instance(), st.Instance())
}

// Append appends the given caps element to these caps. These caps take ownership
// over the given object. If either caps are ANY, the resulting caps will be ANY.
func (c *Caps) Append(caps *Caps) {
	C.gst_caps_append(c.Instance(), caps.Instance())
}

// Size returns the number of structures inside this caps instance.
func (c *Caps) Size() int { return int(C.gst_caps_get_size(c.Instance())) }

// GetStructureAt returns the structure at the given index, or nil if none exists.
func (c *Caps) GetStructureAt(idx int) *Structure {
	st := C.gst_caps_get_structure(c.Instance(), C.guint(idx))
	if st == nil {
		return nil
	}
	return wrapStructure(st)
}

// GetFeaturesAt returns the feature at the given index, or nil if none exists.
func (c *Caps) GetFeaturesAt(idx int) *CapsFeatures {
	feats := C.gst_caps_get_features(c.Instance(), C.guint(idx))
	if feats == nil {
		return nil
	}
	return wrapCapsFeatures(feats)
}

func wrapCaps(caps *C.GstCaps) *Caps { return &Caps{native: caps} }

// CapsFeatures is a go representation of GstCapsFeatures.
type CapsFeatures struct{ native *C.GstCapsFeatures }

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

// Instance returns the native underlying GstCapsFeatures instance.
func (c *CapsFeatures) Instance() *C.GstCapsFeatures {
	return C.toGstCapsFeatures(unsafe.Pointer(c.native))
}

// String implements a stringer on caps features.
func (c *CapsFeatures) String() string {
	return C.GoString(C.gst_caps_features_to_string(c.Instance()))
}

// IsAny returns true if these features match any.
func (c *CapsFeatures) IsAny() bool { return gobool(C.gst_caps_features_is_any(c.Instance())) }

// Equal returns true if the given CapsFeatures are equal to the provided ones.
// If the provided structure is nil, this function immediately returns false.
func (c *CapsFeatures) Equal(feats *CapsFeatures) bool {
	if feats == nil {
		return false
	}
	return gobool(C.gst_caps_features_is_equal(c.Instance(), feats.Instance()))
}

// Go casting of pre-baked caps features
var (
	CapsFeaturesAny = wrapCapsFeatures(C.gst_caps_features_new_any())
)

// Go casting of caps features constants
const (
	CapsFeatureMemorySystemMemory = C.GST_CAPS_FEATURE_MEMORY_SYSTEM_MEMORY
)

func wrapCapsFeatures(features *C.GstCapsFeatures) *CapsFeatures {
	return &CapsFeatures{native: features}
}
