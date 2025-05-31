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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// OpenAIConfig represents the configuration options for Azure OpenAI models.
type OpenAIConfig struct {
	ai.GenerationCommonConfig
	DeploymentName   string            `json:"deploymentName,omitempty"`   // Azure OpenAI deployment name
	MaxTokens        *int32            `json:"maxTokens,omitempty"`        // Maximum number of tokens to generate
	Temperature      *float32          `json:"temperature,omitempty"`      // Controls randomness (0.0 to 2.0)
	TopP             *float32          `json:"topP,omitempty"`             // Nucleus sampling parameter
	PresencePenalty  *float32          `json:"presencePenalty,omitempty"`  // Presence penalty (-2.0 to 2.0)
	FrequencyPenalty *float32          `json:"frequencyPenalty,omitempty"` // Frequency penalty (-2.0 to 2.0)
	LogitBias        map[string]*int32 `json:"logitBias,omitempty"`        // Logit bias modifications (fixed type)
	User             string            `json:"user,omitempty"`             // User identifier
	Seed             *int64            `json:"seed,omitempty"`             // Random seed for deterministic outputs (fixed type)
}

// EmbedConfig contains configuration for embedding requests
type EmbedConfig struct {
	DeploymentName string `json:"deploymentName,omitempty"`
	User           string `json:"user,omitempty"`
}

// defineModel creates and registers a model with Genkit
func defineModel(g *genkit.Genkit, client *azopenai.Client, name string, info ai.ModelInfo) ai.Model {
	return genkit.DefineModel(g, azureOpenAIProvider, name, &info,
		func(ctx context.Context, mr *ai.ModelRequest, cb ai.ModelStreamCallback) (*ai.ModelResponse, error) {
			// Extract config from request
			var cfg OpenAIConfig
			if mr.Config != nil {
				if typedCfg, ok := mr.Config.(*OpenAIConfig); ok {
					cfg = *typedCfg
				}
			}

			if cfg.DeploymentName == "" {
				cfg.DeploymentName = name
				mr.Config = &cfg
			}

			// Convert Genkit request to Azure OpenAI format
			azRequest, err := convertToAzureOpenAIRequest(mr, cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to convert request: %w", err)
			}

			// Handle streaming vs non-streaming
			if cb != nil {
				return handleStreamingRequest(ctx, client, azRequest, cb)
			} else {
				return handleNonStreamingRequest(ctx, client, azRequest)
			}
		})
}

// convertToAzureOpenAIRequest converts a Genkit ModelRequest to Azure OpenAI format
func convertToAzureOpenAIRequest(mr *ai.ModelRequest, cfg OpenAIConfig) (azopenai.ChatCompletionsOptions, error) {
	messages := make([]azopenai.ChatRequestMessageClassification, 0, len(mr.Messages))

	for _, msg := range mr.Messages {
		azMsg, err := convertMessage(msg)
		if err != nil {
			return azopenai.ChatCompletionsOptions{}, err
		}
		messages = append(messages, azMsg)
	}

	deploymentName := cfg.DeploymentName
	if deploymentName == "" {
		return azopenai.ChatCompletionsOptions{}, errors.New("deployment name is required")
	}

	options := azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &deploymentName,
	}

	// Apply configuration options
	if cfg.MaxTokens != nil {
		options.MaxTokens = cfg.MaxTokens
	}
	if cfg.Temperature != nil {
		options.Temperature = cfg.Temperature
	}
	if cfg.TopP != nil {
		options.TopP = cfg.TopP
	}
	if cfg.PresencePenalty != nil {
		options.PresencePenalty = cfg.PresencePenalty
	}
	if cfg.FrequencyPenalty != nil {
		options.FrequencyPenalty = cfg.FrequencyPenalty
	}
	if len(cfg.LogitBias) > 0 {
		options.LogitBias = cfg.LogitBias // Now the types match
	}
	if cfg.User != "" {
		options.User = &cfg.User
	}
	if cfg.Seed != nil {
		options.Seed = cfg.Seed // Now the types match
	}

	// Handle tools if present
	if len(mr.Tools) > 0 {
		tools, err := convertTools(mr.Tools)
		if err != nil {
			return azopenai.ChatCompletionsOptions{}, err
		}
		options.Tools = tools
	}

	return options, nil
}

// convertMessage converts a Genkit message to Azure OpenAI format
func convertMessage(msg *ai.Message) (azopenai.ChatRequestMessageClassification, error) {
	content := extractTextContent(msg.Content)

	switch msg.Role {
	case ai.RoleSystem:
		return &azopenai.ChatRequestSystemMessage{
			Content: azopenai.NewChatRequestSystemMessageContent(content),
		}, nil
	case ai.RoleUser:
		return &azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent(content),
		}, nil
	case ai.RoleModel:
		return &azopenai.ChatRequestAssistantMessage{
			Content: azopenai.NewChatRequestAssistantMessageContent(content), // Fixed type
		}, nil
	case ai.RoleTool:
		// Tool messages need special handling
		return &azopenai.ChatRequestToolMessage{
			Content:    azopenai.NewChatRequestToolMessageContent(content), // Fixed type
			ToolCallID: to.Ptr("tool_call_id"),                             // This should be properly tracked
		}, nil
	default:
		return nil, fmt.Errorf("unsupported role: %s", msg.Role)
	}
}

// extractTextContent extracts text content from message parts
func extractTextContent(parts []*ai.Part) string {
	var textParts []string
	for _, part := range parts {
		if part.IsText() {
			textParts = append(textParts, part.Text)
		}
		// TODO: Handle media parts for multimodal models
	}
	return strings.Join(textParts, "")
}

// convertTools converts Genkit tools to Azure OpenAI format
func convertTools(tools []*ai.ToolDefinition) ([]azopenai.ChatCompletionsToolDefinitionClassification, error) {
	azTools := make([]azopenai.ChatCompletionsToolDefinitionClassification, len(tools))
	for i, tool := range tools {
		// Convert the input schema to JSON bytes
		parametersBytes, err := json.Marshal(tool.InputSchema)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tool parameters: %w", err)
		}

		azTools[i] = &azopenai.ChatCompletionsFunctionToolDefinition{
			Type: to.Ptr("function"),
			Function: &azopenai.ChatCompletionsFunctionToolDefinitionFunction{
				Name:        &tool.Name,
				Description: &tool.Description,
				Parameters:  parametersBytes, // Fixed type
			},
		}
	}
	return azTools, nil
}

// handleStreamingRequest handles streaming chat completions
func handleStreamingRequest(ctx context.Context, client *azopenai.Client, options azopenai.ChatCompletionsOptions, cb ai.ModelStreamCallback) (*ai.ModelResponse, error) {
	resp, err := client.GetChatCompletionsStream(ctx, azopenai.ChatCompletionsStreamOptions{
		Messages:         options.Messages,
		DeploymentName:   options.DeploymentName,
		MaxTokens:        options.MaxTokens,
		Temperature:      options.Temperature,
		TopP:             options.TopP,
		PresencePenalty:  options.PresencePenalty,
		FrequencyPenalty: options.FrequencyPenalty,
		LogitBias:        options.LogitBias,
		User:             options.User,
		Seed:             options.Seed,
		Tools:            options.Tools,
		N:                to.Ptr[int32](1),
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to get chat completions stream: %w", err)
	}
	defer resp.ChatCompletionsStream.Close()

	var fullContent strings.Builder
	var finishReason ai.FinishReason

	for {
		chatCompletion, err := resp.ChatCompletionsStream.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read chat completion: %w", err)
		}

		for _, choice := range chatCompletion.Choices {
			if choice.Delta.Content != nil {
				content := *choice.Delta.Content
				fullContent.WriteString(content)

				// Call the streaming callback
				if cb != nil {
					chunk := &ai.ModelResponseChunk{ // Fixed type
						Content: []*ai.Part{ai.NewTextPart(content)},
						Role:    ai.RoleModel,
					}
					if err := cb(ctx, chunk); err != nil {
						return nil, fmt.Errorf("streaming callback error: %w", err)
					}
				}
			}

			if choice.FinishReason != nil {
				finishReason = convertFinishReason(*choice.FinishReason)
			}
		}
	}

	// Return the final response
	return &ai.ModelResponse{
		Message: &ai.Message{ // Fixed structure
			Content: []*ai.Part{ai.NewTextPart(fullContent.String())},
			Role:    ai.RoleModel,
		},
		FinishReason: finishReason,
	}, nil
}

// handleNonStreamingRequest handles non-streaming chat completions
func handleNonStreamingRequest(ctx context.Context, client *azopenai.Client, options azopenai.ChatCompletionsOptions) (*ai.ModelResponse, error) {
	resp, err := client.GetChatCompletions(ctx, options, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat completions: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from Azure OpenAI")
	}

	choice := resp.Choices[0]
	content := ""
	if choice.Message.Content != nil {
		content = *choice.Message.Content
	}

	finishReason := ai.FinishReasonStop
	if choice.FinishReason != nil {
		finishReason = convertFinishReason(*choice.FinishReason)
	}

	return &ai.ModelResponse{
		Message: &ai.Message{ // Fixed structure
			Content: []*ai.Part{ai.NewTextPart(content)},
			Role:    ai.RoleModel,
		},
		FinishReason: finishReason,
	}, nil
}

// convertFinishReason converts Azure OpenAI finish reason to Genkit format
func convertFinishReason(reason azopenai.CompletionsFinishReason) ai.FinishReason {
	switch reason {
	case azopenai.CompletionsFinishReasonStopped:
		return ai.FinishReasonStop
	case azopenai.CompletionsFinishReasonTokenLimitReached:
		return ai.FinishReasonLength
	case azopenai.CompletionsFinishReasonContentFiltered:
		return ai.FinishReasonBlocked
	case azopenai.CompletionsFinishReasonToolCalls:
		return ai.FinishReasonStop // TODO: Handle tool calls properly
	default:
		return ai.FinishReasonOther
	}
}

// defineEmbedder creates a new embedder for the specified embedding model
func defineEmbedder(g *genkit.Genkit, client *azopenai.Client, name string) ai.Embedder {
	return genkit.DefineEmbedder(g, azureOpenAIProvider, name, func(ctx context.Context, req *ai.EmbedRequest) (*ai.EmbedResponse, error) {
		// Extract configuration from request options
		var config *EmbedConfig
		if opts, ok := req.Options.(*EmbedConfig); ok {
			config = opts
		} else {
			// Use default config with the model name as deployment name
			config = &EmbedConfig{
				DeploymentName: name,
			}
		}

		// Convert input documents to strings
		var input []string
		for _, doc := range req.Input {
			// Extract text content from each document
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

		if len(input) == 0 {
			return nil, fmt.Errorf("no text content found in input documents")
		}

		// Call Azure OpenAI embeddings API
		body := azopenai.EmbeddingsOptions{
			Input:          input,
			DeploymentName: to.Ptr(config.DeploymentName),
		}

		if config.User != "" {
			body.User = to.Ptr(config.User)
		}

		resp, err := client.GetEmbeddings(ctx, body, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get embeddings from Azure OpenAI: %w", err)
		}

		// Convert Azure OpenAI response to Genkit format
		var embeddings []*ai.Embedding
		for _, item := range resp.Data {
			embeddings = append(embeddings, &ai.Embedding{
				Embedding: item.Embedding,
			})
		}

		return &ai.EmbedResponse{
			Embeddings: embeddings,
		}, nil
	})
}
