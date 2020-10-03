package main

import (
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

func createPipeline() (*gst.Pipeline, error) {
	gst.Init(nil)

	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	// Should this actually be a *gst.Element that produces an Appsrc interface like the
	// rust examples?
	src, err := app.NewAppSrc()
	if err != nil {
		return nil, err
	}

	elems, err := gst.NewElementMany("videoconvert", "autovideosink")
	if err != nil {
		return nil, err
	}

	// Place the app source at the top of the element list for linking
	elems = append([]*gst.Element{src.Element}, elems...)

	pipeline.AddMany(elems...)
	gst.ElementLinkMany(elems...)

	// TODO: need to implement video

	return pipeline, nil
}

func main() {

}
