package azopenai_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/herosizy/genkit-go-plugins/azopenai"
)

func ExampleAzureOpenAI_basic() {
	ctx := context.Background()

	// Initialize Genkit
	g, err := genkit.Init(ctx)
	if err != nil {
		fmt.Println("Failed to initialize Genkit:", err)
		return
	}

	// Register the Azure OpenAI plugin
	azurePlugin := &azopenai.AzureOpenAI{
		// APIKey and Endpoint will be read from environment variables
		// AZURE_OPEN_AI_API_KEY and AZURE_OPEN_AI_ENDPOINT
	}

	if err := azurePlugin.Init(ctx, g); err != nil {
		log.Fatal("Failed to initialize Azure OpenAI plugin:", err)
	}

	// Get a reference to a model
	model := azopenai.Model(g, azopenai.Gpt4o)

	// Create a generation request
	request := &ai.ModelRequest{
		Messages: []*ai.Message{
			{
				Role:    ai.RoleSystem,
				Content: []*ai.Part{ai.NewTextPart("You are a helpful assistant.")},
			},
			{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("What is the capital of France?")},
			},
		},
		Config: &azopenai.OpenAIConfig{
			Temperature: to.Ptr(float32(0.7)),
			MaxTokens:   to.Ptr(int32(100)),
		},
	}

	// Generate a response
	response, err := model.Generate(ctx, request, nil)
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		return
	}

	if response.Message != nil && len(response.Message.Content) > 0 {
		fmt.Printf("Response: %s\n", response.Message.Content[0].Text)
	}
}

func ExampleAzureOpenAI_streaming() {
	ctx := context.Background()

	// Initialize Genkit and plugin (same as above)
	g, err := genkit.Init(ctx)
	if err != nil {
		fmt.Println("Failed to initialize Genkit:", err)
		return
	}

	azurePlugin := &azopenai.AzureOpenAI{}
	azurePlugin.Init(ctx, g)

	// Get a reference to a model
	model := azopenai.Model(g, azopenai.Gpt4o)

	// Create a streaming callback
	callback := func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
		for _, part := range chunk.Content {
			if part.IsText() {
				fmt.Print(part.Text)
			}
		}
		return nil
	}

	// Create a generation request
	request := &ai.ModelRequest{
		Messages: []*ai.Message{
			{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("Tell me a short joke.")},
			},
		},
		Config: &azopenai.OpenAIConfig{
			Temperature: to.Ptr(float32(0.8)),
			MaxTokens:   to.Ptr(int32(100)),
		},
	}

	// Generate with streaming
	response, err := model.Generate(ctx, request, callback)
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		return
	}

	fmt.Printf("\nFinish reason: %s\n", response.FinishReason)
}

func ExampleModelRef() {
	ctx := context.Background()

	// Initialize Genkit and plugin (same as above)
	g, err := genkit.Init(ctx)
	if err != nil {
		fmt.Println("Failed to initialize Genkit:", err)
		return
	}

	azurePlugin := &azopenai.AzureOpenAI{}
	azurePlugin.Init(ctx, g)

	// ModelRef is typically used in flows, but for direct generation,
	// you should use the Model() function instead
	model := azopenai.Model(g, "gpt-4o")

	// Create a request with the configuration from ModelRef
	request := &ai.ModelRequest{
		Messages: []*ai.Message{
			{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("Name a creative coffee shop.")},
			},
		},
		Config: &azopenai.OpenAIConfig{
			DeploymentName:   "gpt-4o",
			Temperature:      to.Ptr(float32(0.9)),
			MaxTokens:        to.Ptr(int32(50)),
			PresencePenalty:  to.Ptr(float32(0.5)),
			FrequencyPenalty: to.Ptr(float32(0.5)),
		},
	}

	// Generate using the actual model
	response, err := model.Generate(ctx, request, nil)
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		return
	}

	if response.Message != nil && len(response.Message.Content) > 0 {
		fmt.Printf("Coffee shop name: %s\n", response.Message.Content[0].Text)
	}

	// Note: ModelRef is typically used for creating reusable model configurations
	// that can be passed around and used in flows, not for direct generation
	modelRef := azopenai.ModelRef("gpt-4o", &azopenai.OpenAIConfig{
		DeploymentName: "gpt-4o",
		Temperature:    to.Ptr(float32(0.9)),
	})
	fmt.Printf("ModelRef created for flows: %v\n", modelRef)
}

func ExampleAzureOpenAI_toolCalling() {
	ctx := context.Background()

	// Initialize Genkit and plugin (same as above)
	g, err := genkit.Init(ctx)
	if err != nil {
		fmt.Println("Failed to initialize Genkit:", err)
		return
	}

	azurePlugin := &azopenai.AzureOpenAI{}
	azurePlugin.Init(ctx, g)

	// Get a reference to a model
	model := azopenai.Model(g, "gpt-4o")

	// Define a tool
	tool := &ai.ToolDefinition{
		Name:        "get_weather",
		Description: "Get the current weather for a location",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"location": map[string]any{
					"type":        "string",
					"description": "The city name",
				},
			},
			"required": []string{"location"},
		},
	}

	request := &ai.ModelRequest{
		Messages: []*ai.Message{
			{
				Role:    ai.RoleUser,
				Content: []*ai.Part{ai.NewTextPart("What's the weather in San Francisco?")},
			},
		},
		Tools: []*ai.ToolDefinition{tool},
		Config: &azopenai.OpenAIConfig{
			DeploymentName: "gpt-4o",
		},
	}

	// Generate response (the model might request to call the tool)
	response, err := model.Generate(ctx, request, nil)
	if err != nil {
		log.Printf("Failed to generate response: %v", err)
		return
	}

	fmt.Printf("Response: %v\n", response)
}

// Helper function to demonstrate environment variable usage
func Example() {
	// Set up environment variables for Azure OpenAI
	requiredEnvVars := []string{
		"AZURE_OPEN_AI_API_KEY",
		"AZURE_OPEN_AI_ENDPOINT",
		"AZURE_OPENAI_DEPLOYMENT_NAME", // Optional, but recommended
	}

	fmt.Println("Required environment variables:")
	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			fmt.Printf("❌ %s: not set\n", envVar)
		} else {
			fmt.Printf("✅ %s: configured\n", envVar)
		}
	}

	// Example values:
	fmt.Println("\nExample setup:")
	fmt.Println("export AZURE_OPEN_AI_API_KEY=\"your-api-key-here\"")
	fmt.Println("export AZURE_OPEN_AI_ENDPOINT=\"https://your-resource.openai.azure.com/\"")
	fmt.Println("export AZURE_OPENAI_DEPLOYMENT_NAME=\"gpt-4o\"")
}

func ExampleEmbedder() {
	ctx := context.Background()

	// Initialize Genkit
	g, err := genkit.Init(ctx, genkit.WithPlugins(&azopenai.AzureOpenAI{}))
	if err != nil {
		log.Fatal(err)
	}

	// Get an embedder
	embedder := azopenai.Embedder(g, azopenai.TextEmbedding3Small)

	// Create embed request
	req := &ai.EmbedRequest{
		Input: []*ai.Document{
			ai.DocumentFromText("Hello, world!", nil),
			ai.DocumentFromText("How are you today?", nil),
		},
		Options: &azopenai.EmbedConfig{
			DeploymentName: "text-embedding-3-small", // Your Azure OpenAI deployment name
		},
	}

	// Get embeddings
	resp, err := embedder.Embed(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated %d embeddings\n", len(resp.Embeddings))
	// Output: Generated 2 embeddings
}
