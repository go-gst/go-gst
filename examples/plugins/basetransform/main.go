package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/profile"
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

	pipeline := ret.(gst.Pipeline)

	pipeline.SetState(gst.StatePlaying)

	<-ctx.Done()

	pipeline.BlockSetState(gst.StateNull, gst.ClockTime(time.Second))

	return ctx.Err()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := run(ctx)

	for range 10 {
		runtime.GC()
	}

	if profile.Count() > 0 {
		fmt.Fprintf(os.Stderr, "Memory leak detected: %d objects still tracked\n", profile.Count())
	} else {
		fmt.Fprintln(os.Stderr, "No memory leaks detected")
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
