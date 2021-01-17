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
	glib.NewInt64Param(
		"chunk-size",
		"Chunk size",
		"The size of each chunk uploaded to S3/MinIO in bytes. The final output cannot exceed this value * 10000. Note that buffers will be held in memory until they reach this size.",
		defaultChunkSize, 1024*1024*1024, defaultChunkSize, // Min: 5MB  Max: 1GB  Default: 5MB
		glib.ParameterReadWrite,
	),
	glib.NewUintParam(
		"max-memory-chunks",
		"Maximum Chunks in Memory",
		"The maximum number of chunks to keep in memory at any given point in time. Setting this to a higher value will reduce load on the destination at the expense of increased memory consumption. There needs to be room for the head chunk and at least one more at all times.",
		2, ^uint(0), defaultMaxMemChunks,
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
