package cmn

import "github.com/golang/snappy"

// 使用snappy算法压缩（压缩速度快，占用资源少，压缩比适当，重复多则压缩比大，适用于重复较多的文本压缩）
func Compress(srcBytes []byte) []byte {
	return snappy.Encode(nil, srcBytes)
}

// 解压snappy算法压缩的结果（若解压失败返回原参数）
func UnCompress(snappyEncodedBytes []byte) []byte {
	bt, err := snappy.Decode(nil, snappyEncodedBytes)
	if err != nil {
		Trace(err)
		return snappyEncodedBytes
	}
	return bt
}
