package gst

// #include "gst.go.h"
import "C"

import (
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

func toGObject(data unsafe.Pointer) *glib.Object {
	return &glib.Object{GObject: glib.ToGObject(data)}
}

// gobool provides an easy type conversion between C.gboolean and a go bool.
func gobool(b C.gboolean) bool {
	return int(b) > 0
}

// gboolean converts a go bool to a C.gboolean.
func gboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

// gdateToTime converts a GDate to a time object.
func gdateToTime(gdate *C.GDate) time.Time {
	tm := time.Time{}
	tm.AddDate(int(C.g_date_get_year(gdate)), int(C.g_date_get_month(gdate)), int(C.g_date_get_day(gdate)))
	return tm
}

// gstDateTimeToTime converts a GstDateTime to a time object. If the datetime object could not be parsed,
// an empty time object is returned.
func gstDateTimeToTime(gstdatetime *C.GstDateTime) time.Time {
	dateStr := fmt.Sprintf(
		"%s %s %d:%d:%d %s %d",
		time.Weekday(C.gst_date_time_get_day(gstdatetime)).String(),
		time.Month(C.gst_date_time_get_month(gstdatetime)).String(),
		int(C.gst_date_time_get_hour(gstdatetime)),
		int(C.gst_date_time_get_minute(gstdatetime)),
		int(C.gst_date_time_get_second(gstdatetime)),
		formatOffset(C.gst_date_time_get_time_zone_offset(gstdatetime)),
		int(C.gst_date_time_get_year(gstdatetime)),
	)
	tm, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", dateStr)
	return tm
}

func formatOffset(offset C.gfloat) string {
	if offset < 0 {
		return fmt.Sprintf("-0%d00", int(offset))
	}
	return fmt.Sprintf("+0%d00", int(offset))
}

// goStrings returns a string slice for an array of size argc starting at the address argv.
func goStrings(argc C.int, argv **C.gchar) []string {
	length := int(argc)
	tmpslice := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.gchar)(nil))]*C.gchar)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

func gcharStrings(strs []string) **C.gchar {
	gcharSlc := make([]*C.gchar, len(strs))
	for _, s := range strs {
		cStr := C.CString(s)
		defer C.free(unsafe.Pointer(cStr))
		gcharSlc = append(gcharSlc, cStr)
	}
	return &gcharSlc[0]
}

// newQuarkFromString creates a new GQuark (or returns an existing one) for the given
// string
func newQuarkFromString(str string) C.uint {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	quark := C.g_quark_from_string(cstr)
	return quark
}

func quarkToString(q C.GQuark) string {
	return C.GoString(C.g_quark_to_string(q))
}

func streamSliceToGlist(streams []*Stream) *C.GList {
	var glist C.GList
	wrapped := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	for _, stream := range streams {
		wrapped = wrapped.Append(uintptr(stream.Unsafe()))
	}
	return &glist
}

func glistToStreamSlice(glist *C.GList) []*Stream {
	l := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	out := make([]*Stream, 0)
	l.FreeFull(func(item interface{}) {
		st := item.(*C.GstStream)
		out = append(out, wrapStream(toGObject(unsafe.Pointer(st))))
	})
	return out
}

func glistToPadTemplateSlice(glist *C.GList) []*PadTemplate {
	l := glib.WrapList(uintptr(unsafe.Pointer(&glist)))
	out := make([]*PadTemplate, 0)
	l.FreeFull(func(item interface{}) {
		tmpl := item.(*C.GstPadTemplate)
		out = append(out, FromGstPadTemplateUnsafeNone(unsafe.Pointer(tmpl)))
	})
	return out
}
