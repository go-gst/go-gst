package gst

// #include "gst.go.h"
import "C"
import (
	"runtime"
	"time"
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
)

// Query is a go wrapper around a GstQuery.
type Query struct {
	ptr *C.GstQuery
}

// Type returns the type of the Query.
func (q *Query) Type() QueryType { return QueryType(q.ptr._type) }

// FromGstQueryUnsafeNone wraps the pointer to the given C GstQuery with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstQueryUnsafeNone(query unsafe.Pointer) *Query {
	q := ToGstQuery(query)
	q.Ref()
	runtime.SetFinalizer(q, (*Query).Unref)
	return q
}

// FromGstQueryUnsafeFull wraps the pointer to the given C GstQuery with the go type.
// This is meant for internal usage and is exported for visibility to other packages.
func FromGstQueryUnsafeFull(query unsafe.Pointer) *Query {
	q := ToGstQuery(query)
	runtime.SetFinalizer(q, (*Query).Unref)
	return q
}

// ToGstQuery converts the given pointer into a Message without affecting the ref count or
// placing finalizers.
func ToGstQuery(query unsafe.Pointer) *Query { return wrapQuery((*C.GstQuery)(query)) }

// NewAcceptCapsQuery constructs a new query object for querying if caps are accepted.
func NewAcceptCapsQuery(caps *Caps) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_accept_caps(caps.Instance())))
}

// NewAllocationQuery constructs a new query object for querying the allocation properties.
func NewAllocationQuery(caps *Caps, needPool bool) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_allocation(
		caps.Instance(), gboolean(needPool),
	)))
}

// NewBitrateQuery constructs a new query object for querying the bitrate.
func NewBitrateQuery() *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_bitrate()))
}

// NewBufferingQuery constructs a new query object for querying the buffering status of a stream.
func NewBufferingQuery(format Format) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_buffering(
		C.GstFormat(format),
	)))
}

// NewCapsQuery constructs a new query object for querying the caps.
//
// The CAPS query should return the allowable caps for a pad in the context of the element's state, its link to
// other elements, and the devices or files it has opened. These caps must be a subset of the pad template caps.
// In the NULL state with no links, the CAPS query should ideally return the same caps as the pad template. In
// rare circumstances, an object property can affect the caps returned by the CAPS query, but this is discouraged.
//
// For most filters, the caps returned by CAPS query is directly affected by the allowed caps on other pads. For
// demuxers and decoders, the caps returned by the srcpad's getcaps function is directly related to the stream data.
// Again, the CAPS query should return the most specific caps it reasonably can, since this helps with autoplugging.
//
// The filter is used to restrict the result caps, only the caps matching filter should be returned from the CAPS
// query. Specifying a filter might greatly reduce the amount of processing an element needs to do.
func NewCapsQuery(caps *Caps) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_caps(caps.Instance())))
}

// NewContextQuery constructs a new query object for querying the pipeline-local context.
func NewContextQuery(ctxType string) *Query {
	cName := C.CString(ctxType)
	defer C.free(unsafe.Pointer(cName))
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_context(
		(*C.gchar)(unsafe.Pointer(cName)),
	)))
}

// NewConvertQuery constructs a new convert query object. A convert query is used to ask for a conversion between one
// format and another.
func NewConvertQuery(srcFormat, destFormat Format, value int64) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_convert(
		C.GstFormat(srcFormat), C.gint64(value), C.GstFormat(destFormat),
	)))
}

// NewCustomQuery constructs a new custom query object.
func NewCustomQuery(queryType QueryType, structure *Structure) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_custom(
		C.GstQueryType(queryType),
		structure.Instance(),
	)))
}

// NewDrainQuery constructs a new query object for querying the drain state.
func NewDrainQuery() *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_drain()))
}

// NewDurationQuery constructs a new stream duration query object to query in the given format. A duration query will give the
// total length of the stream.
func NewDurationQuery(format Format) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_duration(C.GstFormat(format))))
}

// NewFormatsQuery constructs a new query object for querying formats of the stream.
func NewFormatsQuery() *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_formats()))
}

// NewLatencyQuery constructs a new latency query object. A latency query is usually performed by sinks to compensate for additional
// latency introduced by elements in the pipeline.
func NewLatencyQuery() *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_latency()))
}

// NewPositionQuery constructs a new query stream position query object. A position query is used to query the current position of playback
// in the streams, in some format.
func NewPositionQuery(format Format) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_position(C.GstFormat(format))))
}

// NewSchedulingQuery constructs a new query object for querying the scheduling properties.
func NewSchedulingQuery() *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_scheduling()))
}

// NewSeekingQuery constructs a new query object for querying seeking properties of the stream.
func NewSeekingQuery(format Format) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_seeking(C.GstFormat(format))))
}

// NewSegmentQuery constructs a new segment query object. A segment query is used to discover information about the currently configured segment
// for playback.
func NewSegmentQuery(format Format) *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_segment(C.GstFormat(format))))
}

// NewURIQuery constructs a new query URI query object. An URI query is used to query the current URI that is used by the source or sink.
func NewURIQuery() *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_new_uri()))
}

// Instance returns the underlying GstQuery instance.
func (q *Query) Instance() *C.GstQuery { return C.toGstQuery(unsafe.Pointer(q.ptr)) }

// AddAllocationMeta adds api with params as one of the supported metadata API to query.
func (q *Query) AddAllocationMeta(api glib.Type, structure *Structure) {
	C.gst_query_add_allocation_meta(q.Instance(), (C.GType)(api), structure.Instance())
}

// AddAllocationParam adds allocator and its params as a supported memory allocator.
func (q *Query) AddAllocationParam(allocator *Allocator, params *AllocationParams) {
	C.gst_query_add_allocation_param(q.Instance(), allocator.Instance(), params.Instance())
}

// AddAllocationPool sets the pool parameters in query.
func (q *Query) AddAllocationPool(pool *BufferPool, size, minBuffers, maxBuffers uint) {
	C.gst_query_add_allocation_pool(
		q.Instance(),
		pool.Instance(),
		C.guint(size),
		C.guint(minBuffers),
		C.guint(maxBuffers),
	)
}

// AddBufferingRange sets the buffering-ranges array field in query. The current last start position of the array should be inferior to start.
func (q *Query) AddBufferingRange(start, stop int64) (ok bool) {
	return gobool(C.gst_query_add_buffering_range(q.Instance(), C.gint64(start), C.gint64(stop)))
}

// AddSchedulingMode adds mode as one of the supported scheduling modes to query.
func (q *Query) AddSchedulingMode(mode PadMode) {
	C.gst_query_add_scheduling_mode(q.Instance(), C.GstPadMode(mode))
}

// Copy copies the given query using the copy function of the parent GstStructure.
func (q *Query) Copy() *Query {
	return FromGstQueryUnsafeFull(unsafe.Pointer(C.gst_query_copy(q.Instance())))
}

// FindAllocationMeta checks if query has metadata api set. When this function returns TRUE, index will contain the index where the requested
// API and the parameters can be found.
func (q *Query) FindAllocationMeta(api glib.Type) (ok bool, index uint) {
	var out C.guint
	gok := C.gst_query_find_allocation_meta(q.Instance(), C.GType(api), &out)
	return gobool(gok), uint(out)
}

// GetNumAllocationMetas retrieves the number of values currently stored in the meta API array of the query's structure.
func (q *Query) GetNumAllocationMetas() uint {
	return uint(C.gst_query_get_n_allocation_metas(q.Instance()))
}

// GetNumAllocationParams retrieves the number of values currently stored in the allocator params array of the query's structure.
//
// If no memory allocator is specified, the downstream element can handle the default memory allocator. The first memory allocator in the query
// should be generic and allow mapping to system memory, all following allocators should be ordered by preference with the preferred one first.
func (q *Query) GetNumAllocationParams() uint {
	return uint(C.gst_query_get_n_allocation_params(q.Instance()))
}

// GetNumAllocationPools retrieves the number of values currently stored in the pool array of the query's structure.
func (q *Query) GetNumAllocationPools() uint {
	return uint(C.gst_query_get_n_allocation_pools(q.Instance()))
}

// GetNumBufferingRanges retrieves the number of values currently stored in the buffered-ranges array of the query's structure.
func (q *Query) GetNumBufferingRanges() uint {
	return uint(C.gst_query_get_n_buffering_ranges(q.Instance()))
}

// GetNumSchedulingModes retrieves the number of values currently stored in the scheduling mode array of the query's structure.
func (q *Query) GetNumSchedulingModes() uint {
	return uint(C.gst_query_get_n_scheduling_modes(q.Instance()))
}

// GetStructure retrieves the structure of a query.
func (q *Query) GetStructure() *Structure {
	return wrapStructure(C.gst_query_get_structure(q.Instance()))
}

// HasSchedulingMode checks if query has scheduling mode set.
func (q *Query) HasSchedulingMode(mode PadMode) bool {
	return gobool(C.gst_query_has_scheduling_mode(q.Instance(), C.GstPadMode(mode)))
}

// HasSchedulingModeWithFlags checks if query has scheduling mode set and flags is set in query scheduling flags.
func (q *Query) HasSchedulingModeWithFlags(mode PadMode, flags SchedulingFlags) bool {
	return gobool(C.gst_query_has_scheduling_mode_with_flags(q.Instance(), C.GstPadMode(mode), C.GstSchedulingFlags(flags)))
}

// ParseAcceptCaps gets the caps from query. The caps remains valid as long as query remains valid.
func (q *Query) ParseAcceptCaps() *Caps {
	caps := (*C.GstCaps)(C.malloc(C.sizeof_GstCaps))
	C.gst_query_parse_accept_caps(q.Instance(), &caps)
	return FromGstCapsUnsafeNone(unsafe.Pointer(caps))
}

// ParseAcceptCapsResult parses the result from the caps query.
func (q *Query) ParseAcceptCapsResult() bool {
	var out C.gboolean
	C.gst_query_parse_accept_caps_result(q.Instance(), &out)
	return gobool(out)
}

// ParseAllocation parses an allocation query.
func (q *Query) ParseAllocation() (caps *Caps, needPool bool) {
	gcaps := (*C.GstCaps)(C.malloc(C.sizeof_GstCaps))
	var needz C.gboolean
	C.gst_query_parse_allocation(q.Instance(), &gcaps, &needz)
	return FromGstCapsUnsafeNone(unsafe.Pointer(gcaps)), gobool(needz)
}

// ParseBitrate gets the results of a bitrate query. See also SetBitrate.
func (q *Query) ParseBitrate() uint {
	var out C.guint
	C.gst_query_parse_bitrate(q.Instance(), &out)
	return uint(out)
}

// ParseBufferingPercent gets the percentage of buffered data. This is a value between 0 and 100. The busy indicator is TRUE when
// the buffering is in progress.
func (q *Query) ParseBufferingPercent() (busy bool, percent int) {
	var gb C.gboolean
	var gp C.gint
	C.gst_query_parse_buffering_percent(q.Instance(), &gb, &gp)
	return gobool(gb), int(gp)
}

// ParseBufferingRange parses a buffering range query.
func (q *Query) ParseBufferingRange() (format Format, start, stop, estimatedTotal int64) {
	var gformat C.GstFormat
	var gstart, gstop, gestimated C.gint64
	C.gst_query_parse_buffering_range(q.Instance(), &gformat, &gstart, &gstop, &gestimated)
	return Format(gformat), int64(gstart), int64(gstop), int64(gestimated)
}

// ParseBufferingStats extracts the buffering stats values from query.
func (q *Query) ParseBufferingStats() (mode BufferingMode, avgIn, avgOut int, bufLeft int64) {
	var gmode C.GstBufferingMode
	var avgi, avgo C.gint
	var gbufleft C.gint64
	C.gst_query_parse_buffering_stats(q.Instance(), &gmode, &avgi, &avgo, &gbufleft)
	return BufferingMode(gmode), int(avgi), int(avgo), int64(gbufleft)
}

// ParseCaps gets the filter from the caps query. The caps remains valid as long as query remains valid.
func (q *Query) ParseCaps() *Caps {
	caps := (*C.GstCaps)(C.malloc(C.sizeof_GstCaps))
	C.gst_query_parse_caps(q.Instance(), &caps)
	return FromGstCapsUnsafeNone(unsafe.Pointer(caps))
}

// ParseCapsResult gets the caps result from query. The caps remains valid as long as query remains valid.
func (q *Query) ParseCapsResult() *Caps {
	caps := (*C.GstCaps)(C.malloc(C.sizeof_GstCaps))
	C.gst_query_parse_caps_result(q.Instance(), &caps)
	return FromGstCapsUnsafeNone(unsafe.Pointer(caps))
}

// ParseContext gets the context from the context query. The context remains valid as long as query remains valid.
func (q *Query) ParseContext() *Context {
	var _ctx *C.GstContext
	ctx := C.makeContextWritable(_ctx)
	C.gst_query_parse_context(q.Instance(), &ctx)
	return FromGstContextUnsafeNone(unsafe.Pointer(ctx))
}

// ParseContextType parses a context type from an existing GST_QUERY_CONTEXT query.
func (q *Query) ParseContextType() (ok bool, ctxType string) {
	tPtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(tPtr))
	gok := C.gst_query_parse_context_type(q.Instance(), (**C.gchar)(unsafe.Pointer(&tPtr)))
	if gobool(gok) {
		return true, C.GoString((*C.char)(unsafe.Pointer(tPtr)))
	}
	return false, ""
}

// ParseConvert parses a convert query answer.
func (q *Query) ParseConvert() (srcformat, destFormat Format, srcValue, destValue int64) {
	var gsrcf, gdestf C.GstFormat
	var gsval, gdval C.gint64
	C.gst_query_parse_convert(q.Instance(), &gsrcf, &gsval, &gdestf, &gdval)
	return Format(gsrcf), Format(gdestf), int64(gsval), int64(gdval)
}

// ParseDuration parses a duration query answer.
func (q *Query) ParseDuration() (format Format, duration int64) {
	var gf C.GstFormat
	var gd C.gint64
	C.gst_query_parse_duration(q.Instance(), &gf, &gd)
	return Format(gf), int64(gd)
}

// ParseLatency parses a latency query answer.
func (q *Query) ParseLatency() (live bool, minLatency, maxLatency time.Duration) {
	var min, max C.GstClockTime
	var gl C.gboolean
	C.gst_query_parse_latency(q.Instance(), &gl, &min, &max)
	return gobool(gl), time.Duration(min), time.Duration(max)
}

// ParseNumFormats parses the number of formats in the formats query.
func (q *Query) ParseNumFormats() uint {
	var out C.guint
	C.gst_query_parse_n_formats(q.Instance(), &out)
	return uint(out)
}

// ParseAllocationMetaAt parses an available query and get the metadata API at index of the metadata API array.
func (q *Query) ParseAllocationMetaAt(idx uint) (api glib.Type, st *Structure) {
	var gs *C.GstStructure
	gtype := C.gst_query_parse_nth_allocation_meta(q.Instance(), C.guint(idx), &gs)
	return glib.Type(gtype), wrapStructure(gs)
}

// ParseAllocationParamAt parses an available query and get the allocator and its params at index of the allocator array.
func (q *Query) ParseAllocationParamAt(idx uint) (*Allocator, *AllocationParams) {
	var alloc *C.GstAllocator
	var params C.GstAllocationParams
	C.gst_query_parse_nth_allocation_param(q.Instance(), C.guint(idx), &alloc, &params)
	return FromGstAllocatorUnsafeFull(unsafe.Pointer(alloc)), wrapAllocationParams(&params)
}

// ParseAllocationPoolAt gets the pool parameters in query.
func (q *Query) ParseAllocationPoolAt(idx uint) (pool *BufferPool, size, minBuffers, maxBuffers uint) {
	var gpool *C.GstBufferPool
	var gs, gmin, gmax C.guint
	C.gst_query_parse_nth_allocation_pool(q.Instance(), C.guint(idx), &gpool, &gs, &gmin, &gmax)
	return FromGstBufferPoolUnsafeFull(unsafe.Pointer(gpool)), uint(gs), uint(gmin), uint(gmax)
}

// ParseBufferingRangeAt parses an available query and get the start and stop values stored at the index of the buffered ranges array.
func (q *Query) ParseBufferingRangeAt(idx uint) (start, stop int64) {
	var gstart, gstop C.gint64
	C.gst_query_parse_nth_buffering_range(q.Instance(), C.guint(idx), &gstart, &gstop)
	return int64(gstart), int64(gstop)
}

// ParseFormatAt parses the format query and retrieve the nth format from it into format. If the list contains less elements than nth,
// format will be set to GST_FORMAT_UNDEFINED.
func (q *Query) ParseFormatAt(idx uint) Format {
	var out C.GstFormat
	C.gst_query_parse_nth_format(q.Instance(), C.guint(idx), &out)
	return Format(out)
}

// ParseSchedulingModeAt parses an available query and get the scheduling mode at index of the scheduling modes array.
func (q *Query) ParseSchedulingModeAt(idx uint) PadMode {
	return PadMode(C.gst_query_parse_nth_scheduling_mode(q.Instance(), C.guint(idx)))
}

// ParsePosition parses a position query, writing the format into format, and the position into cur, if the respective parameters are non-%NULL.
func (q *Query) ParsePosition() (format Format, cur int64) {
	var gf C.GstFormat
	var out C.gint64
	C.gst_query_parse_position(q.Instance(), &gf, &out)
	return Format(gf), int64(out)
}

// ParseScheduling sets the scheduling properties.
func (q *Query) ParseScheduling() (flags SchedulingFlags, minSize, maxSize, align int) {
	var gf C.GstSchedulingFlags
	var gmin, gmax, galign C.gint
	C.gst_query_parse_scheduling(q.Instance(), &gf, &gmin, &gmax, &galign)
	return SchedulingFlags(gf), int(gmin), int(gmax), int(galign)
}

// ParseSeeking parses a seeking query.
func (q *Query) ParseSeeking() (format Format, seekable bool, start, end int64) {
	var gs, ge C.gint64
	var seek C.gboolean
	var f C.GstFormat
	C.gst_query_parse_seeking(q.Instance(), &f, &seek, &gs, &ge)
	return Format(f), gobool(seek), int64(gs), int64(ge)
}

// ParseSegment parses a segment query answer.
func (q *Query) ParseSegment() (rate float64, format Format, start, stop int64) {
	var gs, ge C.gint64
	var f C.GstFormat
	var grate C.gdouble
	C.gst_query_parse_segment(q.Instance(), &grate, &f, &gs, &ge)
	return float64(grate), Format(f), int64(gs), int64(ge)
}

// ParseURI parses a URI query.
func (q *Query) ParseURI() string {
	tPtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(tPtr))
	C.gst_query_parse_uri(q.Instance(), (**C.gchar)(unsafe.Pointer(&tPtr)))
	return C.GoString((*C.char)(unsafe.Pointer(tPtr)))
}

// ParseURIRedirection parses a URI query.
func (q *Query) ParseURIRedirection() string {
	tPtr := C.malloc(C.sizeof_char * 1024)
	defer C.free(unsafe.Pointer(tPtr))
	C.gst_query_parse_uri_redirection(q.Instance(), (**C.gchar)(unsafe.Pointer(&tPtr)))
	return C.GoString((*C.char)(unsafe.Pointer(tPtr)))
}

// ParseURIRedirectionPermanent parses an URI query, and set permanent to TRUE if there is a redirection and it should be considered permanent.
// If a redirection is permanent, applications should update their internal storage of the URI, otherwise they should make all future requests
// to the original URI.
func (q *Query) ParseURIRedirectionPermanent() bool {
	var out C.gboolean
	C.gst_query_parse_uri_redirection_permanent(q.Instance(), &out)
	return gobool(out)
}

// RemoveAllocationMetaAt removes the metadata API at index of the metadata API array.
func (q *Query) RemoveAllocationMetaAt(idx uint) {
	C.gst_query_remove_nth_allocation_meta(q.Instance(), C.guint(idx))
}

// RemoveAllocationParamAt removes the allocation param at index of the allocation param array.
func (q *Query) RemoveAllocationParamAt(idx uint) {
	C.gst_query_remove_nth_allocation_param(q.Instance(), C.guint(idx))
}

// RemoveAllocationPoolAt removes the allocation pool at index of the allocation pool array.
func (q *Query) RemoveAllocationPoolAt(idx uint) {
	C.gst_query_remove_nth_allocation_pool(q.Instance(), C.guint(idx))
}

// SetAcceptCapsResult sets result as the result for the query.
func (q *Query) SetAcceptCapsResult(result bool) {
	C.gst_query_set_accept_caps_result(q.Instance(), gboolean(result))
}

// SetBitrate sets the results of a bitrate query. The nominal bitrate is the average bitrate expected over the length of the stream as advertised
// in file headers (or similar).
func (q *Query) SetBitrate(nominal uint) {
	C.gst_query_set_bitrate(q.Instance(), C.guint(nominal))
}

// SetBufferingPercent sets the percentage of buffered data. This is a value between 0 and 100. The busy indicator is TRUE when the buffering is
// in progress.
func (q *Query) SetBufferingPercent(busy bool, percent int) {
	C.gst_query_set_buffering_percent(q.Instance(), gboolean(busy), C.gint(percent))
}

// SetBufferingRange sets the available query result fields in query.
func (q *Query) SetBufferingRange(format Format, start, stop, estimatedTotal int64) {
	C.gst_query_set_buffering_range(q.Instance(), C.GstFormat(format), C.gint64(start), C.gint64(stop), C.gint64(estimatedTotal))
}

// SetBufferingStats configures the buffering stats values in query.
func (q *Query) SetBufferingStats(mode BufferingMode, avgIn, avgOut int, bufferingLeft int64) {
	C.gst_query_set_buffering_stats(q.Instance(), C.GstBufferingMode(mode), C.gint(avgIn), C.gint(avgOut), C.gint64(bufferingLeft))
}

// SetCapsResult sets the caps result in query.
func (q *Query) SetCapsResult(caps *Caps) {
	C.gst_query_set_caps_result(q.Instance(), caps.Instance())
}

// SetContext answers a context query by setting the requested context.
func (q *Query) SetContext(ctx *Context) {
	C.gst_query_set_context(q.Instance(), ctx.Instance())
}

// SetConvert answers a convert query by setting the requested values.
func (q *Query) SetConvert(srcFormat, destFormat Format, srcValue, destValue int64) {
	C.gst_query_set_convert(q.Instance(), C.GstFormat(srcFormat), C.gint64(srcValue), C.GstFormat(destFormat), C.gint64(destValue))
}

// SetDuration answers a duration query by setting the requested value in the given format.
func (q *Query) SetDuration(format Format, duration int64) {
	C.gst_query_set_duration(q.Instance(), C.GstFormat(format), C.gint64(duration))
}

// SetFormats sets the formats query result fields in query. The number of formats passed must be equal to n_formats.
func (q *Query) SetFormats(formats ...Format) {
	gstFormats := make([]C.GstFormat, len(formats))
	for _, f := range formats {
		gstFormats = append(gstFormats, C.GstFormat(f))
	}
	C.gst_query_set_formatsv(q.Instance(), C.gint(len(formats)), (*C.GstFormat)(unsafe.Pointer(&gstFormats[0])))
}

// SetLatency answers a latency query by setting the requested values in the given format.
func (q *Query) SetLatency(live bool, minLatency, maxLatency time.Duration) {
	C.gst_query_set_latency(q.Instance(), gboolean(live), C.guint64(minLatency.Nanoseconds()), C.guint64(maxLatency.Nanoseconds()))
}

// SetAllocationParamAt sets allocation params in query.
func (q *Query) SetAllocationParamAt(idx uint, allocator *Allocator, params *AllocationParams) {
	C.gst_query_set_nth_allocation_param(q.Instance(), C.guint(idx), allocator.Instance(), params.Instance())
}

// SetAllocationPoolAt sets the pool parameters in query.
func (q *Query) SetAllocationPoolAt(idx uint, pool *BufferPool, size, minBuffers, maxBuffers uint) {
	C.gst_query_set_nth_allocation_pool(q.Instance(), C.guint(idx), pool.Instance(), C.guint(size), C.guint(minBuffers), C.guint(maxBuffers))
}

// SetPosition answers a position query by setting the requested value in the given format.
func (q *Query) SetPosition(format Format, cur int64) {
	C.gst_query_set_position(q.Instance(), C.GstFormat(format), C.gint64(cur))
}

// SetScheduling sets the scheduling properties.
func (q *Query) SetScheduling(flags SchedulingFlags, minSize, maxSize, align int) {
	C.gst_query_set_scheduling(q.Instance(), C.GstSchedulingFlags(flags), C.gint(minSize), C.gint(maxSize), C.gint(align))
}

// SetSeeking sets the seeking query result fields in query.
func (q *Query) SetSeeking(format Format, seekable bool, segmentStart, segmentEnd int64) {
	C.gst_query_set_seeking(q.Instance(), C.GstFormat(format), gboolean(seekable), C.gint64(segmentStart), C.gint64(segmentEnd))
}

// SetSegment answers a segment query by setting the requested values. The normal playback segment of a pipeline is 0 to duration at the default rate of 1.0.
// If a seek was performed on the pipeline to play a different segment, this query will return the range specified in the last seek.
//
// start_value and stop_value will respectively contain the configured playback range start and stop values expressed in format. The values are always between
// 0 and the duration of the media and start_value <= stop_value. rate will contain the playback rate. For negative rates, playback will actually happen from
// stop_value to start_value.
func (q *Query) SetSegment(rate float64, format Format, startValue, stopValue int64) {
	C.gst_query_set_segment(q.Instance(), C.gdouble(rate), C.GstFormat(format), C.gint64(startValue), C.gint64(stopValue))
}

// SetURI answers a URI query by setting the requested URI.
func (q *Query) SetURI(uri string) {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	C.gst_query_set_uri(q.Instance(), (*C.gchar)(unsafe.Pointer(curi)))
}

// SetURIRedirection answers a URI query by setting the requested URI redirection.
func (q *Query) SetURIRedirection(uri string) {
	curi := C.CString(uri)
	defer C.free(unsafe.Pointer(curi))
	C.gst_query_set_uri_redirection(q.Instance(), (*C.gchar)(unsafe.Pointer(curi)))
}

// SetURIRedirectionPermanent answers a URI query by setting the requested URI redirection to permanent or not.
func (q *Query) SetURIRedirectionPermanent(permanent bool) {
	C.gst_query_set_uri_redirection_permanent(q.Instance(), gboolean(permanent))
}

// Ref increases the query ref count by one.
func (q *Query) Ref() *Query {
	C.gst_query_ref(q.Instance())
	return q
}

// Unref decreases the refcount of the query. If the refcount reaches 0, the query will be freed.
func (q *Query) Unref() {
	C.gst_query_unref(q.Instance())
}

// WritableStructure gets the structure of a query. This method should be called with a writable query so that the returned structure is guaranteed to be writable.
func (q *Query) WritableStructure() *Structure {
	return wrapStructure(C.gst_query_writable_structure(q.Instance()))
}
