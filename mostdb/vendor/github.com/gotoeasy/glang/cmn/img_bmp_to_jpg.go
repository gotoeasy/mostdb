package cmn

import (
	"bytes"
	"image/jpeg"

	"golang.org/x/image/bmp"
)

// bmp文件转jpg文件
func BmpToJpg(buf []byte, o *jpeg.Options) []byte {

	img, err := bmp.Decode(bytes.NewReader(buf))
	if err != nil {
		Error(err)
		return buf
	}

	if o == nil {
		o = &jpeg.Options{Quality: 80}
	}

	newBuf := bytes.Buffer{}
	err = jpeg.Encode(&newBuf, img, nil)
	if err != nil {
		Error(err)
		return buf
	}

	return newBuf.Bytes()
}
