package customtransform

import (
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/base"
)

// Register needs to be called after gst.Init() to make the gocustombin available in the standard
// gst element registry. After this call the element can be used like any other gstreamer element
func Register() bool {
	return gst.RegisterElement(
		// no plugin:
		nil,
		// The name of the element
		"gocustomtransform",
		// The rank of the element
		gst.RankNone,
		// The GoElement implementation for the element
		&customBaseTransform{},
		// The base subclass this element extends
		base.ExtendsBaseTransform,
	)
}
