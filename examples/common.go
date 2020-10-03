package examples

import (
	"fmt"

	"github.com/tinyzimmer/go-gst/gst"
)

// Run is used to wrap the given function in a main loop and print any error
func Run(f func() error) {
	mainLoop := gst.NewMainLoop(gst.DefaultMainContext(), false)

	defer mainLoop.Unref()

	go func() {
		if err := f(); err != nil {
			fmt.Println("ERROR!", err)
		}
		mainLoop.Quit()
	}()

	mainLoop.Run()
}

// RunLoop is used to wrap the given function in a main loop and print any error.
// The main loop itself is passed to the function for more control over exiting.
func RunLoop(f func(*gst.MainLoop) error) {
	mainLoop := gst.NewMainLoop(gst.DefaultMainContext(), false)
	defer mainLoop.Unref()

	if err := f(mainLoop); err != nil {
		fmt.Println("ERROR!", err)
	}
}
