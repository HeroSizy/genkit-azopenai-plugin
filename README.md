# 🔥 Azure OpenAI Plugin for Firebase Genkit

> ⚡ **Pure vibe coded SDK**
> 
> 🛠️ Crafted with **Cursor**, **Claude Sonnet 4**, and **yuanyang** ☕️

[![Go Reference](https://pkg.go.dev/badge/github.com/herosizy/genkit-go-plugins.svg)](https://pkg.go.dev/github.com/herosizy/genkit-go-plugins)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/herosizy/genkit-go-plugins)](https://goreportcard.com/report/github.com/herosizy/genkit-go-plugins)

A powerful, production-ready Azure OpenAI integration for [Firebase Genkit](https://firebase.google.com/docs/genkit) that provides seamless access to Azure OpenAI's advanced language models and embedding services.

## ✨ Features

- 🤖 **Complete Model Support** - Access all Azure OpenAI models including GPT-4.1, O-series reasoning models, and GPT-4o variants
- 🧠 **Reasoning Models** - Support for advanced O1, O3, and O4 models for complex problem-solving
- 🌊 **Streaming Support** - Real-time response streaming for interactive applications
- 🛠️ **Tool Calling** - Function calling capabilities for complex AI workflows
- 🔗 **Vector Embeddings** - Support for text-embedding-3-small and text-embedding-3-large
- 🔧 **Flexible Configuration** - Environment variables or programmatic configuration
- 🚀 **Production Ready** - Built with Azure SDK best practices and error handling
- 📝 **Type Safe** - Comprehensive Go type definitions for all configurations

## 🚀 Quick Start

### Installation

```bash
go get github.com/herosizy/genkit-go-plugins/azopenai
```

### Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
    "github.com/firebase/genkit/go/ai"
    "github.com/firebase/genkit/go/genkit"
    "github.com/herosizy/genkit-go-plugins/azopenai"
)

func main() {
    ctx := context.Background()

    // Initialize Genkit
    g, err := genkit.Init(ctx)
    if err != nil {
        log.Fatal("Failed to initialize Genkit:", err)
    }

    // Initialize Azure OpenAI plugin
    azurePlugin := &azopenai.AzureOpenAI{
        // Optional: APIKey and Endpoint (will use env vars if not set)
        // APIKey:   "your-api-key",
        // Endpoint: "https://your-resource.openai.azure.com/",
    }

    if err := azurePlugin.Init(ctx, g); err != nil {
        log.Fatal("Failed to initialize Azure OpenAI plugin:", err)
    }

    // Get a model reference
    model := azopenai.Model(g, azopenai.Gpt4o)

    // Create a simple request
    request := &ai.ModelRequest{
        Messages: []*ai.Message{
            {
                Role:    ai.RoleUser,
                Content: []*ai.Part{ai.NewTextPart("Hello! Tell me a joke.")},
            },
        },
        Config: &azopenai.OpenAIConfig{
            Temperature:    to.Ptr(float32(0.7)),
            MaxTokens:      to.Ptr(int32(100)),
        },
    }

    // Generate response
    response, err := model.Generate(ctx, request, nil)
    if err != nil {
        log.Fatal("Failed to generate response:", err)
    }

    if response.Message != nil && len(response.Message.Content) > 0 {
        fmt.Printf("AI: %s\n", response.Message.Content[0].Text)
    }
}
```

## ⚙️ Configuration

### Environment Variables

Set these environment variables for automatic configuration:

```bash
export AZURE_OPEN_AI_API_KEY="your-azure-openai-api-key"
export AZURE_OPEN_AI_ENDPOINT="https://your-resource.openai.azure.com/"
export AZURE_OPENAI_DEPLOYMENT_NAME="gpt-4o"  # Optional default deployment
```

### Programmatic Configuration

```go
azurePlugin := &azopenai.AzureOpenAI{
    APIKey:   "your-api-key",
    Endpoint: "https://your-resource.openai.azure.com/",
}
```

## 📚 Usage Examples

### Text Generation with Streaming

```go
func generateWithStreaming() {
    // ... initialization code ...

    model := azopenai.Model(g, azopenai.Gpt4o)

    // Streaming callback
    callback := func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
        for _, part := range chunk.Content {
            if part.IsText() {
                fmt.Print(part.Text)
            }
        }
        return nil
    }

    request := &ai.ModelRequest{
        Messages: []*ai.Message{
            {
                Role:    ai.RoleUser,
                Content: []*ai.Part{ai.NewTextPart("Write a short story about AI.")},
            },
        },
        Config: &azopenai.OpenAIConfig{
            Temperature:    to.Ptr(float32(0.8)),
        },
    }

    response, err := model.Generate(ctx, request, callback)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("\nFinish reason: %s\n", response.FinishReason)
}
```

### Tool Calling (Function Calling)

```go
func demonstrateToolCalling() {
    // ... initialization code ...

    model := azopenai.Model(g, azopenai.Gpt4o)

    // Define a weather tool
    weatherTool := &ai.ToolDefinition{
        Name:        "get_weather",
        Description: "Get current weather for a location",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "location": map[string]any{
                    "type":        "string",
                    "description": "The city name",
                },
                "unit": map[string]any{
                    "type":        "string",
                    "enum":        []string{"celsius", "fahrenheit"},
                    "description": "Temperature unit",
                },
            },
            "required": []string{"location"},
        },
    }

    request := &ai.ModelRequest{
        Messages: []*ai.Message{
            {
                Role:    ai.RoleUser,
                Content: []*ai.Part{ai.NewTextPart("What's the weather like in Tokyo?")},
            },
        },
        Tools: []*ai.ToolDefinition{weatherTool},
    }

    response, err := model.Generate(ctx, request, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Handle tool calls in response
    for _, part := range response.Message.Content {
        if part.IsToolRequest() {
            fmt.Printf("Tool called: %s with args: %v\n", 
                part.ToolRequest.Name, part.ToolRequest.Input)
        }
    }
}
```

### Vector Embeddings

```go
func generateEmbeddings() {
    // ... initialization code ...

    embedder := azopenai.Embedder(g, azopenai.TextEmbedding3Small)

    req := &ai.EmbedRequest{
        Input: []*ai.Document{
            ai.DocumentFromText("Machine learning is fascinating", nil),
            ai.DocumentFromText("I love programming in Go", nil),
            ai.DocumentFromText("Azure OpenAI provides powerful AI models", nil),
        },
        Options: &azopenai.EmbedConfig{
            User:           "example-user",
        },
    }

    resp, err := embedder.Embed(ctx, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Generated %d embeddings\n", len(resp.Embeddings))
    for i, embedding := range resp.Embeddings {
        fmt.Printf("Document %d: %d dimensions\n", i, len(embedding.Embedding))
    }
}
```

### Using Model References in Flows

```go
// Create reusable model references
creativeModel := azopenai.ModelRef("gpt-4o", &azopenai.OpenAIConfig{
    DeploymentName:   "gpt-4o",
    Temperature:      to.Ptr(float32(0.9)),
    MaxTokens:        to.Ptr(int32(500)),
    PresencePenalty:  to.Ptr(float32(0.6)),
    FrequencyPenalty: to.Ptr(float32(0.6)),
})

factualModel := azopenai.ModelRef("gpt-4o", &azopenai.OpenAIConfig{
    DeploymentName: "gpt-4o",
    Temperature:    to.Ptr(float32(0.1)),
    MaxTokens:      to.Ptr(int32(200)),
})

// Use in flows
genkit.DefineFlow(g, "creative-writing", func(ctx context.Context, input string) (string, error) {
    response, err := creativeModel.Generate(ctx, &ai.ModelRequest{
        Messages: []*ai.Message{
            {Role: ai.RoleUser, Content: []*ai.Part{ai.NewTextPart(input)}},
        },
    }, nil)
    if err != nil {
        return "", err
    }
    return response.Message.Content[0].Text, nil
})
```

## 🎯 Supported Models

### 🧠 Reasoning Models
Latest o-series models that excel at complex, multi-step tasks:

| Model | Description | Status |
|-------|-------------|--------|
| `o4-mini` | Latest O4 reasoning model (mini) | Preview |
| `o3` | Advanced O3 reasoning model | Preview |
| `o3-mini` | Efficient O3 reasoning model | Preview |
| `o1` | Original O1 reasoning model | Preview |
| `o1-mini` | Compact O1 reasoning model | Preview |
| `o1-pro` | Professional O1 reasoning model | Preview |

### 🚀 Flagship Chat Models
Versatile, high-intelligence flagship models:

| Model | Description | Context Window | Features |
|-------|-------------|----------------|----------|
| `gpt-4.1` | Latest GPT-4.1 model | 128K tokens | Multimodal, Tools |
| `gpt-4.1-mini` | Efficient GPT-4.1 variant | 128K tokens | Multimodal, Tools |
| `gpt-4.1-nano` | Ultra-fast GPT-4.1 variant | 128K tokens | Multimodal, Tools |
| `gpt-4o` | GPT-4 Omni model | 128K tokens | Multimodal, Tools |
| `gpt-4o-mini` | Faster, cost-effective GPT-4 Omni | 128K tokens | Multimodal, Tools |
| `gpt-4o-audio-preview` | GPT-4o with audio capabilities | 128K tokens | Audio, Multimodal |
| `gpt-4o-mini-audio-preview` | GPT-4o Mini with audio | 128K tokens | Audio, Multimodal |
| `chatgpt-4o-latest` | Latest ChatGPT-4o variant | 128K tokens | Multimodal, Tools |

### 🎨 Image Generation Models
Models that can generate and edit images:

| Model | Description | Capabilities |
|-------|-------------|--------------|
| `gpt-image-1` | Latest image generation model | High-quality image generation |
| `dall-e-3` | DALL-E 3 image generation | Advanced image creation |
| `dall-e-2` | DALL-E 2 image generation | Standard image creation |

### 🔗 Embedding Models
Convert text into vector representations:

| Model | Description | Dimensions | Max Input |
|-------|-------------|------------|-----------|
| `text-embedding-3-large` | High-performance embeddings | 3072 | 8191 tokens |
| `text-embedding-3-small` | Efficient embeddings | 1536 | 8191 tokens |

### 📚 Legacy Models
Supported older versions for compatibility:

| Model | Description | Context Window |
|-------|-------------|----------------|
| `gpt-4-turbo` | High-performance GPT-4 | 128K tokens |
| `gpt-4-turbo-preview` | Preview version of GPT-4 Turbo | 128K tokens |
| `gpt-4` | Original GPT-4 model | 8K tokens |
| `gpt-3.5-turbo` | Fast and efficient model | 16K tokens |
| `gpt-3.5-turbo-instruct` | Instruction-following variant | 4K tokens |

## ⚙️ Configuration Options

### OpenAIConfig for Text Generation

```go
type OpenAIConfig struct {
    DeploymentName   string               `json:"deploymentName"`   // Required: Azure deployment name
    MaxTokens        *int32               `json:"maxTokens"`        // Maximum tokens to generate
    Temperature      *float32             `json:"temperature"`      // Randomness (0.0-2.0)
    TopP             *float32             `json:"topP"`             // Nucleus sampling
    PresencePenalty  *float32             `json:"presencePenalty"`  // Presence penalty (-2.0 to 2.0)
    FrequencyPenalty *float32             `json:"frequencyPenalty"` // Frequency penalty (-2.0 to 2.0)
    LogitBias        map[string]*int32    `json:"logitBias"`        // Token bias modifications
    User             string               `json:"user"`             // User identifier
    Seed             *int64               `json:"seed"`             // Deterministic seed
}
```

### EmbedConfig for Embeddings

```go
type EmbedConfig struct {
    DeploymentName string `json:"deploymentName"` // Required: Azure deployment name
    User           string `json:"user"`           // Optional: User identifier
}
```

## 🔧 Advanced Usage

### Custom Model Definition

```go
// Define a custom model with specific capabilities
customModel, err := azurePlugin.DefineModel(g, "my-custom-gpt4", &ai.ModelInfo{
    Label:    "Custom GPT-4",
    Supports: ai.ModelSupports{
        Multiturn:  true,
        Tools:      true,
        SystemRole: true,
        Media:      false,
    },
})
if err != nil {
    log.Fatal(err)
}
```

### Error Handling Best Practices

```go
func robustGeneration(model ai.Model, request *ai.ModelRequest) (*ai.ModelResponse, error) {
    response, err := model.Generate(ctx, request, nil)
    if err != nil {
        // Handle specific Azure OpenAI errors
        if strings.Contains(err.Error(), "rate limit") {
            // Implement exponential backoff
            time.Sleep(time.Second * 2)
            return model.Generate(ctx, request, nil)
        }
        if strings.Contains(err.Error(), "context length") {
            // Reduce request size
            request.Config.(*azopenai.OpenAIConfig).MaxTokens = to.Ptr(int32(500))
            return model.Generate(ctx, request, nil)
        }
        return nil, fmt.Errorf("generation failed: %w", err)
    }
    return response, nil
}
```

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test examples
go test -run ExampleAzureOpenAI_basic ./azopenai
```

## 📋 Prerequisites

- Go 1.21 or later
- Azure OpenAI resource with deployed models
- API key and endpoint from Azure Portal

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/your-username/genkit-go-plugins.git
   cd genkit-go-plugins
   ```

2. **Set up environment**
   ```bash
   cp .example.env .env
   # Edit .env with your Azure OpenAI credentials
   ```

3. **Run tests**
   ```bash
   go test ./...
   ```

4. **Submit a pull request**

### Code Style

- Follow standard Go formatting (`go fmt`)
- Write comprehensive tests for new features
- Update documentation for API changes
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📚 References

### Official Documentation
- **[OpenAI Models Documentation](https://platform.openai.com/docs/models)** - Complete list of available OpenAI models and their capabilities
- **[Genkit Plugin Authoring Guide](https://genkit.dev/go/docs/plugin-authoring-models/)** - Official guide for writing Genkit model plugins
- **[Azure OpenAI SDK for Go](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/ai/azopenai)** - Official Azure SDK for Go that powers this plugin

### Related Resources
- [Azure OpenAI Service Documentation](https://docs.microsoft.com/azure/cognitive-services/openai/)
- [Firebase Genkit Documentation](https://firebase.google.com/docs/genkit)
- [Azure OpenAI REST API Reference](https://docs.microsoft.com/azure/cognitive-services/openai/reference)

## 🙏 Acknowledgments

- [Firebase Genkit](https://firebase.google.com/docs/genkit) team for the excellent framework
- [Azure OpenAI Service](https://azure.microsoft.com/products/ai-services/openai-service) for providing powerful AI models
- All contributors who help improve this plugin

## 📞 Support

- 🐛 **Bug Reports**: [Open an issue](https://github.com/herosizy/genkit-go-plugins/issues)
- 💡 **Feature Requests**: [Start a discussion](https://github.com/herosizy/genkit-go-plugins/discussions)

---

**⭐ If this plugin helps you build amazing AI applications, please give it a star!** 
