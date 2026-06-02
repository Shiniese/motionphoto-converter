# -------------------------- 配置项 --------------------------
# 应用名称
APP_NAME := motionphoto-converter
# 主程序入口路径
MAIN_PATH := ./cmd/motionphoto-converter
# 版权/作者信息
AUTHOR := Shiniese
# 编译输出目录
OUTPUT_DIR := ./build
# 版本号
VERSION := $(shell git describe --tags --always --abbrev=0 --match="v*" 2>/dev/null || echo "v0.0.0")

# Git提交哈希和编译时间，自动注入
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")

# -------------------------- Go编译参数配置 --------------------------
# Go编译注入变量（需要在代码里对应位置定义这些变量）
LDFLAGS := -ldflags "\
  -s -w \
  -X 'main.AppName=$(APP_NAME)' \
  -X 'main.Version=$(VERSION)' \
  -X 'main.GitCommit=$(GIT_COMMIT)' \
  -X 'main.BuildTime=$(BUILD_TIME)' \
  -X 'main.Author=$(AUTHOR)' \
"
# 禁用CGO，保证跨平台二进制纯静态编译无依赖
CGO_ENABLED := 0
GO := CGO_ENABLED=$(CGO_ENABLED) go

# -------------------------- 支持的目标平台（可自行增减） --------------------------
# 格式：操作系统_架构，特殊架构备注：
# - 32位Linux: linux_386
# - 32位Windows: windows_386
# - M1/M2/M3苹果芯片: darwin_arm64
# - Linux ARM64（服务器/树莓派4+）: linux_arm64
# - 树莓派3/2/ARMv7: linux_arm
PLATFORMS := darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64

# -------------------------- 编译规则 --------------------------
# 默认执行：打印帮助信息
.PHONY: all
all: help

# 帮助信息
.PHONY: help
help:
	@echo "Go项目跨平台编译Makefile 用法："
	@echo "  make build          编译当前平台的二进制"
	@echo "  make build-all      编译所有平台的二进制"
	@echo "  make package        编译+打包所有平台（生成zip/tar.gz，带md5校验）"
	@echo "  make release        编译打包+上传到GitHub Release（需要安装gh工具）"
	@echo "  make clean          清理编译产物"
	@echo "  make version        查看当前版本号"

# 查看版本
.PHONY: version
version:
	@echo "当前版本: $(VERSION) (Git commit: $(GIT_COMMIT))"

# 编译当前平台
.PHONY: build
build:
	@mkdir -p $(OUTPUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "编译完成：$(OUTPUT_DIR)/$(APP_NAME)"

# 编译所有平台
.PHONY: build-all
build-all: $(PLATFORMS)

# 单个平台编译规则
.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@$(eval OS := $(word 1, $(subst _, ,$@)))
	@$(eval ARCH := $(word 2, $(subst _, ,$@)))
	@$(eval SUFFIX := $(if $(filter windows,$(OS)),.exe,))  # Make层面判断后缀
	@$(eval OUTPUT_BIN := $(OUTPUT_DIR)/$(APP_NAME)_$(VERSION)_$(OS)_$(ARCH)$(SUFFIX))
	@mkdir -p $(OUTPUT_DIR)
	@echo "正在编译 $(OS)/$(ARCH) 平台..."
	GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(LDFLAGS) -o $(OUTPUT_BIN) $(MAIN_PATH)
	@echo "编译完成：$(OUTPUT_BIN)"

# 打包所有平台（压缩+生成md5）
.PHONY: package
package: build-all
	@echo "正在打包所有编译产物..."
	@cd $(OUTPUT_DIR) && \
	for bin in $(APP_NAME)_$(VERSION)_*; do \
		name=$$(basename $$bin | sed 's/.exe$$//'); \
		if [ "$${bin##*.}" = "exe" ]; then \
			zip $${name}.zip $$bin; \
			md5sum $${name}.zip >> $${name}_md5.txt; \
		else \
			tar -zcf $${name}.tar.gz $$bin; \
			md5sum $${name}.tar.gz >> $${name}_md5.txt; \
		fi; \
		rm -f $$bin; \
	done
	@echo "打包完成，产物在 $(OUTPUT_DIR) 目录"

# 自动上传到GitHub Release（需要先安装GitHub CLI工具：https://cli.github.com/）
.PHONY: release
release: build-all
	@echo "正在上传版本 $(VERSION) 到GitHub Release..."
	gh release create $(VERSION) \
		--title "$(APP_NAME) $(VERSION)" \
		--notes "版本 $(VERSION) 发布，Git提交：$(GIT_COMMIT)，编译时间：$(BUILD_TIME)" \
		$(OUTPUT_DIR)/*
	@echo "上传完成！"

# 清理编译产物
.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	@echo "清理完成"