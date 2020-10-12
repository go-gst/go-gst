# go-gst

Go bindings for the gstreamer C library

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-rounded)](https://pkg.go.dev/github.com/tinyzimmer/go-gst)
[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/tinyzimmer/go-gst)
[![GoReportCard](https://goreportcard.com/badge/github.com/nanomsg/mangos)](https://goreportcard.com/report/github.com/tinyzimmer/go-gst)
![](https://github.com/tinyzimmer/go-gst/workflows/Tests/badge.svg)

See the [godoc.org](https://godoc.org/github.com/tinyzimmer/go-gst) or [pkg.go.dev](https://pkg.go.dev/github.com/tinyzimmer/go-gst) references for documentation and examples.
As the latter requires published tags, see godoc.org for the latest documentation of master at any point in time.

For more examples see the `examples` folder [here](examples/).

**This library still has some bugs and should not be used for any mission critical applications, yet. If you'd like to help out feel free to open a PR.**

## Requirements

For building applications with this library you need the following:

 - `cgo`: You must set `CGO_ENABLED=1` in your environment when building.
 - `gcc` and `pkg-config`
 - `libgstreamer-1.0-dev`: This package name may be different depending on your OS. You need the `gst.h` header files.
   - In some distributions (such as alpine linux) this is in the `gstreamer-dev` package.
 - To use the `pbutils`, `app`, `gstauto/app` packages you will need additional dependencies:
   - `libgstreamer-app-1.0-dev`: This package name may also be different depending on your os. You need the `gstappsink.h` and `gstappsrc.h`
     - In some distributions (such as alpine linux) this is in the `gst-plugins-base-dev` package.
     - In Ubuntu this is in `libgstreamer-plugins-base1.0-0`.
 - You may need platform specific headers also. For example, in alpine linux, you will most likely also need the `musl-dev` package.

For running applications with this library you'll need to have `libgstreamer-1.0` installed. Again, this package may be different depending on your OS.
