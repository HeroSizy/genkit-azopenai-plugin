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
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

const (
	azureOpenAIProvider = "azureopenai"
	labelPrefix         = "Azure OpenAI"

	Provider    = azureOpenAIProvider
	LabelPrefix = labelPrefix
)

// AzureOpenAI is a Genkit plugin for interacting with the Azure OpenAI service.
type AzureOpenAI struct {
	APIKey   string // API key to access the service. If empty, the value of the environment variable AZURE_OPEN_AI_API_KEY will be consulted.
	Endpoint string // Azure OpenAI endpoint. If empty, the value of the environment variable AZURE_OPEN_AI_ENDPOINT will be consulted.

	client  *azopenai.Client // Client for the Azure OpenAI service.
	mu      sync.Mutex       // Mutex to control access.
	initted bool             // Whether the plugin has been initialized.
}

// Name returns the name of the plugin.
func (az *AzureOpenAI) Name() string {
	return azureOpenAIProvider
}

// Init initializes the Azure OpenAI plugin and all known models.
// After calling Init, you may call [DefineModel] to create
// and register any additional generative models
func (az *AzureOpenAI) Init(ctx context.Context, g *genkit.Genkit) (err error) {
	if az == nil {
		az = &AzureOpenAI{}
	}
	az.mu.Lock()
	defer az.mu.Unlock()
	if az.initted {
		return errors.New("plugin already initialized")
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("AzureOpenAI.Init: %w", err)
		}
	}()

	apiKey := az.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("AZURE_OPEN_AI_API_KEY")
		if apiKey == "" {
			return fmt.Errorf("Azure OpenAI requires setting AZURE_OPEN_AI_API_KEY in the environment")
		}
	}

	endpoint := az.Endpoint
	if endpoint == "" {
		endpoint = os.Getenv("AZURE_OPEN_AI_ENDPOINT")
		if endpoint == "" {
			return fmt.Errorf("Azure OpenAI requires setting AZURE_OPEN_AI_ENDPOINT in the environment")
		}
	}

	client, err := azopenai.NewClientWithKeyCredential(endpoint, azcore.NewKeyCredential(apiKey), &azopenai.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Telemetry: policy.TelemetryOptions{
				Disabled: false,
			},
		},
	})
	if err != nil {
		return err
	}
	az.client = client
	az.initted = true

	models, err := listModels()
	if err != nil {
		return err
	}

	// Register all supported models
	for name, modelInfo := range models {
		defineModel(g, az.client, name, modelInfo)
	}

	// Register embedding models
	embeddingModels, err := listEmbedders()
	if err != nil {
		return err
	}
	for _, name := range embeddingModels {
		defineEmbedder(g, az.client, name)
	}

	return nil
}

// DefineModel defines an unknown model with the given name.
// The second argument describes the capability of the model.
// Use [IsDefinedModel] to determine if a model is already defined.
// After [Init] is called, only the known models are defined.
func (az *AzureOpenAI) DefineModel(g *genkit.Genkit, name string, info *ai.ModelInfo) (ai.Model, error) {
	az.mu.Lock()
	defer az.mu.Unlock()
	if !az.initted {
		return nil, errors.New("AzureOpenAI plugin not initialized")
	}
	models, err := listModels()
	if err != nil {
		return nil, err
	}

	var mi ai.ModelInfo
	if info == nil {
		var ok bool
		mi, ok = models[name]
		if !ok {
			return nil, fmt.Errorf("AzureOpenAI.DefineModel: called with unknown model %q and nil ModelInfo", name)
		}
	} else {
		mi = *info
	}

	return defineModel(g, az.client, name, mi), nil
}

// Model returns a reference to the named model.
func Model(g *genkit.Genkit, name string) ai.Model {
	return genkit.LookupModel(g, azureOpenAIProvider, name)
}

// ModelRef creates a model reference that can be used in flows.
func ModelRef(name string, config *OpenAIConfig) ai.ModelRef {
	return ai.NewModelRef(azureOpenAIProvider+"/"+name, config)
}

// DefineModel allows users to define a custom model configuration.
func DefineModel(g *genkit.Genkit, name string, info *ai.ModelInfo) ai.Model {
	return defineModel(g, nil, name, *info)
}

// IsDefinedModel checks if a model is already defined.
func IsDefinedModel(name string) bool {
	model := genkit.LookupModel(nil, azureOpenAIProvider, name)
	return model != nil
}

// Embedder returns an embedder with the given name.
func Embedder(g *genkit.Genkit, name string) ai.Embedder {
	embedder := genkit.LookupEmbedder(g, azureOpenAIProvider, name)
	if embedder == nil {
		panic(fmt.Sprintf("Embedder %q was not found. Make sure you configured the Azure OpenAI plugin and that the embedder is supported.", name))
	}
	return embedder
}

// IsDefinedEmbedder checks if an embedder is supported
func IsDefinedEmbedder(name string) bool {
	embeddingModels, err := listEmbedders()
	if err != nil {
		return false
	}
	for _, model := range embeddingModels {
		if model == name {
			return true
		}
	}
	return false
}

// DefineEmbedder defines an embedder with a given name
func (a *AzureOpenAI) DefineEmbedder(g *genkit.Genkit, name string) (ai.Embedder, error) {
	if !IsDefinedEmbedder(name) {
		return nil, fmt.Errorf("embedder %s is not supported", name)
	}
	return defineEmbedder(g, a.client, name), nil
}

// IsDefinedEmbedder reports whether the named Embedder is defined by this plugin instance.
func (a *AzureOpenAI) IsDefinedEmbedder(g *genkit.Genkit, name string) bool {
	return genkit.LookupEmbedder(g, azureOpenAIProvider, name) != nil
}
