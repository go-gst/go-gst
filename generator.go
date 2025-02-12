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
		},
		Postprocessors: map[string][]girgen.Postprocessor{
			"Gst-1":       {ElementFactoryMakeWithProperties},
			"GstMpegts-1": {GstUseUnstableAPI},
			"GstWebRTC-1": {GstUseUnstableAPI},
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

			// In-out array pointer, not very go like and not correctly handled by girgen, needs custom implementation
			types.AbsoluteFilter("C.gst_audio_get_channel_reorder_map"),

			// re-implemented in go to support a properties map instead of two lists
			types.AbsoluteFilter("C.gst_element_factory_make_with_properties"),
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

func GstUseUnstableAPI(nsgen *girgen.NamespaceGenerator) error {
	fg := nsgen.MakeFile("unstable.go")

	fg.Header().DefineC("GST_USE_UNSTABLE_API // APIs in this package are unstable")

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
	`,
}
