#include <gst/gst.h>
#include <gst/app/gstappsink.h>
#include <gst/app/gstappsrc.h>

/*
	Utilitits
*/
static GObjectClass *
getGObjectClass(void * p) {
	return G_OBJECT_GET_CLASS (p);
}

static int 
sizeOfGCharArray(gchar ** arr) {
	int i;
	for (i = 0 ; 1 ; i = i + 1) {
		if (arr[i] == NULL) { return i; };
	}
}

static gboolean
gstObjectFlagIsSet(GstObject * obj, GstElementFlags flags)
{
	return GST_OBJECT_FLAG_IS_SET (obj, flags);
}

static gboolean
gstElementIsURIHandler(GstElement * elem)
{
	return GST_IS_URI_HANDLER (elem);
}

/*
	Type Castings
*/

static GstUri *
toGstURI(void *p)
{
	return (GST_URI(p));
}

static GstURIHandler *
toGstURIHandler(void *p)
{
	return (GST_URI_HANDLER(p));
}

static GstRegistry *
toGstRegistry(void *p)
{
	return (GST_REGISTRY(p));
}

static GstPlugin *
toGstPlugin(void *p)
{
	return (GST_PLUGIN(p));
}

static GstPluginFeature *
toGstPluginFeature(void *p)
{
	return (GST_PLUGIN_FEATURE(p));
}

static GstObject *
toGstObject(void *p)
{
	return (GST_OBJECT(p));
}

static GstElementFactory *
toGstElementFactory(void *p)
{
	return (GST_ELEMENT_FACTORY(p));
}

static GstElement *
toGstElement(void *p)
{
	return (GST_ELEMENT(p));
}

static GstAppSink *
toGstAppSink(void *p)
{
	return (GST_APP_SINK(p));
}

static GstAppSrc *
toGstAppSrc(void *p)
{
	return (GST_APP_SRC(p));
}

static GstBin *
toGstBin(void *p)
{
	return (GST_BIN(p));
}

static GstBus *
toGstBus(void *p)
{
	return (GST_BUS(p));
}

static GstMessage *
toGstMessage(void *p)
{
	return (GST_MESSAGE(p));
}

static GstPipeline *
toGstPipeline(void *p)
{
	return (GST_PIPELINE(p));
}

static GstPad *
toGstPad(void *p)
{
	return (GST_PAD(p));
}

static GstPadTemplate *
toGstPadTemplate(void *p)
{
	return (GST_PAD_TEMPLATE(p));
}

static GstStructure *
toGstStructure(void *p)
{
	return (GST_STRUCTURE(p));
}

static GstClock *
toGstClock(void *p)
{
	return (GST_CLOCK(p));
}

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

//     /* We're printing pad properties */
//     if (obj == NULL && (owner_type == G_TYPE_OBJECT
//             || owner_type == GST_TYPE_OBJECT || owner_type == GST_TYPE_PAD))
//       continue;

//     g_value_init (&value, param->value_type);

//     n_print ("%s%-20s%s: %s%s%s\n", PROP_NAME_COLOR,
//         g_param_spec_get_name (param), RESET_COLOR, PROP_VALUE_COLOR,
//         g_param_spec_get_blurb (param), RESET_COLOR);

//     push_indent_n (11);

//     first_flag = TRUE;
//     n_print ("%sflags%s: ", PROP_ATTR_NAME_COLOR, RESET_COLOR);
//     readable = ! !(param->flags & G_PARAM_READABLE);
//     if (readable && obj != NULL) {
//       g_object_get_property (obj, param->name, &value);
//     } else {
//       /* if we can't read the property value, assume it's set to the default
//        * (which might not be entirely true for sub-classes, but that's an
//        * unlikely corner-case anyway) */
//       g_param_value_set_default (param, &value);
//     }
//     if (readable) {
//       g_print ("%s%s%s%s", (first_flag) ? "" : ", ", PROP_ATTR_VALUE_COLOR,
//           _("readable"), RESET_COLOR);
//       first_flag = FALSE;
//     }
//     if (param->flags & G_PARAM_WRITABLE) {
//       g_print ("%s%s%s%s", (first_flag) ? "" : ", ", PROP_ATTR_VALUE_COLOR,
//           _("writable"), RESET_COLOR);
//       first_flag = FALSE;
//     }
//     if (param->flags & G_PARAM_DEPRECATED) {
//       g_print ("%s%s%s%s", (first_flag) ? "" : ", ", PROP_ATTR_VALUE_COLOR,
//           _("deprecated"), RESET_COLOR);
//       first_flag = FALSE;
//     }
//     if (param->flags & GST_PARAM_CONTROLLABLE) {
//       g_print (", %s%s%s", PROP_ATTR_VALUE_COLOR, _("controllable"),
//           RESET_COLOR);
//       first_flag = FALSE;
//     }
//     if (param->flags & GST_PARAM_CONDITIONALLY_AVAILABLE) {
//       g_print (", %s%s%s", PROP_ATTR_VALUE_COLOR, _("conditionally available"),
//           RESET_COLOR);
//       first_flag = FALSE;
//     }
//     if (param->flags & GST_PARAM_MUTABLE_PLAYING) {
//       g_print (", %s%s%s", PROP_ATTR_VALUE_COLOR,
//           _("changeable in NULL, READY, PAUSED or PLAYING state"), RESET_COLOR);
//     } else if (param->flags & GST_PARAM_MUTABLE_PAUSED) {
//       g_print (", %s%s%s", PROP_ATTR_VALUE_COLOR,
//           _("changeable only in NULL, READY or PAUSED state"), RESET_COLOR);
//     } else if (param->flags & GST_PARAM_MUTABLE_READY) {
//       g_print (", %s%s%s", PROP_ATTR_VALUE_COLOR,
//           _("changeable only in NULL or READY state"), RESET_COLOR);
//     }
//     if (param->flags & ~KNOWN_PARAM_FLAGS) {
//       g_print ("%s0x%s%0x%s", (first_flag) ? "" : ", ", PROP_ATTR_VALUE_COLOR,
//           param->flags & ~KNOWN_PARAM_FLAGS, RESET_COLOR);
//     }
//     g_print ("\n");

//     switch (G_VALUE_TYPE (&value)) {
//       case G_TYPE_STRING:
//       {
//         const char *string_val = g_value_get_string (&value);

//         n_print ("%sString%s. ", DATATYPE_COLOR, RESET_COLOR);

//         if (string_val == NULL)
//           g_print ("%sDefault%s: %snull%s", PROP_ATTR_NAME_COLOR, RESET_COLOR,
//               PROP_ATTR_VALUE_COLOR, RESET_COLOR);
//         else
//           g_print ("%sDefault%s: %s\"%s\"%s", PROP_ATTR_NAME_COLOR, RESET_COLOR,
//               PROP_ATTR_VALUE_COLOR, string_val, RESET_COLOR);
//         break;
//       }
//       case G_TYPE_BOOLEAN:
//       {
//         gboolean bool_val = g_value_get_boolean (&value);

//         n_print ("%sBoolean%s. %sDefault%s: %s%s%s", DATATYPE_COLOR,
//             RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, bool_val ? "true" : "false", RESET_COLOR);
//         break;
//       }
//       case G_TYPE_ULONG:
//       {
//         GParamSpecULong *pulong = G_PARAM_SPEC_ULONG (param);

//         n_print
//             ("%sUnsigned Long%s. %sRange%s: %s%lu - %lu%s %sDefault%s: %s%lu%s ",
//             DATATYPE_COLOR, RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, pulong->minimum, pulong->maximum,
//             RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, g_value_get_ulong (&value), RESET_COLOR);

//         GST_ERROR ("%s: property '%s' of type ulong: consider changing to "
//             "uint/uint64", G_OBJECT_CLASS_NAME (obj_class),
//             g_param_spec_get_name (param));
//         break;
//       }
//       case G_TYPE_LONG:
//       {
//         GParamSpecLong *plong = G_PARAM_SPEC_LONG (param);

//         n_print ("%sLong%s. %sRange%s: %s%ld - %ld%s %sDefault%s: %s%ld%s ",
//             DATATYPE_COLOR, RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, plong->minimum, plong->maximum, RESET_COLOR,
//             PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//             g_value_get_long (&value), RESET_COLOR);

//         GST_ERROR ("%s: property '%s' of type long: consider changing to "
//             "int/int64", G_OBJECT_CLASS_NAME (obj_class),
//             g_param_spec_get_name (param));
//         break;
//       }
//       case G_TYPE_UINT:
//       {
//         GParamSpecUInt *puint = G_PARAM_SPEC_UINT (param);

//         n_print
//             ("%sUnsigned Integer%s. %sRange%s: %s%u - %u%s %sDefault%s: %s%u%s ",
//             DATATYPE_COLOR, RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, puint->minimum, puint->maximum, RESET_COLOR,
//             PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//             g_value_get_uint (&value), RESET_COLOR);
//         break;
//       }
//       case G_TYPE_INT:
//       {
//         GParamSpecInt *pint = G_PARAM_SPEC_INT (param);

//         n_print ("%sInteger%s. %sRange%s: %s%d - %d%s %sDefault%s: %s%d%s ",
//             DATATYPE_COLOR, RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, pint->minimum, pint->maximum, RESET_COLOR,
//             PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//             g_value_get_int (&value), RESET_COLOR);
//         break;
//       }
//       case G_TYPE_UINT64:
//       {
//         GParamSpecUInt64 *puint64 = G_PARAM_SPEC_UINT64 (param);

//         n_print ("%sUnsigned Integer64%s. %sRange%s: %s%" G_GUINT64_FORMAT " - "
//             "%" G_GUINT64_FORMAT "%s %sDefault%s: %s%" G_GUINT64_FORMAT "%s ",
//             DATATYPE_COLOR, RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, puint64->minimum, puint64->maximum,
//             RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, g_value_get_uint64 (&value), RESET_COLOR);
//         break;
//       }
//       case G_TYPE_INT64:
//       {
//         GParamSpecInt64 *pint64 = G_PARAM_SPEC_INT64 (param);

//         n_print ("%sInteger64%s. %sRange%s: %s%" G_GINT64_FORMAT " - %"
//             G_GINT64_FORMAT "%s %sDefault%s: %s%" G_GINT64_FORMAT "%s ",
//             DATATYPE_COLOR, RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, pint64->minimum, pint64->maximum,
//             RESET_COLOR, PROP_ATTR_NAME_COLOR, RESET_COLOR,
//             PROP_ATTR_VALUE_COLOR, g_value_get_int64 (&value), RESET_COLOR);
//         break;
//       }
//       case G_TYPE_FLOAT:
//       {
//         GParamSpecFloat *pfloat = G_PARAM_SPEC_FLOAT (param);

//         n_print ("%sFloat%s. %sRange%s: %s%15.7g - %15.7g%s "
//             "%sDefault%s: %s%15.7g%s ", DATATYPE_COLOR, RESET_COLOR,
//             PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//             pfloat->minimum, pfloat->maximum, RESET_COLOR, PROP_ATTR_NAME_COLOR,
//             RESET_COLOR, PROP_ATTR_VALUE_COLOR, g_value_get_float (&value),
//             RESET_COLOR);
//         break;
//       }
//       case G_TYPE_DOUBLE:
//       {
//         GParamSpecDouble *pdouble = G_PARAM_SPEC_DOUBLE (param);

//         n_print ("%sDouble%s. %sRange%s: %s%15.7g - %15.7g%s "
//             "%sDefault%s: %s%15.7g%s ", DATATYPE_COLOR, RESET_COLOR,
//             PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//             pdouble->minimum, pdouble->maximum, RESET_COLOR,
//             PROP_ATTR_NAME_COLOR, RESET_COLOR, PROP_ATTR_VALUE_COLOR,
//             g_value_get_double (&value), RESET_COLOR);
//         break;
//       }
//       case G_TYPE_CHAR:
//       case G_TYPE_UCHAR:
//         GST_ERROR ("%s: property '%s' of type char: consider changing to "
//             "int/string", G_OBJECT_CLASS_NAME (obj_class),
//             g_param_spec_get_name (param));
//         /* fall through */
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
