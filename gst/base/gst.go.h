#ifndef __GST_BASE_GO_H__
#define __GST_BASE_GO_H__

#include <gst/base/base.h>
#include <stddef.h>

inline GstBaseSink *     toGstBaseSink      (void *p)    { return GST_BASE_SINK_CAST(p); }
inline GstBaseSrc *      toGstBaseSrc       (void *p)    { return GST_BASE_SRC_CAST(p); }
inline GstCollectPads *  toGstCollectPads   (void *p)    { return GST_COLLECT_PADS(p); }
inline GstPushSrc *      toGstPushSrc       (void *p)    { return GST_PUSH_SRC(p); }

inline GstBaseSinkClass *  toGstBaseSinkClass  (void *p) { return (GstBaseSinkClass *)p; }
inline GstBaseSrcClass *   toGstBaseSrcClass   (void *p) { return (GstBaseSrcClass  *)p; }
inline GstPushSrcClass *   toGstPushSrcClass   (void *p) { return (GstPushSrcClass *)p; }

inline gint64  gstCollectDataDTS  (GstCollectData * gcd)  { return GST_COLLECT_PADS_DTS(gcd); }

#endif