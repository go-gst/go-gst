#include "gst.go.h"

GType                  objectGType            (GObject *obj) { return G_OBJECT_TYPE(obj); };

GstAllocator *         toGstAllocator         (void *p) { return (GST_ALLOCATOR_CAST(p)); }
GstBin *               toGstBin               (void *p) { return (GST_BIN(p)); }
GstBinClass *          toGstBinClass          (void *p) { return (GST_BIN_CLASS(p)); }
GstBufferList *        toGstBufferList        (void *p) { return (GST_BUFFER_LIST(p)); }
GstBufferPool *        toGstBufferPool        (void *p) { return (GST_BUFFER_POOL(p)); }
GstBuffer *            toGstBuffer            (void *p) { return (GST_BUFFER(p)); }
GstBus *               toGstBus               (void *p) { return (GST_BUS(p)); }
GstCapsFeatures *      toGstCapsFeatures      (void *p) { return (GST_CAPS_FEATURES(p)); }
GstCaps *              toGstCaps              (void *p) { return (GST_CAPS(p)); }
GstChildProxy *        toGstChildProxy        (void *p) { return (GST_CHILD_PROXY(p)); }
GstClock *             toGstClock             (void *p) { return (GST_CLOCK(p)); }
GstContext *           toGstContext           (void *p) { return (GST_CONTEXT_CAST(p)); }
GstDevice *            toGstDevice            (void *p) { return (GST_DEVICE_CAST(p)); }
GstElementFactory *    toGstElementFactory    (void *p) { return (GST_ELEMENT_FACTORY(p)); }
GstElementClass *      toGstElementClass      (void *p) { return (GST_ELEMENT_CLASS(p)); }
GstElement *           toGstElement           (void *p) { return (GST_ELEMENT(p)); }
GstEvent *             toGstEvent             (void *p) { return (GST_EVENT(p)); }
GstGhostPad *          toGstGhostPad          (void *p) { return (GST_GHOST_PAD(p)); }
GstMemory *            toGstMemory            (void *p) { return (GST_MEMORY_CAST(p)); }
GstMessage *           toGstMessage           (void *p) { return (GST_MESSAGE(p)); }
GstMeta *              toGstMeta              (void *p) { return (GST_META_CAST(p)); }
GstMiniObject *        toGstMiniObject        (void *p) { return (GST_MINI_OBJECT(p)); }
GstObject *            toGstObject            (void *p) { return (GST_OBJECT(p)); }
GstPad *               toGstPad               (void *p) { return (GST_PAD(p)); }
GstPadTemplate *       toGstPadTemplate       (void *p) { return (GST_PAD_TEMPLATE(p)); }
GstPipeline *          toGstPipeline          (void *p) { return (GST_PIPELINE(p)); }
GstPluginFeature *     toGstPluginFeature     (void *p) { return (GST_PLUGIN_FEATURE(p)); }
GstPlugin *            toGstPlugin            (void *p) { return (GST_PLUGIN(p)); }
GstProxyPad *          toGstProxyPad          (void *p) { return (GST_PROXY_PAD(p)); }
GstQuery *             toGstQuery             (void *p) { return (GST_QUERY(p)); }
GstRegistry *          toGstRegistry          (void *p) { return (GST_REGISTRY(p)); }
GstSample *            toGstSample            (void *p) { return (GST_SAMPLE(p)); }
GstStreamCollection *  toGstStreamCollection  (void *p) { return (GST_STREAM_COLLECTION_CAST(p)); }
GstStream *            toGstStream            (void *p) { return (GST_STREAM_CAST(p)); }
GstStructure *         toGstStructure         (void *p) { return (GST_STRUCTURE(p)); }
GstTagList   *         toGstTagList           (void *p) { return (GST_TAG_LIST(p)); }
GstTask *              toGstTask              (void *p) { return (GST_TASK_CAST(p)); }
GstTaskPool *          toGstTaskPool          (void *p) { return (GST_TASK_POOL_CAST(p)); }
GstURIHandler *        toGstURIHandler        (void *p) { return (GST_URI_HANDLER(p)); }
GstUri *               toGstURI               (void *p) { return (GST_URI(p)); }

/* Buffer Utilities */

GstBuffer * getBufferValue (GValue * val)
{
	return gst_value_get_buffer(val);
}

gboolean bufferIsWritable (GstBuffer * buf)
{
	return (gst_buffer_is_writable(buf));
}

GstBuffer * makeBufferWritable (GstBuffer * buf)
{
	return (gst_buffer_make_writable(buf));
}

GType bufferListType ()
{
	return GST_TYPE_BUFFER_LIST;
}

/* BufferList Utilities */

gboolean bufferListIsWritable (GstBufferList * bufList)
{
	return gst_buffer_list_is_writable(bufList);
}

GstBufferList * makeBufferListWritable (GstBufferList * bufList)
{
	return gst_buffer_list_make_writable(bufList);
}

void addToBufferList (GstBufferList * bufList, GstBuffer * buf)
{
	gst_buffer_list_add(bufList, buf);
}

/* BufferPool utilities */

gboolean bufferPoolIsFlushing (GstBufferPool * pool)
{
	return GST_BUFFER_POOL_IS_FLUSHING(pool);
}

/* Caps utilties */

gboolean capsIsWritable (GstCaps * caps)
{
	return gst_caps_is_writable(caps);
}

GstCaps * makeCapsWritable (GstCaps * caps)
{
	return gst_caps_make_writable(caps);
}

/* Context utilities */

gboolean contextIsWritable (GstContext * ctx)
{
	return gst_context_is_writable(ctx);
}

GstContext * makeContextWritable (GstContext * ctx)
{
	return gst_context_make_writable(ctx);
}

/* Event Utilities */

gboolean eventIsWritable (GstEvent * event)
{
	return gst_event_is_writable(event);
}

GstEvent * makeEventWritable (GstEvent * event)
{
	return gst_event_make_writable(event);
}

/* TOC Utilities */

gboolean         entryTypeIsAlternative (GstTocEntryType * type) { return GST_TOC_ENTRY_TYPE_IS_ALTERNATIVE(type); }
gboolean         entryTypeIsSequence    (GstTocEntryType * type) { return GST_TOC_ENTRY_TYPE_IS_SEQUENCE(type); }
GstTocEntry *    copyTocEntry           (GstTocEntry * entry)    { return gst_toc_entry_copy(entry); }
GstTocEntry *    makeTocEntryWritable   (GstTocEntry * entry)    { return gst_toc_entry_make_writable(entry); }
GstTocEntry *    tocEntryRef            (GstTocEntry * entry)    { return gst_toc_entry_ref(entry); }
void             tocEntryUnref          (GstTocEntry * entry)    { gst_toc_entry_unref(entry); }
GstToc *         copyToc                (GstToc * toc)           { return gst_toc_copy(toc); }
GstToc *         makeTocWritable        (GstToc * toc)           { return gst_toc_make_writable(toc); }
GstToc *         tocRef                 (GstToc * toc)           { return gst_toc_ref(toc); }
void             tocUnref               (GstToc * toc)           { gst_toc_unref(toc); }

/* TagList utilities */

gboolean         tagListIsWritable     (GstTagList * tagList)   { return gst_tag_list_is_writable(tagList); }
GstTagList *     makeTagListWritable   (GstTagList * tagList)   { return gst_tag_list_make_writable(tagList); }

/* Object Utilities */

gboolean        gstElementIsURIHandler  (GstElement * elem)                      { return (GST_IS_URI_HANDLER(elem)); }
gboolean        gstObjectFlagIsSet      (GstObject * obj, GstElementFlags flags) { return (GST_OBJECT_FLAG_IS_SET(obj, flags)); }

/* Element utilities */

GstTocSetter *  toTocSetter (GstElement * elem) { return GST_TOC_SETTER(elem); }
GstTagSetter *  toTagSetter (GstElement *elem)  { return GST_TAG_SETTER(elem); }

/* Sample Utilities */

GstSample * getSampleValue (GValue * val)
{
	return gst_value_get_sample(val);
}


/* Misc */

gpointer glistNext (GList * list)
{
	return g_list_next(list);
}

int sizeOfGCharArray (gchar ** arr)
{
	int i;
	for (i = 0 ; 1 ; i = i + 1) {
		if (arr[i] == NULL) { return i; };
	}
}
