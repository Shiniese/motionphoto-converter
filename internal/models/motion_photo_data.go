package models

// MotionPhotoData 存储解析后的动态照片数据
type MotionPhotoData struct {
	ImageData             []byte
	VideoData             []byte
	StillImageTime        int
	Brand                 MotionPhotoBrand
	VideoOffset           *int
	PresentationTimestamp *float64
}

// NewMotionPhotoData 创建新的 MotionPhotoData 实例
func NewMotionPhotoData(imageData, videoData []byte, brand MotionPhotoBrand, 
	videoOffset *int, presentationTimestamp *float64) *MotionPhotoData {
	return &MotionPhotoData{
		ImageData:             imageData,
		VideoData:             videoData,
		StillImageTime:        0,
		Brand:                 brand,
		VideoOffset:           videoOffset,
		PresentationTimestamp: presentationTimestamp,
	}
}