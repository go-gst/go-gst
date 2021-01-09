// This example demonstrates a filesrc plugin implemented in Go.
//
// Every element in a Gstreamer pipeline is provided by plugins. Some are builtin while
// others are provided by third-parties or distributed privately. The plugins are built
// around the GObject type system.
//
// Go-gst offers loose bindings around the GObject type system to provide the necessary
// functionality to implement these plugins. The example in this code produces an element
// that can read from a file on the local system.
//
// In order to build the plugin for use by GStreamer, you can do the following:
//
//     $ go generate
//     $ go build -o libgstgofilesrc.so -buildmode c-shared .
//
//
//go:generate gst-plugin-gen
//
// +plugin:Name=gofilesrc
// +plugin:Description=File plugins written in go
// +plugin:Version=v0.0.1
// +plugin:License=gst.LicenseLGPL
// +plugin:Source=go-gst
// +plugin:Package=examples
// +plugin:Origin=https://github.com/tinyzimmer/go-gst
// +plugin:ReleaseDate=2021-01-04
//
// +element:Name=gofilesrc
// +element:Rank=gst.RankNone
// +element:Impl=fileSrc
// +element:Subclass=base.ExtendsBaseSrc
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
)

// main is left unimplemented since these files are compiled to c-shared.
func main() {}

// CAT is the log category for the gofilesrc. It is safe to define GStreamer objects as globals
// without calling gst.Init, since in the context of a loaded plugin all initialization has
// already been taken care of by the loading application.
var CAT = gst.NewDebugCategory(
	"gofilesrc",
	gst.DebugColorNone,
	"GoFileSrc Element",
)

// Here we define a list of ParamSpecs that will make up the properties for our element.
// This element only has a single property, the location of the file to read from.
// When getting and setting properties later on, you will reference them by their index in
// this list.
var properties = []*gst.ParameterSpec{
	gst.NewStringParameter(
		"location",                          // The name of the parameter
		"File Location",                     // The long name for the parameter
		"Location of the file to read from", // A blurb about the parameter
		nil,                                 // A default value for the parameter
		gst.ParameterReadWrite,              // Flags for the parameter
	),
}

// Here we declare a private struct to hold our internal state.
type state struct {
	// Whether the element is started or not
	started bool
	// The file the element is reading from
	file *os.File
	// The information about the file retrieved from stat
	fileInfo os.FileInfo
	// The current position in the file
	position uint64
}

// This is another private struct where we hold the parameter values set on our
// element.
type settings struct {
	location string
}

// Finally a structure is defined that implements (at a minimum) the gst.GoElement interface.
// It is possible to signal to the bindings to inherit from other classes or implement other
// interfaces via the registration and TypeInit processes.
type fileSrc struct {
	// The settings for the element
	settings *settings
	// The current state of the element
	state *state
}

// Private methods only used internally by the plugin

// setLocation is a simple method to check the validity of a provided file path and set the
// local value with it.
func (f *fileSrc) setLocation(path string) error {
	if f.state.started {
		return errors.New("Changing the `location` property on a started `GoFileSrc` is not supported")
	}
	f.settings.location = strings.TrimPrefix(path, "file://") // should obviously use url.URL and do actual parsing
	return nil
}

// The ObjectSubclass implementations below are for registering the various aspects of our
// element and its capabilities with the type system. These are the minimum methods that
// should be implemented by an element.

// Every element needs to provide its own constructor that returns an initialized
// gst.GoElement implementation. Here we simply create a new fileSrc with zeroed settings
// and state objects.
func (f *fileSrc) New() gst.GoElement {
	CAT.Log(gst.LevelLog, "Initializing new fileSrc object")
	return &fileSrc{
		settings: &settings{},
		state:    &state{},
	}
}

// The TypeInit method should register any additional interfaces provided by the element.
// In this example we signal to the type system that we also implement the GstURIHandler interface.
func (f *fileSrc) TypeInit(instance *gst.TypeInstance) {
	CAT.Log(gst.LevelLog, "Adding URIHandler interface to type")
	instance.AddInterface(gst.InterfaceURIHandler)
}

// The ClassInit method should specify the metadata for this element and add any pad templates
// and properties.
func (f *fileSrc) ClassInit(klass *gst.ElementClass) {
	CAT.Log(gst.LevelLog, "Initializing gofilesrc class")
	klass.SetMetadata(
		"File Source",
		"Source/File",
		"Read stream from a file",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
	CAT.Log(gst.LevelLog, "Adding src pad template and properties to class")
	klass.AddPadTemplate(gst.NewPadTemplate(
		"src",
		gst.PadDirectionSource,
		gst.PadPresenceAlways,
		gst.NewAnyCaps(),
	))
	klass.InstallProperties(properties)
}

// Object implementations are used during the initialization of an element. The
// methods are called once the object is constructed and its properties are read
// and written to. These and the rest of the methods described below are documented
// in interfaces in the bindings, however only individual methods needs from those
// interfaces need to be implemented. When left unimplemented, the behavior of the parent
// class is inherited.

// SetProperty is called when a `value` is set to the property at index `id` in the
// properties slice that we installed during ClassInit. It should attempt to register
// the value locally or signal any errors that occur in the process.
func (f *fileSrc) SetProperty(self *gst.Object, id uint, value *glib.Value) {
	param := properties[id]
	switch param.Name() {
	case "location":
		var val string
		if value == nil {
			val = ""
		} else {
			val, _ = value.GetString()
		}
		if err := f.setLocation(val); err != nil {
			gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorSettings,
				fmt.Sprintf("Could not set location on object: %s", err.Error()),
				"",
			)
			return
		}
		self.Log(CAT, gst.LevelInfo, fmt.Sprintf("Set `location` to %s", f.settings.location))
	}
}

// GetProperty is called to retrieve the value of the property at index `id` in the properties
// slice provided at ClassInit.
func (f *fileSrc) GetProperty(self *gst.Object, id uint) *glib.Value {
	param := properties[id]
	switch param.Name() {
	case "location":
		if f.settings.location == "" {
			return nil
		}
		val, err := glib.GValue(f.settings.location)
		if err == nil {
			return val
		}
		gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
			fmt.Sprintf("Could not convert %s to GValue", f.settings.location),
			err.Error(),
		)
	}
	return nil
}

// Constructed is called when the type system is done constructing the object. Any finalizations required
// during the initialization process can be performed here. In this example, we set the format on our
// underlying GstBaseSrc to bytes.
func (f *fileSrc) Constructed(self *gst.Object) {
	self.Log(CAT, gst.LevelLog, "Setting format of GstBaseSrc to bytes")
	base.ToGstBaseSrc(self).SetFormat(gst.FormatBytes)
}

// GstBaseSrc implementations are optional methods to implement from the base.GstBaseSrcImpl interface.
// If the method is not overridden by the implementing struct, it will be inherited from the parent class.

// IsSeekable returns that we are, in fact, seekable.
func (f *fileSrc) IsSeekable(*base.GstBaseSrc) bool { return true }

// GetSize will return the total size of the file at the configured location.
func (f *fileSrc) GetSize(self *base.GstBaseSrc) (bool, int64) {
	if !f.state.started {
		return false, 0
	}
	return true, f.state.fileInfo.Size()
}

// Start is called to start this element. In this example, the configured file is opened for reading,
// and any error encountered in the process is posted to the pipeline.
func (f *fileSrc) Start(self *base.GstBaseSrc) bool {
	if f.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "GoFileSrc is already started", "")
		return false
	}

	if f.settings.location == "" {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "File location is not defined", "")
		return false
	}

	stat, err := os.Stat(f.settings.location)
	if err != nil {
		if os.IsNotExist(err) {
			self.ErrorMessage(gst.DomainResource, gst.ResourceErrorOpenRead,
				fmt.Sprintf("%s does not exist", f.settings.location), "")
			return false
		}
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorOpenRead,
			fmt.Sprintf("Could not stat %s, err: %s", f.settings.location, err.Error()), "")
		return false
	}
	if stat.IsDir() {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorOpenRead,
			fmt.Sprintf("%s is a directory", f.settings.location), "")
		return false
	}
	f.state.fileInfo = stat
	self.Log(CAT, gst.LevelDebug, fmt.Sprintf("file stat - name: %s  size: %d  mode: %v  modtime: %v", stat.Name(), stat.Size(), stat.Mode(), stat.ModTime()))

	self.Log(CAT, gst.LevelDebug, fmt.Sprintf("Opening file %s for reading", f.settings.location))
	f.state.file, err = os.Open(f.settings.location)
	if err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorOpenRead,
			fmt.Sprintf("Could not open file %s for reading", f.settings.location), err.Error())
		return false
	}

	f.state.position = 0
	f.state.started = true

	self.StartComplete(gst.FlowOK)

	self.Log(CAT, gst.LevelInfo, "GoFileSrc has started")
	return true
}

// Stop is called to stop the element. The file is closed and the local values are zeroed out.
func (f *fileSrc) Stop(self *base.GstBaseSrc) bool {
	if !f.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "FileSrc is not started", "")
		return false
	}

	if err := f.state.file.Close(); err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorClose, "Failed to close the source file", err.Error())
		return false
	}

	f.state.file = nil
	f.state.position = 0
	f.state.started = false

	self.Log(CAT, gst.LevelInfo, "GoFileSrc has stopped")
	return true
}

// Fill is called to fill a pre-allocated buffer with the data at offset up to the given size.
// Since we declared that we are seekable, we need to support the provided offset not neccesarily matching
// where we currently are in the file. This is why we store the position in the file locally.
func (f *fileSrc) Fill(self *base.GstBaseSrc, offset uint64, size uint, buffer *gst.Buffer) gst.FlowReturn {
	if !f.state.started || f.state.file == nil {
		self.ErrorMessage(gst.DomainCore, gst.CoreErrorFailed, "Not started yet", "")
		return gst.FlowError
	}

	self.Log(CAT, gst.LevelLog, fmt.Sprintf("Request to fill buffer from offset %v with size %v", offset, size))

	if f.state.position != offset {
		self.Log(CAT, gst.LevelDebug, fmt.Sprintf("Seeking to new position at offset %v from previous position at offset %v", offset, f.state.position))
		if _, err := f.state.file.Seek(int64(offset), 0); err != nil {
			self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSeek,
				fmt.Sprintf("Failed to seek to %d in file", offset), err.Error())
			return gst.FlowError
		}
		f.state.position = offset
	}

	bufmap := buffer.Map(gst.MapWrite)
	if bufmap == nil {
		self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, "Failed to map buffer", "")
		return gst.FlowError
	}
	defer buffer.Unmap()

	self.Log(CAT, gst.LevelLog, fmt.Sprintf("Reading %v bytes from offset %v in file into buffer at %v", size, f.state.position, bufmap.Data()))
	if _, err := io.CopyN(bufmap.Writer(), f.state.file, int64(size)); err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorRead,
			fmt.Sprintf("Failed to read %d bytes from file at %d into buffer", size, offset), err.Error())
		return gst.FlowError
	}
	buffer.SetSize(int64(size))

	f.state.position = f.state.position + uint64(size)
	self.Log(CAT, gst.LevelLog, fmt.Sprintf("Incremented current position to %v", f.state.position))

	return gst.FlowOK
}

// URIHandler implementations are the methods required by the GstURIHandler interface.

// GetURI returns the currently configured URI
func (f *fileSrc) GetURI() string { return fmt.Sprintf("file://%s", f.settings.location) }

// GetURIType returns the types of URI this element supports.
func (f *fileSrc) GetURIType() gst.URIType { return gst.URISource }

// GetProtocols returns the protcols this element supports.
func (f *fileSrc) GetProtocols() []string { return []string{"file"} }

// SetURI should set the URI that this element is working on.
func (f *fileSrc) SetURI(uri string) (bool, error) {
	if uri == "file://" {
		return true, nil
	}
	err := f.setLocation(uri)
	if err != nil {
		return false, err
	}
	CAT.Log(gst.LevelInfo, fmt.Sprintf("Set `location` to %s via URIHandler", f.settings.location))
	return true, nil
}
