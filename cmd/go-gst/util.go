package main

import (
	"os"

	"github.com/tinyzimmer/go-gst/gst"
)

func getSrcFile() (*os.File, error) {
	if fromStdin || srcFile != "" {
		if fromStdin {
			if verbose {
				logInfo("file", "Reading media data from stdin")
			}
			return os.Stdin, nil
		}
		if verbose {
			logInfo("file", "Reading media data from", srcFile)
		}
		return os.Open(srcFile)
	}
	// Commands should do internal checking before calling this command
	return nil, nil
}

func getDestFile() (*os.File, error) {
	if toStdout || destFile != "" {
		if toStdout {
			if verbose {
				logInfo("file", "Writing media output to stdout")
			}
			return os.Stdout, nil
		}
		if verbose {
			logInfo("file", "Writing media output to", destFile)
		}
		return os.Create(destFile)
	}
	// Commands should do internal checking before calling this command
	return nil, nil
}

func getCLIFiles() (src, dest *os.File, err error) {
	src, err = getSrcFile()
	if err != nil {
		return nil, nil, err
	}
	dest, err = getDestFile()
	if err != nil {
		src.Close()
		return nil, nil, err
	}
	// Commands should do internal checking before calling this command
	return src, dest, nil
}

func setupVerbosePipelineListeners(gstPipeline *gst.Pipeline, name string) {
	logInfo(name, "Starting message listeners")
	go func() {
		var currentState gst.State
		for msg := range gstPipeline.GetPipelineBus().MessageChan() {

			defer msg.Unref()

			switch msg.Type() {

			case gst.MessageStreamStart:
				logInfo(name, "Stream has started")
			case gst.MessageEOS:
				logInfo(name, "Stream has ended")
			case gst.MessageStateChanged:
				if currentState != gstPipeline.GetState() {
					logInfo(name, "New pipeline state:", gstPipeline.GetState().String())
					currentState = gstPipeline.GetState()
				}
			case gst.MessageInfo:
				info := msg.ParseInfo()
				logInfo(name, info.Message())
				for k, v := range info.Structure().Values() {
					logInfo(name, k, ":", v)
				}
			case gst.MessageWarning:
				info := msg.ParseWarning()
				logInfo(name, "WARNING:", info.Message())
				for k, v := range info.Structure().Values() {
					logInfo(name, k, ":", v)
				}

			case gst.MessageError:
				err := msg.ParseError()
				logInfo(name, "ERROR:", err.Error())
				logInfo(name, "DEBUG:", err.DebugString())

			}
		}
	}()

}
