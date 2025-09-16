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

	pluginf := registry.LookupFeature("vtdec_hw")

	if pluginf == nil {
		fmt.Printf("codec vtdec_hw not found")

		return
	}

	plugin := gst.BasePluginFeature(pluginf)

	plugin.SetRank(higherThanHighRank)

	rank := plugin.Rank()
	fmt.Println("vtdec_hw rank is:", rank)
}
