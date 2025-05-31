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

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/firebase/genkit/go/ai"
)

func TestConvertFinishReason(t *testing.T) {
	tests := []struct {
		name           string
		azureReason    azopenai.CompletionsFinishReason
		expectedReason ai.FinishReason
	}{
		{
			name:           "stop reason",
			azureReason:    azopenai.CompletionsFinishReasonStopped,
			expectedReason: ai.FinishReasonStop,
		},
		{
			name:           "length reason",
			azureReason:    azopenai.CompletionsFinishReasonTokenLimitReached,
			expectedReason: ai.FinishReasonLength,
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

func TestExtractTextContent(t *testing.T) {
	tests := []struct {
		name     string
		parts    []*ai.Part
		expected string
	}{
		{
			name:     "nil parts",
			parts:    nil,
			expected: "",
		},
		{
			name:     "empty parts",
			parts:    []*ai.Part{},
			expected: "",
		},
		{
			name:     "single text part",
			parts:    []*ai.Part{ai.NewTextPart("hello world")},
			expected: "hello world",
		},
		{
			name:     "multiple text parts",
			parts:    []*ai.Part{ai.NewTextPart("hello"), ai.NewTextPart(" "), ai.NewTextPart("world")},
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTextContent(tt.parts)
			if result != tt.expected {
				t.Errorf("extractTextContent() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConvertMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  *ai.Message
		hasError bool
	}{
		{
			name: "user message with text",
			message: &ai.Message{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("Hello")},
			},
			hasError: false,
		},
		{
			name: "system message with text",
			message: &ai.Message{
				Role:    ai.RoleSystem,
				Content: []*ai.Part{ai.NewTextPart("You are a helpful assistant")},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertMessage(tt.message)
			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			_ = result // Just test that it returns something
		})
	}
}

func TestConvertTools(t *testing.T) {
	tests := []struct {
		name     string
		tools    []*ai.ToolDefinition
		hasError bool
	}{
		{
			name:     "nil tools",
			tools:    nil,
			hasError: false,
		},
		{
			name:     "empty tools",
			tools:    []*ai.ToolDefinition{},
			hasError: false,
		},
		{
			name: "single tool",
			tools: []*ai.ToolDefinition{
				{
					Name:        "test_tool",
					Description: "A test tool",
					InputSchema: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"param": map[string]any{
								"type": "string",
							},
						},
					},
				},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertTools(tt.tools)
			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			_ = result // Just test that it returns something
		})
	}
}

func TestConvertToAzureOpenAIRequest(t *testing.T) {
	config := OpenAIConfig{
		DeploymentName: "test-deployment",
	}

	request := &ai.ModelRequest{
		Messages: []*ai.Message{
			{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("Hello")},
			},
		},
	}

	result, err := convertToAzureOpenAIRequest(request, config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	_ = result // Just test that it returns something
}
