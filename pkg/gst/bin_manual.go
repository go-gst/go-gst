package gst

// // #cgo pkg-config: gstreamer-1.0
// // #cgo CFLAGS: -Wno-deprecated-declarations
// // #include <gst/gst.h>
// import "C"

type BinExtManual interface {
	// DebugBinToDotData wraps gst_debug_bin_to_dot_data
	DebugBinToDotData(details DebugGraphDetails) string

	// DebugBinToDotFile wraps gst_debug_bin_to_dot_file
	DebugBinToDotFile(details DebugGraphDetails, fileName string)

	// DebugBinToDotFileWithTs wraps gst_debug_bin_to_dot_file_with_ts
	DebugBinToDotFileWithTs(details DebugGraphDetails, fileName string)

	// AddMany adds many elements at once to the bin
	AddMany(els ...Element) bool

	// RemoveMany removes many elements at once from the bin
	RemoveMany(els ...Element) bool
}

// DebugBinToDotData wraps gst_debug_bin_to_dot_data
func (bin *BinInstance) DebugBinToDotData(details DebugGraphDetails) string {
	return DebugBinToDotData(bin, details)
}

// DebugBinToDotFile wraps gst_debug_bin_to_dot_file
func (bin *BinInstance) DebugBinToDotFile(details DebugGraphDetails, fileName string) {
	DebugBinToDotFile(bin, details, fileName)
}

// DebugBinToDotFileWithTs wraps gst_debug_bin_to_dot_file_with_ts
func (bin *BinInstance) DebugBinToDotFileWithTs(details DebugGraphDetails, fileName string) {
	DebugBinToDotFileWithTs(bin, details, fileName)
}

func (bin *BinInstance) AddMany(els ...Element) bool {

	for _, el := range els {
		if !bin.Add(el) {

			return false
		}
	}

	return true
}

func (bin *BinInstance) RemoveMany(els ...Element) bool {
	for _, el := range els {
		if !bin.Remove(el) {
			return false
		}
	}

	return true
}
