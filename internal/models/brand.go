package models

// MotionPhotoBrand 表示动态照片的品牌类型
type MotionPhotoBrand string

const (
	BrandXiaomi  MotionPhotoBrand = "Xiaomi"
	BrandAndroid MotionPhotoBrand = "Android"
	BrandHuawei  MotionPhotoBrand = "Huawei"
	BrandUnknown MotionPhotoBrand = "Unknown"
)

// DisplayName 返回品牌的显示名称
func (b MotionPhotoBrand) DisplayName() string {
	switch b {
	case BrandAndroid:
		return "Android (Pixel/Samsung)"
	case BrandHuawei:
		return "Huawei"
	case BrandUnknown:
		return "Unknown (MP4 Detection)"
	default:
		return string(b)
	}
}

// AllBrands 返回所有支持的品牌列表
func AllBrands() []MotionPhotoBrand {
	return []MotionPhotoBrand{BrandXiaomi, BrandAndroid, BrandHuawei, BrandUnknown}
}