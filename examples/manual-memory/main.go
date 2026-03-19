// This demonstrates manual memory management in GStreamer using go-gst. This is useful when handling a lot of
// GstBuffer or other frequently created/destroyed objects, where the overhead of Go's garbage collector
// would be too high.
package main

import (
	"runtime"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

func main() {
	gst.Init()

	buffer := gst.NewBuffer()

	// when done:
	gst.UnsafeBufferUnref(buffer)

	// GC will not manage the memory of the buffer anymore

	runtime.GC()
	runtime.GC()
	runtime.GC()

	time.Sleep(5 * time.Second)
}
