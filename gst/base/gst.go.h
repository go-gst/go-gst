#ifndef __GST_BASE_GO_H__
#define __GST_BASE_GO_H__

#include <gst/base/base.h>
#include <stddef.h>

extern GstBaseSink *       toGstBaseSink       (void *p);
extern GstBaseSrc *        toGstBaseSrc        (void *p);
extern GstBaseTransform *  toGstBaseTransform  (void *p);
extern GstCollectPads *    toGstCollectPads    (void *p);
extern GstPushSrc *        toGstPushSrc        (void *p);

extern GstBaseSinkClass *       toGstBaseSinkClass       (void *p);
extern GstBaseSrcClass *        toGstBaseSrcClass        (void *p);
extern GstBaseTransformClass *  toGstBaseTransformClass  (void *p);
extern GstPushSrcClass *        toGstPushSrcClass        (void *p);

extern gint64  gstCollectDataDTS  (GstCollectData * gcd);

#endif
