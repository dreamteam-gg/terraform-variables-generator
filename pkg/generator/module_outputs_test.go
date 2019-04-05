package generator

import (
	"testing"
)

var testTeraformFolder = "./fixtures/module_outputs"

func TestGetModules(t *testing.T) {
	modules := getModules(testTeraformFolder)

	if len(modules) != 2 {
		t.Errorf("Should find 2 modules, got %d", len(modules))
	}
}

func TestParsingModule(t *testing.T) {
	modules := getModules(testTeraformFolder)

	modules[0].parseOutputs()

	if len(modules[0].Outputs) != 1 {
		t.Errorf("Modules should have one output, got %d", len(modules[0].Outputs))
	}

	if modules[0].Outputs[0].Name != "first_module_first_output" {
		t.Errorf("Modules output should be first_output, got %s", modules[0].Outputs[0])
	}
}
