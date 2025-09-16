package gst

// // #cgo pkg-config: gstreamer-1.0
// // #cgo CFLAGS: -Wno-deprecated-declarations
// // #include <gst/gst.h>
// import "C"

type BinExtManual interface {
	// DebugBinToDotData wraps gst_debug_bin_to_dot_data
	//
	// The function takes the following parameters:
	//
	// 	- bin Bin: the top-level pipeline that should be analyzed
	// 	- details DebugGraphDetails: type of #GstDebugGraphDetails to use
	//
	// The function returns the following values:
	//
	// 	- goret string
	//
	// To aid debugging applications one can use this method to obtain the whole
	// network of gstreamer elements that form the pipeline into a dot file.
	// This data can be processed with graphviz to get an image.
	DebugBinToDotData(details DebugGraphDetails) string

	// DebugBinToDotFile wraps gst_debug_bin_to_dot_file
	//
	// The function takes the following parameters:
	//
	//   - bin Bin: the top-level pipeline that should be analyzed
	//   - details DebugGraphDetails: type of #GstDebugGraphDetails to use
	//   - fileName string: output base filename (e.g. "myplayer")
	//
	// To aid debugging applications one can use this method to write out the whole
	// network of gstreamer elements that form the pipeline into a dot file.
	// This file can be processed with graphviz to get an image.
	//
	// ``` shell
	//
	//	dot -Tpng -oimage.png graph_lowlevel.dot
	//
	// ```
	DebugBinToDotFile(details DebugGraphDetails, fileName string)

	// DebugBinToDotFileWithTs wraps gst_debug_bin_to_dot_file_with_ts
	//
	// The function takes the following parameters:
	//
	//   - bin Bin: the top-level pipeline that should be analyzed
	//   - details DebugGraphDetails: type of #GstDebugGraphDetails to use
	//   - fileName string: output base filename (e.g. "myplayer")
	//
	// This works like gst_debug_bin_to_dot_file(), but adds the current timestamp
	// to the filename, so that it can be used to take multiple snapshots.
	DebugBinToDotFileWithTs(details DebugGraphDetails, fileName string)

	// AddMany adds many elements at once to the bin
	AddMany(els ...Element) bool

	// RemoveMany removes many elements at once from the bin
	RemoveMany(els ...Element) bool
}

// DebugBinToDotData wraps gst_debug_bin_to_dot_data
//
// The function takes the following parameters:
//
//   - bin Bin: the top-level pipeline that should be analyzed
//   - details DebugGraphDetails: type of #GstDebugGraphDetails to use
//
// The function returns the following values:
//
//   - goret string
//
// To aid debugging applications one can use this method to obtain the whole
// network of gstreamer elements that form the pipeline into a dot file.
// This data can be processed with graphviz to get an image.
func (bin *BinInstance) DebugBinToDotData(details DebugGraphDetails) string {
	return DebugBinToDotData(bin, details)
}

// DebugBinToDotFile wraps gst_debug_bin_to_dot_file
//
// The function takes the following parameters:
//
//   - bin Bin: the top-level pipeline that should be analyzed
//   - details DebugGraphDetails: type of #GstDebugGraphDetails to use
//   - fileName string: output base filename (e.g. "myplayer")
//
// To aid debugging applications one can use this method to write out the whole
// network of gstreamer elements that form the pipeline into a dot file.
// This file can be processed with graphviz to get an image.
//
// ``` shell
//
//	dot -Tpng -oimage.png graph_lowlevel.dot
//
// ```
func (bin *BinInstance) DebugBinToDotFile(details DebugGraphDetails, fileName string) {
	DebugBinToDotFile(bin, details, fileName)
}

// DebugBinToDotFileWithTs wraps gst_debug_bin_to_dot_file_with_ts
//
// The function takes the following parameters:
//
//   - bin Bin: the top-level pipeline that should be analyzed
//   - details DebugGraphDetails: type of #GstDebugGraphDetails to use
//   - fileName string: output base filename (e.g. "myplayer")
//
// This works like gst_debug_bin_to_dot_file(), but adds the current timestamp
// to the filename, so that it can be used to take multiple snapshots.
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
