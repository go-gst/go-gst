package gst

import "time"

// ClockTimeNone is a constant that represents an infinite or non-existent time.
const ClockTimeNone ClockTime = 0xFFFFFFFFFFFFFFFF

// Common durations type casted from time.Duration to ClockTime.
const (
	NanoSecond  ClockTime = ClockTime(time.Nanosecond)
	MicroSecond           = ClockTime(time.Microsecond)
	Millisecond           = ClockTime(time.Millisecond)
	Second                = ClockTime(time.Second)
	Minute                = ClockTime(time.Minute)
	Hour                  = ClockTime(time.Hour)
)

// ClockTime wraps GstClockTime
//
// A datatype to hold a time, measured in nanoseconds.
type ClockTime uint64

// String returns a string representation of the ClockTime.
func (c ClockTime) String() string {
	return time.Duration(c).String()
}
