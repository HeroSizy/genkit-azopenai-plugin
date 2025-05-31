# Contributing to Azure OpenAI Plugin for Firebase Genkit

Thank you for your interest in contributing to this project! We welcome contributions from everyone.

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- Azure OpenAI service access (for testing)

### Setting up your development environment

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/genkit-go-plugins.git
   cd genkit-go-plugins
   ```

3. **Set up environment variables** for testing:
   ```bash
   export AZURE_OPEN_AI_API_KEY="your-test-api-key"
   export AZURE_OPEN_AI_ENDPOINT="https://your-resource.openai.azure.com/"
   export AZURE_OPENAI_DEPLOYMENT_NAME="gpt-4o"
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Verify your setup**:
   ```bash
   go test ./...
   ```

## 🎯 How to Contribute

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to avoid duplicates
2. **Use a descriptive title** that clearly identifies the problem
3. **Provide detailed information** including:
   - Go version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant code snippets or error messages

### Suggesting Features

We welcome feature suggestions! Please:

1. **Check existing feature requests** in the issues
2. **Open a discussion** first for major features
3. **Provide context** about your use case
4. **Consider backward compatibility**

### Pull Requests

1. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Add tests** for new functionality

4. **Update documentation** if needed

5. **Run the test suite**:
   ```bash
   go test ./...
   go test -race ./...
   go vet ./...
   ```

6. **Commit your changes** with a clear message:
   ```bash
   git commit -m "feat: add support for new model configuration"
   ```

7. **Push to your fork** and create a pull request

## 📝 Coding Standards

### Go Code Style

- **Follow standard Go conventions** (`go fmt`, `go vet`)
- **Use meaningful variable and function names**
- **Add comments for exported functions and types**
- **Keep functions focused and small**
- **Handle errors appropriately**

Example:
```go
// DefineCustomModel creates a new model with the specified configuration.
// It returns an error if the model name is already in use or if the
// configuration is invalid.
func DefineCustomModel(g *genkit.Genkit, name string, config *ModelConfig) (ai.Model, error) {
    if name == "" {
        return nil, fmt.Errorf("model name cannot be empty")
    }
    
    // Implementation here...
}
```

### Testing

- **Write tests for all new functionality**
- **Use table-driven tests** where appropriate
- **Include both positive and negative test cases**
- **Mock external dependencies** when needed

Example:
```go
func TestDefineModel(t *testing.T) {
    tests := []struct {
        name        string
        modelName   string
        config      *ModelConfig
        wantErr     bool
        expectedErr string
    }{
        {
            name:      "valid model",
            modelName: "test-model",
            config:    &ModelConfig{DeploymentName: "test"},
            wantErr:   false,
        },
        {
            name:        "empty model name",
            modelName:   "",
            config:      &ModelConfig{},
            wantErr:     true,
            expectedErr: "model name cannot be empty",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Documentation

- **Update README.md** for new features
- **Add inline documentation** for complex logic
- **Include examples** in your documentation
- **Keep documentation up-to-date** with code changes

## 🔄 Development Workflow

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for adding or updating tests
- `refactor:` for code refactoring
- `perf:` for performance improvements
- `chore:` for maintenance tasks

Examples:
```
feat: add streaming support for embeddings
fix: handle rate limit errors gracefully
docs: update installation instructions
test: add integration tests for tool calling
```

### Branch Naming

Use descriptive branch names:
- `feature/add-streaming-embeddings`
- `fix/rate-limit-handling`
- `docs/update-readme`
- `test/integration-tests`

### Code Review Process

1. **All PRs require review** before merging
2. **Address reviewer feedback** promptly
3. **Keep PRs focused** - one feature/fix per PR
4. **Update your PR** if the main branch has moved forward

## 🧪 Testing Guidelines

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestSpecificFunction ./azopenai

# Run tests verbosely
go test -v ./...
```

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test interaction with Azure OpenAI API
3. **Example Tests**: Ensure examples in documentation work

### Mock Testing

For tests that don't require actual API calls:

```go
type mockClient struct {
    // Mock implementation
}

func (m *mockClient) Generate(ctx context.Context, req *azopenai.ChatCompletionsOptions) (*azopenai.ChatCompletions, error) {
    // Return mock response
}
```

## 📦 Release Process

1. **Version bumping** follows [Semantic Versioning](https://semver.org/)
2. **Update CHANGELOG.md** with new features and fixes
3. **Create a release PR** with version updates
4. **Tag the release** after merging

## 🤝 Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/). Please be respectful and inclusive in all interactions.

## 📞 Getting Help

If you need help:

1. **Check the documentation** and examples first
2. **Search existing issues** for similar problems
3. **Ask in GitHub Discussions** for questions
4. **Create an issue** for bugs or feature requests

## 🎉 Recognition

Contributors will be recognized in:
- The project README
- Release notes
- GitHub contributor graphs

Thank you for helping make this project better! 🚀 