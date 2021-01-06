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
// In order to build the plugin for use by GStreamer, you may do the following:
//
//     $ go build -o libgstgofilesrc.so -buildmode c-shared .
//
//
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/gotk3/gotk3/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
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
	path = strings.TrimPrefix(path, "file://")
	if path == "" {
		f.settings.location = ""
		return nil
	}
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", path)
		}
		return fmt.Errorf("Could not stat %s, err: %s", path, err.Error())
	}
	if stat.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}
	f.settings.location = path
	return nil
}

// The ObjectSubclass implementations below are for registering the various aspects of our
// element and its capabilities with the type system.

// Every element needs to provide its own constructor that returns an initialized
// gst.GoElement implementation. Here we simply create a new fileSrc with zeroed settings
// and state objects.
func (f *fileSrc) New() gst.GoElement {
	return &fileSrc{
		settings: &settings{},
		state:    &state{},
	}
}

// The TypeInit method should register any additional interfaces provided by the element.
// In this example we signal to the type system that we also implement the GstURIHandler interface.
func (f *fileSrc) TypeInit(instance *gst.TypeInstance) {
	instance.AddInterface(gst.InterfaceURIHandler)
}

// The ClassInit method should specify the metadata for this element and add any pad templates
// and properties.
func (f *fileSrc) ClassInit(klass *gst.ElementClass) {
	klass.SetMetadata(
		"File Source",
		"Source/File",
		"Read stream from a file",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
	caps := gst.NewAnyCaps()
	srcPadTemplate := gst.NewPadTemplate(
		"src",
		gst.PadDirectionSource,
		gst.PadPresenceAlways,
		caps,
	)
	klass.AddPadTemplate(srcPadTemplate)
	klass.InstallProperties(properties)
}

// Object implementations are used during the initialization of element. The
// methods are called once the obejct is constructed and its properties are read
// and written to.

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
			gst.ToElement(self).Error(gst.DomainLibrary, gst.LibraryErrorSettings,
				"Could not set location on object",
				err.Error(),
			)
		}
		gst.ToElement(self).Info(gst.DomainLibrary, fmt.Sprintf("Set location to %s", f.settings.location))
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
		gst.ToElement(self).Error(gst.DomainLibrary, gst.LibraryErrorSettings,
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
	stat, err := f.state.file.Stat()
	if err != nil {
		// This should never happen
		self.Error(gst.DomainResource, gst.ResourceErrorFailed,
			"Could not retrieve fileinfo on opened file",
			err.Error(),
		)
		return false, 0
	}
	return true, stat.Size()
}

// Start is called to start this element. In this example, the configured file is opened for reading,
// and any error encountered in the process is posted to the pipeline.
func (f *fileSrc) Start(self *base.GstBaseSrc) bool {
	if f.state.started {
		self.Error(gst.DomainResource, gst.ResourceErrorSettings, "FileSrc is already started", "")
		return false
	}

	if f.settings.location == "" {
		self.Error(gst.DomainResource, gst.ResourceErrorSettings, "File location is not defined", "")
		return false
	}

	var err error
	f.state.file, err = os.OpenFile(f.settings.location, syscall.O_RDONLY, 0444)
	if err != nil {
		self.Error(gst.DomainResource, gst.ResourceErrorOpenRead,
			fmt.Sprintf("Could not open file %s for reading", f.settings.location), err.Error())
		return false
	}
	f.state.position = 0

	f.state.started = true

	self.StartComplete(gst.FlowOK)

	self.Info(gst.DomainResource, "Started")
	return true
}

// Stop is called to stop the element. The file is closed and the local values are zeroed out.
func (f *fileSrc) Stop(self *base.GstBaseSrc) bool {
	if !f.state.started {
		self.Error(gst.DomainResource, gst.ResourceErrorSettings, "FileSrc is not started", "")
		return false
	}

	if err := f.state.file.Close(); err != nil {
		self.Error(gst.DomainResource, gst.ResourceErrorClose, "Failed to close the source file", err.Error())
		return false
	}

	f.state.file = nil
	f.state.position = 0
	f.state.started = false

	self.Info(gst.DomainResource, "Stopped")
	return true
}

// Fill is called to fill a pre-allocated buffer with the data at offset to the given size.
// Since we declared that we are seekable, we need to support the provided offset not neccesarily matching
// where we currently are in the file. This is why we store the position in the file locally.
func (f *fileSrc) Fill(self *base.GstBaseSrc, offset uint64, size uint, buffer *gst.Buffer) gst.FlowReturn {
	if !f.state.started || f.state.file == nil {
		self.Error(gst.DomainCore, gst.CoreErrorFailed, "Not started yet", "")
		return gst.FlowError
	}

	if f.state.position != offset {
		if _, err := f.state.file.Seek(int64(offset), 0); err != nil {
			self.Error(gst.DomainResource, gst.ResourceErrorSeek,
				fmt.Sprintf("Failed to seek to %d in file", offset), err.Error())
			return gst.FlowError
		}
	}

	out := make([]byte, int(size))
	if _, err := f.state.file.Read(out); err != nil && err != io.EOF {
		self.Error(gst.DomainResource, gst.ResourceErrorRead,
			fmt.Sprintf("Failed to read %d bytes from file at %d", size, offset), err.Error())
		return gst.FlowError
	}

	f.state.position = f.state.position + uint64(size)

	bufmap := buffer.Map(gst.MapWrite)
	if bufmap == nil {
		self.Error(gst.DomainLibrary, gst.LibraryErrorFailed, "Failed to map buffer", "")
		return gst.FlowError
	}
	defer buffer.Unmap()

	bufmap.WriteData(out)
	buffer.SetSize(int64(size))

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
	return true, nil
}
