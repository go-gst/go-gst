package gst

/*
#include "gst.go.h"
*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// BufferPool is a go wrapper around a GstBufferPool.
//
// For more information refer to the official documentation:
// https://gstreamer.freedesktop.org/documentation/gstreamer/gstbufferpool.html?gi-language=c
type BufferPool struct{ *Object }

// NewBufferPool returns a new BufferPool instance.
func NewBufferPool() *BufferPool {
	pool := C.gst_buffer_pool_new()
	return FromGstBufferPoolUnsafeFull(unsafe.Pointer(pool))
}

// FromGstBufferPoolUnsafeNone wraps the given unsafe.Pointer in a BufferPool instance. It takes a
// ref and places a runtime finalizer on the resulting object.
func FromGstBufferPoolUnsafeNone(bufferPool unsafe.Pointer) *BufferPool {
	pool := wrapBufferPool(glib.TransferNone(bufferPool))
	return pool
}

// // FromGstBufferPoolUnsafe is an alias to FromGstBufferPoolUnsafeNone.
// func FromGstBufferPoolUnsafe(bufferPool unsafe.Pointer) *BufferPool {
// 	return FromGstBufferPoolUnsafeNone(bufferPool)
// }

// FromGstBufferPoolUnsafeFull wraps the given unsafe.Pointer in a BufferPool instance. It just
// places a runtime finalizer on the resulting object.
func FromGstBufferPoolUnsafeFull(bufferPool unsafe.Pointer) *BufferPool {
	pool := wrapBufferPool(glib.TransferFull(bufferPool))
	return pool
}

// Instance returns the underlying GstBufferPool instance.
func (b *BufferPool) Instance() *C.GstBufferPool { return C.toGstBufferPool(b.Unsafe()) }

// IsFlushing returns true if this BufferPool is currently flushing.
func (b *BufferPool) IsFlushing() bool { return gobool(C.bufferPoolIsFlushing(b.Instance())) }

// BufferPoolAcquireParams represents parameters to an AcquireBuffer call.
type BufferPoolAcquireParams struct {
	Format Format                 // format (GstFormat) – the format of start and stop
	Start  int64                  // start (gint64) – the start position
	Stop   int64                  // stop (gint64) – the stop position
	Flags  BufferPoolAcquireFlags // flags (GstBufferPoolAcquireFlags) – additional flags
}

// AcquireBuffer acquires a buffer from this pool.
func (b *BufferPool) AcquireBuffer(params *BufferPoolAcquireParams) (*Buffer, FlowReturn) {
	var buf *C.GstBuffer
	if params != nil {
		gparams := (*C.GstBufferPoolAcquireParams)(C.malloc(C.sizeof_GstBufferPoolAcquireParams))
		defer C.free(unsafe.Pointer(gparams))
		gparams.format = C.GstFormat(params.Format)
		gparams.start = C.gint64(params.Start)
		gparams.stop = C.gint64(params.Stop)
		gparams.flags = C.GstBufferPoolAcquireFlags(params.Flags)
		ret := C.gst_buffer_pool_acquire_buffer(b.Instance(), &buf, gparams)
		if FlowReturn(ret) != FlowOK {
			return nil, FlowReturn(ret)
		}
		return FromGstBufferUnsafeFull(unsafe.Pointer(buf)), FlowReturn(ret)
	}
	ret := C.gst_buffer_pool_acquire_buffer(b.Instance(), &buf, nil)
	if FlowReturn(ret) != FlowOK {
		return nil, FlowReturn(ret)
	}
	return FromGstBufferUnsafeFull(unsafe.Pointer(buf)), FlowReturn(ret)
}

// GetConfig retrieves a copy of the current configuration of the pool. This configuration can either
// be modified and used for the SetConfig call or it must be freed after usage with Free.
func (b *BufferPool) GetConfig() *BufferPoolConfig {
	st := C.gst_buffer_pool_get_config(b.Instance())
	if st == nil {
		return nil
	}
	return &BufferPoolConfig{Structure: FromGstStructureUnsafe(unsafe.Pointer(st))}
}

// GetOptions retrieves a list of supported bufferpool options for the pool. An option would typically
// be enabled with AddOption.
func (b *BufferPool) GetOptions() []string {
	opts := C.gst_buffer_pool_get_options(b.Instance())
	defer C.g_free((C.gpointer)(unsafe.Pointer(opts)))
	return goStrings(C.sizeOfGCharArray(opts), opts)
}

// HasOption returns true if this BufferPool supports the given option.
func (b *BufferPool) HasOption(opt string) bool {
	gOpt := C.CString(opt)
	defer C.free(unsafe.Pointer(gOpt))
	return gobool(C.gst_buffer_pool_has_option(b.Instance(), (*C.gchar)(gOpt)))
}

// IsActive returns true if this BufferPool is active. A pool can be activated with SetActive.
func (b *BufferPool) IsActive() bool { return gobool(C.gst_buffer_pool_is_active(b.Instance())) }

// ReleaseBuffer releases the given buffer from the pool. The buffer should have previously been
// allocated from pool with AcquireBiffer.
//
// This function is usually called automatically when the last ref on buffer disappears.
func (b *BufferPool) ReleaseBuffer(buf *Buffer) {
	C.gst_buffer_pool_release_buffer(b.Instance(), buf.Ref().Instance())
}

// SetActive can be used to control the active state of pool. When the pool is inactive, new calls to
// AcquireBuffer will return with FlowFlushing.
//
// Activating the bufferpool will preallocate all resources in the pool based on the configuration of the pool.
//
// Deactivating will free the resources again when there are no outstanding buffers. When there are outstanding buffers, they will be freed as soon as they are all returned to the pool.
func (b *BufferPool) SetActive(active bool) (ok bool) {
	return gobool(C.gst_buffer_pool_set_active(b.Instance(), gboolean(active)))
}

// SetConfig sets the configuration of the pool. If the pool is already configured, and the configurations
// haven't changed, this function will return TRUE. If the pool is active, this method will return FALSE and
// active configurations will remain. Buffers allocated form this pool must be returned or else this function
// will do nothing and return FALSE.
//
// config is a GstStructure that contains the configuration parameters for the pool. A default and mandatory set
// of parameters can be configured with gst_buffer_pool_config_set_params, gst_buffer_pool_config_set_allocator
// and gst_buffer_pool_config_add_option.
//
// If the parameters in config can not be set exactly, this function returns FALSE and will try to update as much
// state as possible. The new state can then be retrieved and refined with GetConfig.
//
// This function takes ownership of the given structure.
func (b *BufferPool) SetConfig(cfg *BufferPoolConfig) bool {
	return gobool(C.gst_buffer_pool_set_config(b.Instance(), cfg.Instance()))
}

// SetFlushing enables or disable the flushing state of a pool without freeing or allocating buffers.
func (b *BufferPool) SetFlushing(flushing bool) {
	C.gst_buffer_pool_set_flushing(b.Instance(), gboolean(flushing))
}

// BufferPoolConfig wraps the Structure interface with extra methods for interacting with BufferPool
// configurations.
type BufferPoolConfig struct{ *Structure }

// AddOption enables the option in config. This will instruct the bufferpool to enable the specified option
// on the buffers that it allocates.
//
// The supported options by pool can be retrieved with GetOptions.
func (b *BufferPoolConfig) AddOption(opt string) {
	cOpt := C.CString(opt)
	defer C.free(unsafe.Pointer(cOpt))
	C.gst_buffer_pool_config_add_option(b.Instance(), (*C.gchar)(unsafe.Pointer(cOpt)))
}

// GetAllocator retrieves the allocator and params from config.
func (b *BufferPoolConfig) GetAllocator() (*Allocator, *AllocationParams) {
	var allocator *C.GstAllocator
	var params C.GstAllocationParams
	C.gst_buffer_pool_config_get_allocator(b.Instance(), &allocator, &params)
	var allo *Allocator
	if allocator != nil {
		allo = wrapAllocator(glib.TransferNone(unsafe.Pointer(allocator)))
	}
	return allo, &AllocationParams{ptr: &params}
}

// GetOption retrieves the option at index of the options API array.
func (b *BufferPoolConfig) GetOption(index uint) string {
	return C.GoString(C.gst_buffer_pool_config_get_option(b.Instance(), C.guint(index)))
}

// GetParams retrieves the values from this config. All params return 0 or nil if they could not be fetched.
func (b *BufferPoolConfig) GetParams() (caps *Caps, size, minBuffers, maxBuffers uint) {
	var gcaps *C.GstCaps
	var gsize, gminBuffers, gmaxBuffers C.guint
	ret := gobool(C.gst_buffer_pool_config_get_params(
		b.Instance(),
		&gcaps,
		&gsize,
		&gminBuffers,
		&gmaxBuffers,
	))
	if ret {
		return wrapCaps(gcaps), uint(gsize), uint(gminBuffers), uint(gmaxBuffers)
	}
	return nil, 0, 0, 0
}

// HasOption returns true if this config has the given option.
func (b *BufferPoolConfig) HasOption(opt string) bool {
	cOpt := C.CString(opt)
	defer C.free(unsafe.Pointer(cOpt))
	return gobool(C.gst_buffer_pool_config_has_option(b.Instance(), (*C.gchar)(unsafe.Pointer(cOpt))))
}

// NumOptions retrieves the number of values currently stored in the options array of the config structure.
func (b *BufferPoolConfig) NumOptions() uint {
	return uint(C.gst_buffer_pool_config_n_options(b.Instance()))
}

// SetAllocator sets the allocator and params on config.
//
// One of allocator and params can be nil, but not both. When allocator is nil, the default allocator of
// the pool will use the values in param to perform its allocation. When param is nil, the pool will use
// the provided allocator with its default AllocationParams.
//
// A call to SetConfig on the BufferPool can update the allocator and params with the values that it is able to do.
// Some pools are, for example, not able to operate with different allocators or cannot allocate with the values
// specified in params. Use GetConfig on the pool to get the currently used values.
func (b *BufferPoolConfig) SetAllocator(allocator *Allocator, params *AllocationParams) {
	if allocator == nil && params != nil {
		C.gst_buffer_pool_config_set_allocator(b.Instance(), nil, params.Instance())
		return
	}
	if allocator != nil && params == nil {
		C.gst_buffer_pool_config_set_allocator(b.Instance(), allocator.Instance(), nil)
		return
	}
	if allocator != nil && params != nil {
		C.gst_buffer_pool_config_set_allocator(b.Instance(), allocator.Instance(), params.Instance())
	}
}

// SetParams configures the config with the given parameters.
func (b *BufferPoolConfig) SetParams(caps *Caps, size, minBuffers, maxBuffers uint) {
	if caps == nil {
		C.gst_buffer_pool_config_set_params(
			b.Instance(),
			nil,
			C.guint(size), C.guint(minBuffers), C.guint(maxBuffers),
		)
		return
	}
	C.gst_buffer_pool_config_set_params(
		b.Instance(),
		caps.Instance(),
		C.guint(size), C.guint(minBuffers), C.guint(maxBuffers),
	)
}

// Validate that changes made to config are still valid in the context of the expected parameters. This function is a
// helper that can be used to validate changes made by a pool to a config when SetConfig returns FALSE. This expects
// that caps haven't changed and that min_buffers aren't lower then what we initially expected. This does not check if
// options or allocator parameters are still valid, and won't check if size have changed, since changing the size is valid
// to adapt padding.
func (b *BufferPoolConfig) Validate(caps *Caps, size, minBuffers, maxBuffers uint) bool {
	return gobool(C.gst_buffer_pool_config_validate_params(
		b.Instance(),
		caps.Instance(),
		C.guint(size), C.guint(minBuffers), C.guint(maxBuffers),
	))
}
