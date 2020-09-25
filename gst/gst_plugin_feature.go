package gst

/*
#cgo pkg-config: gstreamer-1.0
#cgo CFLAGS: -Wno-deprecated-declarations -g -Wall
#include <gst/gst.h>
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// PluginFeature wraps the C GstPluginFeature.
type PluginFeature struct{ *Object }

// Instance returns the underlying GstPluginFeature instance
func (p *PluginFeature) Instance() *C.GstPluginFeature { return C.toGstPluginFeature(p.unsafe()) }

// GetPlugin returns the plugin that provides this feature or  nil. Unref after usage.
func (p *PluginFeature) GetPlugin() *Plugin {
	plugin := C.gst_plugin_feature_get_plugin((*C.GstPluginFeature)(p.Instance()))
	if plugin == nil {
		return nil
	}
	return wrapPlugin(glib.Take(unsafe.Pointer(plugin)))
}

// GetPluginName returns the name of the plugin that provides this feature.
func (p *PluginFeature) GetPluginName() string {
	pluginName := C.gst_plugin_feature_get_plugin_name((*C.GstPluginFeature)(p.Instance()))
	if pluginName == nil {
		return ""
	}
	return C.GoString(pluginName)
}
