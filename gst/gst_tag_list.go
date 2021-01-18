package gst

/*
#include "gst.go.h"

extern void goTagForEachFunc (const GstTagList * tagList, const gchar * tag, gpointer user_data);

void cgoTagForEachFunc (const GstTagList * tagList, const gchar * tag, gpointer user_data)
{
	goTagForEachFunc(tagList, tag, user_data);
}

*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// TagList is a go wrapper around a GstTagList. For now, until the rest of the methods are
// implemnented, this struct is primarily used for retrieving serialized copies of the tags.
type TagList struct {
	ptr *C.GstTagList
}

// FromGstTagListUnsafeNone wraps the pointer to the given C GstTagList with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstTagListUnsafeNone(tags unsafe.Pointer) *TagList {
	tl := wrapTagList(C.toGstTagList(tags))
	tl.Ref()
	runtime.SetFinalizer(tl, (*TagList).Unref)
	return tl
}

// FromGstTagListUnsafeFull wraps the pointer to the given C GstTagList with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstTagListUnsafeFull(tags unsafe.Pointer) *TagList {
	tl := wrapTagList(C.toGstTagList(tags))
	runtime.SetFinalizer(tl, (*TagList).Unref)
	return tl
}

// NewEmptyTagList returns a new empty tag list.
//
//   tagList := gst.NewEmptyTagList()
//   fmt.Println(tagList.IsEmpty())
//   // true
//
func NewEmptyTagList() *TagList {
	return FromGstTagListUnsafeFull(unsafe.Pointer(C.gst_tag_list_new_empty()))
}

// NewTagListFromString creates a new tag list from the given string. This is the same format produced
// by the stringer interface on the TagList.
func NewTagListFromString(tags string) *TagList {
	ctags := C.CString(tags)
	defer C.free(unsafe.Pointer(ctags))
	tagList := C.gst_tag_list_new_from_string((*C.gchar)(unsafe.Pointer(ctags)))
	if tagList == nil {
		return nil
	}
	return FromGstTagListUnsafeFull(unsafe.Pointer(tagList))
}

// Instance returns the underlying GstTagList instance.
func (t *TagList) Instance() *C.GstTagList { return C.toGstTagList(unsafe.Pointer(t.ptr)) }

// String implements a stringer on the TagList and serializes it to a string.
func (t *TagList) String() string { return C.GoString(C.gst_tag_list_to_string(t.Instance())) }

// Ref increases the ref count on this TagList by one.
func (t *TagList) Ref() *TagList { return wrapTagList(C.gst_tag_list_ref(t.Instance())) }

// Unref decreses the ref count on this TagList by one. When the ref count reaches zero, the object
// is destroyed.
func (t *TagList) Unref() { C.gst_tag_list_unref(t.Instance()) }

// AddValue adds a value to a given tag using the given merge mode. If the value provided
// cannot be coerced to a GValue, nothing will happen.
//
//   tagList := gst.NewEmptyTagList()
//   tagList.AddValue(gst.TagMergeAppend, gst.TagAlbum, "MyNewAlbum")
//   myAlbum, _ := tagList.GetString(gst.TagAlbum)
//   fmt.Println(myAlbum)
//   // MyNewAlbum
//
func (t *TagList) AddValue(mergeMode TagMergeMode, tag Tag, value interface{}) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	gVal, err := glib.GValue(value)
	if err != nil {
		return
	}
	C.gst_tag_list_add_value(
		t.Instance(),
		C.GstTagMergeMode(mergeMode),
		(*C.gchar)(ctag),
		(*C.GValue)(unsafe.Pointer(gVal.GValue)),
	)
}

// AddValues can be used to add multiple values to a tag with the given merge mode.
// Values that cannot be coerced to C types will be ignored.
func (t *TagList) AddValues(mergeMode TagMergeMode, tag Tag, vals ...interface{}) {
	for _, val := range vals {
		t.AddValue(mergeMode, tag, val)
	}
}

// Copy creates a new TagList as a copy of the old taglist. The new taglist will have a refcount of 1,
// owned by the caller, and will be writable as a result.
//
// Note that this function is the semantic equivalent of a Ref followed by a MakeWritable. If you only want
// to hold on to a reference to the data, you should use Ref.
//
// When you are finished with the taglist, call Unref on it.
func (t *TagList) Copy() *TagList {
	return FromGstTagListUnsafeFull(unsafe.Pointer(C.gst_tag_list_copy(t.Instance())))
}

// TagListForEachFunc is a function that will be called in ForEach. The function may not modify the tag list.
type TagListForEachFunc func(tagList *TagList, tag Tag)

// ForEach calls the given function for each tag inside the tag list. Note that if there is no tag,
// the function won't be called at all.
//
//   tagList := gst.NewEmptyTagList()
//
//   tagList.AddValue(gst.TagMergeAppend, gst.TagAlbumArtist, "tinyzimmer")
//   tagList.AddValue(gst.TagMergeAppend, gst.TagAlbum, "GstreamerInGo")
//
//   tagList.ForEach(func(_ *gst.TagList, tag gst.Tag) {
//       val, _ := tagList.GetString(tag)
//       fmt.Println(tag, ":", val)
//   })
//
//   // album-artist : tinyzimmer
//   // album : GstreamerInGo
//
func (t *TagList) ForEach(f TagListForEachFunc) {
	ptr := gopointer.Save(f)
	defer gopointer.Unref(ptr)
	C.gst_tag_list_foreach(
		t.Instance(),
		C.GstTagForeachFunc(C.cgoTagForEachFunc),
		(C.gpointer)(ptr),
	)
}

// GetBool returns the boolean value at the given tag key. If multiple values are associated with the tag they
// are merged.
func (t *TagList) GetBool(tag Tag) (value, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gboolean
	gok := C.gst_tag_list_get_boolean(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return gobool(gout), gobool(gok)
}

// GetBoolIndex retrieves the bool at the given index in the tag key.
func (t *TagList) GetBoolIndex(tag Tag, idx uint) (value, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gboolean
	gok := C.gst_tag_list_get_boolean_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.guint(idx),
		&gout,
	)
	return gobool(gout), gobool(gok)
}

// GetDate returns the date stored at the given tag key. If there are multiple values, the first one
// is returned.
func (t *TagList) GetDate(tag Tag) (value time.Time, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.GDate
	gok := C.gst_tag_list_get_date(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	if gobool(gok) {
		defer C.g_date_free(gout)
		return gdateToTime(gout), true
	}
	return time.Time{}, false
}

// GetDateIndex returns the date stored at the given index in tag key.
func (t *TagList) GetDateIndex(tag Tag, idx uint) (value time.Time, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.GDate
	gok := C.gst_tag_list_get_date_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.guint(idx),
		&gout,
	)
	if gobool(gok) {
		defer C.g_date_free(gout)
		return gdateToTime(gout), true
	}
	return time.Time{}, false
}

// GetDateTime returns the date and time stored at the given tag key. If there are multiple values, the first one
// is returned.
func (t *TagList) GetDateTime(tag Tag) (value time.Time, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.GstDateTime
	gok := C.gst_tag_list_get_date_time(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	if gobool(gok) {
		defer C.gst_date_time_unref(gout)
		return gstDateTimeToTime(gout), true
	}
	return time.Time{}, false
}

// GetDateTimeIndex returns the date and time stored at the given tag key at the given index.
func (t *TagList) GetDateTimeIndex(tag Tag, idx uint) (value time.Time, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.GstDateTime
	gok := C.gst_tag_list_get_date_time_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.guint(idx),
		&gout,
	)
	if gobool(gok) {
		defer C.gst_date_time_unref(gout)
		return gstDateTimeToTime(gout), true
	}
	return time.Time{}, false
}

// GetFloat64 returns the float at the given tag key, merging multiple values into one if multiple values
// are associated with the tag. This is the equivalent of a C double stored in the tag.
func (t *TagList) GetFloat64(tag Tag) (value float64, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gdouble
	gok := C.gst_tag_list_get_double(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return float64(gout), gobool(gok)
}

// GetFloat64Index returns the float at the index of the given tag key.
func (t *TagList) GetFloat64Index(tag Tag, idx uint) (value float64, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gdouble
	gok := C.gst_tag_list_get_double_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.uint(idx),
		&gout,
	)
	return float64(gout), gobool(gok)
}

// GetFloat32 returns the float at the given tag key, merging multiple values into one if multiple values
// are associated with the tag.
func (t *TagList) GetFloat32(tag Tag) (value float32, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gfloat
	gok := C.gst_tag_list_get_float(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return float32(gout), gobool(gok)
}

// GetFloat32Index returns the float at the index of the given tag key.
func (t *TagList) GetFloat32Index(tag Tag, idx uint) (value float32, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gfloat
	gok := C.gst_tag_list_get_float_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.uint(idx),
		&gout,
	)
	return float32(gout), gobool(gok)
}

// GetInt32 returns the integer at the given tag key, merging multiple values into one if multiple values
// are associated with the tag.
func (t *TagList) GetInt32(tag Tag) (value int32, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gint
	gok := C.gst_tag_list_get_int(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return int32(gout), gobool(gok)
}

// GetInt32Index returns the integer at the index of the given tag key.
func (t *TagList) GetInt32Index(tag Tag, idx uint) (value int32, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gint
	gok := C.gst_tag_list_get_int_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.uint(idx),
		&gout,
	)
	return int32(gout), gobool(gok)
}

// GetInt64 returns the integer at the given tag key, merging multiple values into one if multiple values
// are associated with the tag.
func (t *TagList) GetInt64(tag Tag) (value int64, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gint64
	gok := C.gst_tag_list_get_int64(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return int64(gout), gobool(gok)
}

// GetInt64Index returns the integer at the index of the given tag key.
func (t *TagList) GetInt64Index(tag Tag, idx uint) (value int64, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gint64
	gok := C.gst_tag_list_get_int64_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.uint(idx),
		&gout,
	)
	return int64(gout), gobool(gok)
}

// GetPointer returns the C pointer stored at the given tag key, merging values if there are multiple.
func (t *TagList) GetPointer(tag Tag) (value unsafe.Pointer, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gpointer
	gok := C.gst_tag_list_get_pointer(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return unsafe.Pointer(gout), gobool(gok)
}

// GetPointerIndex returns the C pointer stored at the given tag key index.
func (t *TagList) GetPointerIndex(tag Tag, idx uint) (value unsafe.Pointer, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.gpointer
	gok := C.gst_tag_list_get_pointer_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.guint(idx),
		&gout,
	)
	return unsafe.Pointer(gout), gobool(gok)
}

// GetSample copies the first sample for the given tag in the taglist. Free the sample with Unref when it
// is no longer needed.
func (t *TagList) GetSample(tag Tag) (value *Sample, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.GstSample
	gok := C.gst_tag_list_get_sample(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	if gobool(gok) {
		return FromGstSampleUnsafeFull(unsafe.Pointer(gout)), true
	}
	return nil, false
}

// GetSampleIndex copies the sample for the given index in tag in the taglist. Free the sample with Unref
// when it is no longer needed.
func (t *TagList) GetSampleIndex(tag Tag, idx uint) (value *Sample, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.GstSample
	gok := C.gst_tag_list_get_sample_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.guint(idx),
		&gout,
	)
	if gobool(gok) {
		return FromGstSampleUnsafeFull(unsafe.Pointer(gout)), true
	}
	return nil, false
}

// GetString returns the string for the given tag, possibly merging multiple values into one.
func (t *TagList) GetString(tag Tag) (value string, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.gchar
	gok := C.gst_tag_list_get_string(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	defer C.g_free((C.gpointer)(unsafe.Pointer(gout)))
	return C.GoString(gout), gobool(gok)
}

// GetStringIndex returns the string for the given index in tag.
func (t *TagList) GetStringIndex(tag Tag, idx uint) (value string, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.gchar
	gok := C.gst_tag_list_get_string_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.guint(idx),
		&gout,
	)
	defer C.g_free((C.gpointer)(unsafe.Pointer(gout)))
	return C.GoString(gout), gobool(gok)
}

// GetUint32 returns the unsigned integer at the given tag key, merging multiple values into one if multiple values
// are associated with the tag.
func (t *TagList) GetUint32(tag Tag) (value uint32, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.guint
	gok := C.gst_tag_list_get_uint(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return uint32(gout), gobool(gok)
}

// GetUint32Index returns the unsigned integer at the index of the given tag key.
func (t *TagList) GetUint32Index(tag Tag, idx uint) (value uint32, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.guint
	gok := C.gst_tag_list_get_uint_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.uint(idx),
		&gout,
	)
	return uint32(gout), gobool(gok)
}

// GetUint64 returns the unsigned integer at the given tag key, merging multiple values into one if multiple values
// are associated with the tag.
func (t *TagList) GetUint64(tag Tag) (value uint64, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.guint64
	gok := C.gst_tag_list_get_uint64(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		&gout,
	)
	return uint64(gout), gobool(gok)
}

// GetUint64Index returns the unsigned integer at the index of the given tag key.
func (t *TagList) GetUint64Index(tag Tag, idx uint) (value uint64, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout C.guint64
	gok := C.gst_tag_list_get_uint64_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.uint(idx),
		&gout,
	)
	return uint64(gout), gobool(gok)
}

// GetValueIndex retrieves the GValue at the given index in tag, or nil if none exists.
// Note that this function can also return nil if the stored value cannot be cleanly coerced
// to a go type. It is safer to use the other functions provided when you know the expected
// return type.
func (t *TagList) GetValueIndex(tag Tag, idx uint) interface{} {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	gval := C.gst_tag_list_get_value_index(t.Instance(), (*C.gchar)(unsafe.Pointer(ctag)), C.guint(idx))
	if gval == nil {
		return nil
	}
	val := glib.ValueFromNative(unsafe.Pointer(gval))
	iface, _ := val.GoValue()
	return iface
}

// GetScope returns the scope for this TagList.
func (t *TagList) GetScope() TagScope { return TagScope(C.gst_tag_list_get_scope(t.Instance())) }

// GetTagSize returns the number of tag values at the given tag key.
func (t *TagList) GetTagSize(tagKey string) int {
	cStr := C.CString(tagKey)
	defer C.free(unsafe.Pointer(cStr))
	return int(C.gst_tag_list_get_tag_size(t.Instance(), (*C.gchar)(cStr)))
}

// Insert inserts the tags from the provided list using the given merge mode.
func (t *TagList) Insert(tagList *TagList, mergeMode TagMergeMode) {
	C.gst_tag_list_insert(t.Instance(), tagList.Instance(), C.GstTagMergeMode(mergeMode))
}

// IsEmpty returns true if this tag list is empty.
func (t *TagList) IsEmpty() bool { return gobool(C.gst_tag_list_is_empty(t.Instance())) }

// IsEqual checks if the two tag lists are equal.
func (t *TagList) IsEqual(tagList *TagList) bool {
	return gobool(C.gst_tag_list_is_equal(t.Instance(), tagList.Instance()))
}

// IsWritable returns true if this TagList is writable.
func (t *TagList) IsWritable() bool { return gobool(C.tagListIsWritable(t.Instance())) }

// MakeWritable will return a writable copy of the tag list if it is not already so.
func (t *TagList) MakeWritable() *TagList {
	return FromGstTagListUnsafeFull(unsafe.Pointer(C.makeTagListWritable(t.Instance())))
}

// Merge merges the two tag lists with the given mode.
func (t *TagList) Merge(tagList *TagList, mergeMode TagMergeMode) *TagList {
	return FromGstTagListUnsafeFull(unsafe.Pointer(C.gst_tag_list_merge(
		t.Instance(),
		tagList.Instance(),
		C.GstTagMergeMode(mergeMode),
	)))
}

// NumTags returns the number of key/value pairs in ths TagList
func (t *TagList) NumTags() int { return int(C.gst_tag_list_n_tags(t.Instance())) }

// TagNameAt returns the tag name at the given index.
func (t *TagList) TagNameAt(idx int) string {
	return C.GoString(C.gst_tag_list_nth_tag_name(t.Instance(), C.guint(idx)))
}

// PeekStringIndex peeks at the value that is at the given index for the given tag in the given list.
func (t *TagList) PeekStringIndex(tag Tag, idx uint) (value string, ok bool) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	var gout *C.gchar
	gok := C.gst_tag_list_peek_string_index(
		t.Instance(),
		(*C.gchar)(unsafe.Pointer(ctag)),
		C.guint(idx),
		&gout,
	)
	defer C.g_free((C.gpointer)(unsafe.Pointer(gout)))
	return C.GoString(gout), gobool(gok)
}

// RemoveTag removes the values for the given tag in this list.
func (t *TagList) RemoveTag(tag Tag) {
	ctag := C.CString(string(tag))
	defer C.free(unsafe.Pointer(ctag))
	C.gst_tag_list_remove_tag(t.Instance(), (*C.gchar)(unsafe.Pointer(ctag)))
}

// SetScope sets the scope of this TagList. By default, the scope of a tag list is stream scope.
func (t *TagList) SetScope(scope TagScope) {
	C.gst_tag_list_set_scope(t.Instance(), C.GstTagScope(scope))
}
