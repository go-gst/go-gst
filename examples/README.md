# Go-gst Examples

This directory contains examples of some common use cases of gstreamer using the go bindings.

The common package provided to each example exports two methods.

 - `Run(f)` - This wraps the given function in a goroutine and wraps a GMainLoop around it.
 - `RunLoop(f(loop))` - This simply creates (but does not start) a GMainLoop and passes it to the example to manage.

Each example can be run in one of two ways:

```bash
# For single-file examples
go run <example>/main.go [..args]

# For multiple-file examples (but would also work for single file examples)
cd <example> && go build .
./<example> [..args]
```