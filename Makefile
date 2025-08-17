# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# 根據作業系統動態設定二進制文件名稱
ifeq ($(OS),Windows_NT)
    BINARY_NAME=gk.exe
else
    BINARY_NAME=gk
endif

BINARY_UNIX=$(BINARY_NAME)_unix

# 主要建構目標
all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
ifeq ($(OS),Windows_NT)
	-del /Q gk.exe 2>nul
else
	rm -f gk
endif

run: build
ifeq ($(OS),Windows_NT)
	$(BINARY_NAME)
else
	./$(BINARY_NAME)
endif

# 安裝到 GOPATH/bin
install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) -v ./

# 快速測試 resource 命令
test-resource: build
	./$(BINARY_NAME) resource

# 交叉編譯支援
build-all:
	# Linux AMD64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_NAME)-linux-amd64 -v ./
	# Linux ARM64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o builds/$(BINARY_NAME)-linux-arm64 -v ./
	# macOS AMD64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_NAME)-darwin-amd64 -v ./
	# macOS ARM64 (Apple Silicon)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o builds/$(BINARY_NAME)-darwin-arm64 -v ./
	# Windows AMD64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_NAME)-windows-amd64.exe -v ./

.PHONY: all build build-linux test clean run install test-resource build-all
