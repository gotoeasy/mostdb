package cmn

import (
	"encoding/binary"
	"strings"
	"unicode/utf8"
	"unsafe"
)

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Uint32ToBytes(num uint32) []byte {
	bkey := make([]byte, 4)
	binary.BigEndian.PutUint32(bkey, num)
	return bkey
}

func BytesToUint32(bytes []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(bytes))
}

func Len(str string) int {
	return utf8.RuneCountInString(str)
}

func Left(str string, length int) string {
	if Len(str) <= length {
		return str
	}

	var rs string
	for i, s := range str {
		if i == length {
			break
		}
		rs = rs + string(s)
	}
	return rs
}

func Right(str string, length int) string {
	lenr := Len(str)
	if lenr <= length {
		return str
	}

	var rs string
	start := lenr - length
	for i, s := range str {
		if i >= start {
			rs = rs + string(s)
		}
	}
	return rs
}

func Trim(str string) string {
	return strings.TrimSpace(str)
}

func Startwiths(str string, startstr string) bool {
	lstr := Left(str, Len(startstr))
	return lstr == startstr
}

func Endwiths(str string, endstr string) bool {
	rstr := Right(str, Len(endstr))
	return rstr == endstr
}

func SubString(str string, start int, end int) string {
	srune := []rune(str)
	slen := len(srune)
	if start >= slen || start >= end || start < 0 {
		return ""
	}

	rs := ""
	for i := start; i < slen && i < end; i++ {
		rs += string(srune[i])
	}
	return rs
}

func IndexOf(str string, substr string) int {
	idx := strings.Index(str, substr)
	return utf8.RuneCountInString(str[:idx])
}

func Contains(str string, substr string) bool {
	return IndexOf(str, substr) >= 0
}

func ContainsIngoreCase(str string, substr string) bool {
	return IndexOf(ToLower(str), ToLower(substr)) >= 0
}

func EqualsIngoreCase(str1 string, str2 string) bool {
	return ToLower(str1) == ToLower(str2)
}

func ToLower(str string) string {
	return strings.ToLower(str)
}

func ToUpper(str string) string {
	return strings.ToUpper(str)
}

func Repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

func PadLeft(str string, pad string, length int) string {
	if length < Len(str) {
		return str
	}
	s := Repeat(pad, length) + str
	max := Len(s)
	return SubString(s, max-length, max)
}

func PadRight(str string, pad string, length int) string {
	if length < Len(str) {
		return str
	}
	s := str + Repeat(pad, length)
	return SubString(s, 0, length)
}

func Replace(str string, old string, new string, n int) string {
	return strings.Replace(str, old, new, n)
}

func ReplaceAll(str string, old string, new string) string {
	return strings.ReplaceAll(str, old, new)
}
