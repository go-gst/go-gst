package gst

//#include "gst.go.h"
import "C"

// ProtectionMeta is a go wrapper around C GstProtectionMeta.
type ProtectionMeta struct {
	Meta *Meta
	Info *Structure
}

// GetProtectionMetaInfo retrieves global ProtectionMetaInfo.
func GetProtectionMetaInfo() *MetaInfo {
	return wrapMetaInfo(C.gst_protection_meta_get_info())
}

// FilterProtectionSystemByDecryptors tterates the supplied list of UUIDs
// and checks the GstRegistry for all the decryptors supporting one of the supplied UUIDs.
func FilterProtectionSystemByDecryptors(decryptors []string) []string {
	gArr := gcharStrings(decryptors)
	defer C.g_free((C.gpointer)(gArr))
	avail := C.gst_protection_filter_systems_by_available_decryptors(gArr)
	if avail == nil {
		return nil
	}
	defer C.g_free((C.gpointer)(avail))
	return goStrings(C.sizeOfGCharArray(avail), avail)
}

// SelectProtectionSystem iterates the supplied list of UUIDs and checks the GstRegistry for
// an element that supports one of the supplied UUIDs. If more than one element matches, the
// system ID of the highest ranked element is selected.
func SelectProtectionSystem(decryptors []string) string {
	gArr := gcharStrings(decryptors)
	defer C.g_free((C.gpointer)(gArr))
	avail := C.gst_protection_select_system(gArr)
	if avail == nil {
		return ""
	}
	defer C.g_free((C.gpointer)(avail))
	return C.GoString(avail)
}
