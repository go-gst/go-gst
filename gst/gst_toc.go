package gst

// #include "gst.go.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// TOC is a go representation of a GstToc.
type TOC struct {
	ptr *C.GstToc
}

// FromGstTOCUnsafeNone wraps the pointer to the given C GstToc with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstTOCUnsafeNone(toc unsafe.Pointer) *TOC {
	gotoc := wrapTOC((*C.GstToc)(toc))
	gotoc.Ref()
	runtime.SetFinalizer(gotoc, (*TOC).Unref)
	return gotoc
}

// FromGstTOCUnsafeFull wraps the pointer to the given C GstToc with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstTOCUnsafeFull(toc unsafe.Pointer) *TOC {
	gotoc := wrapTOC((*C.GstToc)(toc))
	runtime.SetFinalizer(gotoc, (*TOC).Unref)
	return gotoc
}

// NewTOC returns a new TOC with the given scope.
func NewTOC(scope TOCScope) *TOC {
	toc := C.gst_toc_new(C.GstTocScope(scope))
	if toc == nil {
		return nil
	}
	return FromGstTOCUnsafeFull(unsafe.Pointer(toc))
}

// Instance returns the underlying GstToc instance.
func (t *TOC) Instance() *C.GstToc { return t.ptr }

// Ref increases the ref count on the TOC by one.
func (t *TOC) Ref() *TOC {
	C.tocRef(t.Instance())
	return t
}

// Unref decreases the ref count on the TOC by one.
func (t *TOC) Unref() {
	C.tocUnref(t.Instance())
}

// MakeWritable returns a writable copy of the TOC if it isn't already,
func (t *TOC) MakeWritable() *TOC {
	return FromGstTOCUnsafeFull(unsafe.Pointer(C.makeTocWritable(t.Instance())))
}

// Copy creates a copy of the TOC.
func (t *TOC) Copy() *TOC {
	return FromGstTOCUnsafeFull(unsafe.Pointer(C.copyToc(t.Instance())))
}

// AppendEntry appends the given TOCEntry to this TOC.
func (t *TOC) AppendEntry(entry *TOCEntry) {
	C.gst_toc_append_entry(t.Instance(), entry.Instance())
}

// Dump dumps the TOC.
func (t *TOC) Dump() {
	C.gst_toc_dump(t.Instance())
}

// FindEntry finds the entry with the given uid.
func (t *TOC) FindEntry(uid string) *TOCEntry {
	cuid := C.CString(uid)
	defer C.free(unsafe.Pointer(cuid))
	entry := C.gst_toc_find_entry(t.Instance(), (*C.gchar)(cuid))
	if entry == nil {
		return nil
	}
	return FromGstTocEntryUnsafeNone(unsafe.Pointer(entry))
}

// GetEntries returns a list of all TOCEntries.
func (t *TOC) GetEntries() []*TOCEntry {
	gList := C.gst_toc_get_entries(t.Instance())

	defer C.g_list_free(gList)
	out := make([]*TOCEntry, 0)

	for {
		entry := C.glistNext(gList)
		if entry == nil {
			break
		}
		out = append(out, FromGstTocEntryUnsafeNone(unsafe.Pointer(entry)))
	}
	return out
}

// GetScope returns the scope of this TOC.
func (t *TOC) GetScope() TOCScope {
	return TOCScope(C.gst_toc_get_scope(t.Instance()))
}

// GetTags returns the TagList for this TOC.
func (t *TOC) GetTags() *TagList {
	tagList := C.gst_toc_get_tags(t.Instance())
	if tagList == nil {
		return nil
	}
	return FromGstTagListUnsafeNone(unsafe.Pointer(tagList))
}

// MergeTags merges the given tags into this TOC's TagList.
func (t *TOC) MergeTags(tagList *TagList, mergeMode TagMergeMode) {
	C.gst_toc_merge_tags(t.Instance(), tagList.Instance(), C.GstTagMergeMode(mergeMode))
}

// SetTags sets tags for the entire TOC.
func (t *TOC) SetTags(tagList *TagList) {
	C.gst_toc_set_tags(t.Instance(), tagList.Ref().Instance())
}

// TOCEntry is a go representation of a GstTocEntry,
type TOCEntry struct {
	ptr *C.GstTocEntry
}

// FromGstTocEntryUnsafeNone wraps the given TOCEntry.
func FromGstTocEntryUnsafeNone(entry unsafe.Pointer) *TOCEntry {
	t := wrapTOCEntry((*C.GstTocEntry)(entry))
	t.Ref()
	runtime.SetFinalizer(t, (*TOCEntry).Unref)
	return t
}

// FromGstTocEntryUnsafeFull wraps the given TOCEntry.
func FromGstTocEntryUnsafeFull(entry unsafe.Pointer) *TOCEntry {
	t := wrapTOCEntry((*C.GstTocEntry)(entry))
	runtime.SetFinalizer(t, (*TOCEntry).Unref)
	return t
}

// NewTOCEntry creates a new TOCEntry with the given UID and type.
func NewTOCEntry(entryType TOCEntryType, uid string) *TOCEntry {
	cuid := C.CString(uid)
	defer C.free(unsafe.Pointer(cuid))
	entry := C.gst_toc_entry_new(
		C.GstTocEntryType(entryType),
		(*C.gchar)(unsafe.Pointer(cuid)),
	)
	if entry == nil {
		return nil
	}
	return FromGstTocEntryUnsafeFull(unsafe.Pointer(entry))
}

// Instance returns the underlying GstTocEntry instance.
func (t *TOCEntry) Instance() *C.GstTocEntry { return t.ptr }

// Ref increases the ref count on the TOCEntry by one.
func (t *TOCEntry) Ref() *TOCEntry {
	C.tocEntryRef(t.Instance())
	return t
}

// Unref decreases the ref count on the TOCEntry by one.
func (t *TOCEntry) Unref() {
	C.tocEntryUnref(t.Instance())
}

// MakeWritable returns a writable copy of the TOCEntry if it is not already so.
func (t *TOCEntry) MakeWritable() *TOCEntry {
	return FromGstTocEntryUnsafeFull(unsafe.Pointer(C.makeTocEntryWritable(t.Instance())))
}

// Copy creates a copy of the TOCEntry
func (t *TOCEntry) Copy() *TOCEntry {
	return FromGstTocEntryUnsafeFull(unsafe.Pointer(C.copyTocEntry(t.Instance())))
}

// AppendSubEntry appends the given entry as a subentry to this one.
func (t *TOCEntry) AppendSubEntry(subEntry *TOCEntry) {
	C.gst_toc_entry_append_sub_entry(t.Instance(), subEntry.Ref().Instance())
}

// GetEntryType returns the type of this TOCEntry
func (t *TOCEntry) GetEntryType() TOCEntryType {
	return TOCEntryType(C.gst_toc_entry_get_entry_type(t.Instance()))
}

// GetEntryTypeString returns a string representation of the entry type.
func (t *TOCEntry) GetEntryTypeString() string {
	return C.GoString(C.gst_toc_entry_type_get_nick(C.GstTocEntryType(t.GetEntryType())))
}

// GetLoop gets the loop type and repeat count for the TOC entry.
func (t *TOCEntry) GetLoop() (bool, TOCLoopType, int) {
	var loopType C.GstTocLoopType
	var repeatCount C.gint
	ok := C.gst_toc_entry_get_loop(t.Instance(), &loopType, &repeatCount)
	return gobool(ok), TOCLoopType(loopType), int(repeatCount)
}

// GetParent gets the parent of this TOCEntry.
func (t *TOCEntry) GetParent() *TOCEntry {
	parent := C.gst_toc_entry_get_parent(t.Instance())
	if parent == nil {
		return nil
	}
	return FromGstTocEntryUnsafeNone(unsafe.Pointer(parent))
}

// GetStartStopTimes gets the start and stop times for the TOCEntry if available.
func (t *TOCEntry) GetStartStopTimes() (ok bool, startTime, stopTime int64) {
	var start, stop C.gint64
	gok := C.gst_toc_entry_get_start_stop_times(t.Instance(), &start, &stop)
	return gobool(gok), int64(start), int64(stop)
}

// GetSubEntries gets all the subentries for this TOCEntry.
func (t *TOCEntry) GetSubEntries() []*TOCEntry {
	gList := C.gst_toc_entry_get_sub_entries(t.Instance())

	defer C.g_list_free(gList)
	out := make([]*TOCEntry, 0)

	for {
		entry := C.glistNext(gList)
		if entry == nil {
			break
		}
		out = append(out, FromGstTocEntryUnsafeNone(unsafe.Pointer(entry)))
	}
	return out
}

// GetTags gets the tags for this entry.
func (t *TOCEntry) GetTags() *TagList {
	tagList := C.gst_toc_entry_get_tags(t.Instance())
	if tagList == nil {
		return nil
	}
	return FromGstTagListUnsafeNone(unsafe.Pointer(tagList))
}

// GetTOC returns the parent TOC of this entry.
func (t *TOCEntry) GetTOC() *TOC {
	toc := C.gst_toc_entry_get_toc(t.Instance())
	if toc == nil {
		return nil
	}
	return FromGstTOCUnsafeNone(unsafe.Pointer(toc))
}

// GetUID returns the uid of this entry.
func (t *TOCEntry) GetUID() string {
	return C.GoString(C.gst_toc_entry_get_uid(t.Instance()))
}

// IsAlternative returns true if this is an alternative entry.
func (t *TOCEntry) IsAlternative() bool {
	return gobool(C.gst_toc_entry_is_alternative(t.Instance()))
}

// IsSequence returns true if this is a sequence entry.
func (t *TOCEntry) IsSequence() bool {
	return gobool(C.gst_toc_entry_is_sequence(t.Instance()))
}

// MergeTags merges the given tags with the given mode.
func (t *TOCEntry) MergeTags(tagList *TagList, mergeMode TagMergeMode) {
	if tagList == nil {
		C.gst_toc_entry_merge_tags(t.Instance(), nil, C.GstTagMergeMode(mergeMode))
		return
	}
	C.gst_toc_entry_merge_tags(t.Instance(), tagList.Instance(), C.GstTagMergeMode(mergeMode))
}

// SetLoop sets the loop type and repeat counts for the entry.
func (t *TOCEntry) SetLoop(loopType TOCLoopType, repeatCount int) {
	C.gst_toc_entry_set_loop(t.Instance(), C.GstTocLoopType(loopType), C.gint(repeatCount))
}

// SetStartStopTimes sets the start and stop times for the TOC entry.
func (t *TOCEntry) SetStartStopTimes(startTime, stopTime int64) {
	C.gst_toc_entry_set_start_stop_times(t.Instance(), C.gint64(startTime), C.gint64(stopTime))
}

// SetTags sets the tags on the TOC entry.
func (t *TOCEntry) SetTags(tagList *TagList) {
	C.gst_toc_entry_set_tags(t.Instance(), tagList.Ref().Instance())
}
