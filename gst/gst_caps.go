package gst

/*
#include "gst.go.h"

extern gboolean goCapsMapFunc (GstCapsFeatures * features, GstStructure * structure, gpointer user_data);

gboolean cgoCapsMapFunc (GstCapsFeatures * features, GstStructure * structure, gpointer user_data)
{
	return goCapsMapFunc(features, structure, user_data);
}

*/
import "C"

import (
	"runtime"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// TypeCaps is the static Glib Type for a GstCaps.
var TypeCaps = glib.Type(C.gst_caps_get_type())

var _ glib.ValueTransformer = &Caps{}

// Caps is a go wrapper around GstCaps.
type Caps struct {
	native *C.GstCaps
}

// FromGstCapsUnsafeNone wraps the pointer to the given C GstCaps with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
// A ref is taken on the caps and finalizer placed on the object.
func FromGstCapsUnsafeNone(caps unsafe.Pointer) *Caps {
	if caps == nil {
		return nil
	}
	gocaps := ToGstCaps(caps)
	gocaps.Ref()
	runtime.SetFinalizer(gocaps, (*Caps).Unref)
	return gocaps
}

// FromGstCapsUnsafeFull wraps the pointer to the given C GstCaps with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
// A finalizer is placed on the object to Unref after leaving scope.
func FromGstCapsUnsafeFull(caps unsafe.Pointer) *Caps {
	if caps == nil {
		return nil
	}
	gocaps := ToGstCaps(caps)
	runtime.SetFinalizer(gocaps, (*Caps).Unref)
	return gocaps
}

// ToGstCaps converts the given pointer into a Caps without affecting the ref count or
// placing finalizers.
func ToGstCaps(caps unsafe.Pointer) *Caps {
	return wrapCaps(C.toGstCaps(caps))
}

// CapsMapFunc represents a function passed to the Caps MapInPlace, ForEach, and FilterAndMapInPlace methods.
type CapsMapFunc func(features *CapsFeatures, structure *Structure) bool

// NewAnyCaps creates a new caps that indicate compatibility with any format.
//
//   caps := gst.NewAnyCaps()
//   fmt.Println(caps.IsAny())
//   // true
//
func NewAnyCaps() *Caps { return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_new_any())) }

// NewEmptyCaps creates a new empty caps object. This is essentially the opposite of
// NewAnyCamps.
//
//   caps := gst.NewEmptyCaps()
//   fmt.Println(caps.IsEmpty())
//   // true
//
func NewEmptyCaps() *Caps { return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_new_empty())) }

// NewEmptySimpleCaps returns a new empty caps object with the given media format.
//
//   caps := gst.NewEmptySimpleCaps("audio/x-raw")
//   fmt.Println(caps.String())
//   // audio/x-raw
//
func NewEmptySimpleCaps(mediaFormat string) *Caps {
	cFormat := C.CString(mediaFormat)
	defer C.free(unsafe.Pointer(cFormat))
	caps := C.gst_caps_new_empty_simple(cFormat)
	return FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// NewFullCaps creates a new caps from the given structures.
func NewFullCaps(structures ...*Structure) *Caps {
	caps := NewEmptyCaps()
	for _, st := range structures {
		caps.AppendStructure(st)
	}
	return caps
}

// NewCapsFromString creates a new Caps object from the given string.
//
//   caps := gst.NewCapsFromString("audio/x-raw, channels=2")
//   fmt.Println(caps.String())
//   // audio/x-raw, channels=(int)2
//
func NewCapsFromString(capsStr string) *Caps {
	cStr := C.CString(capsStr)
	defer C.free(unsafe.Pointer(cStr))
	caps := C.gst_caps_from_string(cStr)
	if caps == nil {
		return nil
	}
	return FromGstCapsUnsafeFull(unsafe.Pointer(caps))
}

// Unsafe returns an unsafe.Pointer to the underlying GstCaps.
func (c *Caps) Unsafe() unsafe.Pointer { return unsafe.Pointer(c.native) }

// ToGValue returns a GValue containing the given caps.
func (c *Caps) ToGValue() (*glib.Value, error) {
	val, err := glib.ValueInit(glib.Type(C.gst_caps_get_type()))
	if err != nil {
		return nil, err
	}
	C.gst_value_set_caps((*C.GValue)(unsafe.Pointer(val.GValue)), c.Instance())
	return val, nil
}

// Ref increases the ref count on these caps by one.
//
// From this point on, until the caller calls Unref or MakeWritable, it is guaranteed that the caps object
// will not change. This means its structures won't change, etc. To use a Caps object, you must always have a
// refcount on it -- either the one made implicitly by e.g. NewSimpleCaps, or via taking one explicitly with
// this function. Note that when a function provided by these bindings returns caps, or they are converted
// through the FromGstCapsUnsafe methods, a ref is automatically taken if necessary and a runtime Finalizer
// is used to remove it.
func (c *Caps) Ref() *Caps {
	C.gst_caps_ref(c.Instance())
	return c
}

// Unref decreases the ref count on these caps by one.
func (c *Caps) Unref() { C.gst_caps_unref(c.Instance()) }

// Instance returns the native GstCaps instance
func (c *Caps) Instance() *C.GstCaps { return C.toGstCaps(c.Unsafe()) }

// MakeWritable returns a writable copy of caps.
func (c *Caps) MakeWritable() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.makeCapsWritable(c.Instance())))
}

// String implements a stringer on a caps instance. This same string can be used for NewCapsFromString.
func (c *Caps) String() string {
	cStr := C.gst_caps_to_string(c.Instance())
	defer C.g_free((C.gpointer)(unsafe.Pointer(cStr)))
	return C.GoString(cStr)
}

// AppendStructure appends the given structure to this caps instance.
func (c *Caps) AppendStructure(st *Structure) {
	C.gst_caps_append_structure(c.Instance(), st.Instance())
}

// AppendStructureFull appends structure with features to caps. The structure is not copied;
// caps becomes the owner of structure.
func (c *Caps) AppendStructureFull(st *Structure, features *CapsFeatures) {
	C.gst_caps_append_structure_full(c.Instance(), st.Instance(), features.Instance())
}

// Append appends the given caps element to these caps. These caps take ownership
// over the given object. If either caps are ANY, the resulting caps will be ANY.
func (c *Caps) Append(caps *Caps) {
	C.gst_caps_append(c.Instance(), caps.Ref().Instance())
}

// CanIntersect tries intersecting these caps with those given and reports whether the result would not be empty.
func (c *Caps) CanIntersect(caps *Caps) bool {
	return gobool(C.gst_caps_can_intersect(c.Instance(), caps.Instance()))
}

// Copy creates a new Caps as a copy of these. The new caps will have a refcount of 1, owned by the caller.
// The structures are copied as well.
//
// Note that this function is the semantic equivalent of a Ref followed by a MakeWritable. If you only want to hold
// on to a reference to the data, you should use Ref.
//
// When you are finished with the caps, call Unref on it.
func (c *Caps) Copy() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_copy(c.Instance())))
}

// CopyNth creates a new GstCaps and appends a copy of the nth structure contained in caps.
func (c *Caps) CopyNth(n uint) *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_copy_nth(c.Instance(), C.guint(n))))
}

// FilterAndMapInPlace calls the provided function once for each structure and caps feature in the Caps.
// In contrast to ForEach, the function may modify the structure and features. In contrast to MapInPlace,
// the structure and features are removed from the caps if FALSE is returned from the function. The caps must be mutable.
//
//   caps := gst.NewCapsFromString("audio/x-raw")
//
//   caps.FilterAndMapInPlace(func(features *gst.CapsFeatures, structure *gst.Structure) bool {
//       if features.Contains(gst.CapsFeatureMemorySystemMemory) {
//           fmt.Println("Removing system memory feature")
//           return false
//       }
//       return true
//   })
//
//   fmt.Println(caps.IsEmpty())
//
//   // Removing system memory feature
//   // true
//
func (c *Caps) FilterAndMapInPlace(f CapsMapFunc) {
	ptr := gopointer.Save(f)
	defer gopointer.Unref(ptr)
	C.gst_caps_filter_and_map_in_place(
		c.Instance(),
		C.GstCapsFilterMapFunc(C.cgoCapsMapFunc),
		(C.gpointer)(ptr),
	)
}

// Fixate modifies the given caps into a representation with only fixed values. First the caps will be truncated and
// then the first structure will be fixated with Structure's Fixate.
//
// This function takes ownership of caps and will call gst_caps_make_writable on it so you must not use caps afterwards
// unless you keep an additional reference to it with Ref.
//
// Note that it is not guaranteed that the returned caps have exactly one structure. If caps are empty caps then then
// returned caps will be the empty too and contain no structure at all.
//
// Calling this function with any caps is not allowed.
func (c *Caps) Fixate() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_fixate(c.Ref().Instance())))
}

// ForEach calls the provided function once for each structure and caps feature in the GstCaps. The function must not
// modify the fields. There is an unresolved bug in this function currently and it is better to use MapInPlace instead.
//
//   caps := gst.NewCapsFromString("audio/x-raw")
//
//   caps.ForEach(func(features *gst.CapsFeatures, structure *gst.Structure) bool {
//       fmt.Println(structure)
//       return true
//   })
//
//   // audio/x-raw;
//
func (c *Caps) ForEach(f CapsMapFunc) bool {
	ptr := gopointer.Save(f)
	defer gopointer.Unref(ptr)
	return gobool(C.gst_caps_foreach(
		c.Instance(),
		C.GstCapsForeachFunc(C.cgoCapsMapFunc),
		(C.gpointer)(ptr),
	))
}

// GetSize returns the number of structures inside this caps instance.
func (c *Caps) GetSize() int { return int(C.gst_caps_get_size(c.Instance())) }

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

// CapsIntersectMode represents the modes of caps intersection.
// See the official documentation for more details:
// https://gstreamer.freedesktop.org/documentation/gstreamer/gstcaps.html?gi-language=c#GstCapsIntersectMode
type CapsIntersectMode int

// Type castings of intersect modes
const (
	CapsIntersectZigZag CapsIntersectMode = C.GST_CAPS_INTERSECT_ZIG_ZAG
	CapsIntersectFirst  CapsIntersectMode = C.GST_CAPS_INTERSECT_FIRST
)

// Intersect creates a new Caps that contains all the formats that are common to both these caps and those given.
// Defaults to CapsIntersectZigZag mode.
func (c *Caps) Intersect(caps *Caps) *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_intersect(c.Instance(), caps.Instance())))
}

// IntersectFull creates a new Caps that contains all the formats that are common to both these caps those given.
// The order is defined by the CapsIntersectMode used.
func (c *Caps) IntersectFull(caps *Caps, mode CapsIntersectMode) *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_intersect_full(c.Instance(), caps.Instance(), C.GstCapsIntersectMode(mode))))
}

// IsAlwaysCompatible returns if this structure is always compatible with another if every media format that is in
// the first is also contained in the second. That is, these caps are a subset of those given.
func (c *Caps) IsAlwaysCompatible(caps *Caps) bool {
	return gobool(C.gst_caps_is_always_compatible(c.Instance(), caps.Instance()))
}

// IsAny returns true if these caps match any media format.
func (c *Caps) IsAny() bool { return gobool(C.gst_caps_is_any(c.Instance())) }

// IsEmpty returns true if these caps are empty.
func (c *Caps) IsEmpty() bool { return gobool(C.gst_caps_is_empty(c.Instance())) }

// IsEqual returns true if the caps given represent the same set as these.
func (c *Caps) IsEqual(caps *Caps) bool {
	return gobool(C.gst_caps_is_equal(c.Instance(), caps.Instance()))
}

// IsEqualFixed tests if the Caps are equal. This function only works on fixed Caps.
func (c *Caps) IsEqualFixed(caps *Caps) bool {
	return gobool(C.gst_caps_is_equal_fixed(c.Instance(), caps.Instance()))
}

// IsFixed returns true if these caps are fixed, that is, they describe exactly one format.
func (c *Caps) IsFixed() bool { return gobool(C.gst_caps_is_fixed(c.Instance())) }

// IsStrictlyEqual checks if the given caps are exactly the same set of caps.
func (c *Caps) IsStrictlyEqual(caps *Caps) bool {
	return gobool(C.gst_caps_is_strictly_equal(c.Instance(), caps.Instance()))
}

// IsSubset checks if all caps represented by these are also represented by those given.
func (c *Caps) IsSubset(caps *Caps) bool {
	return gobool(C.gst_caps_is_subset(c.Instance(), caps.Instance()))
}

// IsSubsetStructure checks if the given structure is a subset of these caps.
func (c *Caps) IsSubsetStructure(structure *Structure) bool {
	return gobool(C.gst_caps_is_subset_structure(c.Instance(), structure.Instance()))
}

// IsSubsetStructureFull checks if the given structure is a subset of these caps with features.
func (c *Caps) IsSubsetStructureFull(structure *Structure, features *CapsFeatures) bool {
	if features == nil {
		return gobool(C.gst_caps_is_subset_structure_full(
			c.Instance(), structure.Instance(), nil,
		))
	}
	return gobool(C.gst_caps_is_subset_structure_full(
		c.Instance(), structure.Instance(), features.Instance(),
	))
}

// IsWritable returns true if these caps are writable.
func (c *Caps) IsWritable() bool {
	return gobool(C.capsIsWritable(c.Instance()))
}

// MapInPlace calls the provided function once for each structure and caps feature in the Caps.
// In contrast to ForEach, the function may modify, but not delete, the structures and features.
// The caps must be mutable.
func (c *Caps) MapInPlace(f CapsMapFunc) bool {
	ptr := gopointer.Save(f)
	defer gopointer.Unref(ptr)
	return gobool(C.gst_caps_map_in_place(
		c.Instance(),
		C.GstCapsMapFunc(C.cgoCapsMapFunc),
		(C.gpointer)(ptr),
	))
}

// Merge appends the structures contained in the given caps if they are not yet expressed by these.
// The structures in the given caps are not copied -- they are transferred to a writable copy of these ones,
// and then those given are freed. If either caps are ANY, the resulting caps will be ANY.
func (c *Caps) Merge(caps *Caps) *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_merge(c.Ref().Instance(), caps.Ref().Instance())))
}

// MergeStructure appends structure to caps if its not already expressed by caps.
func (c *Caps) MergeStructure(structure *Structure) *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_merge_structure(c.Ref().Instance(), structure.Instance())))
}

// MergeStructureFull appends structure with features to the caps if its not already expressed.
func (c *Caps) MergeStructureFull(structure *Structure, features *CapsFeatures) *Caps {
	if features == nil {
		return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_merge_structure_full(
			c.Ref().Instance(), structure.Instance(), nil,
		)))
	}
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_merge_structure_full(
		c.Ref().Instance(), structure.Instance(), features.Instance(),
	)))
}

// Normalize returns a Caps that represents the same set of formats as caps, but contains no lists.
// Each list is expanded into separate GstStructures.
//
// This function takes ownership of caps and will call MakeWritable on it so you must not
// use caps afterwards unless you keep an additional reference to it with Ref.
func (c *Caps) Normalize() *Caps {
	return wrapCaps(C.gst_caps_normalize(c.Instance()))
}

// RemoveStructureAt removes the structure with the given index from the list of structures.
func (c *Caps) RemoveStructureAt(idx uint) {
	C.gst_caps_remove_structure(c.Instance(), C.guint(idx))
}

// SetFeaturesAt sets the CapsFeatures features for the structure at index.
func (c *Caps) SetFeaturesAt(idx uint, features *CapsFeatures) {
	if features == nil {
		C.gst_caps_set_features(c.Instance(), C.guint(idx), nil)
		return
	}
	C.gst_caps_set_features(c.Instance(), C.guint(idx), features.Instance())
}

// SetFeaturesSimple sets the CapsFeatures for all the structures of these caps.
func (c *Caps) SetFeaturesSimple(features *CapsFeatures) {
	if features == nil {
		C.gst_caps_set_features_simple(c.Instance(), nil)
		return
	}
	C.gst_caps_set_features_simple(c.Instance(), features.Instance())
}

// SetValue sets the given field on all structures of caps to the given value. This is a convenience
// function for calling SetValue on all structures of caps. If the value cannot be coerced to a C type,
// then nothing will happen.
func (c *Caps) SetValue(field string, val interface{}) {
	gVal, err := glib.GValue(val)
	if err != nil {
		return
	}
	C.gst_caps_set_value(
		c.Instance(),
		(*C.gchar)(unsafe.Pointer(C.CString(field))),
		(*C.GValue)(unsafe.Pointer(gVal.GValue)),
	)
}

// Simplify converts the given caps into a representation that represents the same set of formats, but in a
// simpler form. Component structures that are identical are merged. Component structures that have values
// that can be merged are also merged.
//
// This function takes ownership of caps and will call MakeWritable on it if necessary, so you must not use
// caps afterwards unless you keep an additional reference to it with Ref.
//
// This method does not preserve the original order of caps.
func (c *Caps) Simplify() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_simplify(c.Ref().Instance())))
}

// StealStructureAt retrieves the structure with the given index from the list of structures contained in caps.
// The caller becomes the owner of the returned structure.
func (c *Caps) StealStructureAt(idx uint) *Structure {
	return wrapStructure(C.gst_caps_steal_structure(c.Instance(), C.guint(idx)))
}

// Subtract subtracts the given caps from these.
func (c *Caps) Subtract(caps *Caps) *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_subtract(c.Instance(), caps.Instance())))
}

// Truncate discards all but the first structure from caps. Useful when fixating.
//
// This function takes ownership of caps and will call gst_caps_make_writable on it if necessary, so you must not
// use caps afterwards unless you keep an additional reference to it with Ref.
//
// Note that it is not guaranteed that the returned caps have exactly one structure. If caps is any or empty caps
// then then returned caps will be the same and contain no structure at all.
func (c *Caps) Truncate() *Caps {
	return FromGstCapsUnsafeFull(unsafe.Pointer(C.gst_caps_truncate(c.Ref().Instance())))
}
