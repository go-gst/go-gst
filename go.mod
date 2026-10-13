module github.com/go-gst/go-gst

go 1.25.0

require github.com/go-gst/go-glib v0.0.3

retract [v1.0.0, v1.4.1] // Handwritten bindings, not as stable as planned and not licensed open enough

replace github.com/go-gst/go-glib => github.com/iameli/go-glib v0.0.4-0.20261013211756-2f69802326f8
