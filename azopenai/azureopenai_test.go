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
	"context"
	"os"
	"strings"
	"testing"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func TestAzureOpenAI_Name(t *testing.T) {
	plugin := &AzureOpenAI{}
	expected := "azureopenai"

	if got := plugin.Name(); got != expected {
		t.Errorf("AzureOpenAI.Name() = %v, want %v", got, expected)
	}
}

func TestAzureOpenAI_Init_Success(t *testing.T) {
	// Set up test environment
	originalAPIKey := os.Getenv("AZURE_OPEN_AI_API_KEY")
	originalEndpoint := os.Getenv("AZURE_OPEN_AI_ENDPOINT")
	defer func() {
		os.Setenv("AZURE_OPEN_AI_API_KEY", originalAPIKey)
		os.Setenv("AZURE_OPEN_AI_ENDPOINT", originalEndpoint)
	}()

	os.Setenv("AZURE_OPEN_AI_API_KEY", "test-api-key")
	os.Setenv("AZURE_OPEN_AI_ENDPOINT", "https://test.openai.azure.com/")

	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}

	// Note: This will fail without real credentials, but tests the validation logic
	err = plugin.Init(ctx, g)
	// We expect this to fail with auth error, not missing env vars
	if err != nil && err.Error() == "Azure OpenAI requires setting AZURE_OPEN_AI_API_KEY in the environment" {
		t.Errorf("Expected auth error, got missing env var error: %v", err)
	}
}

func TestAzureOpenAI_Init_MissingAPIKey(t *testing.T) {
	// Clear environment variables
	originalAPIKey := os.Getenv("AZURE_OPEN_AI_API_KEY")
	originalEndpoint := os.Getenv("AZURE_OPEN_AI_ENDPOINT")
	defer func() {
		os.Setenv("AZURE_OPEN_AI_API_KEY", originalAPIKey)
		os.Setenv("AZURE_OPEN_AI_ENDPOINT", originalEndpoint)
	}()

	os.Unsetenv("AZURE_OPEN_AI_API_KEY")
	os.Setenv("AZURE_OPEN_AI_ENDPOINT", "https://test.openai.azure.com/")

	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}
	err = plugin.Init(ctx, g)

	expectedErrorContains := "Azure OpenAI requires setting AZURE_OPEN_AI_API_KEY in the environment"
	if err == nil || !strings.Contains(err.Error(), expectedErrorContains) {
		t.Errorf("Expected error containing %q, got %v", expectedErrorContains, err)
	}
}

func TestAzureOpenAI_Init_MissingEndpoint(t *testing.T) {
	// Clear environment variables
	originalAPIKey := os.Getenv("AZURE_OPEN_AI_API_KEY")
	originalEndpoint := os.Getenv("AZURE_OPEN_AI_ENDPOINT")
	defer func() {
		os.Setenv("AZURE_OPEN_AI_API_KEY", originalAPIKey)
		os.Setenv("AZURE_OPEN_AI_ENDPOINT", originalEndpoint)
	}()

	os.Setenv("AZURE_OPEN_AI_API_KEY", "test-api-key")
	os.Unsetenv("AZURE_OPEN_AI_ENDPOINT")

	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}
	err = plugin.Init(ctx, g)

	expectedErrorContains := "Azure OpenAI requires setting AZURE_OPEN_AI_ENDPOINT in the environment"
	if err == nil || !strings.Contains(err.Error(), expectedErrorContains) {
		t.Errorf("Expected error containing %q, got %v", expectedErrorContains, err)
	}
}

func TestAzureOpenAI_Init_WithDirectCredentials(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{
		APIKey:   "direct-api-key",
		Endpoint: "https://direct.openai.azure.com/",
	}

	// This will still fail without real credentials, but tests that direct creds are used
	err = plugin.Init(ctx, g)
	// We don't expect missing env var errors since we provided direct credentials
	if err != nil && (strings.Contains(err.Error(), "Azure OpenAI requires setting AZURE_OPEN_AI_API_KEY in the environment") ||
		strings.Contains(err.Error(), "Azure OpenAI requires setting AZURE_OPEN_AI_ENDPOINT in the environment")) {
		t.Errorf("Direct credentials should be used, but got env var error: %v", err)
	}
}

func TestAzureOpenAI_Init_DoubleInit(t *testing.T) {
	originalAPIKey := os.Getenv("AZURE_OPEN_AI_API_KEY")
	originalEndpoint := os.Getenv("AZURE_OPEN_AI_ENDPOINT")
	defer func() {
		os.Setenv("AZURE_OPEN_AI_API_KEY", originalAPIKey)
		os.Setenv("AZURE_OPEN_AI_ENDPOINT", originalEndpoint)
	}()

	os.Setenv("AZURE_OPEN_AI_API_KEY", "test-api-key")
	os.Setenv("AZURE_OPEN_AI_ENDPOINT", "https://test.openai.azure.com/")

	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}

	// First init (will fail due to auth, but sets initted flag)
	plugin.Init(ctx, g)

	// Second init should fail with "already initialized"
	err = plugin.Init(ctx, g)
	expectedErrorContains := "plugin already initialized"
	if err == nil || !strings.Contains(err.Error(), expectedErrorContains) {
		t.Errorf("Expected error containing %q, got %v", expectedErrorContains, err)
	}
}

func TestAzureOpenAI_DefineModel_NotInitialized(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}

	_, err = plugin.DefineModel(g, "test-model", &ai.ModelInfo{
		Label: "Test Model",
	})

	expectedErrorContains := "AzureOpenAI plugin not initialized"
	if err == nil || !strings.Contains(err.Error(), expectedErrorContains) {
		t.Errorf("Expected error containing %q, got %v", expectedErrorContains, err)
	}
}

func TestIsDefinedEmbedder(t *testing.T) {
	tests := []struct {
		name         string
		embedderName string
		want         bool
	}{
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
		{
			name:         "invalid embedder",
			embedderName: "invalid-embedder",
			want:         false,
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

func TestModelConstants(t *testing.T) {
	// Test that exported constants match expected values
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"GPT-4o", Gpt4o, "gpt-4o"},
		{"GPT-4o Mini", Gpt4oMini, "gpt-4o-mini"},
		{"GPT-4", Gpt4, "gpt-4"},
		{"GPT-3.5 Turbo", Gpt35Turbo, "gpt-3.5-turbo"},
		{"Text Embedding 3 Small", TextEmbedding3Small, "text-embedding-3-small"},
		{"Text Embedding 3 Large", TextEmbedding3Large, "text-embedding-3-large"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant %s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestProvider(t *testing.T) {
	expectedProvider := "azureopenai"
	if Provider != expectedProvider {
		t.Errorf("Provider = %q, want %q", Provider, expectedProvider)
	}
}

func TestLabelPrefix(t *testing.T) {
	expectedPrefix := "Azure OpenAI"
	if LabelPrefix != expectedPrefix {
		t.Errorf("LabelPrefix = %q, want %q", LabelPrefix, expectedPrefix)
	}
}
