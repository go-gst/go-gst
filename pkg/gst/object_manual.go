package gst

type ObjectExtManual interface {
	// GetValue wraps gst_object_get_value
	//
	// The function takes the following parameters:
	//
	// 	- propertyName string: the name of the property to get
	// 	- timestamp ClockTime: the time the control-change should be read from
	//
	// The function returns the following values:
	//
	// 	- goret any
	//
	// Gets the value for the given controlled property at the requested time.
	GetValue(propertyName string, timestamp ClockTime) any
}

// GetValue wraps gst_object_get_value
//
// The function takes the following parameters:
//
//   - propertyName string: the name of the property to get
//   - timestamp ClockTime: the time the control-change should be read from
//
// The function returns the following values:
//
//   - goret *gobject.Value
//
// Gets the value for the given controlled property at the requested time.
func (object *ObjectInstance) GetValue(propertyName string, timestamp ClockTime) any {
	panic("not implemented yet")
	// var carg0 *C.GstObject   // in, none, converted
	// var carg1 *C.gchar       // in, none, string, casted *C.gchar
	// var carg2 C.GstClockTime // in, none, casted, alias
	// var cret  *C.GValue      // return, full, converted

	// carg0 = (*C.GstObject)(UnsafeObjectToGlibNone(object))
	// carg1 = (*C.gchar)(unsafe.Pointer(C.CString(propertyName)))
	// defer C.free(unsafe.Pointer(carg1))
	// carg2 = C.GstClockTime(timestamp)

	// cret = C.gst_object_get_value(carg0, carg1, carg2)
	// runtime.KeepAlive(object)
	// runtime.KeepAlive(propertyName)
	// runtime.KeepAlive(timestamp)

	// var goret *gobject.Value

	// goret = gobject.ValueFromNativeOwned(unsafe.Pointer(cret))

	// return goret
}
