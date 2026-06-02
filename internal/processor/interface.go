package processor

import "github.com/Shiniese/motionphoto-converter/internal/models"

// MotionPhotoProcessor 定义动态照片处理器的接口
type MotionPhotoProcessor interface {
	// Brand 返回处理器支持的品牌
	Brand() models.MotionPhotoBrand

	// CanProcess 检测是否能处理给定 XMP 信息的动态照片
	CanProcess(xmpInfo map[string]string) bool

	// ProcessMotionPhoto 处理动态照片数据
	ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult

	// CalculateStillImageTime 计算静帧图片的时间位置
	CalculateStillImageTime(videoDuration float64, presentationTimestamp *float64, frameRate float64) int
}
