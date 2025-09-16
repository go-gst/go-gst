package main

import (
	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstrtp"
)

// twccURI is the URI for the Transport Wide Congestion Control (TWCC) extension
const twccURI = "http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01"

// this example demonstrates how to create a TWCC RTP header extension, also this is a test for the
// typechecking of the signal emit function
func main() {
	gst.Init()

	payloader := gst.ElementFactoryMake("rtpopuspay", "").(gstrtp.RTPBasePayload)

	twcc := gstrtp.RTPHeaderExtensionCreateFromURI(twccURI)
	twcc.SetID(1)
	payloader.EmitAddExtension(twcc)
}
