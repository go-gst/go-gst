package gst

/*
#include "gst.go.h"

extern void       goGstChildProxyChildAdded        (GstChildProxy * parent, GObject * child, const gchar * name);
extern void       goGstChildProxyChildRemoved      (GstChildProxy * parent, GObject * child, const gchar * name);
extern GObject *  goGstChildProxyGetChildByIndex   (GstChildProxy * parent, guint idx);
extern GObject *  goGstChildProxyGetChildByName    (GstChildProxy * parent, const gchar * name);
extern guint      goGstChildProxyGetChildrenCount  (GstChildProxy * parent);

void  setGstChildProxyChildAdded        (gpointer iface)  { ((GstChildProxyInterface*)iface)->child_added = goGstChildProxyChildAdded; }
void  setGstChildProxyChildRemoved      (gpointer iface)  { ((GstChildProxyInterface*)iface)->child_removed = goGstChildProxyChildRemoved; }
void  setGstChildProxyGetChildByIndex   (gpointer iface)  { ((GstChildProxyInterface*)iface)->get_child_by_index = goGstChildProxyGetChildByIndex; }
void  setGstChildProxyGetChildByName    (gpointer iface)  { ((GstChildProxyInterface*)iface)->get_child_by_name = goGstChildProxyGetChildByName; }
void  setGstChildProxyGetChildrenCount  (gpointer iface)  { ((GstChildProxyInterface*)iface)->get_children_count = goGstChildProxyGetChildrenCount; }

*/
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// InterfaceChildProxy represents the GstChildProxy interface. Use this when querying bins
// for elements that implement GstChildProxy, or when signaling that a GoObjectSubclass
// provides this interface.
var InterfaceChildProxy glib.Interface = &interfaceChildProxy{}

type interfaceChildProxy struct{ glib.Interface }

func (i *interfaceChildProxy) Type() glib.Type { return glib.Type(C.GST_TYPE_CHILD_PROXY) }
func (i *interfaceChildProxy) Init(instance *glib.TypeInstance) {
	goobj := instance.GoType

	if _, ok := goobj.(interface {
		ChildAdded(self *ChildProxy, child *glib.Object, name string)
	}); ok {
		C.setGstChildProxyChildAdded((C.gpointer)(instance.GTypeInstance))
	}

	if _, ok := goobj.(interface {
		ChildRemoved(self *ChildProxy, child *glib.Object, name string)
	}); ok {
		C.setGstChildProxyChildRemoved((C.gpointer)(instance.GTypeInstance))
	}

	if _, ok := goobj.(interface {
		GetChildByIndex(self *ChildProxy, idx uint) *glib.Object
	}); ok {
		C.setGstChildProxyGetChildByIndex((C.gpointer)(instance.GTypeInstance))
	}

	if _, ok := goobj.(interface {
		GetChildByName(self *ChildProxy, name string) *glib.Object
	}); ok {
		C.setGstChildProxyGetChildByName((C.gpointer)(instance.GTypeInstance))
	}

	if _, ok := goobj.(interface {
		GetChildrenCount(self *ChildProxy) uint
	}); ok {
		C.setGstChildProxyGetChildrenCount((C.gpointer)(instance.GTypeInstance))
	}
}

// ChildProxyImpl is the reference implementation for a ChildProxy implemented by a Go object.
type ChildProxyImpl interface {
	ChildAdded(self *ChildProxy, child *glib.Object, name string)
	ChildRemoved(self *ChildProxy, child *glib.Object, name string)
	GetChildByIndex(self *ChildProxy, idx uint) *glib.Object
	GetChildByName(self *ChildProxy, name string) *glib.Object
	GetChildrenCount(self *ChildProxy) uint
}

// ChildProxy is an interface that abstracts handling of property sets for
// elements with children. They all have multiple GstPad or some kind of voice
// objects. Another use case are container elements like GstBin. The element
// implementing the interface acts as a parent for those child objects.
//
// Property names are written as "child-name::property-name". The whole naming
// scheme is recursive. Thus "child1::child2::property" is valid too, if "child1"
// and "child2" implement the GstChildProxy interface.
type ChildProxy struct{ ptr *C.GstChildProxy }

// ToChildProxy returns a ChildProxy for the given element. If the element does not implement
// a ChildProxy it returns nil.
func ToChildProxy(elem *Element) *ChildProxy {
	if proxy := C.toGstChildProxy(elem.Unsafe()); proxy != nil {
		return &ChildProxy{proxy}
	}
	return nil
}

// Instance returns the underlying GstChildProxy instance.
func (c *ChildProxy) Instance() *C.GstChildProxy {
	return C.toGstChildProxy(unsafe.Pointer(c.ptr))
}

// ChildAdded emits the "child-added" signal.
func (c *ChildProxy) ChildAdded(child *glib.Object, name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.gst_child_proxy_child_added(
		c.Instance(),
		(*C.GObject)(child.Unsafe()),
		(*C.gchar)(unsafe.Pointer(cname)),
	)
}

// ChildRemoved emits the "child-removed" signal.
func (c *ChildProxy) ChildRemoved(child *glib.Object, name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.gst_child_proxy_child_removed(
		c.Instance(),
		(*C.GObject)(child.Unsafe()),
		(*C.gchar)(unsafe.Pointer(cname)),
	)
}

// Get gets properties of the parent object and its children. This is a direct alias to looping
// over GetProperty and returning the results in the order of the arguments. If any of the results
// returns nil from an allocation error, nil is returned for the entire slice.
func (c *ChildProxy) Get(names ...string) []*glib.Value {
	out := make([]*glib.Value, len(names))
	for i, name := range names {
		val := c.GetProperty(name)
		if val == nil {
			return nil
		}
		out[i] = val
	}
	return out
}

// GetChildByIndex fetches a child by its number. This function can return nil if the object is not
// found. Unref after usage.
func (c *ChildProxy) GetChildByIndex(idx uint) *glib.Object {
	gobj := C.gst_child_proxy_get_child_by_index(c.Instance(), C.guint(idx))
	if gobj == nil {
		return nil
	}
	return glib.TransferFull(unsafe.Pointer(gobj))
}

// GetChildByName fetches a child by name. The virtual method's default implementation uses Object
// together with Object.GetName. If the interface is to be used with GObjects, this method needs
// to be overridden.
//
// This function can return nil if the object is not found. Unref after usage.
func (c *ChildProxy) GetChildByName(name string) *glib.Object {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	gobj := C.gst_child_proxy_get_child_by_name(c.Instance(), (*C.gchar)(unsafe.Pointer(cname)))
	if gobj == nil {
		return nil
	}
	return glib.TransferFull(unsafe.Pointer(gobj))
}

// GetChildrenCount returns the number of child objects the parent contains.
func (c *ChildProxy) GetChildrenCount() uint {
	return uint(C.gst_child_proxy_get_children_count(c.Instance()))
}

// GetProperty gets a single property using the ChildProxy mechanism. The bindings
// take care of freeing the value when it leaves the user's scope. This function
// can return nil if a failure happens trying to allocate GValues.
func (c *ChildProxy) GetProperty(name string) *glib.Value {
	value, err := glib.ValueAlloc()
	if err != nil {
		return nil
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.gst_child_proxy_get_property(c.Instance(), (*C.gchar)(unsafe.Pointer(cname)), (*C.GValue)(unsafe.Pointer(value.GValue)))
	return value
}

// Lookup looks up which object and and parameter would be affected by the given name.
// If ok is false, the targets could not be found and this function returned nil.
// Unref target after usage.
func (c *ChildProxy) Lookup(name string) (ok bool, target *glib.Object, param *glib.ParamSpec) {
	var gtarget *C.GObject
	var gspec *C.GParamSpec
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ok = gobool(C.gst_child_proxy_lookup(
		c.Instance(),
		(*C.gchar)(unsafe.Pointer(cname)),
		&gtarget, &gspec,
	))
	if !ok {
		return
	}
	target = glib.TransferFull(unsafe.Pointer(gtarget))
	param = glib.ToParamSpec(unsafe.Pointer(gspec))
	return
}

// Set takes a map of names to values and applies them using the ChildProxy mechanism.
func (c *ChildProxy) Set(values map[string]*glib.Value) {
	for name, value := range values {
		c.SetProperty(name, value)
	}
}

// SetProperty sets a single property using the ChildProxy mechanism.
func (c *ChildProxy) SetProperty(name string, value *glib.Value) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.gst_child_proxy_set_property(
		c.Instance(),
		(*C.gchar)(unsafe.Pointer(cname)),
		(*C.GValue)(unsafe.Pointer(value.GValue)),
	)
}
