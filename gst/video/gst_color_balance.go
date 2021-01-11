package video

/*
#include <stdlib.h>
#include <gst/gst.h>
#include <gst/video/video.h>

GstColorBalance * toGstColorBalance(GstElement * element)
{
	return GST_COLOR_BALANCE(element);
}
*/
import "C"
import (
	"unsafe"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
)

// ColorBalanceType is an enumeration indicating whether an element implements color
// balancing operations in software or in dedicated hardware. In general, dedicated
// hardware implementations (such as those provided by xvimagesink) are preferred.
type ColorBalanceType int

// Type castings
const (
	ColorBalanceHardware ColorBalanceType = C.GST_COLOR_BALANCE_HARDWARE // (0) – Color balance is implemented with dedicated hardware.
	ColorBalanceSoftware ColorBalanceType = C.GST_COLOR_BALANCE_SOFTWARE // (1) – Color balance is implemented via software processing.
)

// ColorBalanceChannel represents parameters for modifying the color balance implemented by
// an element providing the GstColorBalance interface. For example, Hue or Saturation.
type ColorBalanceChannel struct {
	// A string containing a descriptive name for this channel
	Label string
	// The minimum valid value for this channel.
	MinValue int
	// The maximum valid value for this channel.
	MaxValue int
}

// ColorBalance is an interface implemented by elements which can perform some color balance
// operation on video frames they process. For example, modifying the brightness, contrast,
// hue or saturation.
//
// Example elements are 'xvimagesink' and 'colorbalance'
type ColorBalance interface {
	// Get the ColorBalanceType of this implementation.
	GetBalanceType() ColorBalanceType
	// Retrieve the current value of the indicated channel, between MinValue and MaxValue.
	GetValue(*ColorBalanceChannel) int
	// Retrieve a list of the available channels.
	ListChannels() []*ColorBalanceChannel
	// Sets the current value of the channel to the passed value, which must be between MinValue
	// and MaxValue.
	SetValue(*ColorBalanceChannel, int)
}

// ColorBalanceFromElement checks if the given element implements the ColorBalance interface,
// and if so, returns a usable interface. This currently only supports elements created from the
// C runtime.
func ColorBalanceFromElement(element *gst.Element) ColorBalance {
	if C.toGstColorBalance(fromCoreElement(element)) == nil {
		return nil
	}
	return &gstColorBalance{fromCoreElement(element)}
}

// gstColorBalance implements a ColorBalance interface backed by an element
// from the C runtime.
type gstColorBalance struct{ elem *C.GstElement }

// Instance returns the C GstColorBalance interface.
func (c *gstColorBalance) Instance() *C.GstColorBalance {
	return C.toGstColorBalance(c.elem)
}

// GetBalanceType gets the ColorBalanceType of this implementation.
func (c *gstColorBalance) GetBalanceType() ColorBalanceType {
	return ColorBalanceType(C.gst_color_balance_get_balance_type(c.Instance()))
}

// GetValue retrieve the current value of the indicated channel, between MinValue and MaxValue.
func (c *gstColorBalance) GetValue(channel *ColorBalanceChannel) int {
	cLabel := C.CString(channel.Label)
	defer C.free(unsafe.Pointer(cLabel))
	gcbc := &C.GstColorBalanceChannel{
		label:     (*C.gchar)(cLabel),
		min_value: C.gint(channel.MinValue),
		max_value: C.gint(channel.MaxValue),
	}
	defer C.free(unsafe.Pointer(gcbc))
	return int(C.gst_color_balance_get_value(c.Instance(), gcbc))
}

// ListChannels retrieves a list of the available channels.
func (c *gstColorBalance) ListChannels() []*ColorBalanceChannel {
	gList := C.gst_color_balance_list_channels(c.Instance())
	if gList == nil {
		return nil
	}
	wrapped := glib.WrapList(uintptr(unsafe.Pointer(gList)))
	defer wrapped.Free()
	out := make([]*ColorBalanceChannel, 0)
	wrapped.Foreach(func(item interface{}) {
		channel := (*C.GstColorBalanceChannel)(item.(unsafe.Pointer))
		out = append(out, &ColorBalanceChannel{
			Label:    C.GoString(channel.label),
			MinValue: int(channel.min_value),
			MaxValue: int(channel.max_value),
		})
	})
	return out
}

// SetValue sets the current value of the channel to the passed value, which must be between MinValue
// and MaxValue.
func (c *gstColorBalance) SetValue(channel *ColorBalanceChannel, value int) {
	cLabel := C.CString(channel.Label)
	defer C.free(unsafe.Pointer(cLabel))
	gcbc := &C.GstColorBalanceChannel{
		label:     (*C.gchar)(cLabel),
		min_value: C.gint(channel.MinValue),
		max_value: C.gint(channel.MaxValue),
	}
	defer C.free(unsafe.Pointer(gcbc))
	C.gst_color_balance_set_value(c.Instance(), gcbc, C.gint(value))
}
