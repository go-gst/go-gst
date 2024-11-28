package gst

// #include "gst.go.h"
import "C"
import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

type Tracer struct {
	*Object
}

func TracingGetActiveTracers() []*Tracer {
	cglist := C.gst_tracing_get_active_tracers()

	wrapped := glib.WrapList(unsafe.Pointer(cglist))
	defer wrapped.Free()

	out := make([]*Tracer, 0)
	wrapped.Foreach(func(item interface{}) {
		ctracer := item.(unsafe.Pointer) // item is a *C.GstTracer
		out = append(out, &Tracer{
			Object: wrapObject(glib.Take(ctracer)),
		})
	})
	return out
}
