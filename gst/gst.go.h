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
	Number functions
*/

static GParamSpecUInt *
getParamUInt(GParamSpec * param)
{
	return G_PARAM_SPEC_UINT (param);
}

static GParamSpecInt *
getParamInt(GParamSpec * param)
{
	return G_PARAM_SPEC_INT (param);
}

static GParamSpecUInt64 *
getParamUInt64(GParamSpec * param)
{
	return G_PARAM_SPEC_UINT64 (param);
}

static GParamSpecInt64 *
getParamInt64(GParamSpec * param)
{
	return G_PARAM_SPEC_INT64 (param);
}

static GParamSpecFloat *
getParamFloat(GParamSpec * param)
{
	return G_PARAM_SPEC_FLOAT (param);
}

static GParamSpecDouble *
getParamDouble(GParamSpec * param)
{
	return G_PARAM_SPEC_DOUBLE (param);
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

static GstMiniObject *
toGstMiniObject(void *p)
{
	return (GST_MINI_OBJECT(p));
}
