package gst

/*
#include "gst.go.h"

extern gboolean goPluginInit (GstPlugin * plugin, gpointer user_data);
extern gboolean goGlobalPluginInit (GstPlugin * plugin);

gboolean cgoPluginInit (GstPlugin * plugin, gpointer user_data)
{
	return goPluginInit(plugin, user_data);
}

gboolean cgoGlobalPluginInit(GstPlugin * plugin)
{
	return goGlobalPluginInit(plugin);
}

GstPluginDesc * getPluginMeta (gint major,
					gint minor,
					gchar * name,
					gchar * description,
					GstPluginInitFunc init,
					gchar * version,
					gchar * license,
					gchar * source,
					gchar * package,
					gchar * origin,
					gchar * release_datetime)
{

	GstPluginDesc * desc = malloc ( sizeof (GstPluginDesc) );

	desc->major_version = major;
	desc->minor_version = minor;
	desc->name = name;
	desc->description = description;
	desc->plugin_init = init;
	desc->version = version;
	desc->license = license;
	desc->source = source;
	desc->package = package;
	desc->origin = origin;
	desc->release_datetime = release_datetime;

	return desc;
}

*/
import "C"

import (
	"errors"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
	"github.com/tinyzimmer/go-glib/glib"
)

// PluginMetadata represents the information to include when registering a new plugin
// with gstreamer.
type PluginMetadata struct {
	// The major version number of the GStreamer core that the plugin was compiled for, you can just use VersionMajor here
	MajorVersion Version
	// The minor version number of the GStreamer core that the plugin was compiled for, you can just use VersionMinor here
	MinorVersion Version
	// A unique name of the plugin (ideally prefixed with an application- or library-specific namespace prefix in order to
	// avoid name conflicts in case a similar plugin with the same name ever gets added to GStreamer)
	Name string
	// A description of the plugin
	Description string
	// The function to call when initiliazing the plugin
	Init PluginInitFunc
	// The version of the plugin
	Version string
	// The license for the plugin, must match one of the license constants in this package
	License License
	// The source module the plugin belongs to
	Source string
	// The shipped package the plugin belongs to
	Package string
	// The URL to the provider of the plugin
	Origin string
	// The date of release in ISO 8601 format.
	// See https://gstreamer.freedesktop.org/documentation/gstreamer/gstplugin.html?gi-language=c#GstPluginDesc for more details.
	ReleaseDate string
}

var globalPluginInit PluginInitFunc

// Export will export the PluginMetadata to an unsafe pointer to a GstPluginDesc.
func (p *PluginMetadata) Export() unsafe.Pointer {
	globalPluginInit = p.Init
	desc := C.getPluginMeta(
		C.gint(p.MajorVersion),
		C.gint(p.MinorVersion),
		(*C.gchar)(unsafe.Pointer(&[]byte(p.Name)[0])),
		(*C.gchar)(C.CString(p.Description)),
		(C.GstPluginInitFunc(C.cgoGlobalPluginInit)),
		(*C.gchar)(C.CString(p.Version)),
		(*C.gchar)(C.CString(string(p.License))),
		(*C.gchar)(C.CString(p.Source)),
		(*C.gchar)(C.CString(p.Package)),
		(*C.gchar)(C.CString(p.Origin)),
		(*C.gchar)(C.CString(p.ReleaseDate)),
	)
	return unsafe.Pointer(desc)
}

// PluginInitFunc is a function called by the plugin loader at startup. This function should register
// all the features of the plugin. The function should return true if the plugin is initialized successfully.
type PluginInitFunc func(*Plugin) bool

// Plugin is a go representation of a GstPlugin.
type Plugin struct{ *Object }

// FromGstPluginUnsafeNone wraps the given pointer in a Plugin.
func FromGstPluginUnsafeNone(plugin unsafe.Pointer) *Plugin {
	return &Plugin{wrapObject(glib.TransferNone(plugin))}
}

// FromGstPluginUnsafeFull wraps the given pointer in a Plugin.
func FromGstPluginUnsafeFull(plugin unsafe.Pointer) *Plugin {
	return &Plugin{wrapObject(glib.TransferFull(plugin))}
}

// RegisterPlugin will register a static plugin, i.e. a plugin which is private to an application
// or library and contained within the application or library (as opposed to being shipped as a
// separate module file).
func RegisterPlugin(desc *PluginMetadata, initFunc PluginInitFunc) bool {
	cName := C.CString(desc.Name)
	cDesc := C.CString(desc.Description)
	cVers := C.CString(desc.Version)
	cLics := C.CString(string(desc.License))
	cSrc := C.CString(desc.Source)
	cPkg := C.CString(desc.Package)
	cOrg := C.CString(desc.Origin)
	defer func() {
		for _, ptr := range []*C.char{cName, cDesc, cVers, cLics, cSrc, cPkg, cOrg} {
			C.free(unsafe.Pointer(ptr))
		}
	}()
	fPtr := gopointer.Save(initFunc)
	return gobool(C.gst_plugin_register_static_full(
		C.gint(desc.MajorVersion), C.gint(desc.MinorVersion),
		(*C.gchar)(cName), (*C.gchar)(cDesc),
		C.GstPluginInitFullFunc(C.cgoPluginInit),
		(*C.gchar)(cVers), (*C.gchar)(cLics),
		(*C.gchar)(cSrc), (*C.gchar)(cPkg),
		(*C.gchar)(cOrg), (C.gpointer)(unsafe.Pointer(fPtr)),
	))
}

// LoadPluginByName loads the named plugin and places a ref count on it. The function
// returns nil if the plugin could not be loaded.
func LoadPluginByName(name string) *Plugin {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	plugin := C.gst_plugin_load_by_name((*C.gchar)(unsafe.Pointer(cstr)))
	if plugin == nil {
		return nil
	}
	return FromGstPluginUnsafeFull(unsafe.Pointer(plugin))
}

// LoadPluginFile loads the given plugin and refs it. If an error is returned Plugin will be nil.
func LoadPluginFile(fpath string) (*Plugin, error) {
	cstr := C.CString(fpath)
	defer C.free(unsafe.Pointer(cstr))
	var gerr *C.GError
	plugin := C.gst_plugin_load_file((*C.gchar)(unsafe.Pointer(cstr)), (**C.GError)(&gerr))
	if gerr != nil {
		defer C.g_free((C.gpointer)(gerr))
		return nil, errors.New(C.GoString(gerr.message))
	}
	return FromGstPluginUnsafeFull(unsafe.Pointer(plugin)), nil
}

// Instance returns the underlying GstPlugin instance.
func (p *Plugin) Instance() *C.GstPlugin { return C.toGstPlugin(p.Unsafe()) }

// Description returns the description for this plugin.
func (p *Plugin) Description() string {
	ret := C.gst_plugin_get_description((*C.GstPlugin)(p.Instance()))
	if ret == nil {
		return ""
	}
	return C.GoString(ret)
}

// Filename returns the filename for this plugin.
func (p *Plugin) Filename() string {
	ret := C.gst_plugin_get_filename((*C.GstPlugin)(p.Instance()))
	if ret == nil {
		return ""
	}
	return C.GoString(ret)
}

// Version returns the version for this plugin.
func (p *Plugin) Version() string {
	ret := C.gst_plugin_get_version((*C.GstPlugin)(p.Instance()))
	if ret == nil {
		return ""
	}
	return C.GoString(ret)
}

// License returns the license for this plugin.
func (p *Plugin) License() License {
	ret := C.gst_plugin_get_license((*C.GstPlugin)(p.Instance()))
	if ret == nil {
		return ""
	}
	return License(C.GoString(ret))
}

// Source returns the source module for this plugin.
func (p *Plugin) Source() string {
	ret := C.gst_plugin_get_source((*C.GstPlugin)(p.Instance()))
	if ret == nil {
		return ""
	}
	return C.GoString(ret)
}

// Package returns the binary package for this plugin.
func (p *Plugin) Package() string {
	ret := C.gst_plugin_get_package((*C.GstPlugin)(p.Instance()))
	if ret == nil {
		return ""
	}
	return C.GoString(ret)
}

// Origin returns the origin URL for this plugin.
func (p *Plugin) Origin() string {
	ret := C.gst_plugin_get_origin((*C.GstPlugin)(p.Instance()))
	if ret == nil {
		return ""
	}
	return C.GoString(ret)
}
