![banner](./img/go-gst-banner.png)

# go-gst: Go bindings for the GStreamer C libraries

[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/go-gst/go-gst)
[![GoReportCard](https://goreportcard.com/badge/github.com/go-gst/go-gst)](https://goreportcard.com/report/github.com/go-gst/go-gst)
<!-- ![](https://github.com/go-gst/go-gst/workflows/Tests/badge.svg) -->

See [pkg.go.dev](https://pkg.go.dev/github.com/go-gst/go-gst) references for documentation and examples.

Please make sure that you have followed the [official gstreamer installation instructions](https://gstreamer.freedesktop.org/documentation/installing/index.html?gi-language=c) before attempting to use the bindings or file an issue.

The bindings are not structured in a way to make version matching with GStreamer easy. We use github actions to verify against the latest supported GStreamer version that is supported by the action https://github.com/blinemedical/setup-gstreamer. Newer GStreamer versions will also work. Always try to use the [latest version of GStreamer](https://gstreamer.freedesktop.org/releases/).

## Requirements

For building applications with this library you need the following:

 - `cgo`: You must set `CGO_ENABLED=1` in your environment when building.
 - `gcc` and `pkg-config`
 - GStreamer development files (the method for obtaining these will [differ depending on your OS](https://gstreamer.freedesktop.org/documentation/installing/index.html?gi-language=c))
   - The core `gst` package utilizes GStreamer core
   - Subpackages (e.g. `app`, `video`) will require development files from their corresponding GStreamer packages
     - Look at `pkg_config.go` in the imported package to see which C libraries are needed.

### Windows

Compiling on Windows may require some more dancing around than on macOS or Linux.
First, make sure you have [mingw](https://chocolatey.org/packages/mingw) and [pkgconfig](https://chocolatey.org/packages/pkgconfiglite) installed (links are for the Chocolatey packages).
Next, go to the [GStreamer downloads](https://gstreamer.freedesktop.org/download/) page and download the latest "development installer" for your MinGW architecture. 
When running your applications on another Windows system, they will need to have the "runtime" installed as well.

Finally, to compile the application you'll have to manually set your `PKG_CONFIG_PATH` to where you installed the GStreamer development files.
For example, if you installed GStreamer to `C:\gstreamer`:

```ps
PS> $env:PKG_CONFIG_PATH='C:\gstreamer\1.0\mingw_x86_64\lib\pkgconfig'
PS> go build .
```

## Quickstart

For more examples see the `examples` folder [here](examples/).

```go
// This is the same as the `launch` example. See the godoc and other examples for more 
// in-depth usage of the bindings.
package main

import (
    "fmt"
    "os"
    "strings"

    "github.com/go-gst/go-glib/glib"
    "github.com/go-gst/go-gst/gst"
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
        fmt.Println(err)
        os.Exit(2)
    }

    // Add a message handler to the pipeline bus, printing interesting information to the console.
    pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
        switch msg.Type() {
        case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
            pipeline.BlockSetState(gst.StateNull)
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

If you find any issues with the bindings or spot areas where things can be improved, feel free to open a PR or start an Issue. A few things to note:

 - Compilation times are insanely slow when working within the bindings.
 - There are a lot of quirks that make generators difficult to deal with for these bindings, so currently everything is hand written. If you have a need for a new binding, feel free to open an issue or create a PR. Writing CGo bindings is not as hard as it seems. (Take a look at https://github.com/go-gst/go-gst/pull/53 for inspiration)
 - More examples would be nice.
 - Support for writing GStreamer plugins and custom elements via the bindings is there, but not well documented.
 - go-gst follows semantic versioning, so it should always be forward compatible for minor versions. If we find an issue in a function and the only way to fix it is to change the function signature, we will break it in a minor version. That way you "get forced" to use the fixed version.

Please make sure that you use the latest version of GStreamer before submitting an issue. If you are using an older version of GStreamer, please try to reproduce the issue with the [latest version](https://gstreamer.freedesktop.org/releases/) before submitting an issue.