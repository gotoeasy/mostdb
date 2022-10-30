package cmn

import (
	"bytes"
	"image/jpeg"
	"log"

	"golang.org/x/image/bmp"
)

func BmpToJpg(buf []byte, o *jpeg.Options) []byte {

	img, err := bmp.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Println(err)
		return buf
	}

	if o == nil {
		o = &jpeg.Options{Quality: 80}
	}

	newBuf := bytes.Buffer{}
	err = jpeg.Encode(&newBuf, img, nil)
	if err != nil {
		log.Println(err)
		return buf
	}

	return newBuf.Bytes()
}
