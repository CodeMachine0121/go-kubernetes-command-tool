# Go K8s Tools 🚀

A modern, modular command-line toolkit for Kubernetes operations, built with Go and following Test-Driven Development (TDD) principles.

## 🎯 Project Purpose

Go K8s Tools aims to simplify and automate Kubernetes workflows for DevOps, SRE, and backend engineers. The project provides a robust, extensible foundation for building custom K8s CLI tools, focusing on maintainability, scalability, and developer productivity.

## 🏗️ Project Structure

```
go-k8s-tools/
├── cmd/                    # CLI entry points and main commands
│   └── root.go             # Main command definition
├── internal/               # Private application logic
│   ├── cli/                # CLI command implementations and handlers
│   │   └── terminal_ui_service.go
│   ├── core/               # Core logic and dependency injection
│   │   └── dependency_injection.go
│   ├── k8s/                # Kubernetes business logic
│   │   ├── k8s_client.go
│   │   ├── k8s_service.go
│   │   ├── k8s_service_test.go
│   │   └── resource_structs.go
├── pkg/                    # Public utility library
│   └── utils/
│       └── utilies.go
├── test/                   # Testing infrastructure and test files
│   └── fixtures/           # Test data and mock files
├── main.go                 # Program entry point
├── go.mod / go.sum         # Go modules management
├── Makefile                # Automation scripts
├── README.md               # Project documentation
├── docs/                   # Additional docs
```

## ⚡️ Quick Start

1. **Install dependencies**
   ```bash
   go mod tidy
   ```
2. **Run the main program**
   ```bash
   go run main.go
   ```
3. **Run tests**
   ```bash
   go test ./...
   ```

## 🛠️ Build with Make (for Users)

To build the executable file using Make:

1. **Build the project**
   ```bash
   make build
   ```
   The compiled executable will be generated in the workspace root as `gk.exe` (on Windows) or `gk` (on Unix-like systems).

2. **Run the executable**
   ```bash
   ./gk.exe
   ```
   (On Unix-like systems, use `./gk`)

> 💡 Make sure you have [Make](https://www.gnu.org/software/make/) and [Go](https://go.dev/doc/install) installed on your system.

## 💡 Development Tips

- Place business logic in `internal/`, reusable utilities in `pkg/`
- Cover all major features with tests in `test/` or alongside code
- Use dependency injection (`internal/core/dependency_injection.go`) for modularity
- CLI commands are recommended to use Cobra
- Configuration via Viper, logging via Logrus

## 🤝 Contributing

Feel free to open issues, submit PRs, or contact maintainers to help improve this project!

---

For more details, check the `docs/` folder. Happy coding! 😄
