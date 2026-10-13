// Package ownershiptest observes actual native destruction in ownership tests.
package ownershiptest

/*
#cgo pkg-config: gstreamer-1.0
#include <gst/gst.h>

typedef struct {
	gint refs;
	gint samples;
	gint buffers;
	gint memories;
} ReleaseCounter;

static ReleaseCounter *release_counter_new(void) {
	ReleaseCounter *counter = g_new0(ReleaseCounter, 1);
	counter->refs = 1;
	return counter;
}

static void release_counter_unref(ReleaseCounter *counter) {
	if (g_atomic_int_dec_and_test(&counter->refs)) {
		g_free(counter);
	}
}

static void released_sample(gpointer data, GstMiniObject *object) {
	ReleaseCounter *counter = data;
	g_atomic_int_inc(&counter->samples);
	release_counter_unref(counter);
}

static void released_buffer(gpointer data, GstMiniObject *object) {
	ReleaseCounter *counter = data;
	g_atomic_int_inc(&counter->buffers);
	release_counter_unref(counter);
}

static void released_memory(gpointer data, GstMiniObject *object) {
	ReleaseCounter *counter = data;
	g_atomic_int_inc(&counter->memories);
	release_counter_unref(counter);
}

static void release_counter_watch(ReleaseCounter *counter, GstSample *sample) {
	GstBuffer *buffer = gst_sample_get_buffer(sample);
	GstMemory *memory = gst_buffer_peek_memory(buffer, 0);
	g_atomic_int_add(&counter->refs, 3);
	gst_mini_object_weak_ref(GST_MINI_OBJECT(sample), released_sample, counter);
	gst_mini_object_weak_ref(GST_MINI_OBJECT(buffer), released_buffer, counter);
	gst_mini_object_weak_ref(GST_MINI_OBJECT(memory), released_memory, counter);
}

static gint release_counter_samples(ReleaseCounter *counter) {
	return g_atomic_int_get(&counter->samples);
}

static gint release_counter_buffers(ReleaseCounter *counter) {
	return g_atomic_int_get(&counter->buffers);
}

static gint release_counter_memories(ReleaseCounter *counter) {
	return g_atomic_int_get(&counter->memories);
}
*/
import "C"

import "unsafe"

type Counter struct {
	native *C.ReleaseCounter
}

func NewCounter() *Counter {
	return &Counter{native: C.release_counter_new()}
}

// WatchSample borrows a sample containing one allocated buffer and memory block.
// Native weak notifications do not retain those objects or their Go wrappers.
func (counter *Counter) WatchSample(sample unsafe.Pointer) {
	C.release_counter_watch(counter.native, (*C.GstSample)(sample))
}

func (counter *Counter) Released() (samples, buffers, memories int) {
	return int(C.release_counter_samples(counter.native)), int(C.release_counter_buffers(counter.native)), int(C.release_counter_memories(counter.native))
}

// Close releases the observer. Pending notifications retain their own counter
// references, so a failing test cannot leave callbacks pointing into freed data.
func (counter *Counter) Close() {
	C.release_counter_unref(counter.native)
	counter.native = nil
}
