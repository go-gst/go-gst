package gstauto

import (
	"errors"

	"github.com/tinyzimmer/go-gst/gst"
)

// PipelineElement represents an `GstElement` in a `GstPipeline` when building a Pipeline with `NewPipelineFromConfig`.
// The Name should coorespond to a valid gstreamer plugin name. The data are additional
// fields to set on the element. If SinkCaps is non-nil, they are applied to the sink of this
// element.
type PipelineElement struct {
	Name     string
	SinkCaps *gst.Caps
	Data     map[string]interface{}
}

// GetName returns the name to use when creating Elements from this configuration.
func (p *PipelineElement) GetName() string { return p.Name }

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

// PushPluginToTop pushes a plugin to the top of the list of elements.
func (p *PipelineConfig) PushPluginToTop(elem *PipelineElement) {
	newSlc := []*PipelineElement{elem}
	newSlc = append(newSlc, p.Elements...)
	p.Elements = newSlc
}

// Apply applies this configuration to the given Pipeline.
func (p *PipelineConfig) Apply(pipeline *gst.Pipeline) error {
	// build all the elements
	elementNames := p.ElementNames()
	elements, err := gst.NewElementMany(elementNames...)
	if err != nil {
		return err
	}
	// iterate the element names and add them to the pipeline
	for idx, name := range elementNames {
		// get the current config and element
		currentCfg := p.GetElementByName(name)
		currentElem := elements[idx]

		// Iterate any data with the plugin and set it on the element
		for key, value := range currentCfg.Data {
			if err := currentElem.Set(key, value); err != nil {
				return err
			}
		}

		// Add the element to the pipeline
		if err := pipeline.Add(currentElem); err != nil {
			return err
		}

		// If this is the first element continue
		if idx == 0 {
			continue
		}

		// get the last element in the chain
		lastElemName := elementNames[idx-1]
		lastElem := elements[idx-1]
		lastCfg := p.GetElementByName(lastElemName)

		if lastCfg == nil {
			// this would never happen unless someone is messing with memory,
			// but safety from panic
			continue
		}

		// If there are sink caps on the last element, do a filtered link to this one and continue
		if lastCfg.SinkCaps != nil {
			if err := lastElem.LinkFiltered(currentElem, lastCfg.SinkCaps); err != nil {
				return err
			}
			continue
		}

		// link the last element to this element
		if err := lastElem.Link(currentElem); err != nil {
			return err
		}
	}

	return nil
}

// NewPipelineFromConfig builds a new pipeline from the given PipelineConfig. The plugins provided
// in the configuration will be linked in the order they are given.
func NewPipelineFromConfig(cfg *PipelineConfig) (*gst.Pipeline, error) {
	if cfg.Elements == nil {
		return nil, errors.New("Element cannot be empty in the configuration")
	}

	// create a new empty pipeline instance
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	if err = cfg.Apply(pipeline); err != nil {
		runOrPrintErr(pipeline.Destroy)
		return nil, err
	}

	return pipeline, nil
}
