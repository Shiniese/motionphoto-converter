package processor

import "github.com/Shiniese/motionphoto-converter/internal/models"

// BaseMotionPhotoProcessor 提供基础实现
type BaseMotionPhotoProcessor struct {
	brand models.MotionPhotoBrand
}

// NewBaseMotionPhotoProcessor 创建基础处理器
func NewBaseMotionPhotoProcessor(brand models.MotionPhotoBrand) *BaseMotionPhotoProcessor {
	return &BaseMotionPhotoProcessor{brand: brand}
}

// Brand 返回品牌
func (b *BaseMotionPhotoProcessor) Brand() models.MotionPhotoBrand {
	return b.brand
}

// CanProcess 默认实现，子类需重写
func (b *BaseMotionPhotoProcessor) CanProcess(xmpInfo map[string]string) bool {
	return false
}

// ProcessMotionPhoto 默认实现，子类需重写
func (b *BaseMotionPhotoProcessor) ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult {
	errMsg := "Not implemented"
	return models.NewErrorResult(errMsg)
}

// CalculateStillImageTime 计算静帧时间
func (b *BaseMotionPhotoProcessor) CalculateStillImageTime(videoDuration float64,
	presentationTimestamp *float64, frameRate float64) int {

	if presentationTimestamp == nil {
		// 默认返回中间帧
		return int(videoDuration * frameRate / 2)
	}

	photoTime := *presentationTimestamp / 1_000_000.0 // 转换为秒
	frameTime := 1.0 / frameRate
	frameNumber := int(photoTime / frameTime)

	maxFrame := int(videoDuration*frameRate) - 1
	if frameNumber < 0 {
		return 0
	}
	if frameNumber > maxFrame {
		return maxFrame
	}
	return frameNumber
}
