package gst

// #include "gst.go.h"
import "C"
import "sync"

var registerMutex sync.RWMutex

var registeredTypes = make(map[string]C.GType)
var registeredClasses = make(map[C.gpointer]GoElement)
var globalURIHdlr URIHandler
