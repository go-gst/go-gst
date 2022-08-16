package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"
)

// DeviceProviderFactory is a Go representation of a GstDeviceProviderFactory.
type DeviceProviderFactory struct {
	ptr *C.GstDeviceProviderFactory
}

func FindDeviceProviderByName(factoryName string) *DeviceProvider {
	cFactoryName := C.CString(factoryName)
	defer C.free(unsafe.Pointer(cFactoryName))
	provider := C.gst_device_provider_factory_get_by_name((*C.gchar)(unsafe.Pointer(cFactoryName)))
	if provider == nil {
		return nil
	}
	return &DeviceProvider{ptr: provider}
}
