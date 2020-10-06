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
FROM golang:${GO_VERSION} 

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y cmake \
        build-essential python3 python3-pip glib2.0 glib2.0-dev libssl-dev bison flex liborc-dev ffmpeg libegl1-mesa \
        libavcodec-dev libavfilter-dev libavformat-dev libavresample-dev libavutil-dev libpciaccess-dev libegl1-mesa-dev \
        libpng-tools libpng-dev libjpeg-dev libx264-dev nasm libgtest-dev libnice-dev libdrm-dev gobject-introspection  \
        libgirepository1.0-dev gir1.2-glib-2.0 libcairo-dev gtk3.0-dev liblcms2-dev liblcms2-2 liborc-0.4-dev liborc-0.4 \
    && pip3 install meson ninja \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /gstreamer \
    && cd /gstreamer \
    && git clone https://gitlab.freedesktop.org/gstreamer/gst-build.git \
    && cd gst-build \
    && meson setup build \
    && cd build \
    && meson compile \
    && meson install \ 
    && cd .. \
    && rm -rf build/

RUN mkdir -p /workspace
WORKDIR /workspace