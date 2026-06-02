package motionphoto

import (
	"github.com/Shiniese/motionphoto-converter/internal/models"
	"github.com/Shiniese/motionphoto-converter/internal/processor"
	"github.com/Shiniese/motionphoto-converter/internal/utils"
)

// ProcessMotionPhoto 处理动态照片的主入口函数
func ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult {
	proc := processor.DefaultFactory.GetProcessorWithFallback(xmpInfo, data)
	if proc == nil {
		errMsg := "No suitable processor found for motion photo"
		return models.NewErrorResult(errMsg)
	}

	result := proc.ProcessMotionPhoto(data, xmpInfo)
	if !result.Success || result.Data == nil {
		return result
	}

	// 计算静帧时间（如果有时间戳和时长信息）
	// 注意：实际使用中需要从视频元数据获取 duration 和 frameRate
	// 这里保留接口供调用方传入
	// result.Data.StillImageTime = proc.CalculateStillImageTime(duration, timestamp, frameRate)

	return result
}

// DetectBrand 检测动态照片的品牌
func DetectBrand(data []byte, xmpInfo map[string]string) *models.MotionPhotoBrand {
	proc := processor.DefaultFactory.GetProcessorWithFallback(xmpInfo, data)
	if proc == nil {
		unknown := models.BrandUnknown
		return &unknown
	}
	brand := proc.Brand()
	return &brand
}

// GetSupportedBrands 获取所有支持的品牌列表
func GetSupportedBrands() []models.MotionPhotoBrand {
	return processor.DefaultFactory.GetAllSupportedBrands()
}

func ExtractXMPInfo(data []byte) map[string]string {
	return utils.ExtractXMPInfo(data)
}
