PROJECT_NAME = prometheus-slurm-exporter
SHELL := $.PHONY: build-single
build-single: go/modules/pkg/mod ## 构建指定架构或所有架构的二进制文件 (可选: ARCH=amd64|arm64)
	@if [ -z "$(ARCH)" ]; then \
	    echo "=== No ARCH specified, building for all architectures ==="; \
	    $(foreach arch,$(ARCHS),$(MAKE) ARCH=$(arch) build-single;) \
	    echo "=== Build completed ==="; \
	    echo "Built binaries:"; \
	    $(foreach arch,$(ARCHS),echo "  $(call GOBIN,$(arch))";) \
	else \
	    mkdir -p $(dir $(call GOBIN,$(ARCH))); \
	    echo "Building $(PROJECT_NAME) for $(OS)/$(ARCH)..."; \
	    echo "Compiling Go files (excluding test files):"; \
	    find . -name '*.go' -not -path './vendor/*' -not -path './go/*' -not -name 'test_*.go' -not -name '*_test.go' | head -5 | sed 's/^/  /' && \
	    ([ $$(find . -name '*.go' -not -path './vendor/*' -not -path './go/*' -not -name 'test_*.go' -not -name '*_test.go' | wc -l) -gt 5 ] && echo "  ... and more" || true); \
	    $(GO_BUILD_ENV) go build $(LDFLAGS) -v -o $(call GOBIN,$(ARCH)) $$(find . -name '*.go' -not -path './vendor/*' -not -path './go/*' -not -name 'test_*.go' -not -name '*_test.go'); \
	    echo "✓ Built: $(call GOBIN,$(ARCH))"; \
	fiash) -eu -o pipefail

# 版本信息
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 多架构配置
OS = linux
ARCHS = amd64 arm64
GOPATH = $(shell pwd)/go/modules

# 只包含非测试的 Go 文件
# 排除规则：
# 1. vendor 和 go 目录下的文件
# 2. 以 test_ 开头的文件 (如 test_parsing.go, test_real_commands.go)  
# 3. 以 _test.go 结尾的文件 (如 node_test.go, queue_test.go)
GOFILES = $(shell find . -name '*.go' -not -path './vendor/*' -not -path './go/*' -not -name 'test_*.go' -not -name '*_test.go')

# Go 编译参数
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.Commit=$(COMMIT) -w -s"
GO_BUILD_ENV = CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH)

# 输出目录结构：bin/linux/<arch>/$(PROJECT_NAME)
define GOBIN
bin/$(OS)/$(1)/$(PROJECT_NAME)
endef

# 默认目标
.DEFAULT_GOAL := build

.PHONY: help
help: ## 显示帮助信息
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: verify-files  
verify-files: ## 验证要编译的文件列表，确保测试文件被正确排除
	@echo "=== Verifying Go files to be compiled ==="
	@echo "Files to be compiled:"
	@find . -name '*.go' -not -path './vendor/*' -not -path './go/*' -not -name 'test_*.go' -not -name '*_test.go' | sort
	@echo ""
	@echo "Excluded test files:"
	@find . -name 'test_*.go' -o -name '*_test.go' | grep -v './vendor/' | grep -v './go/' | sort || echo "  (none found)"
	@echo ""
	@echo "Total files to compile: $$(find . -name '*.go' -not -path './vendor/*' -not -path './go/*' -not -name 'test_*.go' -not -name '*_test.go' | wc -l | tr -d ' ')"
	@echo "Total test files excluded: $$(find . -name 'test_*.go' -o -name '*_test.go' | grep -v './vendor/' | grep -v './go/' | wc -l | tr -d ' ')"

.PHONY: build
build: test ## 构建所有架构的二进制文件
	@echo "=== Building for all architectures ==="
	@$(foreach arch,$(ARCHS),$(MAKE) ARCH=$(arch) build-single;)
	@echo "=== Build completed ==="
	@echo "Built binaries:"
	@$(foreach arch,$(ARCHS),echo "  $(call GOBIN,$(arch))";)

.PHONY: build-single
build-single: go/modules/pkg/mod ## 构建指定架构或所有架构的二进制文件 (可选: ARCH=amd64|arm64)
	@if [ -z "$(ARCH)" ]; then \
	    echo "=== No ARCH specified, building for all architectures ==="; \
	    $(foreach arch,$(ARCHS),$(MAKE) ARCH=$(arch) build-single;) \
	    echo "=== Build completed ==="; \
	    echo "Built binaries:"; \
	    $(foreach arch,$(ARCHS),echo "  $(call GOBIN,$(arch))";) \
	else \
	    mkdir -p $(dir $(call GOBIN,$(ARCH))); \
	    echo "Building $(PROJECT_NAME) for $(OS)/$(ARCH)..."; \
	    $(GO_BUILD_ENV) go build $(LDFLAGS) -v -o $(call GOBIN,$(ARCH)) .; \
	    echo "✓ Built: $(call GOBIN,$(ARCH))"; \
	fi

go/modules/pkg/mod: go.mod go.sum
	@echo "Downloading Go modules..."
	go mod download
	@touch go/modules/pkg/mod

.PHONY: test
test: go/modules/pkg/mod ## 运行测试
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: clean
clean: ## 清理构建文件和缓存
	@echo "Cleaning up..."
	go clean -modcache
	rm -rf bin/ go/ coverage.out coverage.html
	@echo "✓ Cleaned"