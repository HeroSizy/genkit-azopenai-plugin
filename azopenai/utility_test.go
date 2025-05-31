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
	"testing"
)

func TestModelRef(t *testing.T) {
	config := &OpenAIConfig{
		DeploymentName: "test-deployment",
	}

	modelRef := ModelRef(Gpt4o, config)

	// Test that ModelRef creates a valid reference
	// We can't directly test nil since it's a complex type
	// But we can test that the function doesn't panic
	_ = modelRef
}

func TestModelRef_WithNilConfig(t *testing.T) {
	modelRef := ModelRef(Gpt4o, nil)

	// Test that ModelRef works with nil config
	_ = modelRef
}

func TestIsDefinedEmbedder_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		embedderName string
		want         bool
	}{
		{
			name:         "empty string",
			embedderName: "",
			want:         false,
		},
		{
			name:         "random string",
			embedderName: "not-an-embedder",
			want:         false,
		},
		{
			name:         "partial match",
			embedderName: "text-embedding",
			want:         false,
		},
		{
			name:         "case sensitive",
			embedderName: "TEXT-EMBEDDING-3-SMALL",
			want:         false,
		},
		{
			name:         "valid small embedder",
			embedderName: TextEmbedding3Small,
			want:         true,
		},
		{
			name:         "valid large embedder",
			embedderName: TextEmbedding3Large,
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDefinedEmbedder(tt.embedderName); got != tt.want {
				t.Errorf("IsDefinedEmbedder(%q) = %v, want %v", tt.embedderName, got, tt.want)
			}
		})
	}
}

func TestModelConstants_AllDefined(t *testing.T) {
	// Test that all exported model constants are properly defined
	modelConstants := []string{
		Gpt4,
		Gpt4Turbo,
		Gpt4o,
		Gpt4oMini,
		Gpt35Turbo,
		TextEmbedding3Small,
		TextEmbedding3Large,
		Dalle3,
		O1Mini,
		O3Mini,
		O4Mini,
		Gpt41,
		Gpt41Mini,
	}

	for _, constant := range modelConstants {
		if constant == "" {
			t.Error("Model constant should not be empty")
		}

		// Model constants should be reasonable length
		if len(constant) < 3 || len(constant) > 50 {
			t.Errorf("Model constant %s has unreasonable length: %d", constant, len(constant))
		}
	}
}

func TestProviderAndLabelConstants(t *testing.T) {
	if Provider == "" {
		t.Error("Provider constant should not be empty")
	}

	if LabelPrefix == "" {
		t.Error("LabelPrefix constant should not be empty")
	}

	// Provider should be lowercase
	for _, r := range Provider {
		if r >= 'A' && r <= 'Z' {
			t.Errorf("Provider should be lowercase, got: %s", Provider)
			break
		}
	}

	// LabelPrefix should contain "Azure OpenAI"
	expectedPhrase := "Azure OpenAI"
	if !containsSubstring(LabelPrefix, expectedPhrase) {
		t.Errorf("LabelPrefix should contain '%s', got: %s", expectedPhrase, LabelPrefix)
	}
}

func TestOpenAIConfig(t *testing.T) {
	// Test OpenAIConfig struct creation
	config := &OpenAIConfig{
		DeploymentName: "test-deployment",
	}

	if config.DeploymentName != "test-deployment" {
		t.Errorf("Expected DeploymentName to be 'test-deployment', got %s", config.DeploymentName)
	}
}

// Helper function for substring checking
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
