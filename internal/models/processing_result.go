package models

// MotionPhotoProcessingResult 表示动态照片处理的结果
type MotionPhotoProcessingResult struct {
	Success      bool
	Data         *MotionPhotoData
	ErrorMessage *string
}

// NewSuccessResult 创建成功的处理结果
func NewSuccessResult(data *MotionPhotoData) *MotionPhotoProcessingResult {
	return &MotionPhotoProcessingResult{
		Success:      true,
		Data:         data,
		ErrorMessage: nil,
	}
}

// NewErrorResult 创建失败的处理结果
func NewErrorResult(errorMessage string) *MotionPhotoProcessingResult {
	return &MotionPhotoProcessingResult{
		Success:      false,
		Data:         nil,
		ErrorMessage: &errorMessage,
	}
}

// Error 返回错误信息（如果存在）
func (r *MotionPhotoProcessingResult) Error() *string {
	return r.ErrorMessage
}