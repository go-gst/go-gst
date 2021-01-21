package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

const (
	accessKeyIDEnvVar     = "MINIO_ACCESS_KEY_ID"
	secretAccessKeyEnvVar = "MINIO_SECRET_ACCESS_KEY"
)

var (
	defaultEndpoint           = "play.min.io"
	defaultUseTLS             = true
	defaultRegion             = "us-east-1"
	defaultInsecureSkipVerify = false
)

type settings struct {
	endpoint           string
	useTLS             bool
	region             string
	bucket             string
	key                string
	accessKeyID        string
	secretAccessKey    string
	insecureSkipVerify bool
	caCertFile         string
	partSize           uint64
}

func (s *settings) safestring() string {
	return fmt.Sprintf("%+v", &settings{
		endpoint:           s.endpoint,
		useTLS:             s.useTLS,
		region:             s.region,
		bucket:             s.bucket,
		key:                s.key,
		insecureSkipVerify: s.insecureSkipVerify,
		caCertFile:         s.caCertFile,
	})
}

func defaultSettings() *settings {
	return &settings{
		endpoint:           defaultEndpoint,
		useTLS:             defaultUseTLS,
		region:             defaultRegion,
		accessKeyID:        os.Getenv(accessKeyIDEnvVar),
		secretAccessKey:    os.Getenv(secretAccessKeyEnvVar),
		insecureSkipVerify: defaultInsecureSkipVerify,
		partSize:           defaultPartSize,
	}
}

func getMinIOClient(settings *settings) (*minio.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	if settings.useTLS {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		if settings.caCertFile != "" {
			certPool := x509.NewCertPool()
			body, err := ioutil.ReadFile(settings.caCertFile)
			if err != nil {
				return nil, err
			}
			certPool.AppendCertsFromPEM(body)
			transport.TLSClientConfig.RootCAs = certPool
		}
		transport.TLSClientConfig.InsecureSkipVerify = settings.insecureSkipVerify
	}
	return minio.New(settings.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.accessKeyID, settings.secretAccessKey, ""),
		Secure: settings.useTLS,
		Region: settings.region,
	})
}

func setProperty(elem *gst.Element, properties []*glib.ParamSpec, settings *settings, id uint, value *glib.Value) {
	prop := properties[id]

	val, err := value.GoValue()
	if err != nil {
		elem.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorSettings,
			fmt.Sprintf("Could not coerce %v to go value", value), err.Error())
	}

	switch prop.Name() {
	case "endpoint":
		settings.endpoint = val.(string)
	case "use-tls":
		settings.useTLS = val.(bool)
	case "tls-skip-verify":
		settings.insecureSkipVerify = val.(bool)
	case "ca-cert-file":
		settings.caCertFile = val.(string)
	case "region":
		settings.region = val.(string)
	case "bucket":
		settings.bucket = val.(string)
	case "key":
		settings.key = val.(string)
	case "access-key-id":
		settings.accessKeyID = val.(string)
	case "secret-access-key":
		settings.secretAccessKey = val.(string)
	case "part-size":
		settings.partSize = val.(uint64)
	}
}

func getProperty(elem *gst.Element, properties []*glib.ParamSpec, settings *settings, id uint) *glib.Value {
	prop := properties[id]

	var localVal interface{}

	switch prop.Name() {
	case "endpoint":
		localVal = settings.endpoint
	case "use-tls":
		localVal = settings.useTLS
	case "tls-skip-verify":
		localVal = settings.insecureSkipVerify
	case "ca-cert-file":
		localVal = settings.caCertFile
	case "region":
		localVal = settings.region
	case "bucket":
		localVal = settings.bucket
	case "key":
		localVal = settings.key
	case "access-key-id":
		localVal = settings.accessKeyID
	case "secret-access-key":
		localVal = "<private>"
	case "part-size":
		localVal = settings.partSize

	default:
		elem.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorSettings,
			fmt.Sprintf("Cannot get invalid property %s", prop.Name()), "")
		return nil
	}

	val, err := glib.GValue(localVal)
	if err != nil {
		elem.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed,
			fmt.Sprintf("Could not convert %v to GValue", localVal),
			err.Error(),
		)
		return nil
	}

	return val
}
