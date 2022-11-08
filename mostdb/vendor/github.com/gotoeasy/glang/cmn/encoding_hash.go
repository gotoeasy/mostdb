package cmn

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
)

// 哈希码
func HashCode(bts []byte) uint32 {
	return crc32.ChecksumIEEE(bts)
}

// 哈希码
func Hash(str string) uint32 {
	var rs uint32 = 53653 // 5381
	r := []rune(str)
	for i := len(r) - 1; i >= 0; i-- {
		rs = (rs * 33) ^ uint32(r[i])
	}
	return rs
}

// 字符串哈希处理后取模(余数)，返回值最大不超过mod值
func HashMod(str string, mod uint32) string {
	return fmt.Sprint(Hash("添油"+str+"加醋") % mod)
}

// Base64编码（同Base64Encode）
func Base64(bts []byte) string {
	return base64.StdEncoding.EncodeToString(bts)
}

// Base64编码
func Base64Encode(bts []byte) string {
	return base64.StdEncoding.EncodeToString(bts)
}

// Base64解码
func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

// MD5
func Md5(bts []byte) string {
	h := md5.New()
	h.Write(bts)
	return hex.EncodeToString(h.Sum(nil))
}
