package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"

	"github.com/go-gst/go-glib/glib"
)

// DeviceProvider is a Go representation of a GstDeviceProvider.
type DeviceProvider struct{ *Object }

// FromGstDeviceProviderUnsafeNone wraps the given device with a ref and finalizer.
func FromGstDeviceProviderUnsafeNone(deviceProvider unsafe.Pointer) *DeviceProvider {
	return &DeviceProvider{wrapObject(glib.TransferNone(deviceProvider))}
}

// FromGstDeviceProviderUnsafeFull wraps the given device with a finalizer.
func FromGstDeviceProviderUnsafeFull(deviceProvider unsafe.Pointer) *DeviceProvider {
	return &DeviceProvider{wrapObject(glib.TransferFull(deviceProvider))}
}

// Instance returns the underlying GstDevice object.
func (d *DeviceProvider) Instance() *C.GstDeviceProvider { return C.toGstDeviceProvider(d.Unsafe()) }

func (d *DeviceProvider) GetDevices() []*Device {
	glist := C.gst_device_provider_get_devices((*C.GstDeviceProvider)(d.Instance()))
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

// GetBus returns the message bus for this pipeline.
func (d *DeviceProvider) GetBus() *Bus {
	cBus := C.gst_device_provider_get_bus((*C.GstDeviceProvider)(d.Instance()))
	bus := FromGstBusUnsafeFull(unsafe.Pointer(cBus))
	return bus
}

func (d *DeviceProvider) Start() bool {
	return gobool(C.gst_device_provider_start((*C.GstDeviceProvider)(d.Instance())))
}

func (d *DeviceProvider) Stop() {
	C.gst_device_provider_stop((*C.GstDeviceProvider)(d.Instance()))
}
