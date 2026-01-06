package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-gst/go-glib/pkg/gobject/v2"
	"github.com/go-gst/go-gst/pkg/gst"
)

// this is a port of the level example from c https://gstreamer.freedesktop.org/documentation/level/index.html?gi-language=c

type LevelMessageStructure struct {
	Timestamp   gst.ClockTime `gst:"timestamp"`
	StreamTime  gst.ClockTime `gst:"stream-time"`
	RunningTime gst.ClockTime `gst:"running-time"`
	Duration    gst.ClockTime `gst:"duration"`
	EndTime     gst.ClockTime `gst:"endtime"`

	Peak  gobject.ValueArray `gst:"peak"`
	Decay gobject.ValueArray `gst:"decay"`
	RMS   gobject.ValueArray `gst:"rms"`
}

func main() {
	gst.Init()

	parsed, err := gst.ParseLaunch("audiotestsrc ! audio/x-raw,channels=2 ! level post-messages=true ! fakesink sync=true")

	if err != nil {
		log.Fatalf("failed to create pipeline: %v", err)
	}

	pipeline := parsed.(gst.Pipeline)

	go func() {
		for m := range pipeline.GetBus().Messages(context.Background()) {
			if m.Type() == gst.MessageElement {
				structure := m.GetStructure()

				if structure.GetName() == "level" {
					var levelMsg LevelMessageStructure
					err := structure.UnmarshalInto(&levelMsg)
					if err != nil {
						log.Printf("failed to unmarshal level message: %v", err)
						continue
					}

					log.Printf("Level message at %#v:", levelMsg)
					log.Printf("  Peak: %v", levelMsg.Peak)
					log.Printf("  Decay: %v", levelMsg.Decay)
					log.Printf("  RMS: %v", levelMsg.RMS)
				}
			}

			fmt.Println(m.String())
		}
	}()

	pipeline.SetState(gst.StatePlaying)

	select {}
}
