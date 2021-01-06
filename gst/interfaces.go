package gst

/*
#include "gst.go.h"

extern void   goClassInit     (gpointer g_class, gpointer class_data);
extern void   goInstanceInit  (GTypeInstance * instance, gpointer g_class);

extern void            goObjectSetProperty  (GObject * object, guint property_id, const GValue * value, GParamSpec *pspec);
extern void            goObjectGetProperty  (GObject * object, guint property_id, GValue * value, GParamSpec * pspec);
extern void            goObjectConstructed  (GObject * object);
extern void            goObjectFinalize     (GObject * object, gpointer klass);

void objectFinalize (GObject * object)
{
	GObjectClass *parent = g_type_class_peek_parent((G_OBJECT_GET_CLASS(object)));
	goObjectFinalize(object, G_OBJECT_GET_CLASS(object));
	parent->finalize(object);
}

void objectConstructed (GObject * object)
{
	GObjectClass *parent = g_type_class_peek_parent((G_OBJECT_GET_CLASS(object)));
	goObjectConstructed(object);
	parent->constructed(object);
}

void cgoClassInit (gpointer g_class, gpointer class_data)
{
	((GObjectClass *)g_class)->set_property = goObjectSetProperty;
	((GObjectClass *)g_class)->get_property = goObjectGetProperty;
	((GObjectClass *)g_class)->constructed = objectConstructed;
	((GObjectClass *)g_class)->finalize = objectFinalize;

	goClassInit(g_class, class_data);
}

void cgoInstanceInit (GTypeInstance * instance, gpointer g_class)
{
	goInstanceInit(instance, g_class);
}

*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	gopointer "github.com/mattn/go-pointer"
)

// Extendable is an interface implemented by extendable classes. It provides
// the methods necessary to setup the vmethods on the object it represents.
type Extendable interface {
	Type() glib.Type
	ClassSize() int64
	InstanceSize() int64
	InitClass(unsafe.Pointer, GoElement)
}

type extendElement struct{}

func (e *extendElement) Type() glib.Type                                { return glib.Type(C.gst_element_get_type()) }
func (e *extendElement) ClassSize() int64                               { return int64(C.sizeof_GstElementClass) }
func (e *extendElement) InstanceSize() int64                            { return int64(C.sizeof_GstElement) }
func (e *extendElement) InitClass(klass unsafe.Pointer, elem GoElement) {}

// ExtendsElement signifies a GoElement that extends a GstElement.
var ExtendsElement Extendable = &extendElement{}

// GoElement is an interface to be implemented by GStreamer elements built using the
// go bindings. The various methods are called throughout the lifecycle of the plugin.
type GoElement interface {
	GoObjectSubclass
	GoObject
}

// privateFromObj returns the actual value of the address we stored in the object's private data.
func privateFromObj(obj unsafe.Pointer) unsafe.Pointer {
	private := C.g_type_instance_get_private((*C.GTypeInstance)(obj), C.objectGType((*C.GObject)(obj)))
	privAddr := (*unsafe.Pointer)(unsafe.Pointer(private))
	return *privAddr
}

// FromObjectUnsafePrivate will return the GoElement addressed in the private data of the given GObject.
func FromObjectUnsafePrivate(obj unsafe.Pointer) GoElement {
	ptr := gopointer.Restore(privateFromObj(obj))
	return ptr.(GoElement)
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
	// should  install any properties and pad templates it has.
	ClassInit(*ElementClass)
}

// GoObject is an interface that abstracts on the GObject. It should be implemented by plugins using
// the gobindings.
type GoObject interface {
	// SetProperty should set the value of the property with the given id. ID is the index+1 of the parameter
	// in the order it was registered.
	SetProperty(obj *Object, id uint, value *glib.Value)
	// GetProperty should retrieve the value of the property with the given id. ID is the index+1 of the parameter
	// in the order it was registered.
	GetProperty(obj *Object, id uint) *glib.Value
	// Constructed is called when the Object has finished setting up.
	Constructed(*Object)
}

type classData struct {
	elem GoElement
	ext  Extendable
}

func gtypeForGoElement(name string, elem GoElement, extendable Extendable) C.GType {
	registerMutex.Lock()
	defer registerMutex.Unlock()
	// fmt.Printf("Checking registration of %v\n", reflect.TypeOf(elem).String())
	if registered, ok := registeredTypes[reflect.TypeOf(elem).String()]; ok {
		return registered
	}
	classData := &classData{
		elem: elem,
		ext:  extendable,
	}
	ptr := gopointer.Save(classData)
	typeInfo := C.GTypeInfo{
		class_size:     C.gushort(extendable.ClassSize()),
		base_init:      nil,
		base_finalize:  nil,
		class_init:     C.GClassInitFunc(C.cgoClassInit),
		class_finalize: nil,
		class_data:     (C.gconstpointer)(ptr),
		instance_size:  C.gushort(extendable.InstanceSize()),
		n_preallocs:    0,
		instance_init:  C.GInstanceInitFunc(C.cgoInstanceInit),
		value_table:    nil,
	}
	gtype := C.g_type_register_static(
		C.GType(extendable.Type()),
		(*C.gchar)(C.CString(name)),
		&typeInfo,
		C.GTypeFlags(0),
	)
	elem.TypeInit(&TypeInstance{gtype: gtype, gotype: elem})
	// fmt.Printf("Registering %v to type %v\n", reflect.TypeOf(elem).String(), gtype)
	registeredTypes[reflect.TypeOf(elem).String()] = gtype
	return gtype
}
