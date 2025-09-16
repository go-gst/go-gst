// This example shows how to use the feature ranking.
package main

import (
	"fmt"

	"github.com/go-gst/go-gst/pkg/gst"
)

func main() {
	gst.Init()

	registry := gst.RegistryGet()

	higherThanHighRank := uint(258)

	plugin := registry.LookupFeature("vtdec_hw")

	plugin.SetRank(higherThanHighRank)

	rank := plugin.GetRank()
	fmt.Println("vtdec_hw rank is:", rank)
}
