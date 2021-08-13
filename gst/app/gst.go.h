#ifndef __GST_APP_GO_H__
#define __GST_APP_GO_H__

#include <gst/app/gstappsink.h>
#include <gst/app/gstappsrc.h>

extern GstAppSink *  toGstAppSink   (void *p);
extern GstAppSrc *   toGstAppSrc    (void *p);

#endif
