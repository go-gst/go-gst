package gstmpegts

// #cgo pkg-config: gstreamer-mpegts-1.0
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <gst/mpegts/mpegts.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// GetSectionType returns the section_type member of the struct
func (section *Section) GetSectionType() SectionType {
	ret := SectionType(section.section.native.section_type)
	runtime.KeepAlive(section)
	return ret
}

// GetSpliceCommandType returns the splice_command_type member of the struct
func (sit *SCTESIT) GetSpliceCommandType() SCTESpliceCommandType {
	ret := SCTESpliceCommandType(sit.sCTESIT.native.splice_command_type)
	runtime.KeepAlive(sit)
	return ret
}

// GetSpliceTimeSpecified returns the splice_time_specified member of the struct
func (sit *SCTESIT) GetSpliceTimeSpecified() bool {
	ret := sit.sCTESIT.native.splice_time_specified != 0
	runtime.KeepAlive(sit)
	return ret
}

// GetSplices returns the splices member of the struct. The events are owned by sit.
func (sit *SCTESIT) GetSplices() []*SCTESpliceEvent {
	splices := sit.sCTESIT.native.splices
	if splices == nil {
		runtime.KeepAlive(sit)
		return nil
	}

	events := unsafe.Slice((**C.GstMpegtsSCTESpliceEvent)(unsafe.Pointer(splices.pdata)), splices.len)
	ret := make([]*SCTESpliceEvent, 0, len(events))
	for _, e := range events {
		if e == nil {
			continue
		}
		event := UnsafeSCTESpliceEventFromGlibBorrow(unsafe.Pointer(e))
		runtime.AddCleanup(event, func(_ *SCTESIT) {}, sit)
		ret = append(ret, event)
	}
	runtime.KeepAlive(sit)
	return ret
}

// GetSpliceEventID returns the splice_event_id member of the struct
func (event *SCTESpliceEvent) GetSpliceEventID() uint32 {
	ret := uint32(event.sCTESpliceEvent.native.splice_event_id)
	runtime.KeepAlive(event)
	return ret
}

// GetSpliceEventCancelIndicator returns the splice_event_cancel_indicator member of the struct
func (event *SCTESpliceEvent) GetSpliceEventCancelIndicator() bool {
	ret := event.sCTESpliceEvent.native.splice_event_cancel_indicator != 0
	runtime.KeepAlive(event)
	return ret
}

// GetOutOfNetworkIndicator returns the out_of_network_indicator member of the struct
func (event *SCTESpliceEvent) GetOutOfNetworkIndicator() bool {
	ret := event.sCTESpliceEvent.native.out_of_network_indicator != 0
	runtime.KeepAlive(event)
	return ret
}

// GetSpliceImmediateFlag returns the splice_immediate_flag member of the struct
func (event *SCTESpliceEvent) GetSpliceImmediateFlag() bool {
	ret := event.sCTESpliceEvent.native.splice_immediate_flag != 0
	runtime.KeepAlive(event)
	return ret
}
