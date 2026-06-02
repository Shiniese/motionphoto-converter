package processor

import (
	"github.com/Shiniese/motionphoto-converter/internal/models"
	"github.com/Shiniese/motionphoto-converter/internal/utils"
)

// UnknownMotionPhotoProcessor 处理未知品牌的动态照片（基于 File Type Box 检测）
type UnknownMotionPhotoProcessor struct {
	*BaseMotionPhotoProcessor
}

// NewUnknownMotionPhotoProcessor 创建未知品牌处理器
func NewUnknownMotionPhotoProcessor() *UnknownMotionPhotoProcessor {
	return &UnknownMotionPhotoProcessor{
		BaseMotionPhotoProcessor: NewBaseMotionPhotoProcessor(models.BrandUnknown),
	}
}

// CanProcess 默认返回 false，作为降级方案使用
func (p *UnknownMotionPhotoProcessor) CanProcess(xmpInfo map[string]string) bool {
	return false
}

// CanProcessByFileTypeBox 通过文件类型框检测是否可处理
func (p *UnknownMotionPhotoProcessor) CanProcessByFileTypeBox(data []byte) bool {
	return utils.FindMP4VideoByFileTypeBox(data) != nil
}

// ProcessMotionPhoto 处理未知品牌动态照片
func (p *UnknownMotionPhotoProcessor) ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult {
	videoInfo := utils.FindMP4VideoByFileTypeBox(data)
	if videoInfo == nil {
		errMsg := "No MP4 video found using File Type Box detection"
		return models.NewErrorResult(errMsg)
	}

	videoStartOffset := videoInfo.Offset
	videoLength := videoInfo.Length

	imageData := data[:videoStartOffset]
	videoData := data[videoStartOffset : videoStartOffset+videoLength]

	motionPhotoData := models.NewMotionPhotoData(
		imageData, videoData, models.BrandUnknown,
		&videoStartOffset, nil, // 未知格式无时间戳
	)

	return models.NewSuccessResult(motionPhotoData)
}
