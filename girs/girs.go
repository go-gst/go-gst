package girfiles_gst

import (
	"embed"

	girfiles_gotk4 "github.com/diamondburned/gotk4/girs"
)

//go:embed *.gir
var girFiles embed.FS

var GirFiles = girfiles_gotk4.ReadGirFiles(girFiles)
