FROM golang:1.15-alpine 

RUN mkdir /gstreamer \
    && cd /gstreamer \
    && apk add --update musl-dev glib glib-dev gcc meson git bison flex \
    && git clone https://gitlab.freedesktop.org/gstreamer/gst-build.git \
    && cd gst-build \
    && meson setup /workspace/include

WORKDIR /workspace
