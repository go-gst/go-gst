package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/gotk3/gotk3/glib"
	"github.com/tinyzimmer/go-gst-launch/gst"
)

func printFieldType(s string) {
	colorGreen.printIndent(24, s)
}

func printFieldName(s string) {
	colorOrange.print(s)
	colorReset.print(": ")
}

func printFieldValue(s string) {
	colorCyan.printf("%s ", s)
}

// ByName implements sort. Interface based on the Name field.
type ByName []*gst.ParameterSpec

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ByValue implements sort. Interface based on the Value field.
type ByValue []*gst.FlagsValue

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func printObjectPropertiesInfo(obj *gst.Object, description string) {
	colorOrange.printf("%s:\n", description)

	// for now this function only handles elements

	props := obj.ListProperties()
	sort.Sort(ByName(props))

	for _, param := range props {
		defer param.Unref()

		colorBlue.printfIndent(2, "%-20s", param.Name)
		colorReset.printf(": %s", param.Blurb)

		colorReset.print("\n")

		colorOrange.printIndent(24, "flags")
		colorReset.print(": ")
		colorCyan.print(param.Flags.GstFlagsString())

		if !param.Flags.Has(gst.ParameterReadable) {
			colorReset.print(" | ")
			colorLightPurple.print("Write only")
		} else if !param.Flags.Has(gst.ParameterWritable) {
			colorReset.print(" | ")
			colorLightPurple.print("Read only")
		}

		colorReset.print("\n")

		// get the value and continue on any error (very unlikely)

		goval, _ := param.DefaultValue.GoValue()

		// skips deprecated types

		switch param.ValueType {

		case glib.TYPE_STRING:
			printFieldType("String. ")
			printFieldName("Default")
			if goval == nil {
				printFieldValue("null")
			} else {
				str, _ := param.DefaultValue.GetString()
				printFieldValue(str)
			}

		case glib.TYPE_BOOLEAN:
			var valStr string
			if goval == nil {
				valStr = "unknown" // edge case
			} else {
				b := goval.(bool)
				valStr = fmt.Sprintf("%t", b)
			}
			printFieldType("Boolean. ")
			printFieldName("Default")
			printFieldValue(valStr)

		case glib.TYPE_UINT:
			var valStr string
			if goval == nil {
				valStr = "0"
			} else {
				v := goval.(uint)
				valStr = fmt.Sprintf("%d", v)
			}
			printFieldType("Unsigned Integer. ")
			printFieldName("Range")
			min, max := param.UIntRange()
			printFieldValue(fmt.Sprintf("%d - %d ", min, max))
			printFieldName("Default")
			printFieldValue(valStr)

		case glib.TYPE_INT:
			var valStr string
			if goval == nil {
				valStr = "0"
			} else {
				v := goval.(int)
				valStr = fmt.Sprintf("%d", v)
			}
			printFieldType("Integer. ")
			printFieldName("Range")
			min, max := param.IntRange()
			printFieldValue(fmt.Sprintf("%d - %d ", min, max))
			printFieldName("Default")
			printFieldValue(valStr)

		case glib.TYPE_UINT64:
			var valStr string
			if goval == nil {
				valStr = "0"
			} else {
				v := goval.(uint64)
				valStr = fmt.Sprintf("%d", v)
			}
			printFieldType("Unsigned Integer64. ")
			printFieldName("Range")
			min, max := param.UInt64Range()
			printFieldValue(fmt.Sprintf("%d - %d ", min, max))
			printFieldName("Default")
			printFieldValue(valStr)

		case glib.TYPE_INT64:
			var valStr string
			if goval == nil {
				valStr = "0"
			} else {
				v := goval.(int64)
				valStr = fmt.Sprintf("%d", v)
			}
			printFieldType("Integer64. ")
			printFieldName("Range")
			min, max := param.Int64Range()
			printFieldValue(fmt.Sprintf("%d - %d ", min, max))
			printFieldName("Default")
			printFieldValue(valStr)

		case glib.TYPE_FLOAT:
			var valStr string
			if goval == nil {
				valStr = "0"
			} else {
				v := goval.(float64)
				valStr = fmt.Sprintf("%15.7g", v)
			}
			printFieldType("Float. ")
			printFieldName("Range")
			min, max := param.FloatRange()
			printFieldValue(fmt.Sprintf("%15.7g - %15.7g ", min, max))
			printFieldName("Default")
			printFieldValue(valStr)

		case glib.TYPE_DOUBLE:
			var valStr string
			if goval == nil {
				valStr = "0"
			} else {
				v := goval.(float64)
				valStr = fmt.Sprintf("%15.7g", v)
			}
			printFieldType("Double. ")
			printFieldName("Range")
			min, max := param.DoubleRange()
			printFieldValue(fmt.Sprintf("%15.7g - %15.7g ", min, max))
			printFieldName("Default")
			printFieldValue(valStr)

		default:
			if param.IsCaps() {

				if caps := param.GetCaps(); caps != nil {
					printCaps(caps, 24)
				}

			} else if param.IsEnum() {

				enumValues := param.GetEnumValues()
				iface, _ := param.DefaultValue.GoValue()
				var curVal string
				if iface == nil {
					curVal = "-1"
				} else {
					curVal = fmt.Sprintf("%v", iface)
				}
				var defaultStr string
				for _, val := range enumValues {
					if curVal == strconv.Itoa(val.Value) {
						defaultStr = val.ValueNick
					}
				}
				printFieldType(fmt.Sprintf(`Enum "%s" `, param.ValueType.Name()))
				printFieldName("Default")
				printFieldValue(fmt.Sprintf(`%s "%s"`, curVal, defaultStr))
				colorReset.print("\n")
				for idx, val := range enumValues {
					w := new(tabwriter.Writer)
					buf := new(bytes.Buffer)
					w.Init(buf, 100, 73, 0, '\t', 0)
					colorOrange.fprintfIndent(w, 27, "(%d)", val.Value)
					colorReset.fprint(w, ": ")
					colorCyan.fprint(w, val.ValueNick)
					colorReset.fprint(w, "\t - ")
					colorLightGray.fprint(w, val.ValueName)
					w.Flush()
					fmt.Print(buf.String())
					if idx < len(enumValues)-1 {
						colorReset.print("\n")
					}
				}

			} else if param.IsFlags() {

				flags := param.GetFlagValues()
				sort.Sort(ByValue(flags))
				flagStr := "+"
				for _, flag := range flags {
					flagStr += fmt.Sprintf(" %s", flag.ValueNick)
				}
				if flagStr == "+" {
					flagStr = "(none)"
				}
				printFieldType(fmt.Sprintf(`Flags "%s" `, param.ValueType.Name()))
				printFieldName("Default")
				printFieldValue(fmt.Sprintf(`0x%08x "%s"`, param.GetDefaultFlags(), flagStr))

				for idx, flag := range flags {
					w := new(tabwriter.Writer)
					buf := new(bytes.Buffer)
					w.Init(buf, 100, 73, 0, '\t', 0)
					colorOrange.fprintfIndent(w, 27, "(%d)", flag.Value)
					colorReset.fprint(w, ": ")
					colorCyan.fprint(w, flag.ValueNick)
					colorReset.fprint(w, "\t - ")
					colorLightGray.fprint(w, flag.ValueName)
					w.Flush()
					fmt.Print(buf.String())
					if idx < len(flags)-1 {
						colorReset.print("\n")
					}
				}

			} else if param.IsObject() {

				colorLightGray.printIndent(24, "Object of type ")
				colorGreen.printf(`"%s"`, param.ValueType.Name())

			} else if param.IsBoxed() {

				colorLightGray.printIndent(24, "Boxed pointer of type ")
				colorGreen.printf(`"%s"`, param.ValueType.Name())
				if param.ValueType.Name() == "GstStructure" {
					structure := gst.StructureFromGValue(param.DefaultValue)
					if structure != nil {
						for key, val := range structure.Values() {
							colorReset.printIndent(26, "(gpointer)       ")
							colorYellow.printf("%15s:", key)
							colorBlue.printf("%v", val)
						}
					}
				}

			} else if param.IsPointer() {

				colorLightGray.printIndent(24, "Pointer of type ")
				colorGreen.printf(`"%s`, param.ValueType.Name())

			} else if param.IsFraction() {

				colorGreen.printIndent(24, "Fraction.")

			} else if param.IsGstArray() {

				colorGreen.printIndent(24, "GstArray.")

			} else {
				colorReset.printIndent(24, "Unknown type ")
				colorGreen.printf(`"%s`, param.ValueType.Name())
			}
		}

		colorReset.print("\n")

	}

	fmt.Println()
}

func printCaps(caps *gst.Caps, indent int) {
	if caps == nil {
		return
	}

	colorReset.print("\n")
	defer func() { colorReset.print("\n") }()

	if caps.IsAny() {
		colorOrange.printIndent(indent, "ANY")
		return
	}

	if caps.IsEmpty() {
		colorOrange.printIndent(indent, "EMPTY")
		return
	}

	for i := 0; i < caps.Size(); i++ {
		structure := caps.GetStructureAt(i)
		features := caps.GetFeaturesAt(i)

		if features != nil && features.IsAny() {
			colorOrange.printIndent(indent+20, structure.Name())
			colorLightGray.print(" (")
			colorGreen.print(features.String())
			colorLightGray.print(")")
		} else {
			colorOrange.printIndent(indent+20, structure.Name())
		}

		colorReset.print("\n")
		for k, v := range structure.Values() {
			colorCyan.printIndent(indent+10, k)
			colorReset.print(": ")
			colorBlue.printf("%v\n", v)
		}
	}
}

//         } else if (GST_IS_PARAM_SPEC_FRACTION (param)) {
//           GstParamSpecFraction *pfraction = GST_PARAM_SPEC_FRACTION (param);

//           n_print ("%sFraction%s. %sRange%s: %s%d/%d - %d/%d%s "
//               "%sDefault%s: %s%d/%d%s ", DATATYPE_COLOR, RESET_COLOR,
//               PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//               pfraction->min_num, pfraction->min_den, pfraction->max_num,
//               pfraction->max_den, RESET_COLOR, PROP_ATTR_NAME_COLOR,
//               RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//               gst_value_get_fraction_numerator (&value),
//               gst_value_get_fraction_denominator (&value), RESET_COLOR);
//         } else if (param->value_type == GST_TYPE_ARRAY) {
//           GstParamSpecArray *parray = GST_PARAM_SPEC_ARRAY_LIST (param);

//           if (parray->element_spec) {
//             n_print ("%sGstValueArray of GValues of type%s %s\"%s\"%s",
//                 PROP_VALUE_COLOR, RESET_COLOR, DATATYPE_COLOR,
//                 g_type_name (parray->element_spec->value_type), RESET_COLOR);
//           } else {
//             n_print ("%sGstValueArray of GValues%s", PROP_VALUE_COLOR,
//                 RESET_COLOR);
//           }
//         } else {
//           n_print ("%sUnknown type %ld%s %s\"%s\"%s", PROP_VALUE_COLOR,
//               (glong) param->value_type, RESET_COLOR, DATATYPE_COLOR,
//               g_type_name (param->value_type), RESET_COLOR);
//         }
//         break;
