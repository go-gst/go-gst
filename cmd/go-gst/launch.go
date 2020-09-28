package main

import (
	"errors"
	"io"
	"os"
	"strings"

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
	defer mainLoop.Quit()

	pipelineString := strings.Join(args, " ")

	logInfo("pipeline", "Creating pipeline")

	pipeliner, err := getPipeline(launchSrc, launchDest, pipelineString)
	if err != nil {
		return err
	}

	if verbose {
		setupVerbosePipelineListeners(pipeliner.Pipeline(), "pipeline")
	}

	logInfo("pipeline", "Starting pipeline")
	if err := pipeliner.Start(); err != nil {
		return err
	}

	defer pipeliner.Close()

	if launchSrc != nil {
		pipelineWriter := pipeliner.(gstauto.WritePipeliner)
		go io.Copy(pipelineWriter, launchSrc)
	}
	if launchDest != nil {
		pipelineReader := pipeliner.(gstauto.ReadPipeliner)
		go io.Copy(launchDest, pipelineReader)
	}

	return mainLoop.RunError()
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
