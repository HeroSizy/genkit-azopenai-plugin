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

// Package azopenai provides Azure OpenAI integration for Firebase Genkit.
//
// This is a pure vibe coded SDK that provides seamless access to Azure OpenAI's
// advanced language models and embedding services through the Firebase Genkit framework.
//
// # Features
//
//   - Complete Model Support: Access all Azure OpenAI models including GPT-4.1, O-series reasoning models, and GPT-4o variants
//   - Reasoning Models: Support for advanced O1, O3, and O4 models for complex problem-solving
//   - Streaming Support: Real-time response streaming for interactive applications
//   - Tool Calling: Function calling capabilities for complex AI workflows
//   - Vector Embeddings: Support for text-embedding-3-small and text-embedding-3-large
//   - Flexible Configuration: Environment variables or programmatic configuration
//   - Production Ready: Built with Azure SDK best practices and error handling
//   - Type Safe: Comprehensive Go type definitions for all configurations
//
// # Quick Start
//
// First, set up your environment variables:
//
//	export AZURE_OPEN_AI_API_KEY="your-azure-openai-api-key"
//	export AZURE_OPEN_AI_ENDPOINT="https://your-resource.openai.azure.com/"
//	export AZURE_OPENAI_DEPLOYMENT_NAME="gpt-4o"  # Optional default deployment
//
// Basic usage example:
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"log"
//
//		"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
//		"github.com/firebase/genkit/go/ai"
//		"github.com/firebase/genkit/go/genkit"
//		"github.com/herosizy/genkit-go-plugins/azopenai"
//	)
//
//	func main() {
//		ctx := context.Background()
//
//		// Initialize Genkit
//		g, err := genkit.Init(ctx)
//		if err != nil {
//			log.Fatal("Failed to initialize Genkit:", err)
//		}
//
//		// Initialize Azure OpenAI plugin
//		azurePlugin := &azopenai.AzureOpenAI{}
//		if err := azurePlugin.Init(ctx, g); err != nil {
//			log.Fatal("Failed to initialize Azure OpenAI plugin:", err)
//		}
//
//		// Get a model reference
//		model := azopenai.Model(g, azopenai.Gpt4o)
//
//		// Create a simple request
//		request := &ai.ModelRequest{
//			Messages: []*ai.Message{
//				{
//					Role:    ai.RoleUser,
//					Content: []*ai.Part{ai.NewTextPart("Hello! Tell me a joke.")},
//				},
//			},
//			Config: &azopenai.OpenAIConfig{
//				DeploymentName: "gpt-4o", // Your Azure deployment name
//				Temperature:    to.Ptr(float32(0.7)),
//				MaxTokens:      to.Ptr(int32(100)),
//			},
//		}
//
//		// Generate response
//		response, err := model.Generate(ctx, request, nil)
//		if err != nil {
//			log.Fatal("Failed to generate response:", err)
//		}
//
//		if response.Message != nil && len(response.Message.Content) > 0 {
//			fmt.Printf("AI: %s\n", response.Message.Content[0].Text)
//		}
//	}
//
// # Supported Models
//
// ## Text Generation Models
//
//   - GPT-4.1 series: gpt-4.1, gpt-4.1-mini, gpt-4.1-nano
//   - GPT-4o series: gpt-4o, gpt-4o-mini, gpt-4o-audio-preview
//   - GPT-4 series: gpt-4, gpt-4-turbo, gpt-4-turbo-preview
//   - GPT-3.5 series: gpt-3.5-turbo, gpt-3.5-turbo-instruct
//   - Reasoning models: o1, o1-mini, o1-pro, o3, o3-mini, o4-mini
//   - Image generation: dall-e-2, dall-e-3, gpt-image-1
//
// ## Embedding Models
//
//   - text-embedding-3-large: High-performance embeddings (3072 dimensions)
//   - text-embedding-3-small: Efficient embeddings (1536 dimensions)
//
// # Advanced Usage
//
// ## Streaming Responses
//
//	callback := func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
//		for _, part := range chunk.Content {
//			if part.IsText() {
//				fmt.Print(part.Text)
//			}
//		}
//		return nil
//	}
//
//	response, err := model.Generate(ctx, request, callback)
//
// ## Tool Calling
//
//	tool := &ai.ToolDefinition{
//		Name:        "get_weather",
//		Description: "Get current weather for a location",
//		InputSchema: map[string]any{
//			"type": "object",
//			"properties": map[string]any{
//				"location": map[string]any{
//					"type":        "string",
//					"description": "The city name",
//				},
//			},
//			"required": []string{"location"},
//		},
//	}
//
//	request := &ai.ModelRequest{
//		Messages: []*ai.Message{
//			{
//				Role:    ai.RoleUser,
//				Content: []*ai.Part{ai.NewTextPart("What's the weather in Tokyo?")},
//			},
//		},
//		Tools: []*ai.ToolDefinition{tool},
//		Config: &azopenai.OpenAIConfig{
//			DeploymentName: "gpt-4o",
//		},
//	}
//
// ## Vector Embeddings
//
//	embedder := azopenai.Embedder(g, azopenai.TextEmbedding3Small)
//
//	req := &ai.EmbedRequest{
//		Input: []*ai.Document{
//			ai.DocumentFromText("Machine learning is fascinating", nil),
//		},
//		Options: &azopenai.EmbedConfig{
//			DeploymentName: "text-embedding-3-small",
//		},
//	}
//
//	resp, err := embedder.Embed(ctx, req)
//
// # Configuration
//
// The OpenAIConfig struct supports comprehensive configuration options:
//
//   - DeploymentName: Azure OpenAI deployment name (required)
//   - MaxTokens: Maximum number of tokens to generate
//   - Temperature: Controls randomness (0.0 to 2.0)
//   - TopP: Nucleus sampling parameter
//   - PresencePenalty: Presence penalty (-2.0 to 2.0)
//   - FrequencyPenalty: Frequency penalty (-2.0 to 2.0)
//   - LogitBias: Token bias modifications
//   - User: User identifier for tracking
//   - Seed: Random seed for deterministic outputs
//
// # Environment Variables
//
//   - AZURE_OPEN_AI_API_KEY: Your Azure OpenAI API key (required)
//   - AZURE_OPEN_AI_ENDPOINT: Your Azure OpenAI endpoint (required)
//   - AZURE_OPENAI_DEPLOYMENT_NAME: Default deployment name (optional)
//
// # Plugin Interface
//
// This package implements the Genkit plugin interface, providing:
//
//   - Model() function: Get references to defined models
//   - ModelRef() function: Create model references for flows
//   - Embedder() function: Get references to embedding models
//   - DefineModel() function: Define custom model configurations
//   - AzureOpenAI struct: Main plugin implementation
//
// For more information and examples, visit:
// https://github.com/herosizy/genkit-go-plugins
package azopenai
