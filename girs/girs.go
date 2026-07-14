package girfiles_gst

import (
	"embed"

	girfiles "github.com/go-gst/go-glib/girs"
)

//go:embed *.gir
var girFiles embed.FS

var GirFiles = girfiles.ReadGirFiles(girFiles)
