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
	"github.com/diamondburned/gotk4/gir/girgen"
	"github.com/diamondburned/gotk4/gir/girgen/strcases"
	"github.com/diamondburned/gotk4/gir/girgen/types"
)

const Module = "github.com/go-gst/go-gst/pkg"

var Data = genmain.Overlay(
	gendata.Main,
	genmain.Data{
		Module: Module,
		Packages: []genmain.Package{
			// {Name: "gst-editing-services-1.0", Namespaces: []string{"GES-1"}},
			{Name: "gstreamer-1.0", Namespaces: []string{"Gst-1"}},
			{Name: "gstreamer-allocators-1.0", Namespaces: []string{"GstAllocators-1"}},
			// {Name: "gstreamer-analytics-1.0", Namespaces: []string{"GstAnalytics-1"}},
			{Name: "gstreamer-app-1.0", Namespaces: []string{"GstApp-1"}},
			{Name: "gstreamer-audio-1.0", Namespaces: []string{"GstAudio-1"}},
			// {Name: "gstreamer-bad-audio-1.0", Namespaces: []string{"GstBadAudio-1"}},
			// {Name: "gstreamer-bad-base-camerabinsrc-1.0", Namespaces: []string{"GstBadBaseCameraBin-1"}},
			{Name: "gstreamer-base-1.0", Namespaces: []string{"GstBase-1"}},
			{Name: "gstreamer-check-1.0", Namespaces: []string{"GstCheck-1"}},
			// {Name: "gstreamer-codecs-1.0", Namespaces: []string{"GstCodecs-1"}},
			{Name: "gstreamer-controller-1.0", Namespaces: []string{"GstController-1"}},
			// {Name: "gstreamer-cuda-1.0", Namespaces: []string{"GstCuda-1"}},
			{Name: "gstreamer-gl-1.0", Namespaces: []string{"GstGL-1"}},
			// {Name: "gstreamer-insertbin-1.0", Namespaces: []string{"GstInsertBin-1"}},
			{Name: "gstreamer-mpegts-1.0", Namespaces: []string{"GstMpegts-1"}},
			// {Name: "gstreamer-mse-1.0", Namespaces: []string{"GstMse-1"}},
			{Name: "gstreamer-net-1.0", Namespaces: []string{"GstNet-1"}},
			{Name: "gstreamer-pbutils-1.0", Namespaces: []string{"GstPbutils-1"}},
			{Name: "gstreamer-play-1.0", Namespaces: []string{"GstPlay-1"}},
			{Name: "gstreamer-player-1.0", Namespaces: []string{"GstPlayer-1"}},
			{Name: "gstreamer-rtp-1.0", Namespaces: []string{"GstRtp-1"}},
			{Name: "gstreamer-rtsp-1.0", Namespaces: []string{"GstRtsp-1"}},
			// {Name: "gstreamer-rtsp-server-1.0", Namespaces: []string{"GstRtspServer-1"}},
			{Name: "gstreamer-sdp-1.0", Namespaces: []string{"GstSdp-1"}},
			{Name: "gstreamer-tag-1.0", Namespaces: []string{"GstTag-1"}},
			// {Name: "gstreamer-transcoder-1.0", Namespaces: []string{"GstTranscoder-1"}},
			// {Name: "gstreamer-va-1.0", Namespaces: []string{"GstVa-1"}},
			// {Name: "gstreamer-validate-1.0", Namespaces: []string{"GstValidate-1"}},
			{Name: "gstreamer-video-1.0", Namespaces: []string{"GstVideo-1"}},
			// {Name: "gstreamer-vulkan-1.0", Namespaces: []string{"GstVulkan-1"}},
			{Name: "gstreamer-webrtc-1.0", Namespaces: []string{"GstWebRTC-1"}},
		},
		PkgExceptions: []string{
			"core",
		},
		PkgGenerated: []string{
			// gst packages:
			// "cudagst",
			// "ges",
			"gst",
			"gstallocators",
			// "gstanalytics",
			"gstapp",
			"gstaudio",
			// "gstbadaudio",
			"gstbase",
			"gstcheck",
			// "gstcodecs",
			"gstcontroller",
			// "gstcuda",
			// "gstdxva",
			"gstgl",
			// "gstglegl",
			// "gstglwayland",
			// "gstglx11",
			// "gstinsertbin",
			"gstmpegts",
			// "gstmse",
			"gstnet",
			"gstpbutils",
			"gstplay",
			"gstplayer",
			"gstrtp",
			"gstrtsp",
			// "gstrtspserver",
			"gstsdp",
			"gsttag",
			// "gsttranscoder",
			// "gstva",
			// "gstvalidate",
			"gstvideo",
			"gstwebrtc",
		},
		Preprocessors: []types.Preprocessor{
			// Enum has a member of same name:
			types.TypeRenamer("Gst-1.BufferCopyFlags", "BufferCopyFlagsType"),

			// a member of the enum is generated twice:
			DedupBitfieldMembers("GstVideo-1.VideoBufferFlags"),
			DedupBitfieldMembers("GstVideo-1.VideoFrameFlags"),
			DedupBitfieldMembers("GstVideo-1.NavigationModifierType"),

			// the member names do not have a GST prefix, which creates collisions:
			FixCutoffEnumMemberNames("GstMpegts-1.DVBTeletextType"),

			// GArray record fields in SDP, not needed anyways because there are accessors
			types.RemoveRecordFields("GstSdp-1.SDPMedia", "fmts", "connections", "bandwidths", "attributes"),
			types.RemoveRecordFields("GstSdp-1.SDPMessage", "emails", "phones", "bandwidths", "times", "zones", "attributes", "medias"),
			types.RemoveRecordFields("GstSdp-1.SDPTime", "repeat"),
			types.RemoveRecordFields("GstSdp-1.MIKEYMessage", "map_info", "payloads"),
			types.RemoveRecordFields("GstSdp-1.MIKEYPayloadKEMAC", "subpayloads"),
			types.RemoveRecordFields("GstSdp-1.MIKEYPayloadSP", "params"),

			MiniObjectExtenderBorrows(),

			// collides with the base src extenders that actually provide a clock instead of returning the provided one
			types.RenameCallable("Gst-1.Element.provide_clock", "ProvidedClock"),

			// collides with base src set caps
			types.RenameCallable("GstApp-1.AppSrc.set_caps", "AppSrcSetCaps"),

			// collides with extending audio base payloader push
			types.RenameCallable("GstRtp-1.RTPBasePayload.push", "PushBuffer"),

			// collides with extending audio base payloader push
			types.RenameCallable("GstAllocators-1.DRMDumbAllocator.alloc", "DRMAlloc"),

			// otherwise clashes with control binding class extension
			types.RenameCallable("Gst-1.Object.get_control_binding", "CurrentControlBinding"),
		},
		Postprocessors: map[string][]girgen.Postprocessor{
			"Gst-1":       {ElementFactoryMakeWithProperties, ElementBlockSetState, BinAddMany, ElementLinkMany, IteratorValues, StructureGoMarshal},
			"GstMpegts-1": {GstUseUnstableAPI},
			"GstWebRTC-1": {GstUseUnstableAPI, FixWebrtcPkgConfig},
		},
		Filters: []types.FilterMatcher{
			// these collide and are not really useful:
			types.AbsoluteFilter("C.gst_structure_new_from_string"),
			types.AbsoluteFilter("C.gst_structure_from_string"),
			// we have goroutines :) (also creates compile errors)
			types.FileFilterNamespace("Gst", "taskpool"),

			// GArray not working here
			types.AbsoluteFilter("C.gst_mpegts_descriptor_parse_dvb_ca_identifier"),
			types.AbsoluteFilter("C.gst_mpegts_descriptor_parse_dvb_frequency_list"),
			types.AbsoluteFilter("C.gst_source_buffer_get_buffered"),

			// Excluded because the additional_info param seems to have changed between my local gstreamer and 1.24.10 (used in github actions)
			types.AbsoluteFilter("C.gst_mpegts_descriptor_parse_registration"),

			// In-out array pointer, not very go like and not correctly handled by girgen, needs custom implementation
			types.AbsoluteFilter("C.gst_audio_get_channel_reorder_map"),

			// re-implemented in go to support a properties map instead of two lists
			types.AbsoluteFilter("C.gst_element_factory_make_with_properties"),

			// FIXME: creates a name clash with the parents GstObject:
			types.AbsoluteFilter("C.gst_control_binding_get_g_value_array"), // vs gst_object_get_g_value_array
			types.AbsoluteFilter("C.gst_control_binding_sync_values"),       // vs gst_object_sync_values
			types.AbsoluteFilter("C.gst_control_binding_get_value"),         // vs gst_object_get_value
		},
		ExtraGoContents: ExtraGoContents,
		SingleFile:      true,
	},
)

func main() {
	genmain.Run(Data)
}

var borrowedTypes = []string{
	"Gst-1.MiniObject", "Gst-1.Structure", "Gst-1.Caps", "Gst-1.Buffer", "Gst-1.BufferList", "Gst-1.Memory", "Gst-1.Message", "Gst-1.Query", "Gst-1.Sample",
}

// gst.MiniObject extenders must not take a ref on these methods, or they are made readonly
func MiniObjectExtenderBorrows() types.Preprocessor {
	return types.PreprocessorFunc(func(repos gir.Repositories) {
		for _, fulltype := range borrowedTypes {
			res := repos.FindFullType(fulltype)
			if res == nil {
				log.Fatalf("fulltype %s not found", fulltype)
			}

			switch typ := res.Type.(type) {
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

func MarkReturnAsBorrowed(fulltype string) types.Preprocessor {
	return types.ModifyCallable(fulltype, func(c *gir.CallableAttrs) {
		c.ReturnValue.TransferOwnership.TransferOwnership = "borrow"
	})
}

func DedupBitfieldMembers(fulltype string) types.Preprocessor {
	return types.PreprocessorFunc(func(repos gir.Repositories) {
		bf := repos.FindFullType(fulltype).Type.(*gir.Bitfield)

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
func FixCutoffEnumMemberNames(fulltype string) types.Preprocessor {
	return types.MapMembers(fulltype, func(member gir.Member) gir.Member {
		newname := strcases.SnakeToGo(true, strings.ToLower(member.CIdentifier))

		nameAttr := xml.Name{Local: "name"}
		for i, attr := range member.Names {
			if attr.Name == nameAttr {
				attr.Value = newname
			}

			member.Names[i] = attr
		}

		member.CIdentifier = "gst_" + member.CIdentifier

		return member
	})
}

// GstUseUnstableAPI adds a define GST_USE_UNSTABLE_API in every file of the namespace
func GstUseUnstableAPI(nsgen *girgen.NamespaceGenerator) error {
	for _, f := range nsgen.Files {
		f.Header().DefineC("GST_USE_UNSTABLE_API // APIs in this package are unstable")
	}

	return nil
}

func FixWebrtcPkgConfig(nsgen *girgen.NamespaceGenerator) error {
	// see: https://gitlab.freedesktop.org/gstreamer/gstreamer/-/merge_requests/8433, remove after release
	nsgen.Namespace().Repository.Packages = append(nsgen.Namespace().Repository.Packages, gir.Package{
		Name: "gstreamer-sdp-1.0",
	})

	for _, f := range nsgen.Files {
		// see https://gitlab.freedesktop.org/gstreamer/gstreamer/-/merge_requests/8470 , remove after release
		f.Header().IncludeC("gst/webrtc/sctptransport.h")
	}

	return nil
}

func IteratorValues(nsgen *girgen.NamespaceGenerator) error {
	fg := nsgen.MakeFile("iteratorvalues.go")

	p := fg.Pen()

	fg.Header().Import("iter")

	p.Line(`
	// Values allows you to access the values from the iterator in a go for loop via function iterators
	func (it *Iterator) Values() iter.Seq[any] {
		return func(yield func(any) bool) {
			for {
				v, ret := it.Next()
				switch ret {
				case IteratorDone:
					return
				case IteratorResync:
					it.Resync()
				case IteratorOK:
					if !yield(v.GoValue()) {
						return
					}

				case IteratorError:
					panic("iterator values failed")
				default:
					panic("iterator values returned unknown state")
				}
			}
		}
	}
	`)

	return nil
}

func StructureGoMarshal(nsgen *girgen.NamespaceGenerator) error {
	fg := nsgen.MakeFile("structuremarshal.go")

	p := fg.Pen()

	fg.Header().NeedsExternGLib()
	fg.Header().Import("reflect")

	p.Line(`
	// MarshalStructure will convert the given go struct into a GstStructure. Currently nested
	// structs are not supported. You can control the mapping of the field names via the tags of the go struct.
	func MarshalStructure(data interface{}) *Structure {
		typeOf := reflect.TypeOf(data)
		valsOf := reflect.ValueOf(data)
		st := NewStructureEmpty(typeOf.Name())
		for i := 0; i < valsOf.NumField(); i++ {
			gval := valsOf.Field(i).Interface()
			
			fieldName, ok := typeOf.Field(i).Tag.Lookup("gst")

			if !ok {
				fieldName = typeOf.Field(i).Name
			}

			st.SetValue(fieldName, coreglib.NewValue(gval))
		}
		return st
	}

	// UnmarshalInto will unmarshal this structure into the given pointer. The object
	// reflected by the pointer must be non-nil. You can control the mapping of the field names via the tags of the go struct.
	func (s *Structure) UnmarshalInto(data interface{}) error {
		rv := reflect.ValueOf(data)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("data is invalid (nil or non-pointer)")
		}

		val := reflect.ValueOf(data).Elem()
		nVal := rv.Elem()
		for i := 0; i < val.NumField(); i++ {
			nvField := nVal.Field(i)

			fieldName, ok := val.Type().Field(i).Tag.Lookup("gst")

			if !ok {
				fieldName = val.Type().Field(i).Name
			}

			val := s.Value(fieldName)

			nvField.Set(reflect.ValueOf(val))
		}

		return nil
	}
	`)

	return nil
}
func BinAddMany(nsgen *girgen.NamespaceGenerator) error {
	fg := nsgen.MakeFile("binaddmany.go")

	p := fg.Pen()

	p.Line(`
	// AddMany repeatedly calls Add for each param
	func (bin *Bin) AddMany(elements... Elementer) bool {
		for _, el := range elements {
			if !bin.Add(el) {
				return false
			}
		}

		return true
	}
	`)

	return nil
}

func ElementLinkMany(nsgen *girgen.NamespaceGenerator) error {
	fg := nsgen.MakeFile("elementlinkmany.go")

	p := fg.Pen()

	p.Line(`
	// LinkMany links the given elements in the order passed
	func LinkMany(elements... Elementer) bool {
		if len(elements) == 0 {
			return true
		}

		current := elements[0].(*Element)

		for _, next := range elements[1:] {
			if !	current.Link(next) {
				return false
			}

			current = next.(*Element)
		}

		return true
	}
	`)

	return nil
}

func ElementBlockSetState(nsgen *girgen.NamespaceGenerator) error {
	fg := nsgen.MakeFile("elementblocksetstate.go")

	p := fg.Pen()

	p.Line(`
	// BlockSetState is a convenience wrapper around calling SetState and State to wait for async state changes. See State for more info.
	func (el *Element) BlockSetState(state State, timeout ClockTime) StateChangeReturn {
		ret := el.SetState(state)

		if ret == StateChangeAsync {
			_, _, ret = el.State(timeout)
		}

		return ret
	}
	`)

	return nil
}

func ElementFactoryMakeWithProperties(nsgen *girgen.NamespaceGenerator) error {
	fg := nsgen.MakeFile("elementfactory.go")
	fg.Header().NeedsExternGLib()

	p := fg.Pen()

	// this is adapted from the autogenerated code and the coreglib.NewObjectWithProperties
	p.Line(`
		// ElementFactoryMakeWithProperties: create a new element of the type defined by
		// the given elementfactory. The supplied list of properties, will be passed at
		// object construction.
		//
		// The function takes the following parameters:
		//
		//   - factoryname: named factory to instantiate.
		//   - names (optional): array of properties names.
		//   - values (optional): array of associated properties values.
		//
		// The function returns the following values:
		//
		//   - element (optional): new Element or NULL if the element couldn't be
		//     created.
		func ElementFactoryMakeWithProperties(factoryname string, properties map[string]any) Elementer {
			var _arg1 *C.gchar  // out
			var _arg2 C.guint
			var _cret *C.GstElement // in

			_arg1 = (*C.gchar)(unsafe.Pointer(C.CString(factoryname)))
			defer C.free(unsafe.Pointer(_arg1))

			var names_ **C.gchar
			var values_ *C.GValue

			if len(properties) > 0 {
				names := make([]*C.char, 0, len(properties))
				values := make([]C.GValue, 0, len(properties))

				for name, value := range properties {
					cname := (*C.char)(C.CString(name))
					defer C.free(unsafe.Pointer(cname))

					gvalue := coreglib.NewValue(value)
					defer runtime.KeepAlive(gvalue)

					names = append(names, cname)
					values = append(values, *(*C.GValue)(unsafe.Pointer(gvalue.Native())))
				}

				names_ = &names[0]
				values_ = &values[0]
			}

			_cret = C.gst_element_factory_make_with_properties(_arg1, _arg2, names_, values_)
			runtime.KeepAlive(factoryname)

			var _element Elementer // out

			if _cret != nil {
				{
					objptr := unsafe.Pointer(_cret)

					object := coreglib.Take(objptr)
					casted := object.WalkCast(func(obj coreglib.Objector) bool {
						_, ok := obj.(Elementer)
						return ok
					})
					rv, ok := casted.(Elementer)
					if !ok {
						panic("no marshaler for " + object.TypeFromInstance().String() + " matching gst.Elementer")
					}
					_element = rv
				}
			}

			return _element
		}`)

	return nil

}

var ExtraGoContents = map[string]string{
	"gst/gst.go": `
		// Init binds to the gst_init() function. Argument parsing is not
		// supported.
		func Init() {
			C.gst_init(nil, nil)
		}

		// ClockTimeNone means infinite timeout or an empty value
		const ClockTimeNone ClockTime = 0xffffffffffffffff // Ideally this would be set to C.GST_CLOCK_TIME_NONE but this causes issues on MacOS and Windows
	`,
}
