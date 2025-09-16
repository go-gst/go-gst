package custombin

import (
	"github.com/diamondburned/gotk4/pkg/gobject/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

// Register needs to be called after gst.Init() to make the gocustombin available in the standard
// gst element registry. After this call the element can be used like any other gstreamer element
func Register() bool {
	registered := gst.RegisterBinSubClass[*customBin](
		"gocustombin",
		classInit,
		nil,
		gst.BinOverrides[*customBin]{
			ElementOverrides: gst.ElementOverrides[*customBin]{
				ObjectOverrides: gst.ObjectOverrides[*customBin]{
					InitiallyUnownedOverrides: gobject.InitiallyUnownedOverrides[*customBin]{
						ObjectOverrides: gobject.ObjectOverrides[*customBin]{
							Constructed: (*customBin).constructed,
						},
					},
				},
			},
		},
		map[string]gobject.SignalDefinition{},
	)

	return gst.ElementRegister(
		// no plugin:
		nil,
		// The name of the element
		"gocustombin",
		// The rank of the element
		uint(gst.RankNone),
		registered,
	)
}
