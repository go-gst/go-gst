package main

//go:generate go run . -o ./pkg/

import (
	"encoding/xml"
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/gotk4/gir"
	"github.com/diamondburned/gotk4/gir/cmd/gir-generate/gendata"
	"github.com/diamondburned/gotk4/gir/cmd/gir-generate/genmain"
	"github.com/diamondburned/gotk4/gir/girgen/strcases"
	"github.com/diamondburned/gotk4/gir/girgen/typesystem"
	girfiles_gst "github.com/go-gst/go-gst/girs"
)

const Module = "github.com/go-gst/go-gst/pkg"

var Data = genmain.Data{
	Module:   Module,
	GirFiles: girfiles_gst.GirFiles,
	Preprocessors: []gir.Preprocessor{
		gir.MustIntrospect("Gst-1.Message.copy"),

		// Bitfield has a member of same name:
		gir.PreprocessorFunc(func(r gir.Repositories) {
			bitfield := r.FindFullType("Gst-1.BufferCopyFlags").(*gir.Bitfield)

			for _, m := range bitfield.Members {
				if m.CIdentifier == "GST_BUFFER_COPY_FLAGS" {
					m.CIdentifier = "GST_BUFFER_COPY_BUFFER_FLAGS"

					m.Doc.String += " (go-gst: renamed from BufferCopyFlags)"
				}
			}
		}),

		// a member of the enum is generated twice:
		DedupBitfieldMembers("GstVideo-1.VideoBufferFlags"),
		DedupBitfieldMembers("GstVideo-1.VideoFrameFlags"),
		DedupBitfieldMembers("GstVideo-1.NavigationModifierType"),

		// the member names do not have a GST prefix, which creates collisions:
		FixCutoffEnumMemberNames("GstMpegts-1.DVBTeletextType"),

		MiniObjectExtenderBorrows(),

		// collides with the base src extenders that actually provide a clock instead of returning the provided one
		gir.RenameCallable("Gst-1.Element.provide_clock", "provided_clock"),

		// collides with base src set caps
		gir.RenameCallable("GstApp-1.AppSrc.set_caps", "app_src_set_caps"),

		// collides with extending audio base payloader push
		gir.RenameCallable("GstRtp-1.RTPBasePayload.push", "push_buffer"),

		// collides with extending audio base payloader push
		gir.RenameCallable("GstAllocators-1.DRMDumbAllocator.alloc", "drm_dumb_alloc"),

		// otherwise clashes with control binding class extension
		gir.RenameCallable("Gst-1.Object.get_control_binding", "current_control_binding"),
		// Collides with GstObject:
		gir.RenameCallable("Gst-1.ControlBinding.sync_values", "sync_control_binding_values"),

		// Collides with method of the same name
		gir.TypeRenamer("GstVideo-1.VideoChromaResample", "VideoChromaResampler"),

		gir.PreprocessorFunc(func(r gir.Repositories) {
			t := r.FindFullType("GstVideo-1.VideoTimeCode").(*gir.Record)

			// FIXME: the get_type function requires gst_init(), thus crashes if called
			// during init
			t.GLibGetType = ""
		}),
	},
	Config: typesystem.Config{
		Namespaces: map[string]typesystem.NamespaceConfig{
			"Gst-1": {
				MinVersion: "1.26",
				MaxVersion: "1.26",
				ManualTypes: []typesystem.Type{
					&typesystem.Alias{
						BaseType: typesystem.BaseType{
							GirName: "ClockTime",
							GoTyp:   "ClockTime",
							CGoTyp:  "C.GstClockTime",
							CTyp:    "GstClockTime",
						},
						AliasedType: typesystem.CouldBeForeign[typesystem.Type]{
							Type: typesystem.Guint64,
						},
					},
					// We create custom structs so we can call await on a channel instead of
					// blocking a thread.
					&typesystem.Record{
						BaseType: typesystem.BaseType{
							GirName: "Promise",
							GoTyp:   "Promise",
							CGoTyp:  "C.GstPromise",
							CTyp:    "GstPromise",
						},
						BaseConversions: typesystem.BaseConversions{
							FromGlibBorrowFunction: "UnsafePromiseFromGlibBorrow",
							FromGlibFullFunction:   "UnsafePromiseFromGlibFull",
							FromGlibNoneFunction:   "UnsafePromiseFromGlibNone",
							ToGlibNoneFunction:     "UnsafePromiseToGlibNone",
							ToGlibFullFunction:     "UnsafePromiseToGlibFull",
						},
					},
				},
				IgnoredDefinitions: []typesystem.IgnoreFunc{
					// Collide and use an out array of values. TODO: manually implement
					typesystem.IgnoreMatching("Object.get_g_value_array"),
					typesystem.IgnoreMatching("ControlBinding.get_g_value_array"),

					// Manually implemented:
					typesystem.IgnoreMatching("Object.get_value"),
					typesystem.IgnoreMatching("ControlBinding.get_value"), // TODO

					typesystem.IgnoreMatching("ElementFactory.make_with_properties"),
					typesystem.IgnoreMatching("Message.parse_property_notify"),
					typesystem.IgnoreMatching("Message.new_property_notify"),
					typesystem.IgnoreMatching("Message.get_stream_status_object"),
					typesystem.IgnoreMatching("Structure.get_value"),
					typesystem.IgnoreMatching("Structure.set_value"),
					typesystem.IgnoreMatching("Structure.id_get_value"),
					typesystem.IgnoreMatching("Structure.id_take_value"),
					typesystem.IgnoreMatching("Structure.take_value"),
					typesystem.IgnoreMatching("ChildProxy.set_property"),
					typesystem.IgnoreMatching("ChildProxy.get_property"),
					typesystem.IgnoreMatching("Iterator.next"),
					typesystem.IgnoreMatching("TagList.get_value_index"),

					// we have bindings for parse_launch(_full), if we need the v variants,
					// then manually implement them
					typesystem.IgnoreMatching("parse_launchv"),
					typesystem.IgnoreMatching("parse_launchv_full"),

					// gobject.NewValue handles this already.
					typesystem.IgnoreMatching("util_set_value_from_string"),

					// Buffer mapping is manually implemented:
					typesystem.IgnoreMatching("Buffer.map"),
					typesystem.IgnoreMatching("Buffer.unmap"),
					typesystem.IgnoreMatching("MapInfo"),

					// Requires a gvalue arg, manually implemented:
					typesystem.IgnoreMatching("TagSetter.add_tag_value"),

					// ParamSpec subclass colliding with constructor:
					typesystem.IgnoreMatching("ParamSpecArray"),
					typesystem.IgnoreMatching("ParamSpecFraction"),
				},
			},
			"GstApp-1": {
				MinVersion: "1.26",
				MaxVersion: "1.26",
			},
			"GstAllocators-1": {
				MinVersion: "1.26",
				MaxVersion: "1.26",
			},
			"GstWebRTC-1": {
				MinVersion: "1.26",
				MaxVersion: "1.26",
			},
			"GstBase-1": {
				MinVersion: "1.26",
				MaxVersion: "1.26",
				IgnoredDefinitions: []typesystem.IgnoreFunc{
					// has unexported free function that crashes the linker when compiling the examples:
					typesystem.IgnoreMatching("TypeFindData"),
				},
			},
			"GstVideo-1": {
				MinVersion: "1.26",
				MaxVersion: "1.26",
				IgnoredDefinitions: []typesystem.IgnoreFunc{
					// must be implemented manually
					typesystem.IgnoreMatching("VideoCodecFrame.set_user_data"),
					typesystem.IgnoreMatching("VideoCodecFrame.get_user_data"),
					// returns a gconstpointer to an array, manually implemented instead
					typesystem.IgnoreMatching("VideoFormat.get_palette"),
				},
			},
			"GstPbutils-1": {
				MinVersion: "1.26",
				MaxVersion: "1.26",
				IgnoredDefinitions: []typesystem.IgnoreFunc{
					// Resolve to ObjectClass:
					typesystem.IgnoreMatching("DiscovererAudioInfoClass"),
					typesystem.IgnoreMatching("DiscovererContainerInfoClass"),
					typesystem.IgnoreMatching("DiscovererInfoClass"),
					typesystem.IgnoreMatching("DiscovererStreamInfoClass"),
					typesystem.IgnoreMatching("DiscovererSubtitleInfoClass"),
					typesystem.IgnoreMatching("DiscovererVideoInfoClass"),
					typesystem.IgnoreMatching("EncodingTargetClass"),
				},
			},
		},
	},
	Postprocessors: []typesystem.PostProcessor{
		typesystem.MarkAsManuallyExtended("Gst-1", "Object"),
		typesystem.MarkAsManuallyExtended("Gst-1", "Element"),
		typesystem.MarkAsManuallyExtended("Gst-1", "Bin"),
		typesystem.MarkAsManuallyExtended("Gst-1", "Bus"),
		typesystem.MarkAsManuallyExtended("Gst-1", "ChildProxy"),
		typesystem.MarkAsManuallyExtended("Gst-1", "TagSetter"),

		// Virtual methods of BaseTransform collide with Element
		func(r *typesystem.Registry) error {
			base := r.FindNamespaceByName("GstBase-1")

			bt := base.FindLocalTypeByGIRName("BaseTransform").(*typesystem.Class)

			bt.FindVirtualMethod("query").ParentName = "ParentQueryBaseTransform"

			return nil
		},
		// Virtual methods of PushSrc collide with BaseSrc
		func(r *typesystem.Registry) error {
			base := r.FindNamespaceByName("GstBase-1")

			pushsrc := base.FindLocalTypeByGIRName("PushSrc").(*typesystem.Class)

			pushsrc.FindVirtualMethod("alloc").ParentName = "ParentAllocPushSrc"
			pushsrc.FindVirtualMethod("fill").ParentName = "ParentFillPushSrc"

			return nil
		},
		// Virtual methods of AudioSink collide with BaseSink
		func(r *typesystem.Registry) error {
			audio := r.FindNamespaceByName("GstAudio-1")

			audioSink := audio.FindLocalTypeByGIRName("AudioSink").(*typesystem.Class)

			audioSink.FindVirtualMethod("prepare").ParentName = "ParentPrepareAudioSink"
			audioSink.FindVirtualMethod("stop").ParentName = "ParentStopAudioSink"

			return nil
		},
		// Virtual methods of RTPBasePayload collide with Element
		func(r *typesystem.Registry) error {
			rtp := r.FindNamespaceByName("GstRtp-1")

			rtpBasePayload := rtp.FindLocalTypeByGIRName("RTPBasePayload").(*typesystem.Class)

			rtpBasePayload.FindVirtualMethod("query").ParentName = "ParentQueryRTPBasePayload"

			return nil
		},
	},
}

func main() {
	genmain.Run(
		gendata.Main,
		Data,
	)
}

var borrowedTypes = []string{
	"Gst-1.MiniObject", "Gst-1.Structure", "Gst-1.Caps", "Gst-1.Buffer", "Gst-1.BufferList", "Gst-1.Memory", "Gst-1.Message", "Gst-1.Query", "Gst-1.Sample",
}

// gst.MiniObject extenders must not take a ref on these methods, or they are made readonly
func MiniObjectExtenderBorrows() gir.Preprocessor {
	return gir.PreprocessorFunc(func(repos gir.Repositories) {
		for _, fulltype := range borrowedTypes {
			res := repos.FindFullType(fulltype)
			if res == nil {
				log.Fatalf("fulltype %s not found", fulltype)
			}

			switch typ := res.(type) {
			case *gir.Record:
				for i, m := range typ.Methods {
					if m.ReturnValue.TransferOwnership.TransferOwnership == "none" && slices.ContainsFunc(borrowedTypes, func(typ string) bool {
						return strings.SplitN(typ, ".", 2)[1] == m.ReturnValue.Type.Name
					}) {
						log.Printf("marking function as borrowing: %s", m.Name)
						typ.Methods[i].ReturnValue.TransferOwnership.TransferOwnership = "borrow"
					}
				}
			default:
				log.Fatalf("unhandled type for %s", fulltype)
			}
		}
	})
}

func MarkReturnAsBorrowed(fulltype string) gir.Preprocessor {
	return gir.ModifyCallable(fulltype, func(c *gir.CallableAttrs) {
		c.ReturnValue.TransferOwnership.TransferOwnership = "borrow"
	})
}

func DedupBitfieldMembers(fulltype string) gir.Preprocessor {
	return gir.PreprocessorFunc(func(repos gir.Repositories) {
		bf := repos.FindFullType(fulltype).(*gir.Bitfield)

		oldmembers := bf.Members

		bf.Members = nil

		seen := make(map[string]struct{})

		for _, m := range oldmembers {
			if _, ok := seen[m.CIdentifier]; ok {
				continue
			}

			seen[m.CIdentifier] = struct{}{}

			bf.Members = append(bf.Members, m)
		}
	})
}

// FixCutoffEnumMemberNames adds a namespace prefix, that will later be cut off by FormatMember. It also regenerates the name from the C identifier
func FixCutoffEnumMemberNames(fulltype string) gir.Preprocessor {
	return gir.MapMembers(fulltype, func(member *gir.Member) {
		newname := strcases.SnakeToGo(true, strings.ToLower(member.CIdentifier))

		nameAttr := xml.Name{Local: "name"}
		for i, attr := range member.Names {
			if attr.Name == nameAttr {
				attr.Value = newname
			}

			member.Names[i] = attr
		}

		member.CIdentifier = "gst_" + member.CIdentifier
	})
}
