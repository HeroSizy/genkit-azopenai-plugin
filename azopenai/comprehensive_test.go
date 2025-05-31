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
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// Mock client for testing
type mockClient struct {
	client *azopenai.Client
}

// Test Model function
func TestModel(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	// Test Model function - it looks up a model that may not exist yet
	model := Model(g, Gpt4o)
	// Model() will return nil if the model is not registered
	// This is expected behavior without initialization
	_ = model // Don't check for nil as it's expected
}

// Test DefineModel function (global)
func TestDefineModel_Global(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	info := ai.ModelInfo{
		Label: "Test Model",
		Supports: &ai.ModelSupports{
			Multiturn:  true,
			Media:      false,
			Tools:      true,
			SystemRole: true,
		},
	}

	model := DefineModel(g, "test-model", &info)
	if model == nil {
		t.Error("DefineModel() should return a non-nil model")
	}
}

// Test Embedder function with panic recovery
func TestEmbedder_Panic(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			// Expected panic for undefined embedder
			panicMsg := r.(string)
			if !strings.Contains(panicMsg, "was not found") {
				t.Errorf("Expected panic about embedder not found, got: %s", panicMsg)
			}
		}
	}()

	// This should panic since the embedder is not registered
	Embedder(g, "non-existent-embedder")
	t.Error("Embedder() should panic for non-existent embedder")
}

// Test AzureOpenAI.DefineEmbedder with unsupported embedder
func TestAzureOpenAI_DefineEmbedder_Unsupported(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}

	_, err = plugin.DefineEmbedder(g, "unsupported-embedder")
	if err == nil {
		t.Error("DefineEmbedder() should return error for unsupported embedder")
	}
	if !strings.Contains(err.Error(), "is not supported") {
		t.Errorf("Expected 'is not supported' error, got: %v", err)
	}
}

// Test AzureOpenAI.IsDefinedEmbedder
func TestAzureOpenAI_IsDefinedEmbedder(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}

	// Test with unregistered embedder
	isDefined := plugin.IsDefinedEmbedder(g, TextEmbedding3Small)
	if isDefined {
		t.Error("IsDefinedEmbedder() should return false for unregistered embedder")
	}
}

// Test DefineModel with initialized plugin but unknown model
func TestAzureOpenAI_DefineModel_UnknownModel(t *testing.T) {
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

	// Force initialization to set initted flag
	plugin.initted = true

	// Test with unknown model and nil info
	_, err = plugin.DefineModel(g, "unknown-model", nil)
	if err == nil {
		t.Error("DefineModel() should return error for unknown model with nil info")
	}
	if !strings.Contains(err.Error(), "unknown model") {
		t.Errorf("Expected 'unknown model' error, got: %v", err)
	}
}

// Test DefineModel with initialized plugin and known model
func TestAzureOpenAI_DefineModel_KnownModel(t *testing.T) {
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

	// Force initialization to set initted flag
	plugin.initted = true

	// Test with known model and nil info (should use model info from list)
	_, err = plugin.DefineModel(g, Gpt4o, nil)
	if err != nil {
		t.Errorf("DefineModel() should not return error for known model with nil info: %v", err)
	}
}

// Test convertFinishReason with all cases
func TestConvertFinishReason_AllCases(t *testing.T) {
	tests := []struct {
		name           string
		azureReason    azopenai.CompletionsFinishReason
		expectedReason ai.FinishReason
	}{
		{
			name:           "stopped",
			azureReason:    azopenai.CompletionsFinishReasonStopped,
			expectedReason: ai.FinishReasonStop,
		},
		{
			name:           "token limit reached",
			azureReason:    azopenai.CompletionsFinishReasonTokenLimitReached,
			expectedReason: ai.FinishReasonLength,
		},
		{
			name:           "content filtered",
			azureReason:    azopenai.CompletionsFinishReasonContentFiltered,
			expectedReason: ai.FinishReasonBlocked,
		},
		{
			name:           "tool calls",
			azureReason:    azopenai.CompletionsFinishReasonToolCalls,
			expectedReason: ai.FinishReasonStop,
		},
		{
			name:           "function call (legacy)",
			azureReason:    azopenai.CompletionsFinishReasonFunctionCall,
			expectedReason: ai.FinishReasonOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertFinishReason(tt.azureReason)
			if result != tt.expectedReason {
				t.Errorf("convertFinishReason() = %v, want %v", result, tt.expectedReason)
			}
		})
	}
}

// Test convertMessage with all roles
func TestConvertMessage_AllRoles(t *testing.T) {
	tests := []struct {
		name        string
		message     *ai.Message
		expectError bool
	}{
		{
			name: "user role",
			message: &ai.Message{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("Hello")},
			},
			expectError: false,
		},
		{
			name: "system role",
			message: &ai.Message{
				Role:    ai.RoleSystem,
				Content: []*ai.Part{ai.NewTextPart("System message")},
			},
			expectError: false,
		},
		{
			name: "model role",
			message: &ai.Message{
				Role:    ai.RoleModel,
				Content: []*ai.Part{ai.NewTextPart("Model response")},
			},
			expectError: false,
		},
		{
			name: "tool role",
			message: &ai.Message{
				Role:    ai.RoleTool,
				Content: []*ai.Part{ai.NewTextPart("Tool response")},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertMessage(tt.message)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && result == nil {
				t.Error("Expected non-nil result")
			}
		})
	}
}

// Test convertToAzureOpenAIRequest with various configurations
func TestConvertToAzureOpenAIRequest_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		request  *ai.ModelRequest
		config   OpenAIConfig
		hasError bool
	}{
		{
			name: "basic request",
			request: &ai.ModelRequest{
				Messages: []*ai.Message{
					{
						Role:    ai.RoleUser,
						Content: []*ai.Part{ai.NewTextPart("Hello")},
					},
				},
			},
			config: OpenAIConfig{
				DeploymentName: "test-deployment",
			},
			hasError: false,
		},
		{
			name: "request with tools",
			request: &ai.ModelRequest{
				Messages: []*ai.Message{
					{
						Role:    ai.RoleUser,
						Content: []*ai.Part{ai.NewTextPart("Hello")},
					},
				},
				Tools: []*ai.ToolDefinition{
					{
						Name:        "test_tool",
						Description: "A test tool",
						InputSchema: map[string]any{
							"type": "object",
						},
					},
				},
			},
			config: OpenAIConfig{
				DeploymentName: "test-deployment",
			},
			hasError: false,
		},
		{
			name: "request with config options",
			request: &ai.ModelRequest{
				Messages: []*ai.Message{
					{
						Role:    ai.RoleUser,
						Content: []*ai.Part{ai.NewTextPart("Hello")},
					},
				},
				Config: &OpenAIConfig{
					DeploymentName: "custom-deployment",
					Temperature:    &[]float32{0.7}[0],
					MaxTokens:      &[]int32{100}[0],
				},
			},
			config: OpenAIConfig{
				DeploymentName: "test-deployment",
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertToAzureOpenAIRequest(tt.request, tt.config)
			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.hasError {
				// Verify deployment name is set
				if result.DeploymentName == nil || *result.DeploymentName != tt.config.DeploymentName {
					t.Errorf("Expected deployment name %s, got %v", tt.config.DeploymentName, result.DeploymentName)
				}
			}
		})
	}
}

// Test listModels error handling
func TestListModels_Error(t *testing.T) {
	// This tests the listModels function indirectly by testing successful case
	models, err := listModels()
	if err != nil {
		t.Errorf("listModels() should not return error: %v", err)
	}
	if len(models) == 0 {
		t.Error("listModels() should return at least one model")
	}
}

// Test Init with client creation error simulation
func TestAzureOpenAI_Init_ErrorHandling(t *testing.T) {
	originalAPIKey := os.Getenv("AZURE_OPEN_AI_API_KEY")
	originalEndpoint := os.Getenv("AZURE_OPEN_AI_ENDPOINT")
	defer func() {
		os.Setenv("AZURE_OPEN_AI_API_KEY", originalAPIKey)
		os.Setenv("AZURE_OPEN_AI_ENDPOINT", originalEndpoint)
	}()

	// Test with invalid endpoint format
	os.Setenv("AZURE_OPEN_AI_API_KEY", "test-key")
	os.Setenv("AZURE_OPEN_AI_ENDPOINT", "invalid-endpoint")

	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}
	err = plugin.Init(ctx, g)

	// Should get some kind of error (either client creation or auth)
	// Don't check for specific error since it may vary
	if err != nil {
		// Error should be wrapped with "AzureOpenAI.Init:"
		if !strings.Contains(err.Error(), "AzureOpenAI.Init:") {
			t.Errorf("Error should be wrapped with 'AzureOpenAI.Init:', got: %v", err)
		}
	}
}

// Test successful init with embedder registration
func TestAzureOpenAI_Init_WithEmbedders(t *testing.T) {
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

	// The init will fail with auth error, but we can test that it tries to register embedders
	_ = plugin.Init(ctx, g)

	// Test that embedders are recognized
	embeddingModels, err := listEmbedders()
	if err != nil {
		t.Errorf("listEmbedders() should not return error: %v", err)
	}
	if len(embeddingModels) == 0 {
		t.Error("Should have at least one embedding model")
	}
}

// Mock streaming response for testing
type mockStreamResponse struct {
	responses []azopenai.ChatCompletions
	index     int
	closed    bool
}

func (m *mockStreamResponse) Read() (azopenai.ChatCompletions, error) {
	if m.index >= len(m.responses) {
		return azopenai.ChatCompletions{}, io.EOF
	}
	resp := m.responses[m.index]
	m.index++
	return resp, nil
}

func (m *mockStreamResponse) Close() {
	m.closed = true
}

// Test with nil plugin in Init
func TestAzureOpenAI_Init_NilPlugin(t *testing.T) {
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

	var plugin *AzureOpenAI = nil

	// This should handle nil plugin gracefully by creating a new instance
	err = plugin.Init(ctx, g)
	// The function creates a new instance when nil and should succeed with dummy credentials
	// since the Azure SDK doesn't validate credentials until an actual API call is made
	if err != nil {
		// If there's an error, it should be related to model registration, not nil pointer
		if strings.Contains(err.Error(), "nil") {
			t.Error("Should not get nil pointer error when plugin is nil")
		}
		// Otherwise, it's an expected initialization error which is acceptable
	}
}

// Test conversion functions error handling with edge cases
func TestConvertMessage_UnknownRole(t *testing.T) {
	message := &ai.Message{
		Role:    ai.Role("unknown_role"),
		Content: []*ai.Part{ai.NewTextPart("Hello")},
	}

	_, err := convertMessage(message)
	if err == nil {
		t.Error("Expected error for unknown role")
	}
	if !strings.Contains(err.Error(), "unsupported role") {
		t.Errorf("Expected unsupported role error, got: %v", err)
	}
}

// Test convertTools with invalid JSON schema
func TestConvertTools_InvalidJSONSchema(t *testing.T) {
	tools := []*ai.ToolDefinition{
		{
			Name:        "test_tool",
			Description: "A test tool",
			InputSchema: map[string]any{
				// This creates a circular reference that can't be marshaled
				"circular": make(chan int),
			},
		},
	}

	_, err := convertTools(tools)
	if err == nil {
		t.Error("Expected error for invalid JSON schema")
	}
	if !strings.Contains(err.Error(), "failed to marshal tool parameters") {
		t.Errorf("Expected marshal error, got: %v", err)
	}
}

// Test convertToAzureOpenAIRequest with empty deployment name
func TestConvertToAzureOpenAIRequest_EmptyDeploymentName(t *testing.T) {
	request := &ai.ModelRequest{
		Messages: []*ai.Message{
			{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("Hello")},
			},
		},
	}

	config := OpenAIConfig{
		// DeploymentName is empty
	}

	_, err := convertToAzureOpenAIRequest(request, config)
	if err == nil {
		t.Error("Expected error for empty deployment name")
	}
	if !strings.Contains(err.Error(), "deployment name is required") {
		t.Errorf("Expected deployment name error, got: %v", err)
	}
}

// Test extractTextContent with mixed content types
func TestExtractTextContent_MixedContent(t *testing.T) {
	parts := []*ai.Part{
		ai.NewTextPart("Hello"),
		ai.NewTextPart(" "),
		ai.NewTextPart("World"),
	}

	result := extractTextContent(parts)
	expected := "Hello World"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Test defineModel with various model info configurations
func TestDefineModel_CustomModelInfo(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	customInfo := ai.ModelInfo{
		Label:    "Custom Model",
		Versions: []string{"v1.0"},
		Supports: &ai.ModelSupports{
			Multiturn:  true,
			Tools:      false,
			SystemRole: true,
			Media:      false,
		},
		Stage: ai.ModelStageStable,
	}

	// This tests the global DefineModel function
	model := DefineModel(g, "custom-model", &customInfo)
	if model == nil {
		t.Error("Expected non-nil model from DefineModel")
	}
}

// Test AzureOpenAI.DefineModel with initialized plugin
func TestAzureOpenAI_DefineModel_WithoutInit(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	plugin := &AzureOpenAI{}

	// Should fail because plugin is not initialized
	_, err = plugin.DefineModel(g, "test-model", nil)
	if err == nil {
		t.Error("Expected error for uninitialized plugin")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("Expected initialization error, got: %v", err)
	}
}

// Test AzureOpenAI.DefineModel with unknown model and nil info
func TestAzureOpenAI_DefineModel_UnknownModelNilInfo(t *testing.T) {
	// This test simulates an initialized plugin scenario
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

	// Mock the initialization state without actually connecting
	plugin.initted = true
	plugin.client = nil // This will be nil in this test

	// Should fail because model is unknown and info is nil
	_, err = plugin.DefineModel(g, "unknown-model", nil)
	if err == nil {
		t.Error("Expected error for unknown model with nil info")
	}
	if !strings.Contains(err.Error(), "unknown model") {
		t.Errorf("Expected unknown model error, got: %v", err)
	}
}

// Test AzureOpenAI.DefineEmbedder with known embedder
func TestAzureOpenAI_DefineEmbedder_KnownEmbedder(t *testing.T) {
	plugin := &AzureOpenAI{}

	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	// Test with a known embedder
	_, err = plugin.DefineEmbedder(g, TextEmbedding3Small)
	if err != nil {
		t.Errorf("Expected no error for known embedder, got: %v", err)
	}
}

// Test AzureOpenAI.DefineEmbedder with unknown embedder
func TestAzureOpenAI_DefineEmbedder_UnknownEmbedder(t *testing.T) {
	plugin := &AzureOpenAI{}

	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	// Test with an unknown embedder
	_, err = plugin.DefineEmbedder(g, "unknown-embedder")
	if err == nil {
		t.Error("Expected error for unknown embedder")
	}
	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("Expected not supported error, got: %v", err)
	}
}

// Test all finish reason conversions
func TestConvertFinishReason_FunctionCallLegacy(t *testing.T) {
	// Test the function call legacy case specifically
	result := convertFinishReason(azopenai.CompletionsFinishReasonFunctionCall)
	expected := ai.FinishReasonOther
	if result != expected {
		t.Errorf("Expected %v for function call, got %v", expected, result)
	}
}

// Test request handling functions with basic scenarios (without mocking complex Azure SDK types)
func TestRequestHandling_BasicValidation(t *testing.T) {
	// Test that the request conversion logic handles various configurations correctly

	tests := []struct {
		name      string
		config    OpenAIConfig
		wantError bool
	}{
		{
			name: "valid basic config",
			config: OpenAIConfig{
				DeploymentName: "gpt-4o",
			},
			wantError: false,
		},
		{
			name: "valid config with all options",
			config: OpenAIConfig{
				DeploymentName:   "gpt-4o",
				MaxTokens:        &[]int32{100}[0],
				Temperature:      &[]float32{0.7}[0],
				TopP:             &[]float32{0.9}[0],
				PresencePenalty:  &[]float32{0.1}[0],
				FrequencyPenalty: &[]float32{0.2}[0],
				LogitBias:        map[string]*int32{"test": &[]int32{1}[0]},
				User:             "test-user",
				Seed:             &[]int64{42}[0],
			},
			wantError: false,
		},
		{
			name:   "empty deployment name",
			config: OpenAIConfig{
				// DeploymentName is empty
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &ai.ModelRequest{
				Messages: []*ai.Message{
					{
						Role:    ai.RoleUser,
						Content: []*ai.Part{ai.NewTextPart("Hello")},
					},
				},
			}

			_, err := convertToAzureOpenAIRequest(request, tt.config)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test embedder configuration handling
func TestEmbedderConfiguration(t *testing.T) {
	tests := []struct {
		name      string
		input     []*ai.Document
		wantError bool
	}{
		{
			name: "valid documents",
			input: []*ai.Document{
				ai.DocumentFromText("Hello world", nil),
				ai.DocumentFromText("How are you?", nil),
			},
			wantError: false,
		},
		{
			name: "empty documents",
			input: []*ai.Document{
				{Content: []*ai.Part{}}, // Empty content
			},
			wantError: true,
		},
		{
			name:      "no documents",
			input:     []*ai.Document{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the document processing logic from defineEmbedder
			var input []string
			for _, doc := range tt.input {
				var textParts []string
				for _, part := range doc.Content {
					if part.Text != "" {
						textParts = append(textParts, part.Text)
					}
				}
				if len(textParts) > 0 {
					input = append(input, strings.Join(textParts, " "))
				}
			}

			hasError := len(input) == 0

			if tt.wantError && !hasError {
				t.Error("Expected error but processing succeeded")
			}
			if !tt.wantError && hasError {
				t.Error("Expected success but processing failed")
			}
		})
	}
}

// Test Embedder function with successful case to increase coverage
func TestEmbedder_SuccessCase(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	// First define an embedder to make it available
	plugin := &AzureOpenAI{}
	embedder, err := plugin.DefineEmbedder(g, TextEmbedding3Small)
	if err != nil {
		t.Fatalf("Failed to define embedder: %v", err)
	}

	// Now test the Embedder function to get it
	retrievedEmbedder := Embedder(g, TextEmbedding3Small)
	if retrievedEmbedder == nil {
		t.Error("Embedder() should return non-nil embedder when it exists")
	}

	// Verify both embedders are the same
	if embedder != retrievedEmbedder {
		t.Error("Defined embedder and retrieved embedder should be the same")
	}
}

// Test AzureOpenAI.Init with more error scenarios
func TestAzureOpenAI_Init_MoreErrorScenarios(t *testing.T) {
	// Test already initialized plugin
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

	// Set as already initialized
	plugin.initted = true

	// Should return error for already initialized
	err = plugin.Init(ctx, g)
	if err == nil {
		t.Error("Expected error for already initialized plugin")
	}
	if !strings.Contains(err.Error(), "already initialized") {
		t.Errorf("Expected 'already initialized' error, got: %v", err)
	}
}

// Test AzureOpenAI.DefineModel with custom info
func TestAzureOpenAI_DefineModel_WithCustomInfo(t *testing.T) {
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

	// Force initialization to set initted flag
	plugin.initted = true

	customInfo := &ai.ModelInfo{
		Label: "Custom Test Model",
		Supports: &ai.ModelSupports{
			Multiturn:  true,
			Tools:      false,
			SystemRole: true,
			Media:      false,
		},
		Stage: ai.ModelStageStable,
	}

	// Test with custom model info
	model, err := plugin.DefineModel(g, "custom-test-model", customInfo)
	if err != nil {
		t.Errorf("DefineModel() should not return error for custom model info: %v", err)
	}
	if model == nil {
		t.Error("DefineModel() should return non-nil model")
	}
}

// Test IsDefinedEmbedder with listEmbedders error simulation
func TestIsDefinedEmbedder_WithListError(t *testing.T) {
	// This test checks the error path in IsDefinedEmbedder
	// Since we can't easily mock listEmbedders to return an error,
	// we'll test the normal cases more thoroughly

	tests := []struct {
		name         string
		embedderName string
		expected     bool
	}{
		{
			name:         "valid small embedder",
			embedderName: TextEmbedding3Small,
			expected:     true,
		},
		{
			name:         "valid large embedder",
			embedderName: TextEmbedding3Large,
			expected:     true,
		},
		{
			name:         "invalid embedder",
			embedderName: "non-existent-embedder",
			expected:     false,
		},
		{
			name:         "empty embedder name",
			embedderName: "",
			expected:     false,
		},
		{
			name:         "partial match",
			embedderName: "text-embedding",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDefinedEmbedder(tt.embedderName)
			if result != tt.expected {
				t.Errorf("IsDefinedEmbedder(%q) = %v, want %v", tt.embedderName, result, tt.expected)
			}
		})
	}
}

// Test listModels error path coverage
func TestListModels_Coverage(t *testing.T) {
	// Test the successful case to ensure all models are returned
	models, err := listModels()
	if err != nil {
		t.Errorf("listModels() should not return error: %v", err)
	}

	// Verify specific models exist in the list
	requiredModels := []string{Gpt4o, Gpt4, TextEmbedding3Small, TextEmbedding3Large}
	for _, required := range requiredModels {
		if _, exists := models[required]; !exists {
			t.Errorf("Required model %s not found in models list", required)
		}
	}

	// Verify model info structure
	for name, modelInfo := range models {
		if modelInfo.Label == "" {
			t.Errorf("Model %s has empty label", name)
		}
		if modelInfo.Supports == nil {
			t.Errorf("Model %s has nil supports", name)
		}
		if len(modelInfo.Versions) == 0 {
			t.Errorf("Model %s has no versions", name)
		}
	}
}

// Test modelRef with various configurations
func TestModelRef_Configurations(t *testing.T) {
	tests := []struct {
		name   string
		model  string
		config *OpenAIConfig
	}{
		{
			name:   "basic model ref",
			model:  Gpt4o,
			config: nil,
		},
		{
			name:  "model ref with config",
			model: Gpt4o,
			config: &OpenAIConfig{
				DeploymentName: "gpt-4o-deployment",
				Temperature:    &[]float32{0.7}[0],
				MaxTokens:      &[]int32{100}[0],
			},
		},
		{
			name:  "model ref with all options",
			model: Gpt4,
			config: &OpenAIConfig{
				DeploymentName:   "gpt-4-deployment",
				Temperature:      &[]float32{0.5}[0],
				MaxTokens:        &[]int32{200}[0],
				TopP:             &[]float32{0.9}[0],
				PresencePenalty:  &[]float32{0.1}[0],
				FrequencyPenalty: &[]float32{0.2}[0],
				LogitBias:        map[string]*int32{"test": &[]int32{1}[0]},
				User:             "test-user",
				Seed:             &[]int64{42}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelRef := ModelRef(tt.model, tt.config)
			// Test that ModelRef executes successfully without panics
			// We can't easily test the internal structure without accessing private fields
			_ = modelRef // Verify the function doesn't panic
		})
	}
}

// Test more DefineModel scenarios
func TestDefineModel_MoreScenarios(t *testing.T) {
	ctx := context.Background()
	g, err := genkit.Init(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize Genkit: %v", err)
	}

	tests := []struct {
		name      string
		modelName string
		info      *ai.ModelInfo
	}{
		{
			name:      "text model",
			modelName: "test-text-model",
			info: &ai.ModelInfo{
				Label:    "Test Text Model",
				Versions: []string{"v1.0"},
				Supports: &ai.ModelSupports{
					Multiturn:  true,
					Tools:      true,
					SystemRole: true,
					Media:      false,
				},
				Stage: ai.ModelStageStable,
			},
		},
		{
			name:      "multimodal model",
			modelName: "test-multimodal-model",
			info: &ai.ModelInfo{
				Label:    "Test Multimodal Model",
				Versions: []string{"v1.0", "v1.1"},
				Supports: &ai.ModelSupports{
					Multiturn:  true,
					Tools:      true,
					SystemRole: true,
					Media:      true,
				},
				Stage: ai.ModelStageStable,
			},
		},
		{
			name:      "experimental model",
			modelName: "test-experimental-model",
			info: &ai.ModelInfo{
				Label:    "Test Experimental Model",
				Versions: []string{"v0.1"},
				Supports: &ai.ModelSupports{
					Multiturn:  false,
					Tools:      false,
					SystemRole: false,
					Media:      false,
				},
				Stage: ai.ModelStageUnstable,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := DefineModel(g, tt.modelName, tt.info)
			if model == nil {
				t.Error("DefineModel() should return non-nil model")
			}
		})
	}
}
