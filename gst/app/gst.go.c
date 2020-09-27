#include <gst/app/gstappsink.h>
#include <gst/app/gstappsrc.h>

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
