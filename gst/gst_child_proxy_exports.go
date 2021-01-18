package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

func wrapParent(parent *C.GstChildProxy) *ChildProxy { return &ChildProxy{ptr: parent} }

//export goGstChildProxyChildAdded
func goGstChildProxyChildAdded(parent *C.GstChildProxy, child *C.GObject, name *C.gchar) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(parent))
	caller := iface.(interface {
		ChildAdded(self *ChildProxy, child *glib.Object, name string)
	})
	caller.ChildAdded(
		wrapParent(parent),
		&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))},
		C.GoString(name),
	)
}

//export goGstChildProxyChildRemoved
func goGstChildProxyChildRemoved(parent *C.GstChildProxy, child *C.GObject, name *C.gchar) {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(parent))
	caller := iface.(interface {
		ChildRemoved(self *ChildProxy, child *glib.Object, name string)
	})
	caller.ChildRemoved(
		wrapParent(parent),
		&glib.Object{GObject: glib.ToGObject(unsafe.Pointer(child))},
		C.GoString(name),
	)
}

//export goGstChildProxyGetChildByIndex
func goGstChildProxyGetChildByIndex(parent *C.GstChildProxy, idx C.guint) *C.GObject {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(parent))
	caller := iface.(interface {
		GetChildByIndex(self *ChildProxy, idx uint) *glib.Object
	})
	obj := caller.GetChildByIndex(wrapParent(parent), uint(idx))
	if obj == nil {
		return nil
	}
	return (*C.GObject)(unsafe.Pointer(obj.GObject))
}

//export goGstChildProxyGetChildByName
func goGstChildProxyGetChildByName(parent *C.GstChildProxy, name *C.gchar) *C.GObject {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(parent))
	caller := iface.(interface {
		GetChildByName(self *ChildProxy, name string) *glib.Object
	})
	obj := caller.GetChildByName(wrapParent(parent), C.GoString(name))
	if obj == nil {
		return nil
	}
	return (*C.GObject)(unsafe.Pointer(obj.GObject))
}

//export goGstChildProxyGetChildrenCount
func goGstChildProxyGetChildrenCount(parent *C.GstChildProxy) C.guint {
	iface := glib.FromObjectUnsafePrivate(unsafe.Pointer(parent))
	caller := iface.(interface {
		GetChildrenCount(self *ChildProxy) uint
	})
	return C.guint(caller.GetChildrenCount(wrapParent(parent)))
}
