package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
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

	ret, err := gst.ParseLaunch("gocustombin ! fakesink sync=true")

	if err != nil {
		return err
	}

	pipeline := ret.(*gst.Pipeline)

	pipeline.UseClock(systemclock)

	bus := pipeline.Bus()

	mainloop := glib.NewMainLoop(glib.MainContextDefault(), false)

	pipeline.SetState(gst.StatePlaying)

	bus.AddWatch(0, func(bus *gst.Bus, msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageStateChanged:
			old, new, _ := msg.ParseStateChanged()
			dot := gst.DebugBinToDotData(&pipeline.Bin, gst.DebugGraphShowVerbose)

			f, err := os.OpenFile(filepath.Join(wd, fmt.Sprintf("pipeline-%s-to-%s.dot", old, new)), os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)

			if err != nil {
				cancel()
				return false
			}

			defer f.Close()

			_, err = f.Write([]byte(dot))

			if err != nil {
				fmt.Println(err)
				cancel()
				return false
			}

		case gst.MessageEos:
			fmt.Println("reached EOS")
			cancel()
			return false
		}

		return true
	})

	go mainloop.Run()

	<-ctx.Done()

	mainloop.Quit()

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
