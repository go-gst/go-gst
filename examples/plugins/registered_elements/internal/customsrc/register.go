package customsrc

import (
	"github.com/go-gst/go-gst/gst"
)

// Register needs to be called after gst.Init() to make the gocustomsrc available in the standard
// gst element registry. After this call the element can be used like any other gstreamer element
func Register() bool {
	return gst.RegisterElement(
		// no plugin:
		nil,
		// The name of the element
		"gocustomsrc",
		// The rank of the element
		gst.RankNone,
		// The GoElement implementation for the element
		&customSrc{},
		// The base subclass this element extends
		gst.ExtendsBin,
	)
}
