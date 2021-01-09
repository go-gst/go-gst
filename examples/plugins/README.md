# Plugins

This directory contains examples of writing GStreamer plugins using `go-gst`. 
The metadata required by GStreamer is generated via `go generate` with the code for the generator contained in this repo
at [`cmd/gst-plugin-gen`](../../cmd/gst-plugin-gen).

The generator assumes the above is compiled and accessible in your PATH as `gst-plugin-gen`. 
You can build and install it to your `GOPATH` by running `make install-plugin-gen` in the root of the repository.