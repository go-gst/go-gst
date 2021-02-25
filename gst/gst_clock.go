package gst

/*
#include "gst.go.h"

extern gboolean goClockCb                 (GstClock * clock, GstClockTime time, GstClockID id, gpointer user_data);
extern void     goGDestroyNotifyFuncNoRun (gpointer user_data);

gboolean cgoClockCb (GstClock * clock, GstClockTime time, GstClockID id, gpointer user_data)
{
	return goClockCb(clock, time, id, user_data);
}

void clockDestroyNotify (gpointer user_data)
{
	goGDestroyNotifyFuncNoRun(user_data);
}

*/
import "C"

import (
	"runtime"
	"time"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// ClockCallback is the prototype of a clock callback function.
type ClockCallback func(clock *Clock, clockTime time.Duration) bool

// ClockID is a go wrapper around a GstClockID.
type ClockID struct {
	ptr C.GstClockID // which is actually just a casted pointer
}

// Instance returns the underlying pointer.
func (c *ClockID) Instance() C.GstClockID { return c.ptr }

// GetClock returns the clock for this ClockID.
func (c *ClockID) GetClock() *Clock {
	clk := C.gst_clock_id_get_clock(c.Instance())
	return FromGstClockUnsafeFull(unsafe.Pointer(clk))
}

// GetTime returns the time for this ClockID
func (c *ClockID) GetTime() time.Duration {
	return time.Duration(C.gst_clock_id_get_time(c.Instance()))
}

// Unschedule cancels an outstanding request with id. This can either be an outstanding async notification or a pending sync notification.
// After this call, id cannot be used anymore to receive sync or async notifications, you need to create a new GstClockID.
func (c *ClockID) Unschedule() {
	C.gst_clock_id_unschedule(c.Instance())
}

// UsesClock returns whether id uses clock as the underlying clock. clock can be nil, in which case the return value indicates whether the
// underlying clock has been freed. If this is the case, the id is no longer usable and should be freed.
func (c *ClockID) UsesClock(clock *Clock) bool {
	return gobool(C.gst_clock_id_uses_clock(c.Instance(), clock.Instance()))
}

// Wait performs a blocking wait on id. id should have been created with NewSingleShotID or NewPeriodicID and should not have been unscheduled
// with a call to Unschedule.
//
// If the jitter argument is not 0 and this function returns ClockOK or ClockEarly, it will contain the difference against the clock and the
// time of id when this method was called. Positive values indicate how late id was relative to the clock (in which case this function will
// return ClockEarly). Negative values indicate how much time was spent waiting on the clock before this function returned.
func (c *ClockID) Wait() (ret ClockReturn, jitter ClockTimeDiff) {
	var gjitter C.GstClockTimeDiff
	ret = ClockReturn(C.gst_clock_id_wait(c.Instance(), &gjitter))
	jitter = ClockTimeDiff(gjitter)
	return
}

// WaitAsync registers a callback on the given ClockID id with the given function and user_data. When passing a ClockID with an invalid time to
// this function, the callback will be called immediately with a time set to ClockTimeNone. The callback will be called when the time of id has been reached.
//
// The callback func can be invoked from any thread, either provided by the core or from a streaming thread. The application should be prepared for this.
//
//   // Example
//
//   pipeline, _ := gst.NewPipelineFromString("fakesrc ! fakesink")
//   defer pipeline.Unref()
//
//   clock := pipeline.GetPipelineClock()
//
//   id := clock.NewSingleShotID(gst.ClockTime(1000000000)) // 1 second
//
//   id.WaitAsync(func(clock *gst.Clock, clockTime time.Duration) bool {
//       fmt.Println("Single shot triggered at", clockTime.Nanoseconds())
//       pipeline.SetState(gst.StateNull)
//       return true
//   })
//
//   pipeline.SetState(gst.StatePlaying)
//   gst.Wait(pipeline)
//
//   // Single shot triggered at 1000000000
func (c *ClockID) WaitAsync(f ClockCallback) ClockReturn {
	ptr := gopointer.Save(f)
	return ClockReturn(C.gst_clock_id_wait_async(
		c.Instance(),
		C.GstClockCallback(C.cgoClockCb),
		(C.gpointer)(unsafe.Pointer(ptr)),
		C.GDestroyNotify(C.clockDestroyNotify),
	))
}

// Ref increaes the ref count on ClockID.
func (c *ClockID) Ref() *ClockID {
	C.gst_clock_id_ref(c.Instance())
	return c
}

// Unref unrefs a ClockID.
func (c *ClockID) Unref() {
	C.gst_clock_id_unref(c.Instance())
}

// Clock is a go wrapper around a GstClock.
type Clock struct{ *Object }

// FromGstClockUnsafeNone takes a pointer to a GstClock and wraps it in a Clock instance.
// A ref is taken on the clock and a finalizer applied.
func FromGstClockUnsafeNone(clock unsafe.Pointer) *Clock {
	return wrapClock(glib.TransferNone(clock))
}

// FromGstClockUnsafeFull takes a pointer to a GstClock and wraps it in a Clock instance.
// A finalizer is set on the returned object.
func FromGstClockUnsafeFull(clock unsafe.Pointer) *Clock {
	return wrapClock(glib.TransferFull(clock))
}

// Instance returns the underlying GstClock instance.
func (c *Clock) Instance() *C.GstClock { return C.toGstClock(c.Unsafe()) }

// AddObservation adds the time master of the master clock and the time slave of the slave clock to the list of observations.
// If enough observations are available, a linear regression algorithm is run on the observations and clock is recalibrated.
//
// If this functions returns TRUE, the float will contain the correlation coefficient of the interpolation. A value of 1.0 means
// a perfect regression was performed. This value can be used to control the sampling frequency of the master and slave clocks.
func (c *Clock) AddObservation(slaveTime, masterTime time.Duration) (bool, float64) {
	var out C.gdouble
	ok := gobool(C.gst_clock_add_observation(
		c.Instance(),
		C.GstClockTime(slaveTime),
		C.GstClockTime(masterTime),
		&out,
	))
	return ok, float64(out)
}

// AddObservationUnapplied adds a clock observation to the internal slaving algorithm the same as AddObservation, and returns the
// result of the master clock estimation, without updating the internal calibration.
//
// The caller can then take the results and call SetCalibration with the values, or some modified version of them.
func (c *Clock) AddObservationUnapplied(slaveTime, masterTime time.Duration) (ok bool, rSquared float64, internalTime, externalTime, rateNum, rateDenom time.Duration) {
	var ginternal, gexternal, grateNum, grateDenom C.GstClockTime
	var grSquared C.gdouble
	ok = gobool(C.gst_clock_add_observation_unapplied(
		c.Instance(),
		C.GstClockTime(slaveTime),
		C.GstClockTime(masterTime),
		&grSquared, &ginternal, &gexternal, &grateNum, &grateDenom,
	))
	return ok, float64(grSquared), time.Duration(ginternal), time.Duration(gexternal), time.Duration(grateNum), time.Duration(grateDenom)
}

// AdjustUnlocked converts the given internal clock time to the external time, adjusting for the rate and reference time set with
// SetCalibration and making sure that the returned time is increasing. This function should be called with the clock's OBJECT_LOCK
// held and is mainly used by clock subclasses.
//
// This function is the reverse of UnadjustUnlocked.
func (c *Clock) AdjustUnlocked(internal time.Duration) time.Duration {
	return time.Duration(C.gst_clock_adjust_unlocked(c.Instance(), C.GstClockTime(internal)))
}

// AdjustWithCalibration converts the given internal_target clock time to the external time, using the passed calibration parameters.
// This function performs the same calculation as AdjustUnlocked when called using the current calibration parameters, but
// doesn't ensure a monotonically increasing result as AdjustUnlocked does.
//
// See: https://gstreamer.freedesktop.org/documentation/gstreamer/gstclock.html#gst_clock_adjust_with_calibration
func (c *Clock) AdjustWithCalibration(internalTarget, cinternal, cexternal, cnum, cdenom time.Duration) time.Duration {
	return time.Duration(C.gst_clock_adjust_with_calibration(
		c.Instance(),
		C.GstClockTime(internalTarget),
		C.GstClockTime(cinternal),
		C.GstClockTime(cexternal),
		C.GstClockTime(cnum),
		C.GstClockTime(cdenom),
	))
}

// GetCalibration gets the internal rate and reference time of clock. See gst_clock_set_calibration for more information.
func (c *Clock) GetCalibration() (internal, external, rateNum, rateDenom time.Duration) {
	var ginternal, gexternal, grateNum, grateDenom C.GstClockTime
	C.gst_clock_get_calibration(c.Instance(), &ginternal, &gexternal, &grateNum, &grateDenom)
	return time.Duration(ginternal), time.Duration(gexternal), time.Duration(grateNum), time.Duration(grateDenom)
}

// GetTime gets the current time of the given clock in nanoseconds or -1 if invalid.
// The time is always monotonically increasing and adjusted according to the current offset and rate.
func (c *Clock) GetTime() time.Duration {
	res := C.gst_clock_get_time(c.Instance())
	if res == gstClockTimeNone {
		return ClockTimeNone
	}
	return time.Duration(res)
}

// GetInternalTime gets the current internal time of the given clock in nanoseconds
// or ClockTimeNone if invalid. The time is returned unadjusted for the offset and the rate.
func (c *Clock) GetInternalTime() time.Duration {
	res := C.gst_clock_get_internal_time(c.Instance())
	if res == gstClockTimeNone {
		return ClockTimeNone
	}
	return time.Duration(res)
}

// GetMaster returns the master clock that this clock is slaved to or nil when the clock
// is not slaved to any master clock.
func (c *Clock) GetMaster() *Clock {
	clock := C.gst_clock_get_master(c.Instance())
	if clock == nil {
		return nil
	}
	return FromGstClockUnsafeFull(unsafe.Pointer(clock))
}

// GetResolution gets the accuracy of the clock. The accuracy of the clock is the granularity
// of the values returned by GetTime.
func (c *Clock) GetResolution() time.Duration {
	return time.Duration(C.gst_clock_get_resolution(c.Instance()))
}

// GetTimeout gets the amount of time that master and slave clocks are sampled.
func (c *Clock) GetTimeout() time.Duration {
	return time.Duration(C.gst_clock_get_timeout(c.Instance()))
}

// IsSynced returns true if the clock is synced.
func (c *Clock) IsSynced() bool { return gobool(C.gst_clock_is_synced(c.Instance())) }

// NewPeriodicID gets an ID from clock to trigger a periodic notification. The periodic notifications
// will start at time start_time and will then be fired with the given interval. ID should be unreffed after usage.
func (c *Clock) NewPeriodicID(startTime, interval time.Duration) *ClockID {
	id := C.gst_clock_new_periodic_id(
		c.Instance(),
		C.GstClockTime(startTime),
		C.GstClockTime(interval),
	)
	clkid := &ClockID{id}
	runtime.SetFinalizer(clkid, (*ClockID).Unref)
	return clkid
}

// NewSingleShotID gets a ClockID from the clock to trigger a single shot notification at the requested time.
// The single shot id should be unreffed after usage.
func (c *Clock) NewSingleShotID(at time.Duration) *ClockID {
	id := C.gst_clock_new_single_shot_id(
		c.Instance(),
		C.GstClockTime(at),
	)
	clkid := &ClockID{id}
	runtime.SetFinalizer(clkid, (*ClockID).Unref)
	return clkid
}

// PeriodicIDReinit reinitializes the provided periodic id to the provided start time and interval. Does not
/// modify the reference count.
func (c *Clock) PeriodicIDReinit(clockID *ClockID, startTime, interval time.Duration) bool {
	return gobool(C.gst_clock_periodic_id_reinit(
		c.Instance(),
		clockID.Instance(),
		C.GstClockTime(startTime),
		C.GstClockTime(interval),
	))
}

// SetCalibration adjusts the rate and time of clock.
// See: https://gstreamer.freedesktop.org/documentation/gstreamer/gstclock.html#gst_clock_set_calibration.
func (c *Clock) SetCalibration(internal, external, rateNum, rateDenom time.Duration) {
	C.gst_clock_set_calibration(
		c.Instance(),
		C.GstClockTime(internal),
		C.GstClockTime(external),
		C.GstClockTime(rateNum),
		C.GstClockTime(rateDenom),
	)
}

// SetMaster sets master as the master clock for clock. clock will be automatically calibrated so that
// GetTime reports the same time as the master clock.
//
// A clock provider that slaves its clock to a master can get the current calibration values with GetCalibration.
//
// Master can be nil in which case clock will not be slaved anymore. It will however keep reporting its time
// adjusted with the last configured rate and time offsets.
func (c *Clock) SetMaster(master *Clock) bool {
	if master == nil {
		return gobool(C.gst_clock_set_master(c.Instance(), nil))
	}
	return gobool(C.gst_clock_set_master(c.Instance(), master.Instance()))
}

// SetResolution sets the accuracy of the clock. Some clocks have the possibility to operate with different accuracy
// at the expense of more resource usage. There is normally no need to change the default resolution of a clock.
// The resolution of a clock can only be changed if the clock has the ClockFlagCanSetResolution flag set.
func (c *Clock) SetResolution(resolution time.Duration) time.Duration {
	return time.Duration(C.gst_clock_set_resolution(c.Instance(), C.GstClockTime(resolution)))
}

// SetSynced sets clock to synced and emits the GstClock::synced signal, and wakes up any thread waiting in WaitForSync.
//
// This function must only be called if ClockFlagNeedsStartupSync is set on the clock, and is intended to be called by
// subclasses only.
func (c *Clock) SetSynced(synced bool) { C.gst_clock_set_synced(c.Instance(), gboolean(synced)) }

// SetTimeout sets the amount of time, in nanoseconds, to sample master and slave clocks
func (c *Clock) SetTimeout(timeout time.Duration) {
	C.gst_clock_set_timeout(c.Instance(), C.GstClockTime(timeout))
}

// SingleShotIDReinit reinitializes the provided single shot id to the provided time. Does not modify the reference count.
func (c *Clock) SingleShotIDReinit(clockID *ClockID, at time.Duration) bool {
	return gobool(C.gst_clock_single_shot_id_reinit(c.Instance(), clockID.Instance(), C.GstClockTime(at)))
}

// UnadjustUnlocked converts the given external clock time to the internal time of clock, using the rate and reference time
// set with SetCalibration. This function should be called with the clock's OBJECT_LOCK held and is mainly used by clock subclasses.
//
// This function is the reverse of AdjustUnlocked.
func (c *Clock) UnadjustUnlocked(external time.Duration) time.Duration {
	return time.Duration(C.gst_clock_unadjust_unlocked(c.Instance(), C.GstClockTime(external)))
}

// UnadjustWithCalibration converts the given external_target clock time to the internal time, using the passed calibration parameters.
// This function performs the same calculation as UnadjustUnlocked when called using the current calibration parameters.
func (c *Clock) UnadjustWithCalibration(externalTarget, cinternal, cexternal, cnum, cdenom time.Duration) time.Duration {
	return time.Duration(C.gst_clock_unadjust_with_calibration(
		c.Instance(),
		C.GstClockTime(externalTarget),
		C.GstClockTime(cinternal),
		C.GstClockTime(cexternal),
		C.GstClockTime(cnum),
		C.GstClockTime(cdenom),
	))
}

// WaitForSync waits until clock is synced for reporting the current time. If timeout is ClockTimeNone it will wait forever, otherwise it
// will time out after timeout nanoseconds.
//
// For asynchronous waiting, the GstClock::synced signal can be used.
//
// This returns immediately with TRUE if ClockFlagNeedsStartupSync is not set on the clock, or if the clock is already synced.
func (c *Clock) WaitForSync(timeout time.Duration) bool {
	return gobool(C.gst_clock_wait_for_sync(c.Instance(), C.GstClockTime(timeout)))
}

// String returns the string representation of this clock value.
func (c *Clock) String() string { return c.GetTime().String() }

// InternalString returns the string representation of this clock's internal value.
func (c *Clock) InternalString() string { return c.GetInternalTime().String() }
