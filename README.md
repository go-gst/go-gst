# go-gst

Go bindings for the GStreamer C libraries

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-rounded)](https://pkg.go.dev/github.com/tinyzimmer/go-gst)
[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/tinyzimmer/go-gst)
[![GoReportCard](https://goreportcard.com/badge/github.com/nanomsg/mangos)](https://goreportcard.com/report/github.com/tinyzimmer/go-gst)
![](https://github.com/tinyzimmer/go-gst/workflows/Tests/badge.svg)

See the [godoc.org](https://godoc.org/github.com/tinyzimmer/go-gst) or [pkg.go.dev](https://pkg.go.dev/github.com/tinyzimmer/go-gst) references for documentation and examples.
As the latter requires published tags, see godoc.org for the latest documentation of master at any point in time.

**This library has not been thoroughly tested and as such is not recommended for mission critical applications yet. If you'd like to try it out and encounter any bugs, feel free to open an Issue or PR. For more information see the [Contributing](#contributing) section.**

Recently almost all memory handling has been moved into the bindings. Some documentation may still reflect the original need to Unref resources, but in most situations that is not the case anymore.

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

    "github.com/tinyzimmer/go-glib/glib"
    "github.com/tinyzimmer/go-gst/gst"
)

func main() {
    // This example expects a simple `gst-launch-1.0` string as arguments
    if len(os.Args) == 1 {
        fmt.Println("Pipeline string cannot be empty")
        os.Exit(1)
    }

    // Initialize GStreamer with the arguments passed to the program. Gstreamer
    // and the bindings will automatically pop off any handled arguments leaving
    // nothing but a pipeline string (unless other invalid args are present).
    gst.Init(&os.Args)

    // Create a main loop. This is only required when utilizing signals via the bindings.
    // In this example, the AddWatch on the pipeline bus requires iterating on the main loop.
    mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

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
}
```

## Contributing

If you find any issues with the bindings or spot areas where things can be improved, feel free to open a PR or start an Issue. Here are a couple of the things on my radar already that I'd be happy to accept help with:

 - Compilation times are insanely slow when working within the bindings. This could be alleviated by further separating aspects of Gstreamer core into their own packages, or removing bindings that would see no use in Go.

 - There are a lot of quirks that make generators difficult to deal with for these bindings. That being said, I'd still like to find a way to start migrating some of them into generated code.

 - The bindings are not structured in a way to make version matching with GStreamer easy. Basically, you need a version compatible with what the bindings were written with (>=1.16).

 - More examples would be nice.

 - Support for writing GStreamer plugins via the bindings is still a work-in-progress. At the very least I need to write more plugings to find more holes. 

    - SWIG could be used to fix the need for global interfaces to be matched to C callbacks (most notably the `URIHandler` currently). The limitation present at the moment is URIHandlers can only be implemented ONCE per plugin.