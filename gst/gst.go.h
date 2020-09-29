#include <gst/gst.h>

extern gboolean   structForEachCb  (GQuark field_id, GValue * value, gpointer user_data);
extern gboolean   goBusFunc        (GstBus * bus, GstMessage * msg, gpointer user_data);

gboolean        structureForEach (GQuark field_id, GValue * value, gpointer user_data);
gboolean        cgoBusFunc       (GstBus * bus, GstMessage * msg, gpointer user_data);

GEnumValue *        getEnumValues      (GParamSpec * p, guint * size);
GFlagsValue *       getParamSpecFlags  (GParamSpec * p, guint * size);

int             sizeOfGCharArray (gchar ** arr);

gboolean        isParamSpecTypeCaps   (GParamSpec * p);
gboolean        isParamSpecEnum       (GParamSpec * p);
gboolean        isParamSpecFlags      (GParamSpec * p);
gboolean        isParamSpecObject     (GParamSpec * p);
gboolean        isParamSpecBoxed      (GParamSpec * p);
gboolean        isParamSpecPointer    (GParamSpec * p);
gboolean        isParamSpecFraction   (GParamSpec * p);
gboolean        isParamSpecGstArray   (GParamSpec * p);

GObjectClass *  getGObjectClass (void * p);

gboolean        gstObjectFlagIsSet      (GstObject * obj, GstElementFlags flags);
gboolean        gstElementIsURIHandler  (GstElement * elem);

/*
	Number functions
*/

GParamSpecUInt *    getParamUInt    (GParamSpec * param);
GParamSpecInt *     getParamInt     (GParamSpec * param);
GParamSpecUInt64 *  getParamUInt64  (GParamSpec * param);
GParamSpecInt64 *   getParamInt64   (GParamSpec * param);
GParamSpecFloat *   getParamFloat   (GParamSpec * param);
GParamSpecDouble *  getParamDouble  (GParamSpec * param);

/*
	Type Castings
*/

GstAllocator *         toGstAllocator         (void *p);
GstUri *               toGstURI               (void *p);
GstURIHandler *        toGstURIHandler        (void *p);
GstRegistry *          toGstRegistry          (void *p);
GstPlugin *            toGstPlugin            (void *p);
GstPluginFeature *     toGstPluginFeature     (void *p);
GstObject *            toGstObject            (void *p);
GstElementFactory *    toGstElementFactory    (void *p);
GstElement *           toGstElement           (void *p);
GstBin *               toGstBin               (void *p);
GstBus *               toGstBus               (void *p);
GstMessage *           toGstMessage           (void *p);
GstPipeline *          toGstPipeline          (void *p);
GstPad *               toGstPad               (void *p);
GstPadTemplate *       toGstPadTemplate       (void *p);
GstStructure *         toGstStructure         (void *p);
GstClock *             toGstClock             (void *p);
GstMiniObject *        toGstMiniObject        (void *p);
GstCaps *              toGstCaps              (void *p);
GstCapsFeatures *      toGstCapsFeatures      (void *p);
GstBuffer *            toGstBuffer            (void *p);
GstBufferPool *        toGstBufferPool        (void *p);
GstSample *            toGstSample            (void *p);
GstDevice *            toGstDevice            (void *p);
GstStreamCollection *  toGstStreamCollection  (void *p);
GstStream *            toGstStream            (void *p);
GstMemory *            toGstMemory            (void *p);