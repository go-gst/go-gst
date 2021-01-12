package main

import (
	"os"
)

const (
	accessKeyIDEnvVar     = "MINIO_ACCESS_KEY_ID"
	secretAccessKeyEnvVar = "MINIO_SECRET_ACCESS_KEY"
)

var (
	defaultEndpoint = "play.min.io"
	defaultUseTLS   = true
	defaultRegion   = "us-east-1"
)

type settings struct {
	endpoint        string
	useTLS          bool
	region          string
	bucket          string
	key             string
	accessKeyID     string
	secretAccessKey string
}

func defaultSettings() *settings {
	return &settings{
		endpoint:        defaultEndpoint,
		useTLS:          defaultUseTLS,
		region:          defaultRegion,
		accessKeyID:     os.Getenv(accessKeyIDEnvVar),
		secretAccessKey: os.Getenv(secretAccessKeyEnvVar),
	}
}
