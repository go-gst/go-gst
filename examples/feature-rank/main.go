// This example shows how to use the feature ranking.
package main

import (
	"fmt"

	"github.com/go-gst/go-gst/examples"
	"github.com/go-gst/go-gst/gst"
)

func start() (error) {
	gst.Init(nil)

	registry := gst.GetRegistry()

	higherThanHighRank := (gst.Rank)(258)

	codec, codecErr := registry.LookupFeature("vtdec_hw")

	if codecErr == nil {
		codec.SetPluginRank(higherThanHighRank)
		rank := codec.GetPluginRank()
		fmt.Println("vtdec_hw rank is:", rank)
	}

	codec, codecErr = registry.LookupFeature("vtdec_hw")

	if codecErr == nil {
		codec.SetPluginRank(gst.RankPrimary)
		rank := codec.GetPluginRank()
		fmt.Println("vtdec_hw rank is now:", rank)
	}

	//add a feature you expect to be available to you here and change it's rank

	return codecErr
}

func main() {
	examples.Run(func() error {
		var err error
		if err = start(); err != nil {
			return err
		}
		return nil
	})
}
