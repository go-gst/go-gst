package main

import (
	"fmt"
	"time"

	"github.com/tinyzimmer/go-gst/gst"
)

func main() {
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
	time.Sleep(time.Minute)
}
