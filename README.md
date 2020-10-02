# go-gst

Go bindings for the gstreamer C library

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-rounded)](https://pkg.go.dev/github.com/tinyzimmer/go-gst)
[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/tinyzimmer/go-gst)
[![GoReportCard](https://goreportcard.com/badge/github.com/nanomsg/mangos)](https://goreportcard.com/report/github.com/tinyzimmer/go-gst)
![](https://github.com/tinyzimmer/go-gst/workflows/Tests/badge.svg)

See the go.dev reference for documentation and examples.

For other examples see the command line implementation [here](cmd/go-gst).

_TODO: Write examples on programatically building the pipeline yourself_

## Quickstart

```go
package main

import (
	"io"
	"os"

	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/gstauto"
)

var srcFile, destFile *os.File

func main() {
	gst.Init(nil)

	pipeline, err := gstauto.NewPipelineReadWriterSimpleFromString("opusenc ! webmmux")
	if err != nil {
		panic(err)
	}
	defer pipeline.Close()

	pipeline.Start()

	// Write RAW audio data to the pipeline
	go io.Copy(pipeline, srcFile)
	// Write opus/webm to a destination file
	go io.Copy(destFile, pipeline)

	gst.Wait(pipeline.Pipeline())
}

```

## Requirements

For building applications with this library you need the following:

 - `cgo`: You must set `CGO_ENABLED=1` in your environment when building.
 - `gcc` and `pkg-config`
 - `libgstreamer-1.0-dev`: This package name may be different depending on your OS. You need the `gst.h` header files.
   - In some distributions (such as alpine linux) this is in the `gstreamer-dev` package.
 - To use the `app` or `gstauto/app` packages you will need additional dependencies:
   - `libgstreamer-app-1.0-dev`: This package name may also be different depending on your os. You need the `gstappsink.h` and `gstappsrc.h`
     - In some distributions (such as alpine linux) this is in the `gst-plugins-base-dev` package.
     - In Ubuntu this is in `libgstreamer-plugins-base1.0-0`.
 - You may need platform specific headers also. For example, in alpine linux, you will most likely also need the `musl-dev` package.

For running applications with this library you'll need to have `libgstreamer-1.0` installed. Again, this package may be different depending on your OS.


## CLI

There is a CLI utility included with this package that demonstrates some of the things you can do.

For now the functionality is limitted to GIF encoing, inspection, and other arbitrary pipelines.
If I extend it further I'll publish releases, but for now, you can retrieve it with `go get`.

```bash
go get github.com/tinyzimmer/go-gst/cmd/go-gst
```

The usage is described below:

```
Go-gst is a CLI utility aiming to implement the core functionality
of the core gstreamer-tools. It's primary purpose is to showcase the functionality of 
the underlying go-gst library.

There are also additional commands showing some of the things you can do with the library,
such as websocket servers reading/writing to/from local audio servers and audio/video/image
encoders/decoders.

Usage:
  go-gst [command]

Available Commands:
  completion  Generate completion script
  gif         Encodes the given video to GIF format
  help        Help about any command
  inspect     Inspect the elements of the given pipeline string
  launch      Run a generic pipeline
  websocket   Run a websocket audio proxy for streaming audio from a pulse server 
              and optionally recording to a virtual mic.

Flags:
  -I, --from-stdin      Write to the pipeline from stdin. If this is specified, then -i is ignored.
  -h, --help            help for go-gst
  -i, --input string    An input file, defaults to the first element in the pipeline.
  -o, --output string   An output file, defaults to the last element in the pipeline.
  -O, --to-stdout       Writes the results from the pipeline to stdout. If this is specified, then -o is ignored.
  -v, --verbose         Verbose output. This is ignored when used with --to-stdout.

Use "go-gst [command] --help" for more information about a command.
```