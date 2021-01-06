// The contents of this file could be generated from markers placed in filesrc.go
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
	Name:         "go-file-plugins",
	Description:  "File plugins written in Go",
	Version:      "v0.0.1",
	License:      gst.LicenseLGPL,
	Source:       "go-gst",
	Package:      "examples",
	Origin:       "https://github.com/tinyzimmer/go-gst",
	ReleaseDate:  "2021-01-04",
	// The init function is called to register elements provided
	// by the plugin.
	Init: func(plugin *gst.Plugin) bool {
		return gst.RegisterElement(
			plugin,
			"gofilesrc",         // The name of the element
			gst.RankNone,        // The rank of the element
			&fileSrc{},          // The GoElement implementation for the element
			base.ExtendsBaseSrc, // The base subclass this element extends
		)
	},
}

// A single method must be exported from the compiled library that provides for GStreamer
// to fetch the description and init function for this plugin. The name of the method
// must match the format gst_plugin_NAME_get_desc, where hyphens are replaced with underscores.

//export gst_plugin_gofilesrc_get_desc
func gst_plugin_gofilesrc_get_desc() unsafe.Pointer { return pluginMeta.Export() }

// main is left unimplemented since these files are compiled to c-shared.
func main() {}
