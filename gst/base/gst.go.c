#include "gst.go.h"

GstBaseSink *       toGstBaseSink       (void *p)    { return GST_BASE_SINK_CAST(p); }
GstBaseSrc *        toGstBaseSrc        (void *p)    { return GST_BASE_SRC_CAST(p); }
GstBaseTransform *  toGstBaseTransform  (void *p)    { return GST_BASE_TRANSFORM(p); }
GstCollectPads *    toGstCollectPads    (void *p)    { return GST_COLLECT_PADS(p); }
GstPushSrc *        toGstPushSrc        (void *p)    { return GST_PUSH_SRC(p); }

GstBaseSinkClass *       toGstBaseSinkClass       (void *p) { return (GstBaseSinkClass *)p; }
GstBaseSrcClass *        toGstBaseSrcClass        (void *p) { return (GstBaseSrcClass  *)p; }
GstBaseTransformClass *  toGstBaseTransformClass  (void *p) { return (GstBaseTransformClass *)p; }
GstPushSrcClass *        toGstPushSrcClass        (void *p) { return (GstPushSrcClass *)p; }

gint64  gstCollectDataDTS  (GstCollectData * gcd)  { return GST_COLLECT_PADS_DTS(gcd); }

