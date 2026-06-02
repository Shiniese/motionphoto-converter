# MotionPhoto Converter (Go)

将各品牌手机拍摄的动态照片（Motion Photo）分离为独立的图片和视频文件。

## 支持的品牌

- 📱 **Xiaomi**: 基于 `GCamera:MicroVideoOffset` 元数据
- 🤖 **Android (Pixel/Samsung)**: 支持 GContainer、Directory Item、GCamera 多种格式
- 🌸 **Huawei**: 基于 MP4 File Type Box (ftyp) 检测
- ❓ **Unknown**: 降级方案，通过文件类型框检测任意嵌入式 MP4

## 项目结构

```
.
├── cmd/
│   └── motionphoto-converter/  # 命令行工具入口
├── internal/
│   ├── models/                  # 数据模型定义
│   ├── processor/               # 各品牌处理器实现
│   └── utils/                   # 通用工具函数
├── pkg/
│   └── motionphoto/             # 公共 API 接口
├── go.mod
└── README.md
```

## 快速开始

```bash
# 克隆项目
git clone <repo-url>
cd motionphoto-converter

# 构建
go build -o motionphoto-converter ./cmd/motionphoto-converter

# 使用
./motionphoto-converter path/to/motion_photo.jpg
```

## 开发指南

### 添加新品牌支持

1. 在 `internal/models/brand.go` 中添加新品牌常量
2. 创建 `internal/processor/<brand>.go` 实现 `MotionPhotoProcessor` 接口
3. 在 `factory.go` 的 `standardProcessors` 中注册新处理器

### 核心接口

```go
type MotionPhotoProcessor interface {
    Brand() models.MotionPhotoBrand
    CanProcess(xmpInfo map[string]string) bool
    ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult
    CalculateStillImageTime(videoDuration float64, presentationTimestamp *float64, frameRate float64) int
}
```

## 致谢

一个偶然的机会，我需要转换动态照片的格式，同时又发现 go 语言项目中转换动态照片的库很少，所以就学习写了一个。书写此项目的过程中帮我理解了 Go 语言的最佳实践，包括接口设计、错误处理和包组织。

- [MotionPhotoConverter](https://github.com/Igloo302/MotionPhotoConverter) - 帮助我理解动态照片的格式和实现