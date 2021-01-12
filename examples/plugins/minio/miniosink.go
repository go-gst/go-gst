package main

import (
	"context"
	"fmt"
	"io"
	"sync"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/base"
)

var sinkCAT = gst.NewDebugCategory(
	"miniosink",
	gst.DebugColorNone,
	"MinIOSink Element",
)

var sinkProperties = []*gst.ParamSpec{
	gst.NewStringParam(
		"endpoint",
		"S3 API Endpoint",
		"The endpoint for the S3 API server",
		&defaultEndpoint,
		gst.ParameterReadWrite,
	),
	gst.NewBoolParam(
		"use-tls",
		"Use TLS",
		"Use HTTPS for API requests",
		defaultUseTLS,
		gst.ParameterReadWrite,
	),
	gst.NewStringParam(
		"region",
		"Bucket region",
		"The region where the bucket is",
		&defaultRegion,
		gst.ParameterReadWrite,
	),
	gst.NewStringParam(
		"bucket",
		"Bucket name",
		"The name of the MinIO bucket",
		nil,
		gst.ParameterReadWrite,
	),
	gst.NewStringParam(
		"key",
		"Object key",
		"The key of the object inside the bucket",
		nil,
		gst.ParameterReadWrite,
	),
	gst.NewStringParam(
		"access-key-id",
		"Access Key ID",
		"The access key ID to use for authentication",
		nil,
		gst.ParameterReadWrite,
	),
	gst.NewStringParam(
		"secret-access-key",
		"Secret Access Key",
		"The secret access key to use for authentication",
		nil,
		gst.ParameterReadWrite,
	),
}

type minioSink struct {
	settings *settings
	state    *sinkstate
}

type sinkstate struct {
	started  bool
	rPipe    io.ReadCloser
	wPipe    io.WriteCloser
	doneChan chan struct{}

	mux sync.Mutex
}

func (m *minioSink) New() gst.GoElement {
	srcCAT.Log(gst.LevelLog, "Creating new minioSink object")
	return &minioSink{
		settings: defaultSettings(),
		state:    &sinkstate{},
	}
}

func (m *minioSink) TypeInit(*gst.TypeInstance) {}

func (m *minioSink) ClassInit(klass *gst.ElementClass) {
	sinkCAT.Log(gst.LevelLog, "Initializing miniosink class")
	klass.SetMetadata(
		"MinIO Sink",
		"Sink/File",
		"Write stream to a MinIO object",
		"Avi Zimmerman <avi.zimmerman@gmail.com>",
	)
	sinkCAT.Log(gst.LevelLog, "Adding sink pad template and properties to class")
	klass.AddPadTemplate(gst.NewPadTemplate(
		"sink",
		gst.PadDirectionSink,
		gst.PadPresenceAlways,
		gst.NewAnyCaps(),
	))
	klass.InstallProperties(sinkProperties)
}

func (m *minioSink) SetProperty(self *gst.Object, id uint, value *glib.Value) {
	prop := sinkProperties[id]

	val, err := value.GoValue()
	if err != nil {
		gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorSettings,
			fmt.Sprintf("Could not coerce %v to go value", value), err.Error())
	}

	switch prop.Name() {
	case "endpoint":
		m.settings.endpoint = val.(string)
	case "use-tls":
		m.settings.useTLS = val.(bool)
	case "region":
		m.settings.region = val.(string)
	case "bucket":
		m.settings.bucket = val.(string)
	case "key":
		m.settings.key = val.(string)
	case "access-key-id":
		m.settings.accessKeyID = val.(string)
	case "secret-access-key":
		m.settings.secretAccessKey = val.(string)
	}

}

func (m *minioSink) GetProperty(self *gst.Object, id uint) *glib.Value {
	prop := sinkProperties[id]

	var localVal interface{}

	switch prop.Name() {
	case "endpoint":
		localVal = m.settings.endpoint
	case "use-tls":
		localVal = m.settings.useTLS
	case "region":
		localVal = m.settings.region
	case "bucket":
		localVal = m.settings.bucket
	case "key":
		localVal = m.settings.key
	case "access-key-id":
		localVal = m.settings.accessKeyID
	case "secret-access-key":
		localVal = m.settings.secretAccessKey

	default:
		gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorSettings,
			fmt.Sprintf("Cannot get invalid property %s", prop.Name()), "")
		return nil
	}

	val, err := glib.GValue(localVal)
	if err != nil {
		gst.ToElement(self).ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
			fmt.Sprintf("Could not convert %v to GValue", localVal),
			err.Error(),
		)
	}

	return val
}

func (m *minioSink) Start(self *base.GstBaseSink) bool {
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

	m.state.doneChan = make(chan struct{})

	client, err := minio.New(m.settings.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.settings.accessKeyID, m.settings.secretAccessKey, ""),
		Secure: m.settings.useTLS,
	})

	if err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed,
			fmt.Sprintf("Failed to connect to MinIO endpoint %s", m.settings.endpoint), err.Error())
		return false
	}

	m.state.rPipe, m.state.wPipe = io.Pipe()

	go func() {
		self.Log(sinkCAT, gst.LevelInfo,
			fmt.Sprintf("Starting PutObject operation to %s/%s/%s", m.settings.endpoint, m.settings.bucket, m.settings.key),
		)
		if _, err := client.PutObject(context.Background(),
			m.settings.bucket, m.settings.key,
			m.state.rPipe, -1,
			minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			}); err != nil {
			self.ErrorMessage(gst.DomainResource, gst.ResourceErrorFailed,
				fmt.Sprintf("Error during PutObject call to %s/%s", m.settings.bucket, m.settings.key), err.Error())
		}
		m.state.doneChan <- struct{}{}
	}()

	m.state.started = true
	self.Log(sinkCAT, gst.LevelInfo, "MinIOSink has started")
	return true
}

func (m *minioSink) Stop(self *base.GstBaseSink) bool {
	if !m.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "MinIOSink is not started", "")
		return false
	}
	if err := m.state.wPipe.Close(); err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorWrite, "Failed to close the write pipe", err.Error())
		return false
	}
	self.Log(sinkCAT, gst.LevelInfo, "Waiting for PutObject operation to complete")
	<-m.state.doneChan
	self.Log(sinkCAT, gst.LevelInfo, "MinIOSink has stopped")
	return true
}

func (m *minioSink) Render(self *base.GstBaseSink, buffer *gst.Buffer) gst.FlowReturn {
	if !m.state.started {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorSettings, "MinIOSink is not started", "")
		return gst.FlowError
	}

	self.Log(sinkCAT, gst.LevelLog, fmt.Sprintf("Rendering buffer %v", buffer))
	if _, err := io.Copy(m.state.wPipe, buffer.Reader()); err != nil {
		self.ErrorMessage(gst.DomainResource, gst.ResourceErrorWrite, "Error copying buffer to write pipe", err.Error())
		return gst.FlowError
	}

	return gst.FlowOK
}
