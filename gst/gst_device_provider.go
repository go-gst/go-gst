package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// DeviceProvider is a Go representation of a GstDeviceProvider.
type DeviceProvider struct {
	ptr *C.GstDeviceProvider
	bus *Bus
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

// GetBus returns the message bus for this pipeline.
func (d *DeviceProvider) GetBus() *Bus {
	if d.bus == nil {
		cBus := C.gst_device_provider_get_bus(d.ptr)
		d.bus = FromGstBusUnsafeFull(unsafe.Pointer(cBus))
	}
	return d.bus
}

func (d *DeviceProvider) Start() bool {
	return gobool(C.gst_device_provider_start(d.ptr))
}

func (d *DeviceProvider) Stop() {
	C.gst_device_provider_stop(d.ptr)
}
