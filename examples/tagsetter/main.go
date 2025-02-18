// This example demonstrates how to set and store metadata using GStreamer.
//
// Some elements support setting tags on a media stream. An example would be
// id3v2mux. The element signals this by implementing The GstTagsetter interface.
// You can query any element implementing this interface from the pipeline, and
// then tell the returned implementation of GstTagsetter what tags to apply to
// the media stream.
//
// This example's pipeline creates a new flac file from the testaudiosrc
// that the example application will add tags to using GstTagsetter.
// The operated pipeline looks like this:
//
//	{audiotestsrc} - {flacenc} - {filesink}
//
// For example for pipelines that transcode a multimedia file, the input
// already has tags. For cases like this, the GstTagsetter has the merge
// setting, which the application can configure to tell the element
// implementing the interface whether to merge newly applied tags to the
// already existing ones, or if all existing ones should replace, etc.
// (More modes of operation are possible, see: gst.TagMergeMode)
// This merge-mode can also be supplied to any method that adds new tags.
package main

import (
	"fmt"

	coreglib "github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/go-gst/go-gst/pkg/gst"
)

func tagsetter() error {
	gst.Init()

	ret, err := gst.ParseLaunch(
		"audiotestsrc wave=white-noise num-buffers=10000 ! flacenc ! filesink location=test.flac",
	)
	if err != nil {
		return err
	}

	pipeline := ret.(*gst.Pipeline)

	// Query the pipeline for elements implementing the GstTagsetter interface.
	// In our case, this will return the flacenc element.
	element := pipeline.ByInterface(gst.GTypeTagSetter)

	// We actually just retrieved a *gst.Element with the above call. We can retrieve
	// the underying TagSetter interface like this.
	tagsetter := element.(*gst.TagSetter)

	// Tell the element implementing the GstTagsetter interface how to handle already existing
	// metadata.
	tagsetter.SetTagMergeMode(gst.TagMergeKeepAll)

	// Set the "title" tag to "Special randomized white-noise".
	//
	// The first parameter gst.TagMergeAppend tells the tagsetter to append this title
	// if there already is one.
	tagsetter.AddTagValue(gst.TagMergeAppend, gst.TAG_TITLE, coreglib.NewValue("Special randomized white-noise"))

	pipeline.SetState(gst.StatePlaying)

	var cont bool
	var pipelineErr error
	for {
		msg := pipeline.Bus().TimedPop(gst.ClockTimeNone)
		if msg == nil {
			break
		}
		if cont, pipelineErr = handleMessage(msg); pipelineErr != nil || !cont {
			pipeline.SetState(gst.StateNull)
			break
		}
	}

	return pipelineErr
}

func handleMessage(msg *gst.Message) (bool, error) {
	switch msg.Type() {
	case gst.MessageTag:
		fmt.Println(msg) // Prirnt our tags
	case gst.MessageEos:
		return false, nil
	case gst.MessageError:
		err, _ := msg.ParseError()
		return false, err
	}
	return true, nil
}

func main() {

}
