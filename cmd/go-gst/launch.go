package main

import (
	"errors"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tinyzimmer/go-gst-launch/gst"
)

func init() {
	rootCmd.AddCommand(launchCmd)
}

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Run a generic pipeline",
	Long:  `Uses the provided pipeline string to encode/decode the data in the pipeline.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("The pipeline string cannot be empty")
		}
		return nil
	},
	RunE: launch,
}

func launch(cmd *cobra.Command, args []string) error {

	src, dest, err := getCLIFiles()
	if err != nil {
		return err
	}

	var flags gst.PipelineFlags
	if src != nil {
		flags = flags | gst.PipelineWrite
	}
	if dest != nil {
		flags = flags | gst.PipelineRead
	}

	pipelineString := strings.Join(args, " ")

	logInfo("pipeline", "Creating pipeline")
	gstPipeline, err := gst.NewPipelineFromLaunchString(pipelineString, flags)
	if err != nil {
		return err
	}

	defer gstPipeline.Close()

	if verbose {
		setupVerbosePipelineListeners(gstPipeline, "pipeline")
	}

	logInfo("pipeline", "Starting pipeline")
	if err := gstPipeline.Start(); err != nil {
		return err
	}

	if src != nil {
		go io.Copy(gstPipeline, src)
	}
	if dest != nil {
		go io.Copy(dest, gstPipeline)
	}

	gst.Wait(gstPipeline)

	return nil
}
