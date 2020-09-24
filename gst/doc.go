/*
Package gst provides wrappers for building gstreamer pipelines and then
reading and/or writing from either end of the pipeline.

It uses cgo to interface with the gstreamer-1.0 C API.

A simple opus/webm encoder created from a launch string could look like this:

  import (
	  "os"
	  "github.com/tinyzimmer/go-gst-launch/gst"
  )

  func main() {
	  gst.Init()
	  encoder, err := gst.NewPipelineFromLaunchString("opusenc ! webmmux", gst.PipelineReadWrite)
	  if err != nil {
		  panic(err)
	  }

	  // You should close even if you don't start the pipeline, since this
	  // will free resources created by gstreamer.
	  defer encoder.Close()

	  if err := encoder.Start() ; err != nil {
	      panic(err)
	  }

	  go func() {
		  encoder.Write(...)  // Write raw audio data to the pipeline
	  }()

	  // don't actually do this - copy encoded audio to stdout
	  if _, err  := io.Copy(os.Stdout, encoder) ; err != nil {
		  panic(err)
	  }
  }

You can accomplish the same thing using the "configuration" functionality provided by NewPipelineFromConfig().
Here is an example that will record from a pulse server and make opus/webm data available on the Reader.

  import (
	"io"
	"os"
	"github.com/tinyzimmer/go-gst-launch/gst"
  )

  func main() {
	  gst.Init()
	  encoder, err := gst.NewPipelineFromConfig(&gst.PipelineConfig{
		  Plugins: []*gst.Plugin{
			  {
				  Name: "pulsesrc",
				  Data: map[string]interface{}{
					  "server": "/run/user/1000/pulse/native",
					  "device": "playback-device.monitor",
				  },
				  SinkCaps: gst.NewRawCaps("S16LE", 24000, 2),
			  },
			  {
				  Name: "opusenc",
			  },
			  {
				  Name: "webmmux",
			  },
		  },
	  }, gst.PipelineRead, nil)
	  if err != nil {
		  panic(err)
	  }

	  defer encoder.Close()

	  if err := encoder.Start() ; err != nil {
		  panic(err)
	  }

	  // Create an output file
	  f, err := os.Create("out.opus")
	  if err != nil {
		  panic(err)
	  }

	  // Copy the data from the pipeline to the file
	  if err := io.Copy(f, encoder) ; err != nil {
		  panic(err)
	  }

  }


There are two channels exported for listening for messages from the pipeline.
An example of listening to messages on a fake pipeline for 10 seconds:

  package main

  import (
	"fmt"
	"time"

	"github.com/tinyzimmer/go-gst-launch/gst"
  )

  func main() {
	  gst.Init()

	  pipeline, err := gst.NewPipelineFromLaunchString("audiotestsrc ! fakesink",  gst.PipelineInternalOnly)
	  if err != nil {
	 	  panic(err)
	  }

	  defer pipeline.Close()

	  go func() {
		  for msg := range pipeline.MessageChan() {
			  fmt.Println("Got message:", msg.TypeName())
		  }
	  }()

	  go func() {
		  for msg := range pipeline.ErrorChan() {
			  fmt.Println("Got error:", err)
		  }
	  }()

	  if err := pipeline.Start(); err != nil {
		  fmt.Println("Pipeline failed to start")
		  return
	  }

	  time.Sleep(time.Second * 10)
  }


The package also exposes some low level functionality for building pipelines
and doing dynamic linking yourself. See the NewPipeline() function for creating an
empty pipeline that you can then build out using the other structs and methods provided
by this package.

*/
package gst
