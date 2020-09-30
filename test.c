#include <gst/gst.h>
#include <stdio.h>

gboolean caps_map_func (GstCapsFeatures * features, GstStructure * structure, gpointer user_data)
{
    printf(gst_caps_features_to_string(features));
    return TRUE;
}

int main () {
    gst_init(NULL, NULL);

    GstCaps * caps = gst_caps_from_string("audio/x-raw");

    gst_caps_filter_and_map_in_place(caps, caps_map_func, NULL);
    gst_caps_foreach(caps, caps_map_func, NULL);

    gst_caps_unref(caps);

    return 0;
}