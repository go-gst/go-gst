package examples

import (
	"fmt"

	"github.com/tinyzimmer/go-glib/glib"
)

// Run is used to wrap the given function in a main loop and print any error
func Run(f func() error) {
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

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
func RunLoop(f func(*glib.MainLoop) error) {
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)
	defer mainLoop.Unref()

	if err := f(mainLoop); err != nil {
		fmt.Println("ERROR!", err)
	}
}
