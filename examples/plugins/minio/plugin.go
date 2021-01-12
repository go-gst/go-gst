package main

import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
)

// The metadata for this plugin
var pluginMeta = &gst.PluginMetadata{
	MajorVersion: gst.VersionMajor,
	MinorVersion: gst.VersionMinor,
	Name:         "miniosrc",
	Description:  "GStreamer plugins for reading and writing from Minio",
	Version:      "v0.0.1",
	License:      gst.LicenseLGPL,
	Source:       "go-gst",
	Package:      "examples",
	Origin:       "https://github.com/tinyzimmer/go-gst",
	ReleaseDate:  "2021-01-11",
	// The init function is called to register elements provided by the plugin.
	Init: func(plugin *gst.Plugin) bool {
		if ok := gst.RegisterElement(
			plugin,
			// The name of the element
			"miniosrc",
			// The rank of the element
			gst.RankNone,
			// The GoElement implementation for the element
			&minioSrc{},
			// The base subclass this element extends
			base.ExtendsBaseSrc,
		); !ok {
			return ok
		}

		if ok := gst.RegisterElement(
			plugin,
			// The name of the element
			"miniosink",
			// The rank of the element
			gst.RankNone,
			// The GoElement implementation for the element
			&minioSink{},
			// The base subclass this element extends
			base.ExtendsBaseSink,
		); !ok {
			return ok
		}

		return true
	},
}

func main() {}

//export gst_plugin_minio_get_desc
func gst_plugin_minio_get_desc() unsafe.Pointer { return pluginMeta.Export() }
