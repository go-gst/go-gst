package gst

import "runtime/pprof"

var padprobesProfile *pprof.Profile

func init() {
	padprobes := "go-gst-active-pad-probes"
	padprobesProfile = pprof.Lookup(padprobes)
	if padprobesProfile == nil {
		padprobesProfile = pprof.NewProfile(padprobes)
	}
}
