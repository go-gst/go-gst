package main

import (
	"fmt"
	"time"

	"github.com/tinyzimmer/go-gst/gst"
)

func main() {
	gst.Init(nil)
	pipeline, _ := gst.NewPipelineFromString("fakesrc ! fakesink")
	defer pipeline.Unref()

	clock := pipeline.GetPipelineClock()

	id := clock.NewSingleShotID(gst.ClockTime(time.Minute.Nanoseconds()))

	go func() {
		id.Wait(gst.ClockTimeDiff(time.Minute.NanoSeconds()))
		fmt.Println("I returned")
	}()

	pipeline.SetState(gst.StatePlaying)
	fmt.Println("I am waiting")
	gst.Wait(pipeline)
}
