package examples

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
)

// Run is used to wrap the given function in a main loop and print any error
func Run(f func() error) {
	mainLoop := glib.MainLoopNew(glib.MainContextDefault(), false)

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
	mainLoop := glib.MainLoopNew(glib.MainContextDefault(), false)

	if err := f(mainLoop); err != nil {
		fmt.Println("ERROR!", err)
	}
}
