package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/gstauto"
)

func init() {
	rootCmd.AddCommand(launchCmd)
}

var launchSrc, launchDest *os.File

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Run a generic pipeline",
	Long:  `Uses the provided gstreamer string to encode/decode the data in the pipeline.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("The pipeline string cannot be empty")
		}
		var err error
		launchSrc, launchDest, err = getCLIFiles()
		return err
	},
	RunE: launch,
}

func launch(cmd *cobra.Command, args []string) error {

	mainLoop := gst.NewMainLoop(nil, false)
	defer mainLoop.Unref()

	pipelineString := strings.Join(args, " ")

	pipeliner, err := getPipeline(launchSrc, launchDest, pipelineString)
	if err != nil {
		return err
	}

	// If the pipeline is not dumping to stdout, dump messages to stdout instead.
	if !toStdout {
		pipeliner.Pipeline().GetBus().AddWatch(func(msg *gst.Message) bool {
			fmt.Println(msg)
			return true
		})
	}

	if err := pipeliner.Start(); err != nil {
		return err
	}

	// If there are src or dest files, spawn off copies of the data
	if launchSrc != nil {
		pipelineWriter := pipeliner.(gstauto.WritePipeliner)
		go io.Copy(pipelineWriter, launchSrc)
	}
	if launchDest != nil {
		pipelineReader := pipeliner.(gstauto.ReadPipeliner)
		go io.Copy(launchDest, pipelineReader)
	}

	// Catch SIGINT and SIGTERM
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Catch SIGINT so the pipeline can close cleanly
		<-sigc
		pipeliner.Pipeline().BlockSetState(gst.StateNull) // Do an extra call to stop the state
		// Increases the likelihood the user will get
		// final messages.
		mainLoop.Quit()
	}()

	go func() {
		// If the pipeline finishes, close the main loop
		gst.Wait(pipeliner.Pipeline())
		mainLoop.Quit()
	}()

	// Block on the main loop until either the pipeline finishes
	// or a signal is received
	mainLoop.Run()

	// Close the pipeline
	return pipeliner.Close()
}

func getPipeline(src, dest *os.File, pipelineString string) (gstauto.Pipeliner, error) {
	if src != nil && dest != nil {
		return gstauto.NewPipelineReadWriterSimpleFromString(pipelineString)
	}
	if src != nil {
		return gstauto.NewPipelineWriterSimpleFromString(pipelineString)
	}
	if dest != nil {
		return gstauto.NewPipelineReaderSimpleFromString(pipelineString)
	}
	return gstauto.NewPipelinerSimpleFromString(pipelineString)
}
