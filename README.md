# Go K8s Tools

A powerful command-line toolkit for Kubernetes operations built with Go, following Test-Driven Development (TDD) principles.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Dependencies](#dependencies)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)

## 🎯 Overview

Go K8s Tools is a comprehensive command-line library designed to simplify Kubernetes operations. Built with modern Go practices, it provides a robust foundation for creating custom Kubernetes tooling with extensive testing capabilities.

## ✨ Features

- **CLI Framework**: Built on Cobra for powerful command-line interfaces
- **Configuration Management**: Flexible configuration using Viper
- **Kubernetes Integration**: Native K8s client-go integration
- **Structured Logging**: Professional logging with Logrus
- **Test-Driven Development**: Comprehensive testing with Ginkgo and Gomega
- **Modular Architecture**: Clean separation of concerns

## 📁 Project Structure

```
go-k8s-tools/
├── cmd/                    # Command-line entry points and main commands
│   └── [CLI commands will be implemented here]
├── internal/               # Private application code (not importable by other projects)
│   ├── cli/               # CLI command implementations and handlers
│   ├── config/            # Configuration management and parsing
│   └── k8s/               # Kubernetes-specific business logic
├── pkg/                    # Public library code (importable by other projects)
│   └── utils/             # Utility functions and helpers
├── test/                   # Testing infrastructure and test files
│   ├── fixtures/          # Test data, mock files, and test fixtures
│   └── integration/       # Integration tests for end-to-end scenarios
├── docs/                   # Project documentation and guides
├── scripts/                # Build scripts, automation, and development tools
├── build/                  # Build artifacts and distribution files
├── .gitignore             # Git ignore patterns for Go projects
├── go.mod                 # Go module definition and dependencies
├── go.sum                 # Go module checksums
├── Makefile               # Build automation and common tasks
└── README.md              # Project documentation (this file)
```

### Directory Descriptions

- **`cmd/`**: Contains the main entry points for different CLI commands. Each subdirectory typically represents a major command or subcommand.
- **`internal/`**: Houses private application code that cannot be imported by external projects, ensuring encapsulation.
- **`pkg/`**: Contains public packages that can be imported and reused by other Go projects.
- **`test/`**: Centralized testing infrastructure supporting unit, integration, and end-to-end tests.

## 📦 Dependencies

### Core Libraries
- **[Cobra](https://github.com/spf13/cobra)** `v1.9.1` - A powerful CLI framework for Go
- **[Viper](https://github.com/spf13/viper)** `v1.20.1` - Configuration management with support for JSON, YAML, TOML, and more
- **[Logrus](https://github.com/sirupsen/logrus)** `v1.9.3` - Structured logging library

### Kubernetes Libraries
- **[client-go](https://github.com/kubernetes/client-go)** `v0.33.4` - Official Kubernetes client library
- **[api](https://github.com/kubernetes/api)** `v0.33.4` - Kubernetes API definitions
- **[apimachinery](https://github.com/kubernetes/apimachinery)** `v0.33.4` - Kubernetes API machinery

### Testing & Development
- **[Testify](https://github.com/stretchr/testify)** - Enhanced testing assertions and mocking
- **[Ginkgo](https://github.com/onsi/ginkgo)** `v2.23.4` - BDD testing framework
- **[Gomega](https://github.com/onsi/gomega)** `v1.38.0` - Matcher library for Ginkgo

## 🚀 Installation

### Prerequisites
- Go 1.21 or later
- Access to a Kubernetes cluster (for integration tests)

### Building from Source
```bash
# Clone the repository
git clone <repository-url>
cd go-k8s-tools

# Download dependencies
go mod download

# Build the project
make build

# Install globally (optional)
make install
```

## 💻 Usage

```bash
# Display help information
./go-k8s-tools --help

# Example command structure (to be implemented)
./go-k8s-tools [command] [subcommand] [flags]
```

## 🛠 Development

### Setting up Development Environment

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd go-k8s-tools
   go mod download
   ```

2. **Run tests**:
   ```bash
   make test
   ```

3. **Run with coverage**:
   ```bash
   make test-coverage
   ```

### TDD Workflow

This project follows Test-Driven Development principles:

1. **Red**: Write failing tests first
2. **Green**: Write minimal code to make tests pass
3. **Refactor**: Improve code while keeping tests green

### Code Organization

- Place unit tests alongside source code with `_test.go` suffix
- Use `test/fixtures/` for test data and mock configurations
- Integration tests go in `test/integration/`
- Follow Go naming conventions and package structure

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test suite
go test ./internal/cli/...

# Run integration tests
make test-integration

# Run BDD tests with Ginkgo
ginkgo run ./...
```

### Test Structure

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test component interactions
- **End-to-End Tests**: Test complete workflows

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes (TDD approach)
4. Implement your feature
5. Ensure all tests pass (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Standards

- Follow Go best practices and idioms
- Maintain test coverage above 80%
- Use meaningful variable and function names
- Add comments for public APIs
- Follow the existing code style

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Kubernetes](https://kubernetes.io/) for the amazing orchestration platform
- [Cobra](https://github.com/spf13/cobra) for the excellent CLI framework
- [Ginkgo](https://github.com/onsi/ginkgo) for the BDD testing framework
- The Go community for fantastic tooling and libraries

---

**Note**: This project is in active development. The API and CLI interface may change as we iterate towards a stable release.