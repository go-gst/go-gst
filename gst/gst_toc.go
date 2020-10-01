package gst

// #include "gst.go.h"
import "C"
import "unsafe"

// TOC is a go representation of a GstToc.
type TOC struct {
	ptr *C.GstToc
}

// NewTOC returns a new TOC with the given scope.
func NewTOC(scope TOCScope) *TOC {
	toc := C.gst_toc_new(C.GstTocScope(scope))
	if toc == nil {
		return nil
	}
	return wrapTOC(toc)
}

// Instance returns the underlying GstToc instance.
func (t *TOC) Instance() *C.GstToc { return t.ptr }

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
	return wrapTOCEntry(entry)
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
		out = append(out, wrapTOCEntry((*C.GstTocEntry)(unsafe.Pointer(entry))))
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
	return wrapTagList(tagList)
}

// MergeTags

// TOCEntry is a go representation of a GstTocEntry,
type TOCEntry struct {
	ptr *C.GstTocEntry
}

// Instance returns the underlying GstTocEntry instance.
func (t *TOCEntry) Instance() *C.GstTocEntry { return t.ptr }
