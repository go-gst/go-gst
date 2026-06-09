package gst

import (
	"testing"
)

func TestRegistryList(t *testing.T) {
	// Initialize GStreamer
	Init(nil)

	// Get the default registry
	registry := GetRegistry()
	if registry == nil {
		t.Fatal("Failed to get registry")
	}

	// Get the plugin list
	plugins := registry.GetPluginList()

	t.Logf("Found %d plugins in the registry:", len(plugins))

	// Verify we have at least some plugins
	if len(plugins) == 0 {
		t.Error("Expected to find at least one plugin in the registry")
		return
	}

	// Print first few plugins as examples
	for i, plugin := range plugins {
		if i >= 5 { // Only show first 5 plugins
			t.Logf("... and %d more plugins", len(plugins)-5)
			break
		}

		name := plugin.GetName()
		version := plugin.Version()
		description := plugin.Description()

		t.Logf("%d. %s (v%s) - %s", i+1, name, version, description)

		// Basic validation that plugin has required fields
		if name == "" {
			t.Errorf("Plugin %d has empty name", i+1)
		}
		if version == "" {
			t.Errorf("Plugin %d (%s) has empty version", i+1, name)
		}
	}
}
