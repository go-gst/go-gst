package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/examples/plugins/registered_elements/internal/common"
	"github.com/go-gst/go-gst/examples/plugins/registered_elements/internal/custombin"
	"github.com/go-gst/go-gst/examples/plugins/registered_elements/internal/customsrc"
	"github.com/go-gst/go-gst/gst"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	wd, err := os.Getwd()

	if err != nil {
		return err
	}

	gst.Init(nil)

	customsrc.Register()
	custombin.Register()

	systemclock := gst.ObtainSystemClock()

	pipeline, err := gst.NewPipelineFromString("gocustombin ! fakesink sync=true")

	if err != nil {
		return err
	}

	pipeline.ForceClock(systemclock.Clock)

	bus := pipeline.GetBus()

	mainloop := glib.NewMainLoop(glib.MainContextDefault(), false)

	pipeline.SetState(gst.StatePlaying)

	bus.AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageStateChanged:
			old, new := msg.ParseStateChanged()
			dot := pipeline.DebugBinToDotData(gst.DebugGraphShowVerbose)

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

		case gst.MessageEOS:
			fmt.Println(msg.String())
			cancel()
			return false
		}

		// the String method is expensive and should not be used in prodution:
		fmt.Println(msg.String())
		return true
	})

	go mainloop.Run()

	go func() {
		<-ctx.Done()

		mainloop.Quit()
	}()

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

	runtime.GC()
	runtime.GC()
	runtime.GC()

	prof := pprof.Lookup("go-glib-reffed-objects")

	prof.WriteTo(os.Stdout, 1)

	// we are creating 3 custom elements in total. If this panics, then the go struct will memory leak
	common.AssertFinalizersCalled(3)
}
