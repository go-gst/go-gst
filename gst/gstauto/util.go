package gstauto

import "github.com/tinyzimmer/go-gst/gst"

func runOrPrintErr(f func() error) {
	if err := f(); err != nil {
		gst.CAT.Log(gst.LevelError, err.Error())
	}
}
