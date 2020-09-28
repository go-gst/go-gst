package gst

// #include "gst.go.h"
import "C"
import "unsafe"

// TagList is a go wrapper around a GstTagList. For now, until the rest of the methods are
// implemnented, this struct is primarily used for retrieving serialized copies of the tags.
type TagList struct {
	ptr *C.GstTagList
}

// Instance returns the underlying GstTagList instance.
func (t *TagList) Instance() *C.GstTagList { return t.ptr }

// String implements a stringer on the TagList and serializes it to a string.
func (t *TagList) String() string { return C.GoString(C.gst_tag_list_to_string(t.Instance())) }

// Ref increases the ref count on this TagList by one.
func (t *TagList) Ref() *TagList { return wrapTagList(C.gst_tag_list_ref(t.Instance())) }

// Unref decreses the ref count on this TagList by one. When the ref count reaches zero, the object
// is destroyed.
func (t *TagList) Unref() { C.gst_tag_list_unref(t.Instance()) }

// Size returns the number of key/value pairs in ths TagList
func (t *TagList) Size() int { return int(C.gst_tag_list_n_tags(t.Instance())) }

// TagNameAt returns the tag name at the given index.
func (t *TagList) TagNameAt(idx int) string {
	return C.GoString(C.gst_tag_list_nth_tag_name(t.Instance(), C.guint(idx)))
}

// NumValuesAt returns the number of tag values at the given tag key.
func (t *TagList) NumValuesAt(tagKey string) int {
	cStr := C.CString(tagKey)
	defer C.free(unsafe.Pointer(cStr))
	return int(C.gst_tag_list_get_tag_size(t.Instance(), (*C.gchar)(cStr)))
}
