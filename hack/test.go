package main

import (
	"fmt"
	"time"

	"github.com/tinyzimmer/go-gst/gst"
)

func wait() {
	gst.Init(nil)

	clock := gst.ObtainSystemClock()

	id := clock.NewSingleShotID(clock.GetTime() + gst.ClockTime(time.Minute.Nanoseconds()))

	go func() {
		res, _ := id.Wait()
		if res != gst.ClockOK {
			panic(res)
		}
		fmt.Println("I waited")
	}()

	fmt.Println("I am waiting")
	time.Sleep(time.Second)
	fmt.Println("Still waiting")
}

func capsWeirdness() {
	gst.Init(nil)

	caps := gst.NewCapsFromString("audio/x-raw")

	caps.ForEach(func(features *gst.CapsFeatures, structure *gst.Structure) bool {
		fmt.Println(features)
		return true
	})

	caps.FilterAndMapInPlace(func(features *gst.CapsFeatures, structure *gst.Structure) bool {
		fmt.Println(features)
		return true
	})

	caps.Unref()
}

func main() {
	wait()
}
