package gst_test

import (
	"reflect"
	"testing"

	"github.com/go-gst/go-gst/pkg/gst"
)

type (
	Simple struct {
		I int64
		S string
	}
	WithPrivate struct {
		I int64
		_ map[string]int
	}
	Tagged struct {
		I int64 `gst:"tagged"`
	}
	Embedded struct {
		Simple
		X int64
		Y int64
	}
	SubStructure struct {
		Sub Simple
		X   int64
		Y   int64
	}
	TaggedSubStructure struct {
		Sub Simple `gst:"taggedsub"`
		X   int64
		Y   int64
	}

	// errors:
	Unsupported struct {
		I int
	}
	PointersSimple struct {
		X *int
		Y *int
	}
	PointersEmbedded struct {
		*Simple
		X *int64
		Y *int64
	}
)

func TestStructureMarshal(t *testing.T) {
	gst.Init()

	runMarshalTest(t, Simple{1, "foo"}, false)
	runMarshalTest(t, WithPrivate{}, false)
	runMarshalTest(t, Tagged{1}, false)
	runMarshalTest(t, Embedded{Simple: Simple{1, "foo"}, X: 2, Y: 3}, false)
	runMarshalTest(t, SubStructure{Sub: Simple{1, "foo"}, X: 2, Y: 3}, false)
	runMarshalTest(t, TaggedSubStructure{Sub: Simple{1, "foo"}, X: 2, Y: 3}, false)

	// errors:
	runMarshalTest(t, Unsupported{}, true)
	runMarshalTest(t, PointersSimple{new(int), new(int)}, true)
	runMarshalTest(t, PointersEmbedded{new(Simple), new(int64), new(int64)}, true)
}

func runMarshalTest[T any](t *testing.T, v T, expectErr bool) {
	t.Helper()

	var zero T

	structure, err := gst.MarshalStructure(v)

	if err != nil && !expectErr {
		t.Fatalf("marshal error: %v", err)
	}

	if expectErr {
		t.Logf("got expected error: %v", err)
		return
	}

	t.Log(structure.String())

	err = structure.UnmarshalInto(&zero)

	if err != nil {
		t.Fatalf("could not unmarshal %v", err)
	}

	if !reflect.DeepEqual(zero, v) {
		t.Fatalf("expected %#v, got %#v", v, zero)
	}
}
