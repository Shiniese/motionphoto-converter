package processor

import (
	"github.com/Shiniese/motionphoto-converter/internal/models"
	"github.com/Shiniese/motionphoto-converter/internal/utils"
)

// HuaweiMotionPhotoProcessor 处理华为动态照片（基于 File Type Box 检测）
type HuaweiMotionPhotoProcessor struct {
	*BaseMotionPhotoProcessor
}

// NewHuaweiMotionPhotoProcessor 创建华为处理器
func NewHuaweiMotionPhotoProcessor() *HuaweiMotionPhotoProcessor {
	return &HuaweiMotionPhotoProcessor{
		BaseMotionPhotoProcessor: NewBaseMotionPhotoProcessor(models.BrandHuawei),
	}
}

// CanProcess 检测是否为华为动态照片
func (p *HuaweiMotionPhotoProcessor) CanProcess(xmpInfo map[string]string) bool {
	make, hasMake := xmpInfo["Make"]
	manufacturer, hasManufacturer := xmpInfo["Manufacturer"]

	return (hasMake && make == "HUAWEI") || (hasManufacturer && manufacturer == "HUAWEI")
}

// CanProcessByFileTypeBox 通过文件类型框检测是否可处理
func (p *HuaweiMotionPhotoProcessor) CanProcessByFileTypeBox(data []byte) bool {
	return utils.FindMP4VideoByFileTypeBox(data) != nil
}

// ProcessMotionPhoto 处理华为动态照片
func (p *HuaweiMotionPhotoProcessor) ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult {
	videoInfo := utils.FindMP4VideoByFileTypeBox(data)
	if videoInfo == nil {
		errMsg := "No MP4 video found using File Type Box detection for Huawei motion photo"
		return models.NewErrorResult(errMsg)
	}

	videoStartOffset := videoInfo.Offset
	videoLength := videoInfo.Length

	imageData := data[:videoStartOffset]
	videoData := data[videoStartOffset : videoStartOffset+videoLength]

	motionPhotoData := models.NewMotionPhotoData(
		imageData, videoData, models.BrandHuawei,
		&videoStartOffset, nil, // 华为格式无时间戳
	)

	return models.NewSuccessResult(motionPhotoData)
}
