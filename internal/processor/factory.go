package processor

import (
	"log/slog"

	"github.com/Shiniese/motionphoto-converter/internal/models"
)

// MotionPhotoProcessorFactory 动态照片处理器工厂
type MotionPhotoProcessorFactory struct{}

var (
	standardProcessors = []MotionPhotoProcessor{
		NewXiaomiMotionPhotoProcessor(),
		NewAndroidMotionPhotoProcessor(),
		NewHuaweiMotionPhotoProcessor(),
	}

	huaweiProcessor  = NewHuaweiMotionPhotoProcessor()
	unknownProcessor = NewUnknownMotionPhotoProcessor()
)

// GetProcessor 根据 XMP 信息获取合适的处理器
func (f *MotionPhotoProcessorFactory) GetProcessor(xmpInfo map[string]string) MotionPhotoProcessor {
	for _, processor := range standardProcessors {
		if processor.CanProcess(xmpInfo) {
			return processor
		}
	}
	return nil
}

// GetProcessorWithFallback 带降级方案的处理器获取（支持文件类型框检测）
func (f *MotionPhotoProcessorFactory) GetProcessorWithFallback(xmpInfo map[string]string, data []byte) MotionPhotoProcessor {
	// 先尝试标准处理器
	processor := f.GetProcessor(xmpInfo)
	if processor != nil {
		return processor
	}

	// 尝试华为处理器（文件类型框检测）
	slog.Info("No standard processor found, trying Huawei processor")
	if huaweiProcessor.CanProcessByFileTypeBox(data) {
		return huaweiProcessor
	}

	// 尝试未知处理器（文件类型框检测）
	slog.Info("No Huawei processor found, trying Unknown processor")
	if unknownProcessor.CanProcessByFileTypeBox(data) {
		return unknownProcessor
	}

	return nil
}

// GetAllSupportedBrands 返回所有支持的品牌列表
func (f *MotionPhotoProcessorFactory) GetAllSupportedBrands() []models.MotionPhotoBrand {
	brands := make([]models.MotionPhotoBrand, 0)
	for _, p := range standardProcessors {
		brands = append(brands, p.Brand())
	}
	brands = append(brands, models.BrandUnknown)
	return brands
}

// DefaultFactory 默认工厂实例
var DefaultFactory = &MotionPhotoProcessorFactory{}
