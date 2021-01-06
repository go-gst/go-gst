package gst

/*
#include "gst.go.h"

void cgoDebugLog (GstDebugCategory * category,
               	  GstDebugLevel level,
                  const gchar * file,
                  const gchar * function,
                  gint line,
                  GObject * object,
				  const gchar * msg)
{
	gst_debug_log(category, level, file, function, line, object, msg);
}

*/
import "C"
import (
	"path"
	"runtime"
)

// DebugColorFlags are terminal style flags you can use when creating your debugging
// categories to make them stand out in debugging output.
type DebugColorFlags int

// Type castings of DebugColorFlags
const (
	DebugColorNone      DebugColorFlags = 0                      // (0) - No color
	DebugColorFgBlack   DebugColorFlags = C.GST_DEBUG_FG_BLACK   // (0) – Use black as foreground color.
	DebugColorFgRed     DebugColorFlags = C.GST_DEBUG_FG_RED     // (1) – Use red as foreground color.
	DebugColorFgGreen   DebugColorFlags = C.GST_DEBUG_FG_GREEN   // (2) – Use green as foreground color.
	DebugColorFgYellow  DebugColorFlags = C.GST_DEBUG_FG_YELLOW  // (3) – Use yellow as foreground color.
	DebugColorFgBlue    DebugColorFlags = C.GST_DEBUG_FG_BLUE    // (4) – Use blue as foreground color.
	DebugColorFgMagenta DebugColorFlags = C.GST_DEBUG_FG_MAGENTA // (5) – Use magenta as foreground color.
	DebugColorFgCyan    DebugColorFlags = C.GST_DEBUG_FG_CYAN    // (6) – Use cyan as foreground color.
	DebugColorFgWhite   DebugColorFlags = C.GST_DEBUG_FG_WHITE   // (7) – Use white as foreground color.
	DebugColorBgBlack   DebugColorFlags = C.GST_DEBUG_BG_BLACK   // (0) – Use black as background color.
	DebugColorBgRed     DebugColorFlags = C.GST_DEBUG_BG_RED     // (16) – Use red as background color.
	DebugColorBgGreen   DebugColorFlags = C.GST_DEBUG_BG_GREEN   // (32) – Use green as background color.
	DebugColorBgYellow  DebugColorFlags = C.GST_DEBUG_BG_YELLOW  // (48) – Use yellow as background color.
	DebugColorBgBlue    DebugColorFlags = C.GST_DEBUG_BG_BLUE    // (64) – Use blue as background color.
	DebugColorBgMagenta DebugColorFlags = C.GST_DEBUG_BG_MAGENTA // (80) – Use magenta as background color.
	DebugColorBgCyan    DebugColorFlags = C.GST_DEBUG_BG_CYAN    // (96) – Use cyan as background color.
	DebugColorBgWhite   DebugColorFlags = C.GST_DEBUG_BG_WHITE   // (112) – Use white as background color.
	DebugColorBold      DebugColorFlags = C.GST_DEBUG_BOLD       // (256) – Make the output bold.
	DebugColorUnderline DebugColorFlags = C.GST_DEBUG_UNDERLINE  // (512) – Underline the output.
)

// DebugColorMode represents how to display colors.
type DebugColorMode int

// Type castings of DebugColorModes
const (
	DebugColorModeOff  DebugColorMode = C.GST_DEBUG_COLOR_MODE_OFF  // (0) – Do not use colors in logs.
	DebugColorModeOn   DebugColorMode = C.GST_DEBUG_COLOR_MODE_ON   // (1) – Paint logs in a platform-specific way.
	DebugColorModeUnix DebugColorMode = C.GST_DEBUG_COLOR_MODE_UNIX // (2) – Paint logs with UNIX terminal color codes no matter what platform GStreamer is running on.
)

// DebugLevel defines the importance of a debugging message. The more important a message is, the
// greater the probability that the debugging system outputs it.
type DebugLevel int

// Type castings of DebugLevels
const (
	LevelNone    DebugLevel = C.GST_LEVEL_NONE    // (0) – No debugging level specified or desired. Used to deactivate debugging output.
	LevelError   DebugLevel = C.GST_LEVEL_ERROR   // (1) – Error messages are to be used only when an error occurred that stops the application from keeping working correctly. An examples is gst_element_error, which outputs a message with this priority. It does not mean that the application is terminating as with g_error.
	LevelWarning DebugLevel = C.GST_LEVEL_WARNING // (2) – Warning messages are to inform about abnormal behaviour that could lead to problems or weird behaviour later on. An example of this would be clocking issues ("your computer is pretty slow") or broken input data ("Can't synchronize to stream.")
	LevelFixMe   DebugLevel = C.GST_LEVEL_FIXME   // (3) – Fixme messages are messages that indicate that something in the executed code path is not fully implemented or handled yet. Note that this does not replace proper error handling in any way, the purpose of this message is to make it easier to spot incomplete/unfinished pieces of code when reading the debug log.
	LevelInfo    DebugLevel = C.GST_LEVEL_INFO    // (4) – Informational messages should be used to keep the developer updated about what is happening. Examples where this should be used are when a typefind function has successfully determined the type of the stream or when an mp3 plugin detects the format to be used. ("This file has mono sound.")
	LevelDebug   DebugLevel = C.GST_LEVEL_DEBUG   // (5) – Debugging messages should be used when something common happens that is not the expected default behavior, or something that's useful to know but doesn't happen all the time (ie. per loop iteration or buffer processed or event handled). An example would be notifications about state changes or receiving/sending of events.
	LevelLog     DebugLevel = C.GST_LEVEL_LOG     // (6) – Log messages are messages that are very common but might be useful to know. As a rule of thumb a pipeline that is running as expected should never output anything else but LOG messages whilst processing data. Use this log level to log recurring information in chain functions and loop functions, for example.
	LevelTrace   DebugLevel = C.GST_LEVEL_TRACE   // (7) – Tracing-related messages. Examples for this are referencing/dereferencing of objects.
	LevelMemDump DebugLevel = C.GST_LEVEL_MEMDUMP // (9) – memory dump messages are used to log (small) chunks of data as memory dumps in the log. They will be displayed as hexdump with ASCII characters.
)

// StackTraceFlags are flags for configuring stack traces
type StackTraceFlags int

// Type castings of StackTraceFlags
const (
	StackTraceShowNone StackTraceFlags = C.GST_STACK_TRACE_SHOW_NONE // (0) – Try to retrieve the minimum information available, which may be none on some platforms (Since: 1.18)
	StackTraceShowFull StackTraceFlags = C.GST_STACK_TRACE_SHOW_FULL // (1) – Try to retrieve as much information as possible, including source information when getting the stack trace
)

// DebugCategory is a struct that describes a category of log messages.
type DebugCategory struct {
	ptr *C.GstDebugCategory
}

// NewDebugCategory initializes a new DebugCategory with the given properties and set
// to the default threshold.
func NewDebugCategory(name string, color DebugColorFlags, description string) *DebugCategory {
	cat := C._gst_debug_category_new(C.CString(name), C.guint(color), C.CString(description))
	return &DebugCategory{ptr: cat}
}

// Log logs the given message using the currently registered debugging handlers. You can optionally
// provide a single object to log the message for. GStreamer will automatically add a newline to the
// end of the message.
func (d *DebugCategory) Log(level DebugLevel, message string, obj ...*Object) {
	var o *C.GObject
	if len(obj) > 0 {
		o = (*C.GObject)(obj[0].Unsafe())
	}
	function, file, line, _ := runtime.Caller(1)
	C.cgoDebugLog(
		d.ptr,
		C.GstDebugLevel(level),
		C.CString(path.Base(file)),
		(C.CString(runtime.FuncForPC(function).Name())),
		C.gint(line),
		o,
		C.CString(message),
	)
}
