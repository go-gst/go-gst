package app

import "fmt"

func runOrPrintErr(f func() error) {
	if err := f(); err != nil {
		fmt.Println("[go-gst/gst/gstauto] Internal Error:", err.Error())
	}
}
