package customsrc

import (
	"github.com/diamondburned/gotk4/pkg/gobject/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

// Register needs to be called after gst.Init() to make the gocustomsrc available in the standard
// gst element registry. After this call the element can be used like any other gstreamer element
func Register() bool {
	registered := gst.RegisterBinSubClass[*customSrc](
		"gocustomsrc",
		classInit,
		nil,
		gst.BinOverrides[*customSrc]{
			ElementOverrides: gst.ElementOverrides[*customSrc]{
				ObjectOverrides: gst.ObjectOverrides[*customSrc]{
					InitiallyUnownedOverrides: gobject.InitiallyUnownedOverrides[*customSrc]{
						ObjectOverrides: gobject.ObjectOverrides[*customSrc]{
							Constructed: (*customSrc).init,
							SetProperty: (*customSrc).setProperty,
							GetProperty: (*customSrc).getProperty,
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
		"gocustomsrc",
		// The rank of the element
		uint(gst.RankNone),
		registered,
	)
}
