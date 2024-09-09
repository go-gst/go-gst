//lint:file-ignore U1000 Ignore all unused code, this is example code

// +plugin:Name=boilerplate
// +plugin:Description=My plugin written in go
// +plugin:Version=v0.0.1
// +plugin:License=gst.LicenseLGPL
// +plugin:Source=go-gst
// +plugin:Package=examples
// +plugin:Origin=https://github.com/go-gst/go-gst
// +plugin:ReleaseDate=2021-01-18
//
// +element:Name=myelement
// +element:Rank=gst.RankNone
// +element:Impl=myelement
// +element:Subclass=gst.ExtendsElement
//
//go:generate gst-plugin-gen
package main

import (
	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

func main() {}

type myelement struct{}

func (g *myelement) New() glib.GoObjectSubclass { return &myelement{} }

func (g *myelement) ClassInit(klass *glib.ObjectClass) {
	// Set the plugin's longname as it is a basic requirement for a GStreamer plugin
	class := gst.ToElementClass(klass)
	class.SetMetadata(
		"Boilerplate",
		"General",
		"An empty element which does nothing",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
}
