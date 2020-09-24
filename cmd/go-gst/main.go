package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tinyzimmer/go-gst-launch/gst"
)

var (
	srcFile, destFile, pipelineStr string
	verbose, fromStdin, toStdout   bool

	rootCmd = &cobra.Command{
		Use:   "go-gst",
		Short: "A command-line audio/video encoder and decoder based on gstreamer",
		Long: `Go-gst is a CLI utility aiming to implement the core functionality
of the core gstreamer-tools. It's primary purpose is to showcase the functionality of 
the underlying go-gst library.

There are also additional commands showing some of the things you can do with the library,
such as websocket servers reading/writing to/from local audio servers and audio/video/image
encoders/decoders.
`,
	}
)

func init() {
	gst.Init()

	rootCmd.PersistentFlags().StringVarP(&srcFile, "input", "i", "", "An input file, defaults to the first element in the pipeline.")
	rootCmd.PersistentFlags().StringVarP(&destFile, "output", "o", "", "An output file, defaults to the last element in the pipeline.")
	rootCmd.PersistentFlags().BoolVarP(&fromStdin, "from-stdin", "I", false, "Write to the pipeline from stdin. If this is specified, then -i is ignored.")
	rootCmd.PersistentFlags().BoolVarP(&toStdout, "to-stdout", "O", false, "Writes the results from the pipeline to stdout. If this is specified, then -o is ignored.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output. This is ignored when used with --to-stdout.")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		log.Println("ERROR:", err.Error())
	}
}
