package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// DeviceMonitor is a Go representation of a GstDeviceMonitor.
type DeviceMonitor struct {
	ptr *C.GstDeviceMonitor
	bus *Bus
}

func NewDeviceMonitor() *DeviceMonitor {
	monitor := C.gst_device_monitor_new()
	if monitor == nil {
		return nil
	}
	return &DeviceMonitor{ptr: monitor}
}

func (d *DeviceMonitor) AddFilter(classes string, caps *Caps) {
	var cClasses *C.gchar
	if classes != "" {
		cClasses = C.CString(classes)
		defer C.free(unsafe.Pointer(cClasses))
	}

	C.gst_device_monitor_add_filter(d.ptr, cClasses, caps.Instance())
	// if caps == nil {
	// 	return nil
	// }
	//should return if we were able to add the filter
}

// GetPipelineBus returns the message bus for this pipeline.
func (d *DeviceMonitor) GetBus() *Bus {
	if d.bus == nil {
		cBus := C.gst_device_monitor_get_bus(d.ptr)
		d.bus = FromGstBusUnsafeFull(unsafe.Pointer(cBus))
	}
	return d.bus
}

func (d *DeviceMonitor) Start() bool {
	return gobool(C.gst_device_monitor_start(d.ptr))
	//should return if we were able to add the filter
}

func (d *DeviceMonitor) Stop() {
	C.gst_device_monitor_stop(d.ptr)
	//should return if we were able to add the filter
}

func (d *DeviceMonitor) GetDevices() []*Device {
	glist := C.gst_device_monitor_get_devices(d.ptr)
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

//https://gstreamer.freedesktop.org/documentation/gstreamer/gstdevicemonitor.html?gi-language=c
//gst_device_monitor_get_providers
//gst_device_monitor_get_show_all_devices
//gst_device_monitor_remove_filter
//gst_device_monitor_set_show_all_devices
