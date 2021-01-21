package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"path"

	minio "github.com/minio/minio-go/v7"
)

type seekWriter struct {
	// The current position in the buffer
	currentPosition int64
	// The size of each part to upload
	partSize int64
	// A map of in memory parts to their content
	parts map[int64][]byte
	// A map of uploaded parts to the checksum at time of upload
	uploadedParts map[int64]string
	// A local reference to the minio client
	client      *minio.Client
	bucket, key string
}

func newSeekWriter(client *minio.Client, partsize int64, bucket, key string) *seekWriter {
	return &seekWriter{
		currentPosition: 0,
		partSize:        partsize,
		parts:           make(map[int64][]byte),
		uploadedParts:   make(map[int64]string),
		client:          client,
		bucket:          bucket, key: key,
	}
}

func (s *seekWriter) Write(p []byte) (int, error) {
	wrote, err := s.buffer(0, p)
	if err != nil {
		return wrote, err
	}
	return wrote, s.flush(false)
}

func (s *seekWriter) Seek(offset int64, whence int) (int64, error) {
	// Only needs to support SeekStart
	s.currentPosition = offset
	return s.currentPosition, nil
}

func (s *seekWriter) Close() error {
	if err := s.flush(true); err != nil {
		return err
	}
	if len(s.uploadedParts) == 0 {
		return nil
	}
	opts := make([]minio.CopySrcOptions, len(s.uploadedParts))
	for i := 0; i < len(opts); i++ {
		opts[i] = minio.CopySrcOptions{
			Bucket: s.bucket,
			Object: s.keyForPart(int64(i)),
		}
	}
	_, err := s.client.ComposeObject(context.Background(), minio.CopyDestOptions{
		Bucket: s.bucket,
		Object: s.key,
	}, opts...)

	if err != nil {
		return err
	}
	for _, opt := range opts {
		if err := s.client.RemoveObject(context.Background(), opt.Bucket, opt.Object, minio.RemoveObjectOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (s *seekWriter) buffer(from int, p []byte) (int, error) {
	currentPart := s.currentPosition / s.partSize
	writeat := s.currentPosition % s.partSize
	lenToWrite := int64(len(p))

	var buf []byte
	var ok bool
	if buf, ok = s.parts[currentPart]; !ok {
		if _, ok := s.uploadedParts[currentPart]; !ok {
			s.parts[currentPart] = make([]byte, writeat+lenToWrite)
			buf = s.parts[currentPart]
		} else {
			var err error
			buf, err = s.fetchRemotePart(currentPart)
			if err != nil {
				return from, err
			}
		}
	}

	if lenToWrite+writeat > s.partSize {
		newbuf := make([]byte, s.partSize)
		copy(newbuf, buf)
		s.parts[currentPart] = newbuf
		buf = newbuf
	} else if lenToWrite+writeat > int64(len(buf)) {
		newbuf := make([]byte, lenToWrite+writeat)
		copy(newbuf, buf)
		s.parts[currentPart] = newbuf
		buf = newbuf
	}

	wrote := copy(buf[writeat:], p)

	s.currentPosition += int64(wrote)

	if int64(wrote) != lenToWrite {
		return s.buffer(from+wrote, p[wrote:])
	}

	return from + wrote, nil
}

func (s *seekWriter) flush(all bool) error {
	for part, buf := range s.parts {
		if all || int64(len(buf)) == s.partSize {
			if err := s.uploadPart(part, buf); err != nil {
				return err
			}
			continue
		}
		if !all {
			continue
		}
		if err := s.uploadPart(part, buf); err != nil {
			return err
		}
	}
	return nil
}

func (s *seekWriter) uploadPart(part int64, data []byte) error {
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		return err
	}
	datasum := fmt.Sprintf("%x", h.Sum(nil))
	if sum, ok := s.uploadedParts[part]; ok && sum == datasum {
		return nil
	}
	_, err := s.client.PutObject(context.Background(),
		s.bucket, s.keyForPart(part),
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		return err
	}
	delete(s.parts, part)
	s.uploadedParts[part] = datasum
	return nil
}

func (s *seekWriter) fetchRemotePart(part int64) ([]byte, error) {
	object, err := s.client.GetObject(context.Background(), s.bucket, s.keyForPart(part), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}
	s.parts[part] = body
	return body, nil
}

func (s *seekWriter) keyForPart(part int64) string {
	if path.Dir(s.key) == "" {
		return fmt.Sprintf("%s_tmp/%d", s.key, part)
	}
	return path.Join(
		path.Dir(s.key),
		fmt.Sprintf("%s_tmp/%d", path.Base(s.key), part),
	)
}
