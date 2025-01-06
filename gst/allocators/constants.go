package allocators

// #include "gst.go.h"
import "C"

// FdMemoryFlags represent flags for wrapped memory
type FdMemoryFlags int

// Type castins of FdMemoryFlags
const (
	FdMemoryFlagNone       FdMemoryFlags = C.GST_FD_MEMORY_FLAG_NONE        // (0) – no flags
	FdMemoryFlagKeepMapped FdMemoryFlags = C.GST_FD_MEMORY_FLAG_KEEP_MAPPED // (1) – once the memory is mapped, keep it mapped until the memory is destroyed
	FdMemoryFlagMapPrivate FdMemoryFlags = C.GST_FD_MEMORY_FLAG_MAP_PRIVATE // (2) – do a private mapping instead of the default shared mapping.
	FdMemoryFlagDontClose  FdMemoryFlags = C.GST_FD_MEMORY_FLAG_DONT_CLOSE  // (4) – don't close the file descriptor when the memory is freed.
)
