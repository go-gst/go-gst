package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"
)

// PluginFeature wraps the C GstPluginFeature.
type PluginFeature struct{ *Object }

// Instance returns the underlying GstPluginFeature instance
func (p *PluginFeature) Instance() *C.GstPluginFeature { return C.toGstPluginFeature(p.Unsafe()) }

// GetPlugin returns the plugin that provides this feature or  nil. Unref after usage.
func (p *PluginFeature) GetPlugin() *Plugin {
	plugin := C.gst_plugin_feature_get_plugin((*C.GstPluginFeature)(p.Instance()))
	if plugin == nil {
		return nil
	}
	return FromGstPluginUnsafeFull(unsafe.Pointer(plugin))
}

// GetPluginName returns the name of the plugin that provides this feature.
func (p *PluginFeature) GetPluginName() string {
	pluginName := C.gst_plugin_feature_get_plugin_name((*C.GstPluginFeature)(p.Instance()))
	if pluginName == nil {
		return ""
	}
	return C.GoString(pluginName)
}

func (p *PluginFeature) SetPluginRank(rank Rank) {
	C.gst_plugin_feature_set_rank((*C.GstPluginFeature)(p.Instance()), C.guint(rank))
}

func (p *PluginFeature) GetPluginRank() int {
	rank := C.gst_plugin_feature_get_rank((*C.GstPluginFeature)(p.Instance()))
	return int(rank)
}
