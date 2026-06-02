package processor

import (
	"strconv"

	"github.com/Shiniese/motionphoto-converter/internal/models"
)

// XiaomiMotionPhotoProcessor 处理小米动态照片
type XiaomiMotionPhotoProcessor struct {
	*BaseMotionPhotoProcessor
}

// NewXiaomiMotionPhotoProcessor 创建小米处理器
func NewXiaomiMotionPhotoProcessor() *XiaomiMotionPhotoProcessor {
	return &XiaomiMotionPhotoProcessor{
		BaseMotionPhotoProcessor: NewBaseMotionPhotoProcessor(models.BrandXiaomi),
	}
}

// CanProcess 检测是否为小米动态照片
func (p *XiaomiMotionPhotoProcessor) CanProcess(xmpInfo map[string]string) bool {
	_, exists := xmpInfo["GCamera:MicroVideoOffset"]
	return exists
}

// ProcessMotionPhoto 处理小米动态照片
func (p *XiaomiMotionPhotoProcessor) ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult {
	microVideoOffsetStr, exists := xmpInfo["GCamera:MicroVideoOffset"]
	if !exists {
		return models.NewErrorResult("Unable to parse video offset for Xiaomi motion photo")
	}

	microVideoOffset, err := strconv.Atoi(microVideoOffsetStr)
	if err != nil {
		return models.NewErrorResult("Unable to parse video offset for Xiaomi motion photo")
	}

	// 提取时间戳
	var presentationTimestamp *float64
	if tsStr, ok := xmpInfo["GCamera:MicroVideoPresentationTimestampUs"]; ok {
		if ts, err := strconv.ParseFloat(tsStr, 64); err == nil {
			presentationTimestamp = &ts
		}
	} else if tsStr, ok := xmpInfo["GCamera:MotionPhotoPresentationTimestampUs"]; ok {
		if ts, err := strconv.ParseFloat(tsStr, 64); err == nil {
			presentationTimestamp = &ts
		}
	}

	// 计算视频起始位置
	videoStartOffset := len(data) - microVideoOffset

	// 提取图片和视频数据
	imageData := data[:videoStartOffset]
	videoData := data[videoStartOffset:]

	motionPhotoData := models.NewMotionPhotoData(
		imageData, videoData, models.BrandXiaomi,
		&videoStartOffset, presentationTimestamp,
	)

	return models.NewSuccessResult(motionPhotoData)
}
