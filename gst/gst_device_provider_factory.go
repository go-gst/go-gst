package gst

// #include "gst.go.h"
import "C"

import (
	"unsafe"
)

func FindDeviceProviderByName(factoryName string) *DeviceProvider {
	cFactoryName := C.CString(factoryName)
	defer C.free(unsafe.Pointer(cFactoryName))
	provider := C.gst_device_provider_factory_get_by_name((*C.gchar)(unsafe.Pointer(cFactoryName)))
	if provider == nil {
		return nil
	}
	return FromGstDeviceProviderUnsafeFull(unsafe.Pointer(provider))
}
