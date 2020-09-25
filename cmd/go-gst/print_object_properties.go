package main

import (
	"fmt"
	"sort"

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
			}
		}

		colorReset.print("\n")

	}

	fmt.Println()
}

func printCaps(caps gst.Caps, indent int) {
	if len(caps) == 0 {
		colorReset.print("\n")
		colorOrange.printIndent(indent, "ANY")
		return
	}
	for _, cap := range caps {
		colorReset.print("\n")
		colorOrange.printfIndent(indent, "%s", cap.Name)
		for k, v := range cap.Data {
			colorReset.print("\n")
			colorOrange.printfIndent(indent+2, "%s", k)
			colorReset.print(": ")
			colorLightGray.print(fmt.Sprint(v))
		}
	}
}

/*
static void
print_caps (const GstCaps * caps, const gchar * pfx)
{
  guint i;

  g_return_if_fail (caps != NULL);

  if (gst_caps_is_any (caps)) {
    n_print ("%s%sANY%s\n", CAPS_TYPE_COLOR, pfx, RESET_COLOR);
    return;
  }
  if (gst_caps_is_empty (caps)) {
    n_print ("%s%sEMPTY%s\n", CAPS_TYPE_COLOR, pfx, RESET_COLOR);
    return;
  }

  for (i = 0; i < gst_caps_get_size (caps); i++) {
    GstStructure *structure = gst_caps_get_structure (caps, i);
    GstCapsFeatures *features = gst_caps_get_features (caps, i);

    if (features && (gst_caps_features_is_any (features) ||
            !gst_caps_features_is_equal (features,
                GST_CAPS_FEATURES_MEMORY_SYSTEM_MEMORY))) {
      gchar *features_string = gst_caps_features_to_string (features);

      n_print ("%s%s%s%s(%s%s%s)\n", pfx, STRUCT_NAME_COLOR,
          gst_structure_get_name (structure), RESET_COLOR,
          CAPS_FEATURE_COLOR, features_string, RESET_COLOR);
      g_free (features_string);
    } else {
      n_print ("%s%s%s%s\n", pfx, STRUCT_NAME_COLOR,
          gst_structure_get_name (structure), RESET_COLOR);
    }
    gst_structure_foreach (structure, print_field, (gpointer) pfx);
  }
}

*/

// /* obj will be NULL if we're printing properties of pad template pads */
// static void
// print_object_properties_info (GObject * obj, GObjectClass * obj_class,
//     const gchar * desc)
// {
//   GParamSpec **property_specs;
//   guint num_properties, i;
//   gboolean readable;
//   gboolean first_flag;

//   property_specs = g_object_class_list_properties (obj_class, &num_properties);
//   g_qsort_with_data (property_specs, num_properties, sizeof (gpointer),
//       (GCompareDataFunc) sort_gparamspecs, NULL);

//   n_print ("%s%s%s:\n", HEADING_COLOR, desc, RESET_COLOR);

//   push_indent ();

//   for (i = 0; i < num_properties; i++) {
//     GValue value = { 0, };
//     GParamSpec *param = property_specs[i];
//     GType owner_type = param->owner_type;

//       default:
//         if (param->value_type == GST_TYPE_CAPS) {
//           const GstCaps *caps = gst_value_get_caps (&value);

//           if (!caps)
//             n_print ("%sCaps%s (NULL)", DATATYPE_COLOR, RESET_COLOR);
//           else {
//             print_caps (caps, "                           ");
//           }
//         } else if (G_IS_PARAM_SPEC_ENUM (param)) {
//           GEnumValue *values;
//           guint j = 0;
//           gint enum_value;
//           const gchar *value_nick = "";

//           values = G_ENUM_CLASS (g_type_class_ref (param->value_type))->values;
//           enum_value = g_value_get_enum (&value);

//           while (values[j].value_name) {
//             if (values[j].value == enum_value)
//               value_nick = values[j].value_nick;
//             j++;
//           }

//           n_print ("%sEnum \"%s\"%s %sDefault%s: %s%d, \"%s\"%s",
//               DATATYPE_COLOR, g_type_name (G_VALUE_TYPE (&value)), RESET_COLOR,
//               PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//               enum_value, value_nick, RESET_COLOR);

//           j = 0;
//           while (values[j].value_name) {
//             g_print ("\n");
//             n_print ("   %s(%d)%s: %s%-16s%s - %s%s%s",
//                 PROP_ATTR_NAME_COLOR, values[j].value, RESET_COLOR,
//                 PROP_ATTR_VALUE_COLOR, values[j].value_nick, RESET_COLOR,
//                 DESC_COLOR, values[j].value_name, RESET_COLOR);
//             j++;
//           }
//           /* g_type_class_unref (ec); */
//         } else if (G_IS_PARAM_SPEC_FLAGS (param)) {
//           GParamSpecFlags *pflags = G_PARAM_SPEC_FLAGS (param);
//           GFlagsValue *vals;
//           gchar *cur;

//           vals = pflags->flags_class->values;

//           cur = flags_to_string (vals, g_value_get_flags (&value));

//           n_print ("%sFlags \"%s\"%s %sDefault%s: %s0x%08x, \"%s\"%s",
//               DATATYPE_COLOR, g_type_name (G_VALUE_TYPE (&value)), RESET_COLOR,
//               PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//               g_value_get_flags (&value), cur, RESET_COLOR);

//           while (vals[0].value_name) {
//             g_print ("\n");
//             n_print ("   %s(0x%08x)%s: %s%-16s%s - %s%s%s",
//                 PROP_ATTR_NAME_COLOR, vals[0].value, RESET_COLOR,
//                 PROP_ATTR_VALUE_COLOR, vals[0].value_nick, RESET_COLOR,
//                 DESC_COLOR, vals[0].value_name, RESET_COLOR);
//             ++vals;
//           }

//           g_free (cur);
//         } else if (G_IS_PARAM_SPEC_OBJECT (param)) {
//           n_print ("%sObject of type%s %s\"%s\"%s", PROP_VALUE_COLOR,
//               RESET_COLOR, DATATYPE_COLOR,
//               g_type_name (param->value_type), RESET_COLOR);
//         } else if (G_IS_PARAM_SPEC_BOXED (param)) {
//           n_print ("%sBoxed pointer of type%s %s\"%s\"%s", PROP_VALUE_COLOR,
//               RESET_COLOR, DATATYPE_COLOR,
//               g_type_name (param->value_type), RESET_COLOR);
//           if (param->value_type == GST_TYPE_STRUCTURE) {
//             const GstStructure *s = gst_value_get_structure (&value);
//             if (s) {
//               g_print ("\n");
//               gst_structure_foreach (s, print_field,
//                   (gpointer) "                           ");
//             }
//           }
//         } else if (G_IS_PARAM_SPEC_POINTER (param)) {
//           if (param->value_type != G_TYPE_POINTER) {
//             n_print ("%sPointer of type%s %s\"%s\"%s.", PROP_VALUE_COLOR,
//                 RESET_COLOR, DATATYPE_COLOR, g_type_name (param->value_type),
//                 RESET_COLOR);
//           } else {
//             n_print ("%sPointer.%s", PROP_VALUE_COLOR, RESET_COLOR);
//           }
//         } else if (param->value_type == G_TYPE_VALUE_ARRAY) {
//           GParamSpecValueArray *pvarray = G_PARAM_SPEC_VALUE_ARRAY (param);

//           if (pvarray->element_spec) {
//             n_print ("%sArray of GValues of type%s %s\"%s\"%s",
//                 PROP_VALUE_COLOR, RESET_COLOR, DATATYPE_COLOR,
//                 g_type_name (pvarray->element_spec->value_type), RESET_COLOR);
//           } else {
//             n_print ("%sArray of GValues%s", PROP_VALUE_COLOR, RESET_COLOR);
//           }
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
//     }
//     if (!readable)
//       g_print (" %sWrite only%s\n", PROP_VALUE_COLOR, RESET_COLOR);
//     else
//       g_print ("\n");

//     pop_indent_n (11);

//     g_value_reset (&value);
//   }
//   if (num_properties == 0)
//     n_print ("%snone%s\n", PROP_VALUE_COLOR, RESET_COLOR);

//   pop_indent ();

//   g_free (property_specs);
// }
