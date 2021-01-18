package gst

// #include "gst.go.h"
import "C"

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Device is a Go representation of a GstDevice.
type Device struct{ *Object }

// FromGstDeviceUnsafeNone wraps the given device with a ref and finalizer.
func FromGstDeviceUnsafeNone(device unsafe.Pointer) *Device {
	return &Device{wrapObject(glib.TransferNone(device))}
}

// FromGstDeviceUnsafeFull wraps the given device with a finalizer.
func FromGstDeviceUnsafeFull(device unsafe.Pointer) *Device {
	return &Device{wrapObject(glib.TransferFull(device))}
}

// Instance returns the underlying GstDevice object.
func (d *Device) Instance() *C.GstDevice { return C.toGstDevice(d.Unsafe()) }

// CreateElement creates a new element with all the required parameters set to use this device.
// If name is empty, one is automatically generated.
func (d *Device) CreateElement(name string) *Element {
	var cName *C.gchar
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
	}
	elem := C.gst_device_create_element(d.Instance(), cName)
	if elem == nil {
		return nil
	}
	return FromGstElementUnsafeNone(unsafe.Pointer(elem))
}

// GetCaps returns the caps that this device supports. Unref after usage.
func (d *Device) GetCaps() *Caps {
	caps := C.gst_device_get_caps(d.Instance())
	if caps == nil {
		return nil
	}
	return FromGstCapsUnsafeNone(unsafe.Pointer(caps))
}

// GetDeviceClass gets the "class" of a device. This is a "/" separated list of classes that
// represent this device. They are a subset of the classes of the GstDeviceProvider that produced
// this device.
func (d *Device) GetDeviceClass() string {
	class := C.gst_device_get_device_class(d.Instance())
	defer C.g_free((C.gpointer)(unsafe.Pointer(class)))
	return C.GoString(class)
}

// GetDisplayName gets the user-friendly name of the device.
func (d *Device) GetDisplayName() string {
	name := C.gst_device_get_display_name(d.Instance())
	defer C.g_free((C.gpointer)(unsafe.Pointer(name)))
	return C.GoString(name)
}

// GetProperties gets the extra properties of the device.
func (d *Device) GetProperties() *Structure {
	st := C.gst_device_get_properties(d.Instance())
	if st == nil {
		return nil
	}
	return wrapStructure(st)
}

// HasClasses checks if device matches all of the given classes.
func (d *Device) HasClasses(classes []string) bool {
	cClasses := C.CString(strings.Join(classes, "/"))
	defer C.free(unsafe.Pointer(cClasses))
	return gobool(C.gst_device_has_classes(d.Instance(), cClasses))
}

// ReconfigureElement tries to reconfigure an existing element to use the device.
// If this function fails, then one must destroy the element and create a new one
// using Device.CreateElement().
//
// Note: This should only be implemented for elements that can change their device
// while in the PLAYING state.
func (d *Device) ReconfigureElement(elem *Element) error {
	if ok := gobool(C.gst_device_reconfigure_element(d.Instance(), elem.Instance())); !ok {
		return fmt.Errorf("Failed to reconfigure element %s", elem.GetName())
	}
	return nil
}
