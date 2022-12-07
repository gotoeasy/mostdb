package cmn

import (
	"strings"
	"unicode/utf8"
)

// 按文字计算字符串长度
func Len(str string) int {
	// return utf8.RuneCountInString(str)
	return len([]rune(str))
}

// 取左文字
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

// 取右文字
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

// 去除两边空格
func Trim(str string) string {
	return strings.TrimSpace(str)
}

// 去除左前缀
func TrimPrefix(str string, prefix string) string {
	return strings.TrimPrefix(str, prefix)
}

// 判断是否空白
func IsBlank(str string) bool {
	return strings.TrimSpace(str) == ""
}

// 判断是否指定前缀
func Startwiths(str string, startstr string) bool {
	lstr := Left(str, Len(startstr))
	return lstr == startstr
}

// 判断是否指定后缀
func Endwiths(str string, endstr string) bool {
	rstr := Right(str, Len(endstr))
	return rstr == endstr
}

// 按文字截取字符串
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

// 查找文字下标
func IndexOf(str string, substr string) int {
	idx := strings.Index(str, substr)
	if idx < 0 {
		return idx
	}
	return utf8.RuneCountInString(str[:idx])
}

// 判断是否包含（区分大小写）
func Contains(str string, substr string) bool {
	return IndexOf(str, substr) >= 0
}

// 判断是否包含（忽略大小写）
func ContainsIngoreCase(str string, substr string) bool {
	return IndexOf(ToLower(str), ToLower(substr)) >= 0
}

// 判断是否相同（忽略大小写）
func EqualsIngoreCase(str1 string, str2 string) bool {
	return ToLower(str1) == ToLower(str2)
}

// 转小写
func ToLower(str string) string {
	return strings.ToLower(str)
}

// 转大写
func ToUpper(str string) string {
	return strings.ToUpper(str)
}

// 重复
func Repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

// 左补足
func PadLeft(str string, pad string, length int) string {
	if length < Len(str) {
		return str
	}
	s := Repeat(pad, length) + str
	max := Len(s)
	return SubString(s, max-length, max)
}

// 右补足
func PadRight(str string, pad string, length int) string {
	if length < Len(str) {
		return str
	}
	s := str + Repeat(pad, length)
	return SubString(s, 0, length)
}

// 替换
func Replace(str string, old string, new string, n int) string {
	return strings.Replace(str, old, new, n)
}

// 全部替换
func ReplaceAll(str string, old string, new string) string {
	return strings.ReplaceAll(str, old, new)
}

// 反转
func Reverse(str string) string {
	r := []rune(str)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// 字符串切割
func Split(str string, sep string) []string {
	return strings.Split(str, sep)
}

// 字符串数组拼接为字符串
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// 字符串去重
func Unique(strs []string) []string {
	m := make(map[string]struct{}, 0)
	newS := make([]string, 0)
	for _, i2 := range strs {
		if _, ok := m[i2]; !ok {
			newS = append(newS, i2)
			m[i2] = struct{}{}
		}
	}
	return newS
}
