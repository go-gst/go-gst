package gst

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// MarshalStructure will convert the given go struct into a GstStructure. Currently nested
// structs are not supported.
func MarshalStructure(data any) (*Structure, error) {
	typeOf := reflect.TypeOf(data)
	valsOf := reflect.ValueOf(data)
	st := NewStructureEmpty(marshalStructureName(typeOf, valsOf))

	err := marshalInto(typeOf, valsOf, st)

	if err != nil {
		return nil, err
	}

	return st, nil
}

func marshalStructureName(typeOf reflect.Type, valsOf reflect.Value) string {
	// TODO: allow the user to implement an interface that returns the structure name
	_ = valsOf
	return typeOf.Name()
}

func marshalInto(typeOf reflect.Type, valsOf reflect.Value, st *Structure) error {
	for i := 0; i < valsOf.NumField(); i++ {
		field := typeOf.Field(i)
		fieldVal := valsOf.Field(i)

		if !field.IsExported() {
			// skip private fields when marshaling
			continue
		}

		fieldName := marshalInfoFromField(field).name

		if field.Type.Kind() == reflect.Struct {
			var fieldTargetStructure *Structure
			if field.Anonymous {
				// embedded field, marshal into current structure
				fieldTargetStructure = st
			} else {
				// Struct field with struct type, create a recursive gst.Structure
				fieldTargetStructure = NewStructureEmpty(marshalStructureName(typeOf, valsOf))
			}

			err := marshalInto(field.Type, fieldVal, fieldTargetStructure)

			if err != nil {
				return fmt.Errorf("cannot marshal field %s: %w", field.Name, err)
			}

			if !field.Anonymous {
				// set the value after filling it
				st.SetValue(fieldName, fieldTargetStructure)
			}

			continue
		}

		if !supportedStructureMarshalPrimitive(field) {
			return fmt.Errorf("cannot marshal field to gst.Structure: %s, unsupported type: %s", field.Name, field.Type.String())
		}

		gval := fieldVal.Interface()

		st.SetValue(fieldName, gval)
	}

	return nil
}

// UnmarshalInto will unmarshal this structure into the given pointer. The object
// reflected by the pointer must be non-nil.
func (s *Structure) UnmarshalInto(data any) error {
	valsOf := reflect.ValueOf(data)
	if valsOf.Kind() != reflect.Pointer || valsOf.IsNil() {
		return errors.New("data is invalid (nil or non-pointer)")
	}

	typeOf := reflect.TypeOf(data).Elem()
	valsOf = valsOf.Elem()

	if valsOf.Kind() != reflect.Struct {
		return errors.New("cannot unmarshal into data: data is not pointer to struct")
	}

	return unmarshalInto(typeOf, valsOf, s)
}

func unmarshalInto(typeOf reflect.Type, valsOf reflect.Value, s *Structure) error {
	for i := 0; i < valsOf.NumField(); i++ {
		field := typeOf.Field(i)
		fieldVal := valsOf.Field(i)

		if !field.IsExported() {
			// private field
			continue
		}

		fieldName := marshalInfoFromField(field).name

		if field.Type.Kind() == reflect.Struct {
			var fieldTargetStructure *Structure
			if field.Anonymous {
				// embedded field, unmarshal into current structure
				fieldTargetStructure = s
			} else {
				// Struct field with struct type, expect a structre key that returns a structure

				substructure, ok := s.GetValue(fieldName).(*Structure)

				if !ok {
					continue
				}

				fieldTargetStructure = substructure
			}

			err := unmarshalInto(field.Type, fieldVal, fieldTargetStructure)

			if err != nil {
				return fmt.Errorf("error unmarshaling struct field %s: %w", field.Name, err)
			}

			continue
		}

		val := s.GetValue(fieldName)

		if val == nil {
			// leave the field as a zero value
			continue
		}

		rv := reflect.ValueOf(val)

		if !rv.CanConvert(field.Type) {
			return fmt.Errorf("cannot convert value %#v of type %T to %s", val, val, field.Type.String())
		}

		fieldVal.Set(rv.Convert(field.Type))
	}

	return nil
}

func supportedStructureMarshalPrimitive(field reflect.StructField) bool {
	switch field.Type.Kind() {
	case reflect.Int, reflect.Uint:
		// must use concrete bit size
		return false

	case reflect.Bool:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return true
	case reflect.String:
		return true
	case reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return true
	default:
		return false
	}
}

func marshalInfoFromField(field reflect.StructField) gstMarshalFieldInfo {
	fieldTag, ok := field.Tag.Lookup("gst")

	parsed := parseStructureTags(fieldTag)

	if !ok {
		parsed.name = field.Name
	}

	return parsed
}

type gstMarshalFieldInfo struct {
	name string
	kv   map[string]string
}

// parseStructureTags parses a struct tag in the form of:
//
//	foobar,key=value,key2,key3=value3
//
// where the first part is the name and the others are key value pairs (with optional values)
func parseStructureTags(tags string) gstMarshalFieldInfo {
	if tags == "" {
		return gstMarshalFieldInfo{}
	}

	parts := strings.Split(tags, ",")

	parsed := gstMarshalFieldInfo{
		name: parts[0],
	}

	if len(parts) > 1 {
		parsed.kv = make(map[string]string)
		for _, kv := range parts[1:] {
			if kv == "" {
				continue
			}
			parts := strings.Split(kv, "=")

			key := parts[0]
			var value string
			if len(parts) > 1 {
				value = parts[1]
			}
			parsed.kv[key] = value
		}
	}

	return parsed
}
