// This example demonstrates the use of GStreamer's pad probe APIs.
//
// Probes are callbacks that can be installed by the application and will notify
// the application about the states of the dataflow. Those are mostly used for
// changing pipelines dynamically at runtime or for inspecting/modifying buffers or events
//
//                    |-[probe]
//                   /
//   {audiotestsrc} - {fakesink}
package main

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
)

func padProbes(mainLoop *gst.MainLoop) error {
	gst.Init(nil)

	// Parse the pipeline we want to probe from a static in-line string.
	// Here we give our audiotestsrc a name, so we can retrieve that element
	// from the resulting pipeline.
	pipeline, err := gst.NewPipelineFromString(
		"audiotestsrc name=src ! audio/x-raw,format=S16LE,channels=1 ! fakesink",
	)

	// Get the audiotestsrc element from the pipeline that GStreamer
	// created for us while parsing the launch syntax above.
	//
	// TODO: There are some discrepancies still between methods that check the nil
	// value and return an error, versus those that will instead just return nil.
	// Need to settle on one way or the other.
	src, err := pipeline.GetElementByName("src")
	if err != nil {
		return err
	}
	// Get the audiotestsrc's src-pad.
	srcPad := src.GetStaticPad("src")
	if srcPad == nil {
		return errors.New("src pad on src element was nil")
	}

	// Add a probe handler on the audiotestsrc's src-pad.
	// This handler gets called for every buffer that passes the pad we probe.
	srcPad.AddProbe(gst.PadProbeTypeBuffer, func(self *gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
		// Interpret the data sent over the pad as a buffer. We know to expect this because of
		// the probe mask defined above.
		buffer := info.GetBuffer()

		// At this point, buffer is only a reference to an existing memory region somewhere.
		// When we want to access its content, we have to map it while requesting the required
		// mode of access (read, read/write).
		// This type of abstraction is necessary, because the buffer in question might not be
		// on the machine's main memory itself, but rather in the GPU's memory.
		// So mapping the buffer makes the underlying memory region accessible to us.
		// See: https://gstreamer.freedesktop.org/documentation/plugin-development/advanced/allocation.html
		mapInfo := buffer.Map(gst.MapRead)
		defer buffer.Unmap()

		// We know what format the data in the memory region has, since we requested
		// it by setting the fakesink's caps. So what we do here is interpret the
		// memory region we mapped as an array of signed 16 bit integers.
		samples := mapInfo.AsInt16Slice()
		if len(samples) == 0 {
			return gst.PadProbeOK
		}

		// For each buffer (= chunk of samples) calculate the root mean square.
		var square float64
		for _, i := range samples {
			square += float64(i * i)
		}
		rms := math.Sqrt(square / float64(len(samples)))
		fmt.Println("rms:", rms)

		return gst.PadProbeOK
	})

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Block on messages coming in from the bus instead of using the main loop
	for {
		msg := pipeline.GetPipelineBus().TimedPop(time.Duration(-1))
		if msg == nil {
			break
		}
		if err := handleMessage(msg); err != nil {
			return err
		}
	}

	return pipeline.Destroy()
}

func handleMessage(msg *gst.Message) error {
	defer msg.Unref()
	switch msg.Type() {
	case gst.MessageEOS:
		return errors.New("end-of-stream")
	case gst.MessageError:
		return msg.ParseError()
	}
	return nil
}

func main() {
	examples.RunLoop(padProbes)
}
