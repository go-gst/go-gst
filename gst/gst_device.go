package gst

// #include "gst.go.h"
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

// Device is a Go representation of a GstDevice.
type Device struct{ *Object }

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
	return wrapElement(glib.Take(unsafe.Pointer(elem)))
}

// Caps returns the caps that this device supports. Unref after usage.
func (d *Device) Caps() *Caps {
	return wrapCaps(C.gst_device_get_caps(d.Instance()))
}

// DeviceClass gets the "class" of a device. This is a "/" separated list of classes that
// represent this device. They are a subset of the classes of the GstDeviceProvider that produced
// this device.
func (d *Device) DeviceClass() string {
	class := C.gst_device_get_device_class(d.Instance())
	defer C.g_free((C.gpointer)(unsafe.Pointer(class)))
	return C.GoString(class)
}

// DisplayName gets the user-friendly name of the device.
func (d *Device) DisplayName() string {
	name := C.gst_device_get_display_name(d.Instance())
	defer C.g_free((C.gpointer)(unsafe.Pointer(name)))
	return C.GoString(name)
}

// Properties gets the extra properties of the device.
func (d *Device) Properties() *Structure {
	return wrapStructure(C.gst_device_get_properties(d.Instance()))
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
		return fmt.Errorf("Failed to reconfigure element %s", elem.Name())
	}
	return nil
}
