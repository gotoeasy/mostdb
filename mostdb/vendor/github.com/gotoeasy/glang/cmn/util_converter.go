package cmn

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"unsafe"
)

// 字符串(10进制无符号整数形式)转int，超过int最大值会丢失精度，转换失败时返回默认值
func StringToInt(s string, defaultVal int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// float64 转 int64
func Float64ToInt64(f float64) int64 {
	return int64(f)
}

// int64 转 string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// uint64 转 string
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// int 转 string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// 字符串(10进制无符号整数形式)转uint32，超过uint32最大值会丢失精度，转换失败时返回默认值
func StringToUint32(s string, defaultVal uint32) uint32 {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return defaultVal
	}
	return uint32(v & 0xFFFFFFFF)
}

// int 转 []byte
func IntToBytes(intNum int) []byte {
	data := int64(intNum)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

// int64 转 []byte
func Int64ToBytes(i int64) []byte {
	bkey := make([]byte, 8)
	binary.BigEndian.PutUint64(bkey, uint64(i))
	return bkey
}

// uint32 转 string
func Uint32ToString(num uint32) string {
	return BytesToString(Uint32ToBytes(num))
}

// uint32 转 []byte
func Uint32ToBytes(num uint32) []byte {
	bkey := make([]byte, 4)
	binary.BigEndian.PutUint32(bkey, num)
	return bkey
}

// uint16 转 []byte
func Uint16ToBytes(num uint16) []byte {
	bkey := make([]byte, 2)
	binary.BigEndian.PutUint16(bkey, num)
	return bkey
}

// uint64 转 []byte
func Uint64ToBytes(num uint64) []byte {
	bkey := make([]byte, 8)
	binary.BigEndian.PutUint64(bkey, num)
	return bkey
}

// []byte 转 uint32
func BytesToUint32(bytes []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(bytes))
}

// []byte 转 uint64
func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

// string 转 []byte
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// []byte 转 string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// string 转 bool
func StringToBool(s string, defaultVal bool) bool {
	lower := ToLower(s)
	if lower == "true" {
		return true
	}
	if lower == "false" {
		return false
	}
	return defaultVal
}

// bool 转 string
func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
