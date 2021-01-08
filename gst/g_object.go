package gst

/*
#include "gst.go.h"

extern void   goObjectSetProperty  (GObject * object, guint property_id, const GValue * value, GParamSpec *pspec);
extern void   goObjectGetProperty  (GObject * object, guint property_id, GValue * value, GParamSpec * pspec);
extern void   goObjectConstructed  (GObject * object);
extern void   goObjectFinalize     (GObject * object, gpointer klass);

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

void  setGObjectClassSetProperty  (void * klass)  { ((GObjectClass *)klass)->set_property = goObjectSetProperty; }
void  setGObjectClassGetProperty  (void * klass)  { ((GObjectClass *)klass)->get_property = goObjectGetProperty; }
void  setGObjectClassConstructed  (void * klass)  { ((GObjectClass *)klass)->constructed = objectConstructed; }
void  setGObjectClassFinalize     (void * klass)  { ((GObjectClass *)klass)->finalize = objectFinalize; }

*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// GoObject is an interface that abstracts on the GObject. In almost all cases at least SetProperty and GetProperty
// should be implemented by elements built from the go bindings.
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

// ExtendsObject signifies a GoElement that extends a GObject. It is the base Extendable
// that all other implementations derive from.
var ExtendsObject Extendable = &extendObject{}

type extendObject struct{}

func (e *extendObject) Type() glib.Type     { return glib.Type(C.g_object_get_type()) }
func (e *extendObject) ClassSize() int64    { return int64(C.sizeof_GObjectClass) }
func (e *extendObject) InstanceSize() int64 { return int64(C.sizeof_GObject) }

func (e *extendObject) InitClass(klass unsafe.Pointer, elem GoElement) {
	C.setGObjectClassFinalize(klass)

	if _, ok := elem.(interface {
		SetProperty(obj *Object, id uint, value *glib.Value)
	}); ok {
		C.setGObjectClassSetProperty(klass)
	}
	if _, ok := elem.(interface {
		GetProperty(obj *Object, id uint) *glib.Value
	}); ok {
		C.setGObjectClassGetProperty(klass)
	}
	if _, ok := elem.(interface {
		Constructed(*Object)
	}); ok {
		C.setGObjectClassConstructed(klass)
	}
}
