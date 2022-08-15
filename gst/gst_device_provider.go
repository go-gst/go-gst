package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// DeviceMonitor is a Go representation of a GstDeviceMonitor.
type DeviceProvider struct {
	ptr *C.GstDeviceProvider
}

func (d *DeviceProvider) GetDevices() []*Device {
	glist := C.gst_device_provider_get_devices(d.ptr)
	if glist == nil {
		return nil
	}
	goList := glib.WrapList(uintptr(unsafe.Pointer(glist)))
	out := make([]*Device, 0)
	goList.Foreach(func(item interface{}) {
		pt := item.(unsafe.Pointer)
		out = append(out, FromGstDeviceUnsafeFull(pt))
	})
	return out
}
