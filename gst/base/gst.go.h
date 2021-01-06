#ifndef __GST_BASE_GO_H__
#define __GST_BASE_GO_H__

#include <gst/base/gstbasesrc.h>

inline GstBaseSrc *      toGstBaseSrc       (void *p)    { return GST_BASE_SRC_CAST(p); };
inline GstBaseSrcClass * toGstBaseSrcClass  (void *p)    { return (GstBaseSrcClass  *)p; };

#endif