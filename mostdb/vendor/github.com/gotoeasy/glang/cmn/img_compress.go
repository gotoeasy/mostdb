package cmn

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

// 图片有损压缩（压缩失败时原样保存至目标文件）
// srcFile：压缩前文件
// distFile：压缩后文件
// maxWidth：压缩后的最大宽度
// maxHeight：压缩后的最大高度
// o：压缩比例（nil时为默认80%）
func ImgResize(srcFile string, distFile string, maxWidth uint, maxHeight uint, o *jpeg.Options) {
	by, err := ReadFileBytes(srcFile)
	if err != nil {
		Warn(err)
		WriteFileBytes(distFile, by)
		return
	}

	nby := ImgCompress(by, maxWidth, maxHeight, o)
	err = WriteFileBytes(distFile, nby)
	if err != nil {
		Warn(err)
	}
}

// 图片有损压缩
// maxWidth：压缩后的最大宽度
// maxHeight：压缩后的最大高度
// o：压缩比例（nil时为默认80%）
func ImgCompress(buf []byte, maxWidth uint, maxHeight uint, o *jpeg.Options) []byte {

	// 文件压缩
	img, layout, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		Warn(err)
		return buf
	}

	if o == nil {
		o = &jpeg.Options{Quality: 80}
	}

	// 修改图片的大小(最大3072，百度接口base64编码后小于4M，分辨率不高于4096*4096)
	set := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3) // Lanczos3 算法文件最大，图片最清晰，NearestNeighbor 最差
	newBuf := bytes.Buffer{}
	switch layout {
	case "png":
		err = png.Encode(&newBuf, set)
	case "jpeg", "jpg":
		err = jpeg.Encode(&newBuf, set, o)
	default:
		Warn("暂不支持该文件压缩")
		return buf
	}
	if err != nil {
		Warn(err)
		return buf
	}
	if newBuf.Len() < len(buf) {
		buf = newBuf.Bytes()
	}
	return buf
}
