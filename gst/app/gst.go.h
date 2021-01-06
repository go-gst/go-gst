#ifndef __GST_APP_GO_H__
#define __GST_APP_GO_H__

#include <gst/app/gstappsink.h>
#include <gst/app/gstappsrc.h>

inline GstAppSink *  toGstAppSink   (void *p) { return (GST_APP_SINK(p)); }
inline GstAppSrc *   toGstAppSrc    (void *p) { return (GST_APP_SRC(p)); }

#endif