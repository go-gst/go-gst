package gstmpegts_test

import (
	"bytes"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
	"github.com/go-gst/go-gst/pkg/gstmpegts"
)

func init() {
	gst.Init()
	gstmpegts.Initialize()
}

func awaitGC() {
	type tmp struct{ v string }

	setup := make(chan struct{})
	done := make(chan struct{})

	go func() {
		v := &tmp{"foo"}

		runtime.SetFinalizer(v, func(v *tmp) {
			close(done)
		})

		close(setup)
	}()

	<-setup
	runtime.GC()
	<-done
	runtime.GC()
	time.Sleep(1 * time.Second)
}

func isWritable(section *gstmpegts.Section) bool {
	writable := gst.UnsafeMiniObjectFromGlibBorrow(gstmpegts.UnsafeSectionToGlibNone(section)).IsWritable()
	runtime.KeepAlive(section)
	return writable
}

func newSpliceOutSection(eventID uint32) *gstmpegts.Section {
	sit := gstmpegts.NewScteSpliceOut(eventID, gst.ClockTime(5*time.Second), gst.ClockTime(30*time.Second))
	return gstmpegts.SectionFromScteSit(sit, 500)
}

// A parsed section must drop only the reference it was given, leaving the message's reference intact.
func TestParsedSectionUnref(t *testing.T) {
	section := newSpliceOutSection(7)
	func() {
		msg := gstmpegts.MessageNewMpegtsSection(gst.ElementFactoryMake("identity", ""), section)
		if isWritable(section) {
			t.Fatal("section is writable while the message holds a reference")
		}

		for range 10 {
			if gstmpegts.MessageParseMpegtsSection(msg) == nil {
				t.Fatal("message has no section")
			}
		}
		awaitGC()
		runtime.KeepAlive(msg)
	}()
	awaitGC()

	if !isWritable(section) {
		t.Fatal("section still has extra references after the message and parsed sections were released")
	}
}

func TestScteSitGetters(t *testing.T) {
	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(prev)

	section := newSpliceOutSection(7)
	if got := section.GetSectionType(); got != gstmpegts.SectionScteSit {
		t.Fatalf("section type is %v, want %v", got, gstmpegts.SectionScteSit)
	}

	sit := section.GetScteSit()
	if sit == nil {
		t.Fatal("section has no SCTE SIT")
	}
	if got := sit.GetSpliceCommandType(); got != gstmpegts.GstMtsScteSpliceCommandInsert {
		t.Fatalf("splice command type is %v, want %v", got, gstmpegts.GstMtsScteSpliceCommandInsert)
	}

	splices := sit.GetSplices()
	if len(splices) != 1 {
		t.Fatalf("got %d splices, want 1", len(splices))
	}
	e := splices[0]
	if e.GetSpliceEventID() != 7 || !e.GetOutOfNetworkIndicator() || e.GetSpliceEventCancelIndicator() || e.GetSpliceImmediateFlag() {
		t.Fatalf("unexpected splice event: id=%d out=%v cancel=%v immediate=%v",
			e.GetSpliceEventID(), e.GetOutOfNetworkIndicator(), e.GetSpliceEventCancelIndicator(), e.GetSpliceImmediateFlag())
	}

	if buf.Len() != 0 {
		t.Fatalf("getters logged: %s", buf.String())
	}
}
