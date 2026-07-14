package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-gst/go-gst/examples/plugins/registered_elements/internal/custombin"
	"github.com/go-gst/go-gst/examples/plugins/registered_elements/internal/customsrc"
	"github.com/go-gst/go-gst/pkg/gst"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	gst.Init()

	customsrc.Register()
	custombin.Register()

	systemclock := gst.SystemClockObtain()

	log.Printf("using system clock: %s", systemclock.GetName())

	// After registering the elements, we can use them in the parsed pipeline strings as if they were
	// provided by a plugin.
	ret, err := gst.ParseLaunch("gocustombin ! fakesink sync=true")

	if err != nil {
		return err
	}

	pipeline := ret.(gst.Pipeline)

	pipeline.UseClock(systemclock)

	bus := pipeline.GetBus()

	go func() {
		for msg := range bus.Messages(ctx) {
			switch msg.Type() {
			case gst.MessageStateChanged:
				old, new, _ := msg.ParseStateChanged()
				dot := pipeline.DebugBinToDotData(gst.DebugGraphShowVerbose)

				f, err := os.OpenFile(filepath.Join(wd, fmt.Sprintf("pipeline-%s-to-%s.dot", old, new)), os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)

				if err != nil {
					cancel()
					return
				}

				defer f.Close()

				_, err = f.Write([]byte(dot))

				if err != nil {
					fmt.Println(err)
					cancel()
					return
				}

			case gst.MessageEOS:
				fmt.Println("reached EOS")
				cancel()
				return
			default:
				fmt.Println("got message:", msg.String())
			}
		}
	}()

	pipeline.SetState(gst.StatePlaying)

	<-ctx.Done()

	pipeline.BlockSetState(gst.StateNull, gst.ClockTime(time.Second))

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
