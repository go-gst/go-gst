// This example shows how to use the appsrc element.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

func createPipeline() (*gst.Pipeline, error) {
	gst.Init(nil)

	// Create a pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	// Create the elements
	elems, err := gst.NewElementMany("appsrc", "autoaudiosink")
	if err != nil {
		return nil, err
	}

	// Add the elements to the pipeline and link them
	pipeline.AddMany(elems...)
	gst.ElementLinkMany(elems...)

	// Get the app sourrce from the first element returned
	src := app.SrcFromElement(elems[0])

	// We are instructing downstream elements that we are producing raw signed 16-bit integers.
	src.SetCaps(gst.NewCapsFromString(
		"audio/x-raw, format=S16LE, layout=interleaved, channels=1, rate=44100",
	))

	// Add a callback for whene the sink requests a sample
	i := 1
	src.SetCallbacks(&app.SourceCallbacks{
		NeedDataFunc: func(src *app.Source, _ uint) {
			// Stop after 10 samples
			if i == 10 {
				src.EndStream()
				return
			}

			fmt.Println("Producing sample", i)

			sinWave := newSinWave(44100, 440.0, 1.0, time.Second)

			// Allocate a new buffer with the sin wave
			buffer := gst.NewBufferFromBytes(sinWave)

			// Set the presentation timestamp on thee buffer
			pts := time.Second * time.Duration(i)
			buffer.SetPresentationTimestamp(pts)
			buffer.SetDuration(time.Second)

			// Push tehe buffer onto the src
			src.PushBuffer(buffer)

			i++
		},
	})

	return pipeline, nil
}

func newSinWave(sampleRate int64, freq, vol float64, duration time.Duration) []byte {
	numSamples := duration.Milliseconds() * (sampleRate / 1000.0)
	buf := new(bytes.Buffer)
	for i := int64(0); i < numSamples; i++ {
		data := vol * math.Sin(2.0*math.Pi*freq*(1/float64(sampleRate)))
		binary.Write(buf, binary.LittleEndian, data)
	}
	return buf.Bytes()
}

func handleMessage(msg *gst.Message) error {
	defer msg.Unref() // Messages are a good candidate for trying out runtime finalizers

	switch msg.Type() {
	case gst.MessageEOS:
		return app.ErrEOS
	case gst.MessageError:
		return msg.ParseError()
	}

	return nil
}

func mainLoop(pipeline *gst.Pipeline) error {

	defer pipeline.Destroy() // Will stop and unref the pipeline when this function returns

	// Start the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Retrieve the bus from the pipeline
	bus := pipeline.GetPipelineBus()

	// Loop over messsages from the pipeline
	for {
		msg := bus.TimedPop(time.Duration(-1))
		if msg == nil {
			break
		}
		if err := handleMessage(msg); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	examples.Run(func() error {
		var pipeline *gst.Pipeline
		var err error
		if pipeline, err = createPipeline(); err != nil {
			return err
		}
		return mainLoop(pipeline)
	})
}
