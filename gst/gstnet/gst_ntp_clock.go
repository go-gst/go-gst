package gstnet

// #include "gst.go.h"
import "C"
import (
	"context"
	"errors"
	"time"
	"unsafe"

	"github.com/go-gst/go-gst/gst"
)

// NTP timestamps are relative to 1. Jan 1900, so we need an offset for 70 Years to be Unix TS compatible
const NTPTimeToUnixEpoch = gst.ClockTime(2208988800 * time.Second)

// NTPClock wraps GstClock
type NTPClock struct{ *gst.Clock }

// ObtainNTPClock returns the default NTPClock. The refcount of the clock will be
// increased so you need to unref the clock after usage.
func ObtainNTPClock(ctx context.Context, name, address string, port int) (*NTPClock, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	currentSystemTime := time.Now().UnixNano() + int64(NTPTimeToUnixEpoch)

	ntpC := &NTPClock{gst.FromGstClockUnsafeFull(unsafe.Pointer(C.gst_ntp_clock_new(cname, caddr, C.gint(port), C.GstClockTime(currentSystemTime))))}

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout reached while trying to sync")
		default:
			if ntpC.IsSynced() {
				return ntpC, nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// get the current time of the clock
//
// if you want to access the clocks actual time value, use the underlying NTPClock.Clock.GetTime() instead
func (ntp *NTPClock) GetTime() time.Time {
	// nanos since 1. Jan 1900
	ntpNanos := NTPClockTime(ntp.Clock.GetTime())

	return ntpNanos.ToDate()
}

type NTPClockTime gst.ClockTime

func (ct NTPClockTime) ToDate() time.Time {
	return time.Unix(0, int64(ct)-int64(NTPTimeToUnixEpoch))
}

func NewNTPClockTime(t time.Time) NTPClockTime {
	return NTPClockTime(t.UnixNano() + int64(NTPTimeToUnixEpoch))
}
