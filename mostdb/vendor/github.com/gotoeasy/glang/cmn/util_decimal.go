package cmn

import (
	"github.com/shopspring/decimal"
)

// float64 转 string
func Float64ToString(num float64) string {
	return decimal.NewFromFloat(num).String()
}

// float64 转 string，四舍五入保留或补足指定位数小数
func FormatRound(num float64, digit int32) string {
	return decimal.NewFromFloat(num).StringFixed(digit)
}

// float64 转 string，整数部分3位一撇，四舍五入保留或补足指定位数小数
func FormatAmountRound(num float64, digit int32) string {
	s := FormatRound(num, digit)
	ary := Split(s, ".")

	// 整数部分三位一撇
	a := Split(ary[0], "")
	cnt := 0
	rs := ""
	for i := len(a) - 1; i >= 0; i-- {
		if cnt == 3 {
			rs = "," + rs
			cnt = 0
		}
		rs = a[i] + rs
		cnt++
	}

	// 有小数时拼接小数部分
	if len(ary) == 2 {
		rs = rs + "." + ary[1]
	}
	return Replace(rs, "-,", "-", 1) // 特殊情况处理（如 -,123,456.789 => -123,456.789）
}

// 四舍五入保留指定位数(0-16)的小数
func Round(num float64, digit int32) float64 {
	if digit < 0 {
		digit = 0
	}
	if digit > 16 {
		digit = 16
	}
	rs, _ := decimal.NewFromFloat(num).Round(digit).Float64()
	return rs
}

// 四舍五入保留1位小数
func Round1(num float64) float64 {
	rs, _ := decimal.NewFromFloat(num).Round(1).Float64()
	return rs
}

// 四舍五入保留2位小数
func Round2(num float64) float64 {
	rs, _ := decimal.NewFromFloat(num).Round(2).Float64()
	return rs
}
