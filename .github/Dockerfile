# The intention of this image is to provide an up-to-date installation of gstreamer for
# CI. However, it can also be used as an image in multi-stage builds using this library.
# 
# Compile the binary from this image, and then copy it into a fresh alpine image with
# just the needed libraries installed. For example:
#
#   FROM ghcr.io/tinyzimmer/go-gst:1.15 as builder
#   COPY src /workspace/src
#   RUN go build -o /workspace/app /workspace/src
#
#   FROM ubuntu
#   RUN apt-get update && apt-get install -y gstreamer gst-plugins-good
#   COPY --from=builder /workspace/app /app
#   ENTRYPOINT ["/app"]
#
ARG GO_VERSION=1.15
FROM ubuntu:20.10

RUN mkdir -p /build \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
        golang git make curl \
        libgstreamer1.0 libgstreamer1.0-dev \
        libgstreamer-plugins-bad1.0-dev libgstreamer-plugins-base1.0-dev \
        gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad \
        gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools