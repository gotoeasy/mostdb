package cmn

import (
	"github.com/disintegration/imaging"
)

// 高斯模糊处理
// srcFile：处理前文件
// distFile：处理后文件
// sigma：模糊比例（通常可设定为5）
func ImgBlur(srcFile string, distFile string, sigma float64) error {
	src, err := imaging.Open(srcFile)
	if err != nil {
		return err
	}

	distImg := imaging.Blur(src, sigma)
	return imaging.Save(distImg, distFile)
}
