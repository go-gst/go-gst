package gst

// #include "gst.go.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Registry is a go representation of a GstRegistry.
type Registry struct{ *Object }

// FromGstRegistryUnsafeNone wraps the given GstRegistry pointer.
func FromGstRegistryUnsafeNone(registry unsafe.Pointer) *Registry {
	return &Registry{wrapObject(glib.TransferNone(registry))}
}

// FromGstRegistryUnsafeFull wraps the given GstRegistry pointer.
func FromGstRegistryUnsafeFull(registry unsafe.Pointer) *Registry {
	return &Registry{wrapObject(glib.TransferFull(registry))}
}

// GetRegistry returns the default global GstRegistry.
func GetRegistry() *Registry {
	registry := C.gst_registry_get()
	return FromGstRegistryUnsafeNone(unsafe.Pointer(registry))
}

// Instance returns the underlying GstRegistry instance.
func (r *Registry) Instance() *C.GstRegistry { return C.toGstRegistry(r.Unsafe()) }

// FindPlugin retrieves the plugin by the given name. Unref after usage.
func (r *Registry) FindPlugin(name string) (*Plugin, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	plugin := C.gst_registry_find_plugin((*C.GstRegistry)(r.Instance()), (*C.gchar)(cName))
	if plugin == nil {
		return nil, fmt.Errorf("No plugin named %s found", name)
	}
	return FromGstPluginUnsafeFull(unsafe.Pointer(plugin)), nil
}

// LookupFeature looks up the given plugin feature by name. Unref after usage.
func (r *Registry) LookupFeature(name string) (*PluginFeature, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	feat := C.gst_registry_lookup_feature((*C.GstRegistry)(r.Instance()), (*C.gchar)(cName))
	if feat == nil {
		return nil, fmt.Errorf("No feature named %s found", name)
	}
	return wrapPluginFeature(glib.TransferFull(unsafe.Pointer(feat))), nil
}
