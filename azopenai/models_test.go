// Copyright 2025 herosizy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package azopenai

import (
	"strings"
	"testing"
)

func TestListModels(t *testing.T) {
	models, err := listModels()
	if err != nil {
		t.Fatalf("listModels() returned error: %v", err)
	}

	// Should have a good number of models
	if len(models) < 5 {
		t.Errorf("Expected at least 5 models, got %d", len(models))
	}

	// Check that key models are present
	expectedModels := []string{
		Gpt4o,
		Gpt4oMini,
		Gpt4,
		Gpt35Turbo,
		Gpt4Turbo,
	}

	for _, expectedModel := range expectedModels {
		if _, exists := models[expectedModel]; !exists {
			t.Errorf("Expected model %s not found in models list", expectedModel)
		}
	}

	// Verify model info structure
	for name, model := range models {
		if name == "" {
			t.Error("Model name should not be empty")
		}
		if model.Label == "" {
			t.Error("Model label should not be empty")
		}
		if !strings.Contains(model.Label, "Azure OpenAI") {
			t.Errorf("Model label should contain 'Azure OpenAI', got: %s", model.Label)
		}
		if model.Supports == nil {
			t.Errorf("Model %s should have Supports defined", name)
		}
	}
}

func TestListEmbedders(t *testing.T) {
	embedders, err := listEmbedders()
	if err != nil {
		t.Fatalf("listEmbedders() returned error: %v", err)
	}

	// Should have embedder models
	if len(embedders) < 2 {
		t.Errorf("Expected at least 2 embedders, got %d", len(embedders))
	}

	// Check that key embedders are present
	expectedEmbedders := []string{
		TextEmbedding3Small,
		TextEmbedding3Large,
	}

	embedderSet := make(map[string]bool)
	for _, embedder := range embedders {
		embedderSet[embedder] = true
	}

	for _, expectedEmbedder := range expectedEmbedders {
		if !embedderSet[expectedEmbedder] {
			t.Errorf("Expected embedder %s not found in embedders list", expectedEmbedder)
		}
	}

	// Verify embedder names are not empty
	for _, embedder := range embedders {
		if embedder == "" {
			t.Error("Embedder name should not be empty")
		}
	}
}

func TestModelConstants_Comprehensive(t *testing.T) {
	// Test all major model categories that exist
	reasoningModels := []string{O1Mini, O3Mini, O4Mini}
	flagshipModels := []string{Gpt41, Gpt4o, Gpt4oMini}
	legacyModels := []string{Gpt4, Gpt4Turbo, Gpt35Turbo}
	imageModels := []string{Dalle3}
	embeddingModels := []string{TextEmbedding3Small, TextEmbedding3Large}

	allModels := append(reasoningModels, flagshipModels...)
	allModels = append(allModels, legacyModels...)
	allModels = append(allModels, imageModels...)
	allModels = append(allModels, embeddingModels...)

	for _, model := range allModels {
		if model == "" {
			t.Error("Model constant should not be empty")
		}
		// Model names should be lowercase with dashes or numbers
		if !isValidModelName(model) {
			t.Errorf("Model %s should follow valid naming convention", model)
		}
	}
}

func TestModelCategorization(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		category string
	}{
		{"Reasoning Model", O1Mini, "reasoning"},
		{"Flagship Model", Gpt4o, "flagship"},
		{"Legacy Model", Gpt4, "legacy"},
		{"Image Model", Dalle3, "image"},
		{"Embedding Model", TextEmbedding3Small, "embedding"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the model constant is properly defined
			if tt.model == "" {
				t.Errorf("Model constant for %s should not be empty", tt.category)
			}
		})
	}
}

func TestTextModelCapabilities(t *testing.T) {
	if !TextModel.Multiturn {
		t.Error("TextModel should support multiturn")
	}
	if !TextModel.Tools {
		t.Error("TextModel should support tools")
	}
	if !TextModel.ToolChoice {
		t.Error("TextModel should support tool choice")
	}
	if !TextModel.SystemRole {
		t.Error("TextModel should support system role")
	}
	if TextModel.Media {
		t.Error("TextModel should not support media")
	}
}

func TestMultimodalModelCapabilities(t *testing.T) {
	if !MultimodalModel.Multiturn {
		t.Error("MultimodalModel should support multiturn")
	}
	if !MultimodalModel.Tools {
		t.Error("MultimodalModel should support tools")
	}
	if !MultimodalModel.ToolChoice {
		t.Error("MultimodalModel should support tool choice")
	}
	if !MultimodalModel.SystemRole {
		t.Error("MultimodalModel should support system role")
	}
	if !MultimodalModel.Media {
		t.Error("MultimodalModel should support media")
	}
}

// Helper functions for tests
func isValidModelName(s string) bool {
	if s == "" {
		return false
	}
	// Check if it contains only valid characters (lowercase, numbers, dashes, dots)
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.') {
			return false
		}
	}
	return true
}
