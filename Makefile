.PHONY: build test clean install deps example-simple example-webhook example-api

# 变量定义
GO := go
PROJECT_NAME := kook-go-sdk
BUILD_DIR := build

# 默认目标
all: deps test build

# 安装依赖
deps:
	$(GO) mod download
	$(GO) mod tidy

# 运行测试
test:
	$(GO) test -v ./...

# 构建主程序
build:
	mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(PROJECT_NAME) .

# 构建示例程序
build-examples:
	mkdir -p $(BUILD_DIR)/examples
	$(GO) build -o $(BUILD_DIR)/examples/simple_bot ./examples/simple_bot
	$(GO) build -o $(BUILD_DIR)/examples/webhook_bot ./examples/webhook_bot
	$(GO) build -o $(BUILD_DIR)/examples/api_usage ./examples/api_usage

# 运行简单机器人示例
example-simple:
	$(GO) run examples/simple_bot/main.go

# 运行Webhook机器人示例
example-webhook:
	$(GO) run examples/webhook_bot/main.go

# 运行API使用示例
example-api:
	$(GO) run examples/api_usage/main.go

# 运行主程序测试
run:
	$(GO) run main.go

# 格式化代码
fmt:
	$(GO) fmt ./...

# 检查代码
vet:
	$(GO) vet ./...

# 静态检查
lint:
	golangci-lint run

# 清理构建文件
clean:
	rm -rf $(BUILD_DIR)

# 安装到GOPATH
install:
	$(GO) install .

# 生成文档
docs:
	godoc -http=:6060

# 完整测试（包括格式检查）
test-full: fmt vet test

# 帮助信息
help:
	@echo "可用的命令："
	@echo "  deps          - 安装依赖"
	@echo "  test          - 运行测试"
	@echo "  build         - 构建主程序"
	@echo "  build-examples - 构建示例程序"
	@echo "  run           - 运行主程序测试"
	@echo "  example-simple - 运行简单机器人示例"
	@echo "  example-webhook - 运行Webhook机器人示例"
	@echo "  example-api   - 运行API使用示例"
	@echo "  fmt           - 格式化代码"
	@echo "  vet           - 检查代码"
	@echo "  lint          - 静态检查"
	@echo "  clean         - 清理构建文件"
	@echo "  install       - 安装到GOPATH"
	@echo "  docs          - 生成文档服务器"
	@echo "  test-full     - 完整测试"
	@echo "  help          - 显示此帮助信息" 