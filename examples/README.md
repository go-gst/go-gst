# Go-gst Examples

This directory contains examples of some common use cases of gstreamer using the go bindings.

Each example can be run in one of two ways:

```bash
# For single-file examples
go run <example>/main.go [..args]

# For multiple-file examples (but would also work for single file examples)
cd <example> && go build .
./<example> [..args]
```

See the plugins subdirectory to learn how to write custom elements in `go-gst`