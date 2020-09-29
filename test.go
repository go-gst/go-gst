package main

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

func main() {
	gst.Init(nil)

	metaInfo := gst.RegisterMeta(glib.TypeFromName("GstObject"), "my-meta", 1024, &gst.MetaInfoCallbackFuncs{
		InitFunc: func(params interface{}, buffer *gst.Buffer) bool {
			paramStr := params.(string)
			fmt.Println("Buffer initialized with params:", paramStr)
			return true
		},
		FreeFunc: func(buffer *gst.Buffer) {
			fmt.Println("Buffer was destroyed")
		},
	})
	buf := gst.NewEmptyBuffer()
	buf.AddMeta(metaInfo, "hello world")
	buf.Unref()
}
