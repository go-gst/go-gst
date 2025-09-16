package gst

import (
	"errors"
	"reflect"
)

// MarshalStructure will convert the given go struct into a GstStructure. Currently nested
// structs are not supported.
func MarshalStructure(data interface{}) *Structure {
	typeOf := reflect.TypeOf(data)
	valsOf := reflect.ValueOf(data)
	st := NewStructureEmpty(typeOf.Name())
	for i := 0; i < valsOf.NumField(); i++ {
		gval := valsOf.Field(i).Interface()

		// TODO: if the value is a struct then recursively MarshalStructure

		fieldName := typeOf.Field(i).Name
		st.SetValue(fieldName, gval)
	}
	return st
}

// UnmarshalInto will unmarshal this structure into the given pointer. The object
// reflected by the pointer must be non-nil.
func (s *Structure) UnmarshalInto(data any) error {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("data is invalid (nil or non-pointer)")
	}

	val := reflect.ValueOf(data).Elem()
	nVal := rv.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)

		fieldName, ok := val.Type().Field(i).Tag.Lookup("gst")

		if !ok {
			fieldName = val.Type().Field(i).Name
		}

		val := s.GetValue(fieldName)

		// TODO: if val is a structure do a recursive UnmarshalInto

		nvField.Set(reflect.ValueOf(val))
	}

	return nil
}
