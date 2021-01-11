package gst

/*
#include "gst.go.h"

extern void   goClassInit     (gpointer g_class, gpointer class_data);
extern void   goInstanceInit  (GTypeInstance * instance, gpointer g_class);

void  cgoClassInit     (gpointer g_class, gpointer class_data)       { goClassInit(g_class, class_data); }
void  cgoInstanceInit  (GTypeInstance * instance, gpointer g_class)  { goInstanceInit(instance, g_class); }
*/
import "C"

import (
	"reflect"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// GoElement is an interface to be implemented by GStreamer elements built using the
// go bindings. Select methods from other interfaces can be overridden and declared via
// the Extendable property provided at plugin registration.
//
// Typically, at the very least, an element will want to implement methods from the Element
// Extendable (and by extension the GoObject).
type GoElement interface{ GoObjectSubclass }

// Extendable is an interface implemented by extendable classes. It provides
// the methods necessary to setup the vmethods on the object it represents.
type Extendable interface {
	// Type should return the type of the extended object
	Type() glib.Type
	// ClasSize should return the size of the extended class
	ClassSize() int64
	// InstanceSize should return the size of the object itself
	InstanceSize() int64
	// InitClass should take a pointer to a new subclass and a GoElement and override any
	// methods implemented by the GoElement in the subclass.
	InitClass(unsafe.Pointer, GoElement)
}

// GoObjectSubclass is an interface that abstracts on the GObjectClass. It should be implemented
// by plugins using the go bindings.
type GoObjectSubclass interface {
	// New should return a new instantiated GoElement ready to be used.
	New() GoElement
	// TypeInit is called after the GType is registered and right before ClassInit. It is when the
	// element should add any interfaces it plans to implement.
	TypeInit(*TypeInstance)
	// ClassInit is called on the element after registering it with GStreamer. This is when the element
	// should install its properties and pad templates.
	ClassInit(*ElementClass)
}

// FromObjectUnsafePrivate will return the GoElement addressed in the private data of the given GObject.
func FromObjectUnsafePrivate(obj unsafe.Pointer) GoElement {
	ptr := gopointer.Restore(privateFromObj(obj))
	return ptr.(GoElement)
}

type classData struct {
	elem GoElement
	ext  Extendable
}

func gtypeForGoElement(name string, elem GoElement, extendable Extendable) C.GType {
	registerMutex.Lock()
	defer registerMutex.Unlock()
	if registered, ok := registeredTypes[reflect.TypeOf(elem).String()]; ok {
		return registered
	}
	classData := &classData{
		elem: elem,
		ext:  extendable,
	}
	ptr := gopointer.Save(classData)

	typeInfo := (*C.GTypeInfo)(C.malloc(C.sizeof_GTypeInfo))
	defer C.free(unsafe.Pointer(typeInfo))

	typeInfo.base_init = nil
	typeInfo.base_finalize = nil
	typeInfo.class_size = C.gushort(extendable.ClassSize())
	typeInfo.class_finalize = nil
	typeInfo.class_init = C.GClassInitFunc(C.cgoClassInit)
	typeInfo.class_data = (C.gconstpointer)(ptr)
	typeInfo.instance_size = C.gushort(extendable.InstanceSize())
	typeInfo.n_preallocs = 0
	typeInfo.instance_init = C.GInstanceInitFunc(C.cgoInstanceInit)
	typeInfo.value_table = nil

	gtype := C.g_type_register_static(
		C.GType(extendable.Type()),
		(*C.gchar)(C.CString(name)),
		typeInfo,
		C.GTypeFlags(0),
	)
	elem.TypeInit(&TypeInstance{gtype: gtype, gotype: elem})
	registeredTypes[reflect.TypeOf(elem).String()] = gtype
	return gtype
}

// privateFromObj returns the actual value of the address we stored in the object's private data.
func privateFromObj(obj unsafe.Pointer) unsafe.Pointer {
	private := C.g_type_instance_get_private((*C.GTypeInstance)(obj), C.objectGType((*C.GObject)(obj)))
	privAddr := (*unsafe.Pointer)(unsafe.Pointer(private))
	return *privAddr
}
