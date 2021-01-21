package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
)

var sinkCAT = gst.NewDebugCategory(
	"miniosink",
	gst.DebugColorNone,
	"MinIOSink Element",
)

type minioSink struct {
	settings *settings
	state    *sinkstate

	writer *seekWriter
	mux    sync.Mutex
}

type sinkstate struct {
	started bool
}

func (m *minioSink) New() glib.GoObjectSubclass {
	srcCAT.Log(gst.LevelLog, "Creating new minioSink object")
	return &minioSink{
		settings: defaultSettings(),
		state:    &sinkstate{},
	}
}

func (m *minioSink) ClassInit(klass *glib.ObjectClass) {
	class := gst.ToElementClass(klass)
	sinkCAT.Log(gst.LevelLog, "Initializing miniosink class")
	class.SetMetadata(
		"MinIO Sink",
		"Sink/File",
		"Write stream to a MinIO object",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
	sinkCAT.Log(gst.LevelLog, "Adding sink pad template and properties to class")
	class.AddPadTemplate(gst.NewPadTemplate(
		"sink",
		gst.PadDirectionSink,
		gst.PadPresenceAlways,
		gst.NewAnyCaps(),
	))
	class.InstallProperties(sinkProperties)
}

func (m *minioSink) Constructed(obj *glib.Object) { base.ToGstBaseSink(obj).SetSync(false) }

func (m *minioSink) SetProperty(self *glib.Object, id uint, value *glib.Value) {
	setProperty(gst.ToElement(self), sinkProperties, m.settings, id, value)
}

func (m *minioSink) GetProperty(self *glib.Object, id uint) *glib.Value {
	return getProperty(gst.ToElement(self), sinkProperties, m.settings, id)
}

func (m *minioSink) Query(self *base.GstBaseSink, query *gst.Query) bool {
	switch query.Type() {

	case gst.QuerySeeking:
		self.Log(sinkCAT, gst.LevelDebug, "Answering seeking query")
		query.SetSeeking(gst.FormatTime, true, 0, -1)
		return true

	}
	return false
}

func (m *minioSink) Start(self *base.GstBaseSink) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings,
			"MinIOSink is already started", "")
		return false
	}

	if m.settings.bucket == "" {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings,
			"No bucket configured on the miniosink", "")
		return false
	}

	if m.settings.key == "" {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings,
			"No bucket configured on the miniosink", "")
		return false
	}

	self.Log(sinkCAT, gst.LevelDebug, m.settings.safestring())

	if strings.HasPrefix(m.settings.accessKeyID, "env:") {
		spl := strings.Split(m.settings.accessKeyID, "env:")
		m.settings.accessKeyID = os.Getenv(spl[len(spl)-1])
	}

	if strings.HasPrefix(m.settings.secretAccessKey, "env:") {
		spl := strings.Split(m.settings.secretAccessKey, "env:")
		m.settings.secretAccessKey = os.Getenv(spl[len(spl)-1])
	}

	self.Log(sinkCAT, gst.LevelInfo, fmt.Sprintf("Creating new MinIO client for %s", m.settings.endpoint))
	client, err := getMinIOClient(m.settings)
	if err != nil {
		self.Log(sinkCAT, gst.LevelError, err.Error())
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed,
			fmt.Sprintf("Failed to connect to MinIO endpoint %s", m.settings.endpoint), err.Error())
		return false
	}

	self.Log(sinkCAT, gst.LevelInfo, "Initializing new MinIO writer")
	m.writer = newSeekWriter(client, int64(m.settings.partSize), m.settings.bucket, m.settings.key)

	m.state.started = true
	self.Log(sinkCAT, gst.LevelInfo, "MinIOSink has started")
	return true
}

func (m *minioSink) Stop(self *base.GstBaseSink) bool {
	self.Log(sinkCAT, gst.LevelInfo, "Stopping MinIOSink")
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "MinIOSink is not started", "")
		return false
	}

	m.writer = nil
	m.state.started = false

	self.Log(sinkCAT, gst.LevelInfo, "MinIOSink has stopped")
	return true
}

func (m *minioSink) Render(self *base.GstBaseSink, buffer *gst.Buffer) gst.FlowReturn {
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "MinIOSink is not started", "")
		return gst.FlowError
	}

	self.Log(sinkCAT, gst.LevelTrace, fmt.Sprintf("Rendering buffer %v", buffer))

	if _, err := m.writer.Write(buffer.Bytes()); err != nil {
		self.Log(sinkCAT, gst.LevelError, err.Error())
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorWrite, fmt.Sprintf("Failed to write data to minio buffer: %s", err.Error()), "")
		return gst.FlowError
	}

	return gst.FlowOK
}

func (m *minioSink) Event(self *base.GstBaseSink, event *gst.Event) bool {

	switch event.Type() {

	case gst.EventTypeSegment:
		segment := event.ParseSegment()

		if segment.GetFormat() == gst.FormatBytes {
			if uint64(m.writer.currentPosition) != segment.GetStart() {
				m.mux.Lock()
				self.Log(sinkCAT, gst.LevelInfo, fmt.Sprintf("Seeking to %d", segment.GetStart()))
				if _, err := m.writer.Seek(int64(segment.GetStart()), io.SeekStart); err != nil {
					self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed, err.Error(), "")
					m.mux.Unlock()
					return false
				}
				m.mux.Unlock()
			} else {
				self.Log(sinkCAT, gst.LevelDebug, "Ignored SEGMENT, no seek needed")
			}
		} else {
			self.Log(sinkCAT, gst.LevelDebug, fmt.Sprintf("Ignored SEGMENT event of format %s", segment.GetFormat().String()))
		}

	case gst.EventTypeFlushStop:
		self.Log(sinkCAT, gst.LevelInfo, "Flushing contents of writer and seeking back to start")
		if m.writer.currentPosition != 0 {
			m.mux.Lock()
			if err := m.writer.flush(true); err != nil {
				self.ErrorMessage(gst.DomainResource, gst.ResourceErrorWrite, err.Error(), "")
				m.mux.Unlock()
				return false
			}
			if _, err := m.writer.Seek(0, io.SeekStart); err != nil {
				self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed, err.Error(), "")
				m.mux.Unlock()
				return false
			}
			m.mux.Unlock()
		}

	case gst.EventTypeEOS:
		self.Log(sinkCAT, gst.LevelInfo, "Received EOS, closing MinIO writer")
		m.mux.Lock()
		if err := m.writer.Close(); err != nil {
			self.Log(sinkCAT, gst.LevelError, err.Error())
			self.ErrorMessage(gst.DomainResource, gst.ResourceErrorClose, fmt.Sprintf("Failed to close MinIO writer: %s", err.Error()), "")
			m.mux.Unlock()
			return false
		}
		m.mux.Unlock()
	}

	return self.ParentEvent(event)

}
