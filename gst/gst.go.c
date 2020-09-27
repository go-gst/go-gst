#include <gst/gst.h>
#include <gst/app/gstappsink.h>
#include <gst/app/gstappsrc.h>
#include "gst.go.h"

/*
	Utilitits
*/

gboolean isParamSpecTypeCaps (GParamSpec * p)
{
	return p->value_type == GST_TYPE_CAPS;
}

gboolean isParamSpecEnum (GParamSpec * p)
{
	return G_IS_PARAM_SPEC_ENUM(p);
}

gboolean isParamSpecFlags (GParamSpec * p)
{
	return G_IS_PARAM_SPEC_FLAGS(p);
}

gboolean isParamSpecObject (GParamSpec * p)
{
	return G_IS_PARAM_SPEC_OBJECT(p);
}

gboolean isParamSpecBoxed (GParamSpec * p)
{
	return G_IS_PARAM_SPEC_BOXED(p);
}

gboolean isParamSpecPointer (GParamSpec * p)
{
	return G_IS_PARAM_SPEC_POINTER(p);
}

gboolean isParamSpecFraction (GParamSpec * p)
{
	return GST_IS_PARAM_SPEC_FRACTION(p);
}

gboolean isParamSpecGstArray (GParamSpec * p)
{
	return p->value_type == GST_TYPE_ARRAY;
}

GEnumValue * getEnumValues (GParamSpec * p, guint * size)
{
	GEnumValue * values;
    values = G_ENUM_CLASS (g_type_class_ref (p->value_type))->values;
	guint i = 0;
	while (values[i].value_name) {
    	++i;
    }
	*size = i;
	return values;
}

GFlagsValue * getParamSpecFlags (GParamSpec * p, guint * size)
{
	GParamSpecFlags *pflags = G_PARAM_SPEC_FLAGS (p);
	GFlagsValue *vals = pflags->flags_class->values;
	guint i = 0;
	while (vals[i].value_name) {
    	++i;
    }
	*size = i;
	return vals;
}

gboolean structureForEach (GQuark field_id, GValue * value, gpointer user_data)
{
	return structForEachCb(field_id, value, user_data);
}

GObjectClass * getGObjectClass (void * p) {
	return G_OBJECT_GET_CLASS (p);
}

int sizeOfGCharArray (gchar ** arr) {
	int i;
	for (i = 0 ; 1 ; i = i + 1) {
		if (arr[i] == NULL) { return i; };
	}
}

gboolean gstObjectFlagIsSet (GstObject * obj, GstElementFlags flags)
{
	return GST_OBJECT_FLAG_IS_SET (obj, flags);
}

gboolean gstElementIsURIHandler (GstElement * elem)
{
	return GST_IS_URI_HANDLER (elem);
}

/*
	Number functions
*/

GParamSpecUInt * getParamUInt (GParamSpec * param)
{
	return G_PARAM_SPEC_UINT (param);
}

GParamSpecInt * getParamInt (GParamSpec * param)
{
	return G_PARAM_SPEC_INT (param);
}

GParamSpecUInt64 * getParamUInt64 (GParamSpec * param)
{
	return G_PARAM_SPEC_UINT64 (param);
}

GParamSpecInt64 * getParamInt64 (GParamSpec * param)
{
	return G_PARAM_SPEC_INT64 (param);
}

GParamSpecFloat * getParamFloat (GParamSpec * param)
{
	return G_PARAM_SPEC_FLOAT (param);
}

GParamSpecDouble * getParamDouble(GParamSpec * param)
{
	return G_PARAM_SPEC_DOUBLE (param);
}

/*
	Type Castings
*/

GstUri *
toGstURI(void *p)
{
	return (GST_URI(p));
}

GstURIHandler *
toGstURIHandler(void *p)
{
	return (GST_URI_HANDLER(p));
}

GstRegistry *
toGstRegistry(void *p)
{
	return (GST_REGISTRY(p));
}

GstPlugin *
toGstPlugin(void *p)
{
	return (GST_PLUGIN(p));
}

GstPluginFeature *
toGstPluginFeature(void *p)
{
	return (GST_PLUGIN_FEATURE(p));
}

GstObject *
toGstObject(void *p)
{
	return (GST_OBJECT(p));
}

GstElementFactory *
toGstElementFactory(void *p)
{
	return (GST_ELEMENT_FACTORY(p));
}

GstElement *
toGstElement(void *p)
{
	return (GST_ELEMENT(p));
}

GstAppSink *
toGstAppSink(void *p)
{
	return (GST_APP_SINK(p));
}

GstAppSrc *
toGstAppSrc(void *p)
{
	return (GST_APP_SRC(p));
}

GstBin *
toGstBin(void *p)
{
	return (GST_BIN(p));
}

GstBus *
toGstBus(void *p)
{
	return (GST_BUS(p));
}

GstMessage *
toGstMessage(void *p)
{
	return (GST_MESSAGE(p));
}

GstPipeline *
toGstPipeline(void *p)
{
	return (GST_PIPELINE(p));
}

GstPad *
toGstPad(void *p)
{
	return (GST_PAD(p));
}

GstPadTemplate *
toGstPadTemplate(void *p)
{
	return (GST_PAD_TEMPLATE(p));
}

GstStructure *
toGstStructure(void *p)
{
	return (GST_STRUCTURE(p));
}

GstClock *
toGstClock(void *p)
{
	return (GST_CLOCK(p));
}

GstMiniObject *
toGstMiniObject(void *p)
{
	return (GST_MINI_OBJECT(p));
}

GstCaps *
toGstCaps(void *p)
{
	return (GST_CAPS(p));
}

GstCapsFeatures *
toGstCapsFeatures(void *p)
{
	return (GST_CAPS_FEATURES(p));
}

GstBuffer *
toGstBuffer(void *p)
{
	return (GST_BUFFER(p));
}

GstBufferPool *
toGstBufferPool(void *p)
{
	return (GST_BUFFER_POOL(p));
}
