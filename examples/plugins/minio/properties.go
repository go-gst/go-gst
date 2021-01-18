package main

import (
	"github.com/tinyzimmer/go-glib/glib"
)

// Even though there is overlap in properties, they have to be declared twice.
// This is because the GType system doesn't allow for GObjects to share pointers
// to the exact same GParamSpecs.

var sinkProperties = []*glib.ParamSpec{
	glib.NewStringParam(
		"endpoint",
		"S3 API Endpoint",
		"The endpoint for the S3 API server",
		&defaultEndpoint,
		glib.ParameterReadWrite,
	),
	glib.NewBoolParam(
		"use-tls",
		"Use TLS",
		"Use HTTPS for API requests",
		defaultUseTLS,
		glib.ParameterReadWrite,
	),
	glib.NewBoolParam(
		"tls-skip-verify",
		"Disable TLS Verification",
		"Don't verify the signature of the MinIO server certificate",
		defaultInsecureSkipVerify,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"ca-cert-file",
		"PEM CA Cert Bundle",
		"A file containing a PEM certificate bundle to use to verify the MinIO certificate",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"region",
		"Bucket region",
		"The region where the bucket is",
		&defaultRegion,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"bucket",
		"Bucket name",
		"The name of the MinIO bucket",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"key",
		"Object key",
		"The key of the object inside the bucket",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"access-key-id",
		"Access Key ID",
		"The access key ID to use for authentication",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"secret-access-key",
		"Secret Access Key",
		"The secret access key to use for authentication",
		nil,
		glib.ParameterReadWrite,
	),
}

var srcProperties = []*glib.ParamSpec{
	glib.NewStringParam(
		"endpoint",
		"S3 API Endpoint",
		"The endpoint for the S3 API server",
		&defaultEndpoint,
		glib.ParameterReadWrite,
	),
	glib.NewBoolParam(
		"use-tls",
		"Use TLS",
		"Use HTTPS for API requests",
		defaultUseTLS,
		glib.ParameterReadWrite,
	),
	glib.NewBoolParam(
		"tls-skip-verify",
		"Disable TLS Verification",
		"Don't verify the signature of the MinIO server certificate",
		defaultInsecureSkipVerify,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"ca-cert-file",
		"PEM CA Cert Bundle",
		"A file containing a PEM certificate bundle to use to verify the MinIO certificate",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"region",
		"Bucket region",
		"The region where the bucket is",
		&defaultRegion,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"bucket",
		"Bucket name",
		"The name of the MinIO bucket",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"key",
		"Object key",
		"The key of the object inside the bucket",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"access-key-id",
		"Access Key ID",
		"The access key ID to use for authentication. Use env: prefix to denote an environment variable.",
		nil,
		glib.ParameterReadWrite,
	),
	glib.NewStringParam(
		"secret-access-key",
		"Secret Access Key",
		"The secret access key to use for authentication. Use env: prefix to denote an environment variable.",
		nil,
		glib.ParameterReadWrite,
	),
}
