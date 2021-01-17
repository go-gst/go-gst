// This example demonstrates a src element that reads from objects in a minio bucket.
// Since minio implements the S3 API this plugin could also be used for S3 buckets by
// setting the correct endpoints and credentials.
//
// By default this plugin will use the credentials set in the environment at MINIO_ACCESS_KEY_ID
// and MINIO_SECRET_ACCESS_KEY however these can also be set on the element directly.
//
//
// In order to build the plugin for use by GStreamer, you can do the following:
//
//     $ go build -o libgstminio.so -buildmode c-shared .
//
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
	Name:         "minio-plugins",
	Description:  "GStreamer plugins for reading and writing from Minio",
	Version:      "v0.0.1",
	License:      gst.LicenseLGPL,
	Source:       "gst-pipeline-operator",
	Package:      "plugins",
	Origin:       "https://github.com/tinyzimmer/gst-pipeline-operator",
	ReleaseDate:  "2021-01-12",
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
