package gst

type ObjectExtManual interface {
	// GetValue wraps gst_object_get_value
	GetValue(propertyName string, timestamp ClockTime) any
}

// GetValue wraps gst_object_get_value
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
