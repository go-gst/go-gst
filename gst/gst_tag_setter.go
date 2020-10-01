package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// TagSetter is an interface that elements can implement to provide Tag writing capabilities.
type TagSetter interface {
	// Returns the current list of tags the setter uses. The list should not be modified or freed.
	GetTagList() *TagList
	// Adds the given tag/value pair using the given merge mode.
	// TODO: Either an additional function or this one should be modified to accept go types.
	AddTagValue(mergeMode TagMergeMode, tagKey string, tagValue *glib.Value)
	// Merges a tag list with the given merge mode
	MergeTags(*TagList, TagMergeMode)
	// Resets the internal tag list. Elements should call this from within the state-change handler.
	ResetTags()
	// Queries the mode by which tags inside the setter are overwritten by tags from events
	GetTagMergeMode() TagMergeMode
	// Sets the given merge mode that is used for adding tags from events to tags specified by this interface.
	// The default is TagMergeKeep, which keeps the tags set with this interface and discards tags from events.
	SetTagMergeMode(TagMergeMode)
}

// gstTocSetter implements a TagSetter that is backed by an Element from the C runtime.
type gstTagSetter struct {
	ptr *C.GstElement
}

// Instance returns the underlying TagSetter instance.
func (t *gstTagSetter) Instance() *C.GstTagSetter { return C.toTagSetter(t.ptr) }

func (t *gstTagSetter) GetTagList() *TagList {
	tagList := C.gst_tag_setter_get_tag_list(t.Instance())
	if tagList == nil {
		return nil
	}
	return wrapTagList(tagList)
}

func (t *gstTagSetter) AddTagValue(mergeMode TagMergeMode, tagKey string, tagValue *glib.Value) {
	ckey := C.CString(tagKey)
	defer C.free(unsafe.Pointer(ckey))
	C.gst_tag_setter_add_tag_value(
		t.Instance(),
		C.GstTagMergeMode(mergeMode),
		(*C.gchar)(unsafe.Pointer(ckey)),
		(*C.GValue)(tagValue.GetPointer()),
	)
}

func (t *gstTagSetter) MergeTags(tagList *TagList, mergeMode TagMergeMode) {
	C.gst_tag_setter_merge_tags(t.Instance(), tagList.Instance(), C.GstTagMergeMode(mergeMode))
}

func (t *gstTagSetter) ResetTags() {
	C.gst_tag_setter_reset_tags(t.Instance())
}

func (t *gstTagSetter) GetTagMergeMode() TagMergeMode {
	return TagMergeMode(C.gst_tag_setter_get_tag_merge_mode(t.Instance()))
}

func (t *gstTagSetter) SetTagMergeMode(mergeMode TagMergeMode) {
	C.gst_tag_setter_set_tag_merge_mode(t.Instance(), C.GstTagMergeMode(mergeMode))
}
