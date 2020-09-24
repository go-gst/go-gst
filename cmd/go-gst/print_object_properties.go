package main

import (
	"fmt"

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
	colorCyan.print(s)
}

func printObjectPropertiesInfo(obj *gst.Object, description string) {
	colorOrange.printf("%s:\n", description)

	// for now this function only handles elements

	for _, param := range obj.ListProperties() {
		colorBlue.printfIndent(2, "%-20s", param.Name)
		colorReset.printf(": %s", param.Blurb)

		colorReset.print("\n")

		colorOrange.printIndent(24, "flags")
		colorReset.print(": ")
		colorCyan.print(param.Flags.GstFlagsString())

		colorReset.print("\n")

		switch param.ValueType {

		case glib.TYPE_STRING:
			printFieldType("String. ")
			printFieldName("Default")
			if param.DefaultValue == nil {
				printFieldValue("null")
			} else {
				str, _ := param.DefaultValue.GetString()
				printFieldValue(str)
			}

		case glib.TYPE_BOOLEAN:
			val, err := param.DefaultValue.GoValue()
			var valStr string
			if err != nil {
				valStr = "unknown" // edge case
			} else {
				b := val.(bool)
				valStr = fmt.Sprintf("%t", b)
			}
			printFieldType("Boolean. ")
			printFieldName("Default")
			printFieldValue(valStr)

		case glib.TYPE_ULONG:
			printFieldType("Unsigned Long. ")

		case glib.TYPE_LONG:
			printFieldType("Long. ")

		case glib.TYPE_UINT:
			printFieldType("Unsigned Integer. ")

		case glib.TYPE_INT:
			printFieldType("Integer. ")

		case glib.TYPE_UINT64:
			printFieldType("Unsigned Integer64. ")

		case glib.TYPE_INT64:
			printFieldType("Integer64. ")

		case glib.TYPE_FLOAT:
			printFieldType("Float. ")

		case glib.TYPE_DOUBLE:
			printFieldType("Double. ")

		default:

		}

		colorReset.print("\n")

	}

	fmt.Println()
}
