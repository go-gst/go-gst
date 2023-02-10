#ifndef __GST_GO_H__
#define __GST_GO_H__

#include <stdlib.h>
#include <gst/gst.h>
#include <gst/base/base.h>

typedef struct _PadDestroyNotifyInfo {
	gpointer pad_ptr;
	gpointer func_map_ptr;
} PadDestroyNotifyInfo;

/*
	Type Castings
*/

extern GType                  objectGType            (GObject *obj);

extern GstAllocator *         toGstAllocator         (void *p);
extern GstBin *               toGstBin               (void *p);
extern GstBinClass *          toGstBinClass          (void *p);
extern GstBufferList *        toGstBufferList        (void *p);
extern GstBufferPool *        toGstBufferPool        (void *p);
extern GstBuffer *            toGstBuffer            (void *p);
extern GstBus *               toGstBus               (void *p);
extern GstCapsFeatures *      toGstCapsFeatures      (void *p);
extern GstCaps *              toGstCaps              (void *p);
extern GstChildProxy *        toGstChildProxy        (void *p);
extern GstClock *             toGstClock             (void *p);
extern GstContext *           toGstContext           (void *p);
extern GstDevice *            toGstDevice            (void *p);
extern GstElementFactory *    toGstElementFactory    (void *p);
extern GstElementClass *      toGstElementClass      (void *p);
extern GstElement *           toGstElement           (void *p);
extern GstEvent *             toGstEvent             (void *p);
extern GstGhostPad *          toGstGhostPad          (void *p);
extern GstMemory *            toGstMemory            (void *p);
extern GstMessage *           toGstMessage           (void *p);
extern GstMeta *              toGstMeta              (void *p);
extern GstMiniObject *        toGstMiniObject        (void *p);
extern GstObject *            toGstObject            (void *p);
extern GstPad *               toGstPad               (void *p);
extern GstPadTemplate *       toGstPadTemplate       (void *p);
extern GstPipeline *          toGstPipeline          (void *p);
extern GstPluginFeature *     toGstPluginFeature     (void *p);
extern GstPlugin *            toGstPlugin            (void *p);
extern GstProxyPad *          toGstProxyPad          (void *p);
extern GstQuery *             toGstQuery             (void *p);
extern GstRegistry *          toGstRegistry          (void *p);
extern GstSample *            toGstSample            (void *p);
extern GstStreamCollection *  toGstStreamCollection  (void *p);
extern GstStream *            toGstStream            (void *p);
extern GstStructure *         toGstStructure         (void *p);
extern GstTagList   *         toGstTagList           (void *p);
extern GstTask *              toGstTask              (void *p);
extern GstTaskPool *          toGstTaskPool          (void *p);
extern GstURIHandler *        toGstURIHandler        (void *p);
extern GstUri *               toGstURI               (void *p);

/* Buffer Utilities */

extern GstBuffer * getBufferValue (GValue * val);

extern gboolean bufferIsWritable (GstBuffer * buf);

extern GstBuffer * makeBufferWritable (GstBuffer * buf);

extern GType bufferListType ();

/* BufferList Utilities */

extern gboolean bufferListIsWritable (GstBufferList * bufList);

extern GstBufferList * makeBufferListWritable (GstBufferList * bufList);

extern void addToBufferList (GstBufferList * bufList, GstBuffer * buf);

/* BufferPool utilities */

extern gboolean bufferPoolIsFlushing (GstBufferPool * pool);

/* Caps utilties */

extern gboolean capsIsWritable (GstCaps * caps);

extern GstCaps * makeCapsWritable (GstCaps * caps);

/* Context utilities */

extern gboolean contextIsWritable (GstContext * ctx);

extern GstContext * makeContextWritable (GstContext * ctx);

/* Event Utilities */

extern gboolean eventIsWritable (GstEvent * event);

extern GstEvent * makeEventWritable (GstEvent * event);

/* TOC Utilities */

extern gboolean         entryTypeIsAlternative (GstTocEntryType * type);
extern gboolean         entryTypeIsSequence    (GstTocEntryType * type);
extern GstTocEntry *    copyTocEntry           (GstTocEntry * entry);
extern GstTocEntry *    makeTocEntryWritable   (GstTocEntry * entry);
extern GstTocEntry *    tocEntryRef            (GstTocEntry * entry);
extern void             tocEntryUnref          (GstTocEntry * entry);
extern GstToc *         copyToc                (GstToc * toc);
extern GstToc *         makeTocWritable        (GstToc * toc);
extern GstToc *         tocRef                 (GstToc * toc);
extern void             tocUnref               (GstToc * toc);

/* TagList utilities */

extern gboolean         tagListIsWritable     (GstTagList * tagList);
extern GstTagList *     makeTagListWritable   (GstTagList * tagList);

/* Object Utilities */

extern gboolean        gstElementIsURIHandler  (GstElement * elem);
extern gboolean        gstObjectFlagIsSet      (GstObject * obj, GstElementFlags flags);

/* Element utilities */

extern GstTocSetter *  toTocSetter (GstElement * elem);
extern GstTagSetter *  toTagSetter (GstElement *elem);


/* Misc */

extern gpointer glistNext (GList * list);

extern int sizeOfGCharArray (gchar ** arr);

/* Sample Utilities */

extern GstSample * getSampleValue (GValue * val);

#endif
