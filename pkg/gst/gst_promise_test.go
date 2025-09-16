package gst_test

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/go-gst/go-gst/pkg/gst"
)

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

func TestPromise(t *testing.T) {
	gst.Init()

	prom := gst.NewPromise()

	reply := gst.StructureFromString("foo/bar")
	errchan := make(chan error)

	go func() {
		res, err := prom.Await(context.Background())

		if err != nil {
			errchan <- err
		}

		// even though we don't use the promise, the result structure should be still accessible.
		// the returned structure is owned by the promise, so the promise must not get GC'ed until
		// we don't use the structure anymore
		awaitGC()

		if res.GetName() != reply.GetName() {
			errchan <- errors.New("name mismatch")
		}

		runtime.GC()
		runtime.GC()
		runtime.GC()

		close(errchan)
	}()

	prom.Reply(reply)

	err := <-errchan

	if err != nil {
		t.FailNow()
	}

	awaitGC()

	runtime.KeepAlive(prom)

	awaitGC()
}
