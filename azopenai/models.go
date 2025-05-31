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
	"github.com/firebase/genkit/go/ai"
)

const (
	//
	// Reasoning models
	// o-series models that excel at complex, multi-step tasks.
	o4Mini = "o4-mini"
	o3     = "o3"
	o3Mini = "o3-mini"
	o1     = "o1"
	o1Mini = "o1-mini"
	o1Pro  = "o1-pro"

	// Flagship chat models
	// Our versatile, high-intelligence flagship models.
	gpt41          = "gpt-4.1"
	gpt41Mini      = "gpt-4.1-mini"
	gpt41Nano      = "gpt-4.1-nano"
	gpt4o          = "gpt-4o"
	gpt4oMini      = "gpt-4o-mini"
	gpt4oAudio     = "gpt-4o-audio-preview"
	gpt4oMiniAudio = "gpt-4o-mini-audio-preview"
	chatgpt4o      = "chatgpt-4o-latest"

	// Image generation models
	// Models that can generate and edit images, given a natural language prompt.
	gptImage1 = "gpt-image-1"
	dalle3    = "dall-e-3"
	dalle2    = "dall-e-2"

	// Embeddings
	// A set of models that can convert text into vector representations.
	textEmbedding3Large = "text-embedding-3-large"
	textEmbedding3Small = "text-embedding-3-small"

	// Older GPT models
	// Supported older versions of our general purpose and chat models.
	gpt35Turbo         = "gpt-3.5-turbo"
	gpt35TurboInstruct = "gpt-3.5-turbo-instruct"
	gpt4               = "gpt-4"
	gpt4Turbo          = "gpt-4-turbo"
	gpt4TurboPreview   = "gpt-4-turbo-preview"

	// Exported  model constants for external use
	Gpt35Turbo          = gpt35Turbo
	Gpt35TurboInstruct  = gpt35TurboInstruct
	Gpt4                = gpt4
	Gpt4Turbo           = gpt4Turbo
	Gpt4TurboPreview    = gpt4TurboPreview
	Gpt41               = gpt41
	Gpt41Mini           = gpt41Mini
	Gpt41Nano           = gpt41Nano
	Gpt4o               = gpt4o
	Gpt4oMini           = gpt4oMini
	Gpt4oAudio          = gpt4oAudio
	Gpt4oMiniAudio      = gpt4oMiniAudio
	O4Mini              = o4Mini
	O3                  = o3
	O3Mini              = o3Mini
	O1                  = o1
	O1Mini              = o1Mini
	O1Pro               = o1Pro
	Chatgpt4o           = chatgpt4o
	GptImage1           = gptImage1
	Dalle3              = dalle3
	Dalle2              = dalle2
	TextEmbedding3Large = textEmbedding3Large
	TextEmbedding3Small = textEmbedding3Small
)

var (
	// List of supported Azure OpenAI models
	azureOpenAIModels = []string{
		gpt4,
		gpt4Turbo,
		gpt4TurboPreview,
		gpt4o,
		gpt4oMini,
		gpt35Turbo,
		gpt35TurboInstruct,
		textEmbedding3Large,
		textEmbedding3Small,
		gpt41,
		gpt41Mini,
		o4Mini,
	}

	// Model capabilities for text models
	TextModel = ai.ModelSupports{
		Multiturn:  true,
		Tools:      true,
		ToolChoice: true,
		SystemRole: true,
		Media:      false,
	}

	// Model capabilities for multimodal models
	MultimodalModel = ai.ModelSupports{
		Multiturn:  true,
		Tools:      true,
		ToolChoice: true,
		SystemRole: true,
		Media:      true,
	}

	// supportedAzureOpenAIModels maps model names to their capabilities
	supportedAzureOpenAIModels = map[string]ai.ModelInfo{
		gpt4: {
			Label: "GPT-4",
			Versions: []string{
				"gpt-4",
			},
			Supports: &TextModel,
			Stage:    ai.ModelStageStable,
		},
		gpt4Turbo: {
			Label: "GPT-4 Turbo",
			Versions: []string{
				"gpt-4-turbo",
			},
			Supports: &MultimodalModel,
			Stage:    ai.ModelStageStable,
		},
		gpt4TurboPreview: {
			Label: "GPT-4 Turbo Preview",
			Versions: []string{
				"gpt-4-turbo-preview",
			},
			Supports: &MultimodalModel,
			Stage:    ai.ModelStageUnstable,
		},
		gpt4o: {
			Label: "GPT-4o",
			Versions: []string{
				"gpt-4o-2024-05-13",
				"gpt-4o-2024-08-06",
			},
			Supports: &MultimodalModel,
			Stage:    ai.ModelStageStable,
		},
		gpt4oMini: {
			Label: "GPT-4o Mini",
			Versions: []string{
				"gpt-4o-mini-2024-07-18",
			},
			Supports: &MultimodalModel,
			Stage:    ai.ModelStageStable,
		},
		gpt35Turbo: {
			Label: "GPT-3.5 Turbo",
			Versions: []string{
				"gpt-3.5-turbo-0613",
				"gpt-3.5-turbo-16k",
				"gpt-3.5-turbo-16k-0613",
				"gpt-3.5-turbo-1106",
				"gpt-3.5-turbo-0125",
			},
			Supports: &TextModel,
			Stage:    ai.ModelStageStable,
		},
		gpt35TurboInstruct: {
			Label: "GPT-3.5 Turbo Instruct",
			Versions: []string{
				"gpt-3.5-turbo-instruct-0914",
			},
			Supports: &TextModel,
			Stage:    ai.ModelStageStable,
		},
		textEmbedding3Large: {
			Label: "Text Embedding 3 Large",
			Versions: []string{
				"text-embedding-3-large",
			},
			Supports: &ai.ModelSupports{
				Multiturn:  false,
				Tools:      false,
				ToolChoice: false,
				SystemRole: false,
				Media:      false,
			},
			Stage: ai.ModelStageStable,
		},
		textEmbedding3Small: {
			Label: "Text Embedding 3 Small",
			Versions: []string{
				"text-embedding-3-small",
			},
			Supports: &ai.ModelSupports{
				Multiturn:  false,
				Tools:      false,
				ToolChoice: false,
				SystemRole: false,
				Media:      false,
			},
			Stage: ai.ModelStageStable,
		},
		gpt41: {
			Label: "GPT-4.1",
			Versions: []string{
				"gpt-4.1",
			},
			Supports: &MultimodalModel,
			Stage:    ai.ModelStageUnstable,
		},
		gpt41Mini: {
			Label: "GPT-4.1 Mini",
			Versions: []string{
				"gpt-4.1-mini",
			},
			Supports: &MultimodalModel,
			Stage:    ai.ModelStageUnstable,
		},
		o4Mini: {
			Label: "O4 Mini",
			Versions: []string{
				"o4-mini",
			},
			Supports: &MultimodalModel,
			Stage:    ai.ModelStageUnstable,
		},
	}
)

// listModels returns a map of supported models and their capabilities
func listModels() (map[string]ai.ModelInfo, error) {
	models := make(map[string]ai.ModelInfo, len(azureOpenAIModels))
	for _, name := range azureOpenAIModels {
		m, ok := supportedAzureOpenAIModels[name]
		if !ok {
			continue // Skip unknown models
		}
		models[name] = ai.ModelInfo{
			Label:    labelPrefix + " - " + m.Label,
			Versions: m.Versions,
			Supports: m.Supports,
			Stage:    m.Stage,
		}
	}
	return models, nil
}

// listEmbedders returns the list of supported embedding models
func listEmbedders() ([]string, error) {
	return []string{
		textEmbedding3Large,
		textEmbedding3Small,
	}, nil
}
