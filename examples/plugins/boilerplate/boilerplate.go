//go:generate gst-plugin-gen
//
// +plugin:Name=boilerplate
// +plugin:Description=My plugin written in go
// +plugin:Version=v0.0.1
// +plugin:License=gst.LicenseLGPL
// +plugin:Source=go-gst
// +plugin:Package=examples
// +plugin:Origin=https://github.com/tinyzimmer/go-gst
// +plugin:ReleaseDate=2021-01-18
//
// +element:Name=myelement
// +element:Rank=gst.RankNone
// +element:Impl=myelement
// +element:Subclass=gst.ExtendsElement
//
package main

import "github.com/tinyzimmer/go-glib/glib"

func main() {}

type myelement struct{}

func (g *myelement) New() glib.GoObjectSubclass { return &myelement{} }

func (g *myelement) ClassInit(klass *glib.ObjectClass) {}
