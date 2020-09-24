package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/gotk3/gotk3/glib"
	"github.com/spf13/cobra"
	"github.com/tinyzimmer/go-gst-launch/gst"
)

func init() {
	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect the elements of the given pipeline string",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must specify an object to inspect")
		}
		return nil
	},
	RunE: inspect,
}

func inspect(cmd *cobra.Command, args []string) error {

	name := args[0]

	// load the registry
	registry := gst.GetRegistry()
	// get the factory for the element
	factory := gst.Find(name)

	if factory == nil {
		return errors.New("Could not get details for factory")
	}
	defer factory.Unref()

	// assume it's an element for now, can implement more later
	elem, err := gst.NewElement(name)
	if err != nil {
		return err
	}
	defer elem.Unref()

	var maxLevel int

	// dump info about the element

	printFactoryDetails(registry, factory)
	printPluginDetails(registry, factory)
	printHierarchy(elem.TypeFromInstance(), 0, &maxLevel)
	printInterfaces(elem)
	printPadTemplates(elem)
	printClockingInfo(elem)
	printURIHandlerInfo(elem)
	printPadInfo(elem)
	printElementPropertiesInfo(elem)
	printSignalInfo(elem)
	printChildrenInfo(elem)
	printPresentList(elem)

	return nil
}

func printFactoryDetails(registry *gst.Registry, factory *gst.ElementFactory) {

	// initialize tabwriter
	w := new(tabwriter.Writer)
	buf := new(bytes.Buffer)

	w.Init(
		buf,
		40,  // minwidth
		30,  // tabwith
		0,   // padding
		' ', // padchar
		0,   // flags
	)

	colorOrange.fprint(w, "Factory Details:\n")
	for _, key := range factory.GetMetadataKeys() {
		colorBlue.fprintfIndent(w, 2, "%s \t ", strings.Title(key))
		colorLightGray.fprint(w, factory.GetMetadata(key))
		colorReset.fprint(w, "\n")
	}

	w.Flush()
	fmt.Print(buf.String())
	fmt.Println()
}

func printPluginDetails(registry *gst.Registry, factory *gst.ElementFactory) {

	// initialize tabwriter
	w := new(tabwriter.Writer)
	buf := new(bytes.Buffer)

	w.Init(
		buf,
		40,  // minwidth
		30,  // tabwith
		0,   // padding
		' ', // padchar
		0,   // flags
	)

	pluginFeature, err := registry.LookupFeature(factory.Name())
	if err != nil {
		return
	}
	plugin := pluginFeature.GetPlugin()
	if plugin == nil {
		return
	}

	defer pluginFeature.Unref()
	defer plugin.Unref()

	colorOrange.fprint(w, "Plugin Details:\n")

	colorBlue.fprintIndent(w, 2, "Name \t ")
	colorLightGray.fprintf(w, "%s\n", pluginFeature.GetPluginName())

	colorBlue.fprintIndent(w, 2, "Description \t ")
	colorLightGray.fprintf(w, "%s\n", plugin.Description())

	colorBlue.fprintIndent(w, 2, "Filename \t ")
	colorLightGray.fprintf(w, "%s\n", plugin.Filename())

	colorBlue.fprintIndent(w, 2, "Version \t ")
	colorLightGray.fprintf(w, "%s\n", plugin.Version())

	colorBlue.fprintIndent(w, 2, "License \t ")
	colorLightGray.fprintf(w, "%s\n", plugin.License())

	colorBlue.fprintIndent(w, 2, "Source module \t ")
	colorLightGray.fprintf(w, "%s\n", plugin.Source())

	colorBlue.fprintIndent(w, 2, "Binary package \t ")
	colorLightGray.fprintf(w, "%s\n", plugin.Package())

	colorBlue.fprintIndent(w, 2, "Origin URLs \t ")
	colorLightGray.fprintf(w, "%s\n", plugin.Origin())

	w.Flush()
	fmt.Print(buf.String())

	fmt.Println()
}

func printHierarchy(gtype glib.Type, level int, maxLevel *int) {
	parent := gtype.Parent()

	*maxLevel = *maxLevel + 1
	level++

	if parent > 0 {
		printHierarchy(parent, level, maxLevel)
	}

	for i := 1; i < *maxLevel-level; i++ {
		colorReset.print("      ")
	}

	if *maxLevel-level > 0 {
		colorLightPurple.print(" +----")
	}

	colorGreen.printf("%s\n", gtype.Name())

}

func printInterfaces(elem *gst.Element) {
	fmt.Println()

	if ifaces := elem.Interfaces(); len(ifaces) > 0 {
		colorOrange.print("Implemented Interfaces:")
		for _, iface := range ifaces {
			colorGreen.printfIndent(2, "%s\n", iface)
		}
	}
}

func printPadTemplates(elem *gst.Element) {
	fmt.Println()

	tmpls := elem.GetPadTemplates()
	if len(tmpls) == 0 {
		return
	}
	colorOrange.print("Pad templates:\n")
	for _, tmpl := range tmpls {
		colorBlue.printfIndent(2, "%s template", strings.ToUpper(tmpl.Name()))
		colorReset.print(": ")
		colorBlue.printf("'%s'\n", strings.ToLower(tmpl.Direction().String()))

		colorBlue.printIndent(4, "Availability")
		colorReset.print(": ")
		colorLightGray.print(strings.Title(tmpl.Presence().String()))
		colorReset.print("\n")

		colorBlue.printIndent(4, "Capabilities")
		colorReset.print(": ")

		caps := tmpl.Caps()
		if len(caps) == 0 {
			colorOrange.printIndent(6, "ANY")
		} else {
			printCaps(&caps, 6)
		}
	}
	fmt.Println()
	fmt.Println()
}

func printClockingInfo(elem *gst.Element) {
	if !elem.Has(gst.ElementFlagRequireClock) && !elem.Has(gst.ElementFlagProvideClock) {
		colorLightGray.print("Element has no clocking capabilities.\n")
		return
	}
	fmt.Printf("%sClocking Interactions:%s\n", colorOrange, colorReset)

	if elem.Has(gst.ElementFlagRequireClock) {
		colorLightGray.printIndent(2, "element requires a clock\n")
	}

	if elem.Has(gst.ElementFlagProvideClock) {
		clock := elem.GetClock()
		if clock == nil {
			colorLightGray.printIndent(2, "selement is supposed to provide a clock but returned NULL%s\n")
		} else {
			defer clock.Unref()
			colorLightGray.printIndent(2, "element provides a clock: ")
			colorCyan.printf(clock.Name())
		}
	}

	fmt.Println()
}

func printURIHandlerInfo(elem *gst.Element) {
	if !elem.IsURIHandler() {
		colorLightGray.print("Element has no URI handling capabilities.\n")
		fmt.Println()
	}

	fmt.Println()
	colorOrange.print("URI handling capabilities:\n")
	colorLightGray.printfIndent(2, "Element can act as %s.\n", strings.ToLower(elem.GetURIType().String()))

	protos := elem.GetURIProtocols()

	if len(protos) == 0 {
		fmt.Println()
		return
	}

	colorLightGray.printIndent(2, "Supported URI protocols:\n")

	for _, proto := range protos {
		colorCyan.printfIndent(4, "%s\n", proto)
	}

	fmt.Println()
}

func printPadInfo(elem *gst.Element) {

	colorOrange.print("Pads:\n")
	pads := elem.GetPads()
	if len(pads) == 0 {
		colorCyan.printIndent(2, "none\n")
		return
	}
	for _, pad := range elem.GetPads() {
		defer pad.Unref()

		colorBlue.printIndent(2, strings.ToUpper(pad.Direction().String()))
		colorReset.print(": ")
		colorLightGray.printf("'%s'\n", pad.Name())

		if tmpl := pad.Template(); tmpl != nil {
			defer tmpl.Unref()
			colorBlue.printIndent(4, "Pad Template")
			colorReset.print(": ")
			colorLightGray.printf("'%s'\n", tmpl.Name())
		}

		if caps := pad.CurrentCaps(); caps != nil {
			colorBlue.printIndent(2, "Capabilities:\n")
			printCaps(&caps, 4)
		}
	}

	fmt.Println()
}

func printElementPropertiesInfo(elem *gst.Element) {
	printObjectPropertiesInfo(elem.Object, "Element Properties")
}

func printSignalInfo(elem *gst.Element)   {}
func printChildrenInfo(elem *gst.Element) {}
func printPresentList(elem *gst.Element)  {}

func printCaps(caps *gst.Caps, indent int) {
	for _, cap := range *caps {
		colorReset.print("\n")
		colorOrange.printfIndent(indent, "%s", cap.Name)
		for k, v := range cap.Data {
			colorReset.print("\n")
			colorOrange.printfIndent(indent+2, "%s", k)
			colorReset.print(": ")
			colorLightGray.print(fmt.Sprint(v))
		}
	}
	fmt.Println()
}
