package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	minio "github.com/minio/minio-go/v7"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
)

var srcCAT = gst.NewDebugCategory(
	"miniosrc",
	gst.DebugColorNone,
	"MinIOSrc Element",
)

type minioSrc struct {
	settings *settings
	state    *srcstate
}

type srcstate struct {
	started bool
	object  *minio.Object
	objInfo minio.ObjectInfo

	mux sync.Mutex
}

func (m *minioSrc) New() glib.GoObjectSubclass {
	srcCAT.Log(gst.LevelLog, "Creating new minioSrc object")
	return &minioSrc{
		settings: defaultSettings(),
		state:    &srcstate{},
	}
}

func (m *minioSrc) ClassInit(klass *glib.ObjectClass) {
	class := gst.ToElementClass(klass)
	srcCAT.Log(gst.LevelLog, "Initializing miniosrc class")
	class.SetMetadata(
		"MinIO Source",
		"Source/File",
		"Read stream from a MinIO object",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
	srcCAT.Log(gst.LevelLog, "Adding src pad template and properties to class")
	class.AddPadTemplate(gst.NewPadTemplate(
		"src",
		gst.PadDirectionSource,
		gst.PadPresenceAlways,
		gst.NewAnyCaps(),
	))
	class.InstallProperties(srcProperties)
}

func (m *minioSrc) SetProperty(self *glib.Object, id uint, value *glib.Value) {
	setProperty(gst.ToElement(self), srcProperties, m.settings, id, value)
}

func (m *minioSrc) GetProperty(self *glib.Object, id uint) *glib.Value {
	return getProperty(gst.ToElement(self), srcProperties, m.settings, id)
}

func (m *minioSrc) Constructed(self *glib.Object) {
	base.ToGstBaseSrc(self).Log(srcCAT, gst.LevelLog, "Setting format of GstBaseSrc to bytes")
	base.ToGstBaseSrc(self).SetFormat(gst.FormatBytes)
}

func (m *minioSrc) IsSeekable(*base.GstBaseSrc) bool { return true }

func (m *minioSrc) GetSize(self *base.GstBaseSrc) (bool, int64) {
	if !m.state.started {
		return false, 0
	}
	return true, m.state.objInfo.Size
}

func (m *minioSrc) Start(self *base.GstBaseSrc) bool {

	if m.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed, "MinIOSrc is already started", "")
		return false
	}

	if m.settings.bucket == "" {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed, "No source bucket defined", "")
		return false
	}

	if m.settings.key == "" {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed, "No object key defined", "")
		return false
	}

	m.state.mux.Lock()

	if strings.HasPrefix(m.settings.accessKeyID, "env:") {
		spl := strings.Split(m.settings.accessKeyID, "env:")
		m.settings.accessKeyID = os.Getenv(spl[len(spl)-1])
	}

	if strings.HasPrefix(m.settings.secretAccessKey, "env:") {
		spl := strings.Split(m.settings.secretAccessKey, "env:")
		m.settings.secretAccessKey = os.Getenv(spl[len(spl)-1])
	}

	client, err := getMinIOClient(m.settings)
	if err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed,
			fmt.Sprintf("Failed to connect to MinIO endpoint %s", m.settings.endpoint), err.Error())
		m.state.mux.Unlock()
		return false
	}

	self.Log(srcCAT, gst.LevelInfo, fmt.Sprintf("Requesting %s/%s from %s", m.settings.bucket, m.settings.key, m.settings.endpoint))
	m.state.object, err = client.GetObject(context.Background(), m.settings.bucket, m.settings.key, minio.GetObjectOptions{})
	if err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorOpenRead,
			fmt.Sprintf("Failed to retrieve object %q from bucket %q", m.settings.key, m.settings.bucket), err.Error())
		m.state.mux.Unlock()
		return false
	}

	self.Log(srcCAT, gst.LevelInfo, "Getting HEAD for object")
	m.state.objInfo, err = m.state.object.Stat()
	if err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorOpenRead,
			fmt.Sprintf("Failed to stat object %q in bucket %q: %s", m.settings.key, m.settings.bucket, err.Error()), "")
		m.state.mux.Unlock()
		return false
	}
	self.Log(srcCAT, gst.LevelInfo, fmt.Sprintf("%+v", m.state.objInfo))

	m.state.started = true
	m.state.mux.Unlock()

	self.StartComplete(gst.FlowOK)

	self.Log(srcCAT, gst.LevelInfo, "MinIOSrc has started")
	return true
}

func (m *minioSrc) Stop(self *base.GstBaseSrc) bool {
	self.Log(srcCAT, gst.LevelInfo, "Stopping MinIOSrc")
	m.state.mux.Lock()
	defer m.state.mux.Unlock()

	if !m.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "MinIOSrc is not started", "")
		return false
	}

	if err := m.state.object.Close(); err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorClose, "Failed to close the bucket object", err.Error())
		return false
	}

	m.state.object = nil
	m.state.started = false

	self.Log(srcCAT, gst.LevelInfo, "MinIOSrc has stopped")
	return true
}

func (m *minioSrc) Fill(self *base.GstBaseSrc, offset uint64, size uint, buffer *gst.Buffer) gst.FlowReturn {

	if !m.state.started || m.state.object == nil {
		self.ErrorMessage(gst.DomainCore, gst.CoreErrorFailed, "MinIOSrc is not started yet", "")
		return gst.FlowError
	}

	self.Log(srcCAT, gst.LevelLog, fmt.Sprintf("Request to fill buffer from offset %v with size %v", offset, size))

	m.state.mux.Lock()
	defer m.state.mux.Unlock()

	data := make([]byte, size)
	read, err := m.state.object.ReadAt(data, int64(offset))
	if err != nil && err != io.EOF {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorRead,
			fmt.Sprintf("Failed to read %d bytes from object at offset %d", size, offset), err.Error())
		return gst.FlowError
	}

	if read < int(size) {
		self.Log(srcCAT, gst.LevelDebug, fmt.Sprintf("Only read %d bytes from object, trimming", read))
		trim := make([]byte, read)
		copy(trim, data)
		data = trim
	}

	bufmap := buffer.Map(gst.MapWrite)
	if bufmap == nil {
		self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, "Failed to map buffer", "")
		return gst.FlowError
	}
	defer buffer.Unmap()

	bufmap.WriteData(data)
	buffer.SetSize(int64(read))

	return gst.FlowOK
}
