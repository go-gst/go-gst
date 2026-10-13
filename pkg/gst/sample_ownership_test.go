package gst

import (
	"runtime"
	"testing"
	"time"

	"github.com/go-gst/go-gst/pkg/gst/internal/ownershiptest"
)

//go:noinline
func allocateUnreachableSamples(t *testing.T, counter *ownershiptest.Counter, count int) []*uint64 {
	neighbors := make([]*uint64, 0, count*2)
	for range count {
		// Keep pointer-free tiny allocations alive on both sides of the owner.
		neighbors = append(neighbors, new(uint64))
		buffer := NewBufferAllocate(nil, 64, nil)
		if buffer == nil {
			t.Fatal("could not allocate native buffer")
		}
		s := NewSample(buffer, nil, nil, nil)
		if s == nil {
			UnsafeBufferUnref(buffer)
			t.Fatal("could not allocate native sample")
		}
		counter.WatchSample(UnsafeSampleToGlibNone(s))
		// The sample owns the remaining buffer reference. Do not manually unref
		// the sample: its generated owner must release all three native objects.
		UnsafeBufferUnref(buffer)
		neighbors = append(neighbors, new(uint64))
		runtime.KeepAlive(s)
	}
	return neighbors
}

func TestSampleFinalizerWithTinyNeighbors(t *testing.T) {
	Init()

	const count = 1000
	counter := ownershiptest.NewCounter()
	defer counter.Close()
	neighbors := allocateUnreachableSamples(t, counter, count)

	deadline := time.Now().Add(5 * time.Second)
	var samples, buffers, memories int
	for {
		runtime.GC()
		samples, buffers, memories = counter.Released()
		if samples == count && buffers == count && memories == count {
			break
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	// An opaque *C.GstSample owner can share a noscan tiny allocation with
	// these live integers, preventing its finalizer from running at all.
	runtime.KeepAlive(neighbors)
	if samples != count || buffers != count || memories != count {
		t.Fatalf("native releases with live tiny neighbors: samples=%d buffers=%d memories=%d, want %d each", samples, buffers, memories, count)
	}
}

func TestSampleExplicitReleaseClearsFinalizer(t *testing.T) {
	Init()
	counter := ownershiptest.NewCounter()
	defer counter.Close()
	buffer := NewBufferAllocate(nil, 64, nil)
	s := NewSample(buffer, nil, nil, nil)
	counter.WatchSample(UnsafeSampleToGlibNone(s))
	UnsafeBufferUnref(buffer)
	UnsafeSampleUnref(s)
	runtime.KeepAlive(s)
	s = nil
	for range 3 {
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
	}
	if samples, buffers, memories := counter.Released(); samples != 1 || buffers != 1 || memories != 1 {
		t.Fatalf("explicit native releases: samples=%d buffers=%d memories=%d, want 1 each", samples, buffers, memories)
	}
}
