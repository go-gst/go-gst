package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// MpegtsSection is a go representation of a GstMpegtsSection
type MpegtsSection struct {
	section *C.GstMpegtsSection
}

// FromGstMpegtsSectionUnsafeFull wraps the given unsafe.Pointer in a MpegtsSection. No ref is taken
// and a finalizer is placed on the resulting object.
func FromGstMpegtsSectionUnsafeFull(section unsafe.Pointer) *MpegtsSection {
	gosection := ToGstMpegtsSection(section)
	runtime.SetFinalizer(gosection, (*MpegtsSection).Unref)
	return gosection
}

// ToGstMpegtsSection converts the given pointer into a MpegtsSection without affecting the ref count or
// placing finalizers.
func ToGstMpegtsSection(section unsafe.Pointer) *MpegtsSection {
	return wrapMpegtsSection((*C.GstMpegtsSection)(section))
}

// Instance returns the underlying GstMpegtsSection instance.
func (m *MpegtsSection) Instance() *C.GstMpegtsSection {
	return C.toGstMpegtsSection(unsafe.Pointer(m.section))
}

// Unref will call `gst_mpegts_section_unref` on the underlying GstMpegtsSection, freeing it from memory.
func (m *MpegtsSection) Unref() { C.mpegtsSectionUnref(m.Instance()) }

// Ref will increase the ref count on this MpegtsSection. This increases the total amount of times
// Unref needs to be called before the object is freed from memory. It returns the underlying
// MpegtsSection object for convenience.
func (m *MpegtsSection) Ref() *MpegtsSection {
	C.mpegtsSectionRef(m.Instance())
	return m
}

func (m *MpegtsSection) SectionType() MpegtsSectionType {
	return MpegtsSectionType(m.Instance().section_type)
}

func (m *MpegtsSection) GetSCTESIT() *MpegtsSCTESIT {
	scteSit := C.gst_mpegts_section_get_scte_sit(m.Instance())
	if scteSit == nil {
		return nil
	}

	ret := ToGstMpegtsSCTESIT(unsafe.Pointer(scteSit))

	// Take a reference on the underlying GstMpegtsSection to ensure that the parsed table stays vaild until MpegtsSCTESIT gets finalized
	ret.section = m

	return ret
}

// MpegtsSCTESIT is a go representation of a SCTE SIT MpegTS section
type MpegtsSCTESIT struct {
	scteSit *C.GstMpegtsSCTESIT
	section *MpegtsSection // keep a reference to the underlying MpegTSSection object to make sure the GstMpegtsSCTESIT doesn't get freed as it is not independently reference counted
}

// ToGstMpegtsSCTESIT converts the given pointer into a MpegtsSCTESIT without affecting the ref count or
// placing finalizers (GstMpegtsSCTESIT is not a reference counted object)
func ToGstMpegtsSCTESIT(scteSit unsafe.Pointer) *MpegtsSCTESIT {
	return wrapMpegtsSCTESIT((*C.GstMpegtsSCTESIT)(scteSit))
}

// Instance returns the underlying GstMpegtsSCTESIT instance.
func (m *MpegtsSCTESIT) Instance() *C.GstMpegtsSCTESIT {
	return m.scteSit
}

func (m *MpegtsSCTESIT) SpliceCommandType() MpegtsSCTESpliceCommandType {
	return MpegtsSCTESpliceCommandType(m.Instance().splice_command_type)
}

func (m *MpegtsSCTESIT) SpliceTimeSpecified() bool {
	return gobool(m.Instance().splice_time_specified)
}

func (m *MpegtsSCTESIT) SpliceTime() uint64 {
	return uint64(m.Instance().splice_time)
}

func (m *MpegtsSCTESIT) Splices() []*MpegtsSCTESpliceEvent {
	if m.Instance().splices == nil {
		return nil
	}

	ret := []*MpegtsSCTESpliceEvent{}
	for i := uint(0); i < uint(m.Instance().splices.len); i++ {
		ptr := *(**C.GstMpegtsSCTESpliceEvent)(unsafe.Pointer(uintptr(unsafe.Pointer(m.Instance().splices.pdata)) + unsafe.Sizeof(*m.Instance().splices.pdata)*uintptr(i)))
		obj := ToGstMpegtsSCTESpliceEvent(unsafe.Pointer(ptr))
		obj.scteSit = m

		ret = append(ret, obj)
	}

	return ret
}

// MpegtsSCTESpliceEvent is a go representation of a SCTE Splice event
type MpegtsSCTESpliceEvent struct {
	spliceEv *C.GstMpegtsSCTESpliceEvent
	scteSit  *MpegtsSCTESIT // keep a reference to the underlying MpegtsSCTESIT to make sure the GstMpegtsSCTESIT doesn't get freed as it is not independently reference counted
}

// ToMpegtsSCTESpliceEvent converts the given pointer into a MpegtsSCTESpliceEvent without affecting the ref count or
// placing finalizers (GstMpegtsSCTESpliceEvent is not a reference counted object)
func ToGstMpegtsSCTESpliceEvent(spliceEv unsafe.Pointer) *MpegtsSCTESpliceEvent {
	return wrapMpegtsSCTESpliceEvent((*C.GstMpegtsSCTESpliceEvent)(spliceEv))
}

// Instance returns the underlying GstMpegtsSCTESIT instance.
func (ev *MpegtsSCTESpliceEvent) Instance() *C.GstMpegtsSCTESpliceEvent {
	return ev.spliceEv
}

func (ev *MpegtsSCTESpliceEvent) SpliceEventId() uint32 {
	return uint32(ev.Instance().splice_event_id)
}

func (ev *MpegtsSCTESpliceEvent) SpliceEventCancelIndicator() bool {
	return gobool(ev.Instance().splice_event_cancel_indicator)
}

func (ev *MpegtsSCTESpliceEvent) OutOfNetworkIndicator() bool {
	return gobool(ev.Instance().out_of_network_indicator)
}

func (ev *MpegtsSCTESpliceEvent) SpliceImmediateFlag() bool {
	return gobool(ev.Instance().splice_immediate_flag)
}

func (ev *MpegtsSCTESpliceEvent) ProgramSpliceFlag() bool {
	return gobool(ev.Instance().program_splice_flag)
}

func (ev *MpegtsSCTESpliceEvent) ProgramSpliceTimeSpecified() bool {
	return gobool(ev.Instance().program_splice_time_specified)
}

func (ev *MpegtsSCTESpliceEvent) ProgramSpliceTime() uint64 {
	return uint64(ev.Instance().program_splice_time)
}
