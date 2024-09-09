// This example demonstrates a filesink plugin implemented in Go.
//
// Every element in a Gstreamer pipeline is provided by plugins. Some are builtin while
// others are provided by third-parties or distributed privately. The plugins are built
// around the GObject type system.
//
// Go-gst offers loose bindings around the GObject type system to provide the necessary
// functionality to implement these plugins. The example in this code produces an element
// that can write to a file on the local system.
//
// In order to build the plugin for use by GStreamer, you can do the following:
//
//	$ go generate
//	$ go build -o libgstgofilesink.so -buildmode c-shared .
//
// +plugin:Name=gofilesink
// +plugin:Description=File plugins written in go
// +plugin:Version=v0.0.1
// +plugin:License=gst.LicenseLGPL
// +plugin:Source=go-gst
// +plugin:Package=examples
// +plugin:Origin=https://github.com/go-gst/go-gst
// +plugin:ReleaseDate=2021-01-04
//
// +element:Name=gofilesink
// +element:Rank=gst.RankNone
// +element:Impl=FileSink
// +element:Subclass=base.ExtendsBaseSink
// +element:Interfaces=gst.InterfaceURIHandler
//
//go:generate gst-plugin-gen
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/base"
)

// main is left unimplemented since these files are compiled to c-shared.
func main() {}

// CAT is the log category for the gofilesink. It is safe to define GStreamer objects as globals
// without calling gst.Init, since in the context of a loaded plugin all initialization has
// already been taken care of by the loading application.
var CAT = gst.NewDebugCategory(
	"gofilesink",
	gst.DebugColorNone,
	"GoFileSink Element",
)

// Here we define a list of ParamSpecs that will make up the properties for our element.
// This element only has a single property, the location of the file to write to.
// When getting and setting properties later on, you will reference them by their index in
// this list.
var properties = []*glib.ParamSpec{
	glib.NewStringParam(
		"location",                      // The name of the parameter
		"File Location",                 // The long name for the parameter
		"Location to write the file to", // A blurb about the parameter
		nil,                             // A default value for the parameter
		glib.ParameterReadWrite,         // Flags for the parameter
	),
}

// Here we declare a private struct to hold our internal state.
type state struct {
	// Whether the element is started or not
	started bool
	// The file the element is writing to
	file *os.File
	// The current position in the file
	position uint64
}

// This is another private struct where we hold the parameter values set on our
// element.
type settings struct {
	location string
}

// Finally a structure is defined that implements (at a minimum) the glib.GoObject interface.
// It is possible to signal to the bindings to inherit from other classes or implement other
// interfaces via the registration and TypeInit processes.
type FileSink struct {
	// The settings for the element
	settings *settings
	// The current state of the element
	state *state
}

// setLocation is a simple method to check the validity of a provided file path and set the
// local value with it.
func (f *FileSink) setLocation(path string) error {
	if f.state.started {
		return errors.New("changing the `location` property on a started `GoFileSink` is not supported")
	}
	f.settings.location = strings.TrimPrefix(path, "file://") // should obviously use url.URL and do actual parsing
	return nil
}

// The ObjectSubclass implementations below are for registering the various aspects of our
// element and its capabilities with the type system. These are the minimum methods that
// should be implemented by an element.

// Every element needs to provide its own constructor that returns an initialized glib.GoObjectSubclass
// implementation. Here we simply create a new fileSink with zeroed settings and state objects.
func (f *FileSink) New() glib.GoObjectSubclass {
	CAT.Log(gst.LevelLog, "Initializing new fileSink object")
	return &FileSink{
		settings: &settings{},
		state:    &state{},
	}
}

// The ClassInit method should specify the metadata for this element and add any pad templates
// and properties.
func (f *FileSink) ClassInit(klass *glib.ObjectClass) {
	CAT.Log(gst.LevelLog, "Initializing gofilesink class")
	class := gst.ToElementClass(klass)
	class.SetMetadata(
		"File Sink",
		"Sink/File",
		"Write stream to a file",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
	CAT.Log(gst.LevelLog, "Adding sink pad template and properties to class")
	class.AddPadTemplate(gst.NewPadTemplate(
		"sink",
		gst.PadDirectionSink,
		gst.PadPresenceAlways,
		gst.NewAnyCaps(),
	))
	class.InstallProperties(properties)
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
func (f *FileSink) SetProperty(self *glib.Object, id uint, value *glib.Value) {
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
		gst.ToElement(self).Log(CAT, gst.LevelInfo, fmt.Sprintf("Set `location` to %s", f.settings.location))
	}
}

// GetProperty is called to retrieve the value of the property at index `id` in the properties
// slice provided at ClassInit.
func (f *FileSink) GetProperty(self *glib.Object, id uint) *glib.Value {
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

// GstBaseSink implementations are optional methods to implement from the base.GstBaseSinkImpl interface.
// If the method is not overridden by the implementing struct, it will be inherited from the parent class.

// Start is called to start the filesink. Open the file for writing and set the internal state.
func (f *FileSink) Start(self *base.GstBaseSink) bool {
	if f.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "GoFileSink is already started", "")
		return false
	}

	if f.settings.location == "" {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "No location configured on the filesink", "")
		return false
	}

	destFile := f.settings.location

	var err error
	f.state.file, err = os.Create(destFile)
	if err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorOpenWrite,
			fmt.Sprintf("Could not open %s for writing", destFile), err.Error())
		return false
	}

	self.Log(CAT, gst.LevelDebug, fmt.Sprintf("Opened file %s for writing", destFile))

	f.state.started = true
	self.Log(CAT, gst.LevelInfo, "GoFileSink has started")
	return true
}

// Stop is called to stop the element. Set the internal state and close the file.
func (f *FileSink) Stop(self *base.GstBaseSink) bool {
	if !f.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "GoFileSink is not started", "")
		return false
	}

	if err := f.state.file.Close(); err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorWrite, "Failed to close the destination file", err.Error())
		return false
	}

	self.Log(CAT, gst.LevelInfo, "GoFileSink has stopped")
	return true
}

// Render is called when a buffer is ready to be written to the file.
func (f *FileSink) Render(self *base.GstBaseSink, buffer *gst.Buffer) gst.FlowReturn {
	if !f.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "GoFileSink is not started", "")
		return gst.FlowError
	}

	self.Log(CAT, gst.LevelTrace, fmt.Sprintf("Rendering buffer at %v", buffer.Instance()))
	newPos, err := io.Copy(f.state.file, buffer.Reader())
	if err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorWrite, "Error copying buffer to file", err.Error())
		return gst.FlowError
	}

	f.state.position += uint64(newPos)
	self.Log(CAT, gst.LevelTrace, fmt.Sprintf("New position in file: %v", f.state.position))

	return gst.FlowOK
}

// URIHandler implementations are the methods required by the GstURIHandler interface.

// GetURI returns the currently configured URI
func (f *FileSink) GetURI() string { return fmt.Sprintf("file://%s", f.settings.location) }

// GetURIType returns the types of URI this element supports.
func (f *FileSink) GetURIType() gst.URIType { return gst.URISource }

// GetProtocols returns the protcols this element supports.
func (f *FileSink) GetProtocols() []string { return []string{"file"} }

// SetURI should set the URI that this element is working on.
func (f *FileSink) SetURI(uri string) (bool, error) {
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
