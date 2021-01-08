#ifndef __GST_BASE_GO_H__
#define __GST_BASE_GO_H__

#include <gst/base/gstbasesrc.h>
#include <gst/base/gstbasesink.h>

inline GstBaseSink *     toGstBaseSink      (void *p)    { return GST_BASE_SINK_CAST(p); }
inline GstBaseSrc *      toGstBaseSrc       (void *p)    { return GST_BASE_SRC_CAST(p); }

inline GstBaseSinkClass *  toGstBaseSinkClass  (void *p) { return (GstBaseSinkClass *)p; }
inline GstBaseSrcClass *   toGstBaseSrcClass   (void *p) { return (GstBaseSrcClass  *)p; }

#endif