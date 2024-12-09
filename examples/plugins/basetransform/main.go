package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-gst/go-gst/examples/plugins/basetransform/internal/customtransform"
	"github.com/go-gst/go-gst/gst"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	gst.Init(nil)

	customtransform.Register()

	pipeline, err := gst.NewPipelineFromString("audiotestsrc ! gocustomtransform ! fakesink")

	if err != nil {
		return err
	}

	pipeline.SetState(gst.StatePlaying)

	<-ctx.Done()

	pipeline.BlockSetState(gst.StateNull)

	gst.Deinit()

	return ctx.Err()
}

func main() {
	ctx := context.Background()

	err := run(ctx)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
