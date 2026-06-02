package processor

import (
	"strconv"
	"strings"

	"github.com/Shiniese/motionphoto-converter/internal/models"
)

// AndroidMotionPhotoProcessor 处理 Android (Pixel/Samsung) 动态照片
type AndroidMotionPhotoProcessor struct {
	*BaseMotionPhotoProcessor
}

// NewAndroidMotionPhotoProcessor 创建 Android 处理器
func NewAndroidMotionPhotoProcessor() *AndroidMotionPhotoProcessor {
	return &AndroidMotionPhotoProcessor{
		BaseMotionPhotoProcessor: NewBaseMotionPhotoProcessor(models.BrandAndroid),
	}
}

// CanProcess 检测是否为 Android 动态照片
func (p *AndroidMotionPhotoProcessor) CanProcess(xmpInfo map[string]string) bool {
	_, hasGContainer := xmpInfo["GContainer:ItemLength"]
	_, hasDirectoryItem := xmpInfo["Directory Item Length"]
	_, hasGCamera := xmpInfo["GCamera:MotionPhoto"]

	return hasGContainer || hasDirectoryItem || hasGCamera
}

// ProcessMotionPhoto 处理 Android 动态照片
func (p *AndroidMotionPhotoProcessor) ProcessMotionPhoto(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult {
	// 优先尝试 GContainer:ItemLength 格式
	if itemLengthStr, ok := xmpInfo["GContainer:ItemLength"]; ok {
		return p.processGContainerFormat(data, xmpInfo, itemLengthStr)
	}

	// 尝试 Directory Item Length 格式
	if dirItemLengthStr, ok := xmpInfo["Directory Item Length"]; ok {
		return p.processDirectoryItemFormat(data, xmpInfo, dirItemLengthStr)
	}

	// 尝试 GCamera 格式
	if _, ok := xmpInfo["GCamera:MotionPhoto"]; ok {
		return p.processGCameraFormat(data, xmpInfo)
	}

	return models.NewErrorResult("Missing length information for Android motion photo")
}

func (p *AndroidMotionPhotoProcessor) processGContainerFormat(data []byte, xmpInfo map[string]string, itemLengthStr string) *models.MotionPhotoProcessingResult {
	lengthComponents := parseLengthString(itemLengthStr)
	if len(lengthComponents) < 2 {
		return models.NewErrorResult("Unable to parse GContainer length information for Android motion photo")
	}

	// imageLength := lengthComponents[0]
	videoLength := lengthComponents[1]

	// 提取时间戳
	var presentationTimestamp *float64
	if tsStr, ok := xmpInfo["GCamera:MotionPhotoPresentationTimestampUs"]; ok {
		if ts, err := strconv.ParseFloat(tsStr, 64); err == nil {
			presentationTimestamp = &ts
		}
	}

	videoStartOffset := len(data) - videoLength
	imageData := data[:videoStartOffset]
	videoData := data[videoStartOffset:]

	motionPhotoData := models.NewMotionPhotoData(
		imageData, videoData, models.BrandAndroid,
		&videoStartOffset, presentationTimestamp,
	)

	return models.NewSuccessResult(motionPhotoData)
}

func (p *AndroidMotionPhotoProcessor) processDirectoryItemFormat(data []byte, xmpInfo map[string]string, lengthStr string) *models.MotionPhotoProcessingResult {
	lengthComponents := parseLengthString(lengthStr)
	if len(lengthComponents) < 2 {
		return models.NewErrorResult("Unable to parse Directory Item length information for Android motion photo")
	}

	// 提取时间戳
	var presentationTimestamp *float64
	if tsStr, ok := xmpInfo["Motion Photo Presentation Timestamp Us"]; ok {
		if ts, err := strconv.ParseFloat(tsStr, 64); err == nil {
			presentationTimestamp = &ts
		}
	} else if tsStr, ok := xmpInfo["GCamera:MotionPhotoPresentationTimestampUs"]; ok {
		if ts, err := strconv.ParseFloat(tsStr, 64); err == nil {
			presentationTimestamp = &ts
		}
	}

	// 检查是否有 padding 信息（部分三星设备）
	if paddingStr, ok := xmpInfo["Directory Item Padding"]; ok {
		return p.processWithPadding(data, lengthComponents, paddingStr, presentationTimestamp)
	}

	return p.processWithoutPadding(data, lengthComponents[1], presentationTimestamp)
}

func (p *AndroidMotionPhotoProcessor) processWithPadding(data []byte, lengthComponents []int, paddingStr string, presentationTimestamp *float64) *models.MotionPhotoProcessingResult {
	paddingComponents := parseLengthString(paddingStr)
	if len(paddingComponents) < 2 {
		return models.NewErrorResult("Unable to parse padding information for Android motion photo")
	}

	// imageLength := lengthComponents[0]
	videoLength := lengthComponents[1]
	imagePadding := paddingComponents[0]
	videoPadding := paddingComponents[1]

	imageEndOffset := len(data) - videoLength - videoPadding
	videoStartOffset := imageEndOffset + imagePadding

	imageData := data[:imageEndOffset]
	videoData := data[videoStartOffset : videoStartOffset+videoLength]

	motionPhotoData := models.NewMotionPhotoData(
		imageData, videoData, models.BrandAndroid,
		&videoStartOffset, presentationTimestamp,
	)

	return models.NewSuccessResult(motionPhotoData)
}

func (p *AndroidMotionPhotoProcessor) processWithoutPadding(data []byte, videoLength int, presentationTimestamp *float64) *models.MotionPhotoProcessingResult {
	videoStartOffset := len(data) - videoLength
	imageData := data[:videoStartOffset]
	videoData := data[videoStartOffset:]

	motionPhotoData := models.NewMotionPhotoData(
		imageData, videoData, models.BrandAndroid,
		&videoStartOffset, presentationTimestamp,
	)

	return models.NewSuccessResult(motionPhotoData)
}

func (p *AndroidMotionPhotoProcessor) processGCameraFormat(data []byte, xmpInfo map[string]string) *models.MotionPhotoProcessingResult {
	if lengthInfo, ok := xmpInfo["Directory Item Length"]; ok {
		lengths := strings.Split(lengthInfo, ", ")
		if len(lengths) >= 2 {
			if videoLength, err := strconv.Atoi(strings.TrimSpace(lengths[1])); err == nil && videoLength > 0 {
				videoStartOffset := len(data) - videoLength
				if videoStartOffset > 0 && videoStartOffset < len(data) {
					var presentationTimestamp *float64
					if tsStr, ok := xmpInfo["GCamera:MotionPhotoPresentationTimestampUs"]; ok {
						if ts, err := strconv.ParseFloat(tsStr, 64); err == nil {
							presentationTimestamp = &ts
						}
					}

					imageData := data[:videoStartOffset]
					videoData := data[videoStartOffset:]

					motionPhotoData := models.NewMotionPhotoData(
						imageData, videoData, models.BrandAndroid,
						&videoStartOffset, presentationTimestamp,
					)

					return models.NewSuccessResult(motionPhotoData)
				}
			}
		}
	}

	return models.NewErrorResult("Unable to find video data length information in GCamera format Android motion photo")
}

// parseLengthString 解析 "0, 1234567" 格式的字符串
func parseLengthString(s string) []int {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		if val, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			result = append(result, val)
		}
	}
	return result
}
