package motionphoto

import (
	"github.com/Shiniese/motionphoto-converter/internal/models"
	"github.com/Shiniese/motionphoto-converter/internal/processor"
	"github.com/Shiniese/motionphoto-converter/internal/utils"
)

// ProcessMotionPhoto 处理动态照片数据，提取静帧图像及相关元数据。
// 该函数首先从原始数据中提取 XMP 信息，然后根据信息选择合适的处理器进行处理。
// 如果未找到合适的处理器或处理失败，将返回错误结果。
//
// 参数:
//
//	data - 动态照片的原始字节数据。
//
// 返回值:
//
//	*models.MotionPhotoProcessingResult - 处理结果，包含成功状态、静帧数据或错误信息。
func ProcessMotionPhoto(data []byte) *models.MotionPhotoProcessingResult {
	xmpInfo := ExtractXMPInfo(data)
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

// DetectBrand 从提供的字节数据中检测运动照片的品牌信息。
//
// 参数:
//
//	data - 包含运动照片数据的字节切片，用于提取品牌标识。
//
// 返回值:
//
//	models.MotionPhotoBrand - 检测到的运动照片品牌类型。如果无法识别品牌或处理失败，则返回 BrandUnknown。
func DetectBrand(data []byte) models.MotionPhotoBrand {
	xmpInfo := ExtractXMPInfo(data)
	proc := processor.DefaultFactory.GetProcessorWithFallback(xmpInfo, data)
	if proc == nil {
		unknown := models.BrandUnknown
		return unknown
	}
	brand := proc.Brand()
	return brand
}

// GetSupportedBrands 获取所有支持的运动照片品牌列表。
//
// 返回值:
//   - []models.MotionPhotoBrand: 包含所有支持的运动照片品牌的切片。
func GetSupportedBrands() []models.MotionPhotoBrand {
	return processor.DefaultFactory.GetAllSupportedBrands()
}

// ExtractXMPInfo 从给定的字节数据中提取 XMP（Extensible Metadata Platform）信息。
//
// 参数:
//
//	data - 包含 XMP 元数据的原始字节切片。
//
// 返回值:
//
//	一个映射，键为 XMP 属性名，值为对应的属性值。如果未找到 XMP 信息，则返回空映射。
func ExtractXMPInfo(data []byte) map[string]string {
	return utils.ExtractXMPInfo(data)
}
