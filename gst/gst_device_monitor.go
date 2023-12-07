package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
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

func (d *DeviceMonitor) AddFilter(classes string, caps *Caps) uint {
	var cClasses *C.gchar
	if classes != "" {
		cClasses = C.CString(classes)
		defer C.free(unsafe.Pointer(cClasses))
	}

	filterId := C.gst_device_monitor_add_filter(d.ptr, cClasses, caps.Instance())
	return uint(filterId)
}

func (d *DeviceMonitor) RemoveFilter(filterId uint) bool {
	return gobool(C.gst_device_monitor_remove_filter(d.ptr, C.guint(filterId)))
}

// GetBus returns the message bus for this pipeline.
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
	goList := glib.WrapList(unsafe.Pointer(glist))
	out := make([]*Device, 0)
	goList.Foreach(func(item interface{}) {
		pt := item.(unsafe.Pointer)
		out = append(out, FromGstDeviceUnsafeFull(pt))
	})
	return out
}

func (d *DeviceMonitor) SetShowAllDevices(show bool) {
	C.gst_device_monitor_set_show_all_devices(d.ptr, gboolean(show))
}

func (d *DeviceMonitor) GetShowAllDevices() bool {
	return gobool(C.gst_device_monitor_get_show_all_devices(d.ptr))
}

//gst_device_monitor_get_providers
