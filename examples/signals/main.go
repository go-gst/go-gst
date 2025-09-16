package main

import (
	"log"
	"log/slog"
	"runtime"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

func mkObject() {
	el := gst.ElementFactoryMake("decodebin", "").(*gst.Bin)

	h := el.ConnectPadAdded(func(newPad *gst.Pad) {
		el.SetName("foobar")
		// log.Println("handler")
	})

	runtime.AddCleanup(el, func(_ struct{}) {
		log.Println("garbage collected element")
	}, struct{}{})

	// runtime.AddCleanup(pipeline, func(el *gst.Bin) {
	// 	el.HandlerDisconnect(h)
	// }, el)

	_ = h

	return
}

func main() {
	gst.Init()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	log.Println("go:")
	for range 5 {
		log.Println("loop")
		mkObject()

		runtime.GC()
		runtime.GC()
		runtime.GC()
	}

	runtime.GC()
	time.Sleep(2 * time.Second)
	runtime.GC()
}
