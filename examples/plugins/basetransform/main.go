package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-gst/go-gst/examples/plugins/basetransform/internal/customtransform"
	"github.com/go-gst/go-gst/pkg/gst"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	gst.Init()

	customtransform.Register()

	ret, err := gst.ParseLaunch("audiotestsrc ! gocustomtransform ! fakesink")

	if err != nil {
		return err
	}

	pipeline := ret.(*gst.Pipeline)

	pipeline.SetState(gst.StatePlaying)

	<-ctx.Done()

	pipeline.BlockSetState(gst.StateNull, gst.ClockTime(time.Second))

	return ctx.Err()
}

func main() {
	ctx := context.Background()

	err := run(ctx)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
