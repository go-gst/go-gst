package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/gstauto"
	"golang.org/x/net/websocket"
)

var (
	websocketHost                                                    string
	websocketPort                                                    int
	pulseServer, pulseMonitor, encoding, micName, micFifo, micFormat string
	micSampleRate, micChannels                                       int
)

func init() {
	user, err := user.Current()
	var defaultPulseServer, defaultPulseMonitor string
	if err == nil {
		defaultPulseServer = fmt.Sprintf("/run/user/%s/pulse/native", user.Uid)
	}
	defaultMonitor, err := exec.Command("/bin/sh", "-c", "pactl list sources | grep Name | head -n1  | cut -d ' ' -f2").Output()
	if err == nil {
		defaultPulseMonitor = strings.TrimSpace(string(defaultMonitor))
	}
	websocketCmd.PersistentFlags().StringVarP(&websocketHost, "host", "H", "127.0.0.1", "The host to listen on for websocket connections.")
	websocketCmd.PersistentFlags().IntVarP(&websocketPort, "port", "P", 8080, "The port to listen on for websocket connections.")
	websocketCmd.PersistentFlags().StringVarP(&pulseServer, "pulse-server", "p", defaultPulseServer, "The path to the PulseAudio socket.")
	websocketCmd.PersistentFlags().StringVarP(&pulseMonitor, "pulse-monitor", "d", defaultPulseMonitor, "The monitor device to connect to on the Pulse server. The default device is selected if omitted.")
	websocketCmd.PersistentFlags().StringVarP(&encoding, "encoding", "e", "", `The audio encoding to send to websocket connections. The options are:

	opus (default)
		Serves audio data in webm/opus.
		The MediaSource can consume this format by specifying "audio/webm".
	
	vorbis
		Serves audio data in ogg/vorbis.
`)
	websocketCmd.PersistentFlags().StringVarP(&micFifo, "mic-path", "m", "", "A mic FIFO to write received audio data to, by default, nothing is done with received data.")
	websocketCmd.PersistentFlags().StringVarP(&micName, "mic-name", "n", "virtmic", "The name of the mic fifo device in pulse audio.")
	websocketCmd.PersistentFlags().StringVarP(&micFormat, "mic-format", "f", "S16LE", "The audio format pulse audio expects on the fifo.")
	websocketCmd.PersistentFlags().IntVarP(&micSampleRate, "mic-sample-rate", "r", 16000, "The sample rate pulse audio expects on the fifo.")
	websocketCmd.PersistentFlags().IntVarP(&micChannels, "mic-channels", "c", 1, "The number of channels pulse audio expects on the fifo.")

	rootCmd.AddCommand(websocketCmd)
}

var websocketCmd = &cobra.Command{
	Use: "websocket",
	Short: `Run a websocket audio proxy for streaming audio from a pulse server 
              and optionally recording to a virtual mic.`,
	Long: `Starts a websocket server with the given configurations.

This currently only works with PulseAudio or an input file, but may be expanded to be cross-platform.

This command may be expanded to include video support via RFB or RTP.

To use with the mic support, you should first set up a virtual device with something like:

	pactl load-module module-pipe-source source_name=virtmic file=/tmp/mic.fifo format=s16le rate=16000 channels=1
	
And then you can run this command with --mic-path /tmp/mic.fifo. The received data will be expected to
be in the same format as specified with --encoding.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if pulseServer == "" {
			return errors.New("Could not determine pulse server, you should use --pulse-server")
		}
		switch encoding {
		case "opus":
		case "vorbis":
		case "":
			encoding = "opus"
		default:
			return fmt.Errorf("Not a valid audio encoder: %s", encoding)
		}
		return nil
	},
	RunE: websocketProxy,
}

func websocketProxy(cmd *cobra.Command, args []string) error {
	addr := fmt.Sprintf("%s:%d", websocketHost, websocketPort)
	server := &http.Server{
		Handler: &websocket.Server{
			Handshake: func(*websocket.Config, *http.Request) error { return nil },
			Handler:   handleWebsocketConnection,
		},
		Addr: addr,
	}
	logInfo("websocket", "Listening on", addr)
	return server.ListenAndServe()
}

func handleWebsocketConnection(wsconn *websocket.Conn) {
	defer func() {
		wsconn.Close()
		logInfo("websocket", "Connection to", wsconn.Request().RemoteAddr, "closed")
	}()

	logInfo("websocket", "New connection from", wsconn.Request().RemoteAddr)
	wsconn.PayloadType = websocket.BinaryFrame

	playbackPipeline, err := newPlaybackPipeline()
	if err != nil {
		logInfo("websocket", "ERROR:", err.Error())
		return
	}

	logInfo("websocket", "Starting playback pipeline")

	if err = playbackPipeline.Start(); err != nil {
		logInfo("websocket", "ERROR:", err.Error())
		return
	}

	if verbose {
		setupVerbosePipelineListeners(playbackPipeline.Pipeline(), "playback")
	}

	var recordingPipeline gstauto.WritePipeliner
	var sinkPipeline gstauto.Pipeliner
	if micFifo != "" {
		recordingPipeline, err := newRecordingPipeline()
		if err != nil {
			logInfo("websocket", "Could not open pipeline for recording:", err.Error())
			return
		}
		defer recordingPipeline.Close()
		sinkPipeline, err := newSinkPipeline()
		if err != nil {
			logInfo("websocket", "Could not open null sink pipeling. Disabling recording.")
			return
		}
		defer sinkPipeline.Close()
	}

	if recordingPipeline != nil && sinkPipeline != nil {
		logInfo("websocket", "Starting recording pipeline")
		if err = recordingPipeline.Start(); err != nil {
			logInfo("websocket", "Could not start recording pipeline")
			return
		}
		logInfo("websocket", "Starting sink pipeline")
		if err = sinkPipeline.Start(); err != nil {
			logInfo("websocket", "Could not start sink pipeline")
			return
		}

		if verbose {
			setupVerbosePipelineListeners(sinkPipeline.Pipeline(), "mic-null-sink")
		}

		var runMicFunc func()
		runMicFunc = func() {
			if verbose {
				setupVerbosePipelineListeners(recordingPipeline.Pipeline(), "recorder")
			}
			go io.Copy(recordingPipeline, wsconn)
			go func() {
				var lastState gst.State
				for msg := range recordingPipeline.Pipeline().GetBus().MessageChan() {
					defer msg.Unref()
					switch msg.Type() {
					case gst.MessageStateChanged:
						if lastState == gst.StatePlaying && recordingPipeline.Pipeline().GetState() != gst.StatePlaying {
							var nerr error
							recordingPipeline.Close()
							recordingPipeline, nerr = newRecordingPipeline()
							if nerr != nil {
								logInfo("websocket", "Could not create new recording pipeline, stopping input stream")
								return
							}
							logInfo("websocket", "Restarting recording pipeline")
							if nerr = recordingPipeline.Start(); nerr != nil {
								logInfo("websocket", "Could not start new recording pipeline, stopping input stream")
							}
							runMicFunc()
							return
						}
						lastState = recordingPipeline.Pipeline().GetState()
					}
				}
			}()
		}
		runMicFunc()
	}

	defer playbackPipeline.Close()

	go func() {
		io.Copy(wsconn, playbackPipeline)
		// signal the pipeline to do a clean close. This will cause the
		// wait below to break.
		playbackPipeline.Pipeline().SetState(gst.StateNull)
	}()

	if srcFile != "" {
		srcFile, err := getSrcFile()
		if err != nil {
			return
		}
		defer srcFile.Close()
		// stat, err := srcFile.Stat()
		// if err != nil {
		// 	return
		// }
		go io.Copy(playbackPipeline.(gstauto.ReadWritePipeliner), srcFile)
	}

	gst.Wait(playbackPipeline.Pipeline())
}

func newPlaybackPipelineFromString() (gstauto.ReadWritePipeliner, error) {
	pipelineString := "decodebin ! audioconvert ! audioresample"

	switch encoding {
	case "opus":
		pipelineString = fmt.Sprintf("%s ! cutter ! opusenc ! webmmux", pipelineString)
	case "vorbis":
		pipelineString = fmt.Sprintf("%s ! vorbisenc ! oggmux", pipelineString)
	}

	if verbose {
		logInfo("playback", "Using pipeline string", pipelineString)
	}

	return gstauto.NewPipelineReadWriterSimpleFromString(pipelineString)
}

func newPlaybackPipeline() (gstauto.ReadPipeliner, error) {

	if srcFile != "" {
		return newPlaybackPipelineFromString()
	}

	cfg := &gstauto.PipelineConfig{Elements: []*gstauto.PipelineElement{}}

	pulseSrc := &gstauto.PipelineElement{
		Name:     "pulsesrc",
		Data:     map[string]interface{}{"server": pulseServer},
		SinkCaps: gst.NewRawCaps("S16LE", 24000, 2),
	}

	if pulseMonitor != "" {
		pulseSrc.Data["device"] = pulseMonitor
	}

	cfg.Elements = append(cfg.Elements, pulseSrc)

	switch encoding {
	case "opus":
		cfg.Elements = append(cfg.Elements, &gstauto.PipelineElement{Name: "cutter"})
		cfg.Elements = append(cfg.Elements, &gstauto.PipelineElement{Name: "opusenc"})
		cfg.Elements = append(cfg.Elements, &gstauto.PipelineElement{Name: "webmmux"})
	case "vorbis":
		cfg.Elements = append(cfg.Elements, &gstauto.PipelineElement{Name: "vorbisenc"})
		cfg.Elements = append(cfg.Elements, &gstauto.PipelineElement{Name: "oggmux"})
	}

	return gstauto.NewPipelineReaderSimpleFromConfig(cfg)
}

func newRecordingPipeline() (gstauto.WritePipeliner, error) {
	return gstauto.NewPipelineWriterSimpleFromString(newPipelineStringFromOpts())
}

func newPipelineStringFromOpts() string {
	return fmt.Sprintf(
		"decodebin ! audioconvert ! audioresample ! audio/x-raw, format=%s, rate=%d, channels=%d ! filesink location=%s append=true",
		micFormat,
		micSampleRate,
		micChannels,
		micFifo,
	)
}

func newSinkPipeline() (gstauto.Pipeliner, error) {
	cfg := &gstauto.PipelineConfig{
		Elements: []*gstauto.PipelineElement{
			{
				Name:     "pulsesrc",
				Data:     map[string]interface{}{"server": pulseServer, "device": micName},
				SinkCaps: gst.NewRawCaps(micFormat, micSampleRate, micChannels),
			},
			{
				Name: "fakesink",
				Data: map[string]interface{}{"sync": false},
			},
		},
	}
	return gstauto.NewPipelinerSimpleFromConfig(cfg)
}
