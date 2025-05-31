# Azure OpenAI Plugin for Firebase Genkit

> ⚡ **Pure vibe coded SDK**
> 
> 🛠️ Crafted with **Cursor**, **Claude Sonnet 4**, and **[yuanyang](https://en.wikipedia.org/wiki/Yuenyeung)** ☕️

[![Go Reference](https://pkg.go.dev/badge/github.com/HeroSizy/genkit-go-plugins.svg)](https://pkg.go.dev/github.com/HeroSizy/genkit-go-plugins)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/HeroSizy/genkit-go-plugins)](https://goreportcard.com/report/github.com/HeroSizy/genkit-go-plugins)
[![CI](https://github.com/HeroSizy/genkit-go-plugins/actions/workflows/ci.yml/badge.svg)](https://github.com/HeroSizy/genkit-go-plugins/actions/workflows/ci.yml)
[![Go Coverage](https://github.com/HeroSizy/genkit-go-plugins/wiki/coverage.svg)](https://raw.githack.com/wiki/HeroSizy/genkit-go-plugins/coverage.html)

A production-ready Azure OpenAI integration for [Firebase Genkit](https://firebase.google.com/docs/genkit) that provides seamless access to Azure OpenAI's advanced language models and embedding services.

## 🚀 Quick Start

### Installation

```bash
go get github.com/HeroSizy/genkit-go-plugins/azopenai
```

### Basic Usage

```go
package main

import (
    "context"
    "log"

    "github.com/firebase/genkit/go/genkit"
    "github.com/HeroSizy/genkit-go-plugins/azopenai"
)

func main() {
    ctx := context.Background()
    
    // Initialize Genkit
    g, err := genkit.Init(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Initialize Azure OpenAI plugin
    azurePlugin := &azopenai.AzureOpenAI{}
    if err := azurePlugin.Init(ctx, g); err != nil {
        log.Fatal(err)
    }

    // Ready to use Azure OpenAI models!
    model := azopenai.Model(g, azopenai.Gpt4o)
}
```

## 📦 Packages

### [@azopenai](./azopenai/)

The main Azure OpenAI plugin package providing:

- 🤖 **Complete Model Support** - GPT-4.1, O-series reasoning models, GPT-4o variants
- 🧠 **Reasoning Models** - O1, O3, O4 models for complex problem-solving  
- 🌊 **Streaming Support** - Real-time response streaming
- 🛠️ **Tool Calling** - Function calling capabilities
- 🔗 **Vector Embeddings** - text-embedding-3-small/large support
- 🚀 **Production Ready** - Built with Azure SDK best practices

**[📖 View detailed documentation →](./azopenai/)**

## ⚙️ Environment Setup

```bash
export AZURE_OPEN_AI_API_KEY="your-azure-openai-api-key"
export AZURE_OPEN_AI_ENDPOINT="https://your-resource.openai.azure.com/"
export AZURE_OPENAI_DEPLOYMENT_NAME="gpt-4o"  # Optional
```

## 🔧 Development

### Requirements

- Go 1.21 or later
- Azure OpenAI resource with deployed models

### Development Commands

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Build
make build
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Firebase Genkit](https://firebase.google.com/docs/genkit) team for the excellent framework
- [Azure OpenAI Service](https://azure.microsoft.com/products/ai-services/openai-service) for providing powerful AI models

## 📞 Support

- 🐛 **Bug Reports**: [Open an issue](https://github.com/HeroSizy/genkit-go-plugins/issues)
- 💡 **Feature Requests**: [Start a discussion](https://github.com/HeroSizy/genkit-go-plugins/discussions)

---

**⭐ If this plugin helps you build amazing AI applications, please give it a star!** 
