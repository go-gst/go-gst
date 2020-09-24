package gst

import "fmt"

// PipelineConfig represents a list of elements and their configurations
// to be used with NewPipelineFromConfig.
type PipelineConfig struct {
	Elements []*PipelineElement
}

// GetElementByName returns the Element configuration for the given name.
func (p *PipelineConfig) GetElementByName(name string) *PipelineElement {
	for _, elem := range p.Elements {
		if name == elem.GetName() {
			return elem
		}
	}
	return nil
}

// ElementNames returns a string slice of the names of all the plugins.
func (p *PipelineConfig) ElementNames() []string {
	names := make([]string, 0)
	for _, elem := range p.Elements {
		names = append(names, elem.GetName())
	}
	return names
}

// pushPluginToTop pushes a plugin to the top of the list.
func (p *PipelineConfig) pushPluginToTop(elem *PipelineElement) {
	newSlc := []*PipelineElement{elem}
	newSlc = append(newSlc, p.Elements...)
	p.Elements = newSlc
}

// PipelineElement represents an `GstElement` in a `GstPipeline` when building a Pipeline with `NewPipelineFromConfig`.
// The Name should coorespond to a valid gstreamer plugin name. The data are additional
// fields to set on the element. If SinkCaps is non-nil, they are applied to the sink of this
// element.
type PipelineElement struct {
	Name     string
	SinkCaps Caps
	Data     map[string]interface{}
}

// GetName returns the name to use when creating Elements from this configuration.
func (p *PipelineElement) GetName() string { return p.Name }

// NewPipelineFromConfig builds a new pipeline from the given PipelineConfig. The plugins provided
// in the configuration will be linked in the order they are given.
// If using PipelineWrite, you can optionally pass a Caps object to filter between the write-buffer
// and the start of the pipeline.
func NewPipelineFromConfig(cfg *PipelineConfig, flags PipelineFlags, caps Caps) (pipeline *Pipeline, err error) {
	// create a new empty pipeline instance
	pipeline, err = NewPipeline(flags)
	if err != nil {
		return nil, err
	}
	// if any error happens while setting up the pipeline, immediately free it
	defer func() {
		if err != nil {
			if cerr := pipeline.Close(); cerr != nil {
				fmt.Println("Failed to close pipeline:", err)
			}
		}
	}()

	if cfg.Elements == nil {
		cfg.Elements = make([]*PipelineElement, 0)
	}

	if flags.has(PipelineWrite) {
		if flags.has(PipelineUseGstApp) {
			cfg.pushPluginToTop(&PipelineElement{
				Name: "appsrc",
				Data: map[string]interface{}{
					"block":        true,  // TODO: make these all configurable
					"emit-signals": false, // https://gstreamer.freedesktop.org/documentation/app/appsrc.html?gi-language=c
					"is-live":      true,
					"max-bytes":    200000,
					// "size": 0, // If this is known we should specify it
				},
				SinkCaps: caps,
			})
		} else {
			cfg.pushPluginToTop(&PipelineElement{
				Name: "fdsrc",
				Data: map[string]interface{}{
					"fd": pipeline.writerFd(),
				},
				SinkCaps: caps,
			})
		}
	}

	if flags.has(PipelineRead) {
		if flags.has(PipelineUseGstApp) {
			cfg.Elements = append(cfg.Elements, &PipelineElement{
				Name: "appsink",
				Data: map[string]interface{}{
					"emit-signals": false,
				},
			})
		} else {
			cfg.Elements = append(cfg.Elements, &PipelineElement{
				Name: "fdsink",
				Data: map[string]interface{}{
					"fd": pipeline.readerFd(),
				},
			})
		}
	}

	// retrieve a list of the plugin names
	pluginNames := cfg.ElementNames()

	// build all the elements
	var elements map[int]*Element
	elements, err = NewElementMany(pluginNames...)
	if err != nil {
		return
	}

	// iterate the plugin names and add them to the pipeline
	for idx, name := range pluginNames {
		// get the current plugin and element
		currentPlugin := cfg.GetElementByName(name)
		currentElem := elements[idx]

		// Iterate any data with the plugin and set it on the element
		for key, value := range currentPlugin.Data {
			if err = currentElem.Set(key, value); err != nil {
				return
			}
		}

		// Add the element to the pipeline
		if err = pipeline.Add(currentElem); err != nil {
			return
		}

		// If this is the first element continue
		if idx == 0 {
			continue
		}

		// get the last element in the chain
		lastPluginName := pluginNames[idx-1]
		lastElem := elements[idx-1]
		lastPlugin := cfg.GetElementByName(lastPluginName)

		if lastPlugin == nil {
			// this should never happen, since only used internally,
			// but safety from panic
			continue
		}

		// If this is the second element and we are configuring writing
		// call link on the last element
		if idx == 1 && flags.has(PipelineWrite) {
			pipeline.LinkWriterTo(lastElem)
			if flags.has(PipelineUseGstApp) {
				pipeline.appSrc = wrapAppSrc(lastElem)
			}
		}

		// If this is the last element and we are configuring reading
		// call link on the element
		if idx == len(pluginNames)-1 && flags.has(PipelineRead) {
			pipeline.LinkReaderTo(currentElem)
			if flags.has(PipelineUseGstApp) {
				pipeline.appSink = wrapAppSink(currentElem)
			}
		}

		// If there are sink caps on the last element, do a filtered link to this one and continue
		if lastPlugin.SinkCaps != nil {
			if err = lastElem.LinkFiltered(currentElem, lastPlugin.SinkCaps); err != nil {
				return
			}
			continue
		}

		// link the last element to this element
		if err = lastElem.Link(currentElem); err != nil {
			return
		}
	}

	pipeline.pipelineFromHelper = true

	return
}
