# 应用程序的主文件路径
APP_DIR := tests/uat4
APP_MAIN := $(APP_DIR)/load.go

# 输出目录
BUILD_DIR := build

# 应用程序名称
APP_NAME := load

# 默认目标
.PHONY: all
all: build

# 本地编译（当前平台）
.PHONY: build
build:
	@echo "Building for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(APP_MAIN)
	go build -o $(BUILD_DIR)/plow
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME)"

# 交叉编译（Linux）
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux $(APP_MAIN)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/plow 
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME)-linux"

# 交叉编译（Windows）
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME).exe $(APP_MAIN)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/plow.exe
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME).exe"

# 交叉编译（macOS）
.PHONY: build-macos
build-macos:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-macos $(APP_MAIN)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/plow
	@echo "Build completed: $(BUILD_DIR)/$(APP_NAME)-macos"

# 清理生成的文件
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed."

# 帮助信息
.PHONY: help
help:
	@echo "Makefile Usage:"
	@echo "  make          - 编译当前平台的程序"
	@echo "  make build    - 同上"
	@echo "  make build-linux - 交叉编译为 Linux 平台"
	@echo "  make build-windows - 交叉编译为 Windows 平台"
	@echo "  make build-macos - 交叉编译为 macOS 平台"
	@echo "  make clean    - 清理生成的文件"