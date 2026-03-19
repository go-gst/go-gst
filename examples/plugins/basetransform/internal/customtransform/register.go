package customtransform

import (
	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstbase"
)

// Register needs to be called after gst.Init() to make the gocustombin available in the standard
// gst element registry. After this call the element can be used like any other gstreamer element
func Register() bool {
	registered := gstbase.RegisterBaseTransformSubClass[*customBaseTransform](
		"gocustomtransform",
		classInit,
		nil, // no constructor
		gstbase.BaseTransformOverrides[*customBaseTransform]{},
		nil, // no signals
	)

	return gst.ElementRegister(
		// no plugin:
		nil,
		// The name of the element
		"gocustomtransform",
		// The rank of the element
		uint(gst.RankNone),
		registered,
	)
}
