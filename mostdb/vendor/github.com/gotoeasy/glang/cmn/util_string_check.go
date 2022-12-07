package cmn

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// 判断是否数值（123、123.456、-123.456都认为是数值）
func IsNumber(s string) bool {
	dotCount := 0
	length := len(s)
	if length == 0 {
		return false
	}

	for i := 0; i < length; i++ {
		if s[i] == '-' && i == 0 {
			continue
		}
		if s[i] == '.' {
			dotCount++
			if i > 0 && i < length-1 {
				continue
			} else {
				return false
			}
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return dotCount <= 1
}

// 判断是否半角数字
func IsDigit(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}

	for i := 0; i < length; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// 判断是否半角字母
func IsAlpha(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}

	for i := 0; i < length; i++ {
		if !(s[i] > 'a' && s[i] < 'z' || s[i] > 'A' && s[i] < 'Z') {
			return false
		}
	}
	return true
}

// 判断是否半角英数（不含符号）
func IsAlphaDigit(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}

	for i := 0; i < length; i++ {
		if !(s[i] > 'a' && s[i] < 'z' || s[i] > 'A' && s[i] < 'Z' || s[i] < '0' || s[i] > '9') {
			return false
		}
	}
	return true
}

// 判断是否手机号
func IsMobile(phone string) bool {
	if len([]rune(phone)) != 11 {
		return false
	}
	reg, err := regexp.Compile(`^1([38][0-9]|4[5679]|5[^4]|6[2567]|7[0-8]|9[0-35-9])\d{8}$`)
	if err != nil {
		return false
	}
	return reg.Match([]byte(phone))
}

// 判断是否身份证号码，若按标准校验失败也返回false
func IsIdCard(idCard string) bool {
	// 计算规则参考“中国国家标准化管理委员会”官方文档：http://www.gb688.cn/bzgk/gb/newGbInfo?hcno=080D6FBF2BB468F9007657F26D60013E
	if len([]rune(idCard)) != 18 {
		return false
	}
	// a1与对应的校验码对照表，其中key表示a1，value表示校验码，value中的10表示校验码X
	var a1Map = map[int]int{
		0:  1,
		1:  0,
		2:  10,
		3:  9,
		4:  8,
		5:  7,
		6:  6,
		7:  5,
		8:  4,
		9:  3,
		10: 2,
	}

	var idStr = strings.ToUpper(idCard)
	var reg, err = regexp.Compile(`^[0-9]{17}[0-9X]$`)
	if err != nil {
		return false
	}
	if !reg.Match([]byte(idStr)) {
		return false
	}
	var sum int
	var signChar = ""
	for index, c := range idStr {
		var i = 18 - index
		if i != 1 {
			if v, err := strconv.Atoi(string(c)); err == nil {
				// 计算加权因子
				var weight = int(math.Pow(2, float64(i-1))) % 11
				sum += v * weight
			} else {
				return false
			}
		} else {
			signChar = string(c)
		}
	}
	var a1 = a1Map[sum%11]
	var a1Str = fmt.Sprintf("%d", a1)
	if a1 == 10 {
		a1Str = "X"
	}
	return a1Str == signChar
}
