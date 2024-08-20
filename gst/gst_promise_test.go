package gst

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"
)

//go:noinline
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
	initOnce.Do(func() {
		Init(nil)
	})

	prom := NewPromise()
	cprom := prom.Instance()

	reply := NewStructure("foo/bar")
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

		if res.Name() != reply.Name() {
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

	if cprom.parent.refcount != 1 {
		panic("refcount too high")
	}

	runtime.KeepAlive(prom)

	awaitGC()
}

var initOnce sync.Once

func TestPromiseMarshal(t *testing.T) {
	initOnce.Do(func() {
		Init(nil)
	})

	prom := NewPromise()

	gv, err := prom.ToGValue()

	if err != nil {
		t.Fatal(err)
	}

	receivedPromI, err := marshalPromise(unsafe.Pointer(gv.GValue))

	if err != nil {
		t.Fatal(err)
	}

	receivedProm, ok := receivedPromI.(*Promise)

	if !ok {
		t.Fatal("could not cast")
	}

	// Awaiting received promise should error immediately
	_, err = receivedProm.Await(context.Background())

	if err == nil {
		t.FailNow()
	}
}
