# go-gst

Go bindings for the GStreamer C libraries

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-rounded)](https://pkg.go.dev/github.com/tinyzimmer/go-gst)
[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/tinyzimmer/go-gst)
[![GoReportCard](https://goreportcard.com/badge/github.com/nanomsg/mangos)](https://goreportcard.com/report/github.com/tinyzimmer/go-gst)
![](https://github.com/tinyzimmer/go-gst/workflows/Tests/badge.svg)

See the [godoc.org](https://godoc.org/github.com/tinyzimmer/go-gst) or [pkg.go.dev](https://pkg.go.dev/github.com/tinyzimmer/go-gst) references for documentation and examples.
As the latter requires published tags, see godoc.org for the latest documentation of master at any point in time.

**This library has not been thoroughly tested and as such is not recommended for mission critical applications yet. If you'd like to try it out and encounter any bugs, feel free to open an Issue or PR.**

## Requirements

For building applications with this library you need the following:

 - `cgo`: You must set `CGO_ENABLED=1` in your environment when building.
 - `gcc` and `pkg-config`
 - GStreamer development files (the method for obtaining these will differ depending on your OS)
   - The core `gst` package utilizes GStreamer core
   - Subpackages (e.g. `app`, `video`) will require development files from their corresponding GStreamer packages
     - Look at `pkg_config.go` in the imported package to see which C libraries are needed.

## Quickstart

For more examples see the `examples` folder [here](examples/).

```go
// This is the same as the `launch` example. See the godoc and other examples for more 
// in-depth usage of the bindings.
package main

import (
    "fmt"
    "os"

    "github.com/tinyzimmer/go-gst/gst"
)

func main() {
    // This example expects a simple `gst-launch-1.0` string as arguments
    if len(os.Args) == 1 {
        fmt.Println("Pipeline string cannot be empty")
        os.Exit(1)
    }

    // Initialize GStreamer
    gst.Init(nil)

    // Create a main loop. This is only required when utilizing signals via the bindings.
    // In this example, the AddWatch on the pipeline bus requires iterating on the main loop.
    mainLoop := gst.NewMainLoop(gst.DefaultMainContext(), false)
    defer mainLoop.Unref()

    // Build a pipeline string from the cli arguments
    pipelineString := strings.Join(os.Args[1:], " ")

    /// Let GStreamer create a pipeline from the parsed launch syntax on the cli.
    pipeline, err := gst.NewPipelineFromString(pipelineString)
    if err != nil {
        fmt.Println("Pipeline string cannot be empty")
        os.Exit(2)
    }

    // Add a message handler to the pipeline bus, printing interesting information to the console.
    pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
        switch msg.Type() {
        case gst.MessageEOS: // When end-of-stream is received stop the main loop
            mainLoop.Quit()
        case gst.MessageError: // Error messages are always fatal
            err := msg.ParseError()
            fmt.Println("ERROR:", err.Error())
            if debug := err.DebugString(); debug != "" {
                fmt.Println("DEBUG:", debug)
            }
            mainLoop.Quit()
        default:
            // All messages implement a Stringer. However, this is
            // typically an expensive thing to do and should be avoided.
            fmt.Println(msg)
        }
        return true
    })

    // Start the pipeline
    pipeline.SetState(gst.StatePlaying)

    // Block and iterate on the main loop
    mainLoop.Run()
    
    // Destroy the pipeline
    if err := pipeline.Destroy() ; err != nil {
        fmt.Println("Error destroying the pipeline:", err)
    }
}
```