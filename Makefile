
build-cmd:
	cd cmd/go-gst && go build -o ../../dist/go-gst

ARGS ?=
run-cmd: build-cmd
	dist/go-gst $(ARGS)
