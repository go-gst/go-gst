package gst

import (
	"unsafe"

	"github.com/go-gst/go-glib/pkg/core/userdata"
	"github.com/go-gst/go-glib/pkg/glib/v2"
	"github.com/go-gst/go-glib/pkg/gobject/v2"
)

// #include <gst/gst.h>
import "C"

//export _goglib_gst1_StructureForeachFunc
func _goglib_gst1_StructureForeachFunc(carg1 C.GQuark, carg2 *C.GValue, carg3 C.gpointer) (cret C.gboolean) {
	var fn StructureForeachFunc
	{
		v := userdata.Load(unsafe.Pointer(carg3))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(StructureForeachFunc)
	}

	var fieldId glib.Quark   // in, none, casted, alias
	var value *gobject.Value // in, none, converted
	var goret bool           // return

	fieldId = glib.Quark(carg1)
	value = gobject.ValueFromNative(unsafe.Pointer(carg2))

	goret = fn(glib.QuarkToString(fieldId), value.GoValue())

	if goret {
		cret = C.TRUE
	}

	return cret
}
