package cmn

import (
	"math/rand"
	"time"
)

// 随机数
func RandomInt(min, max int) int {
	if min == max {
		return min
	}
	if max < min {
		min, max = max, min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

// 随机数
func RandomUint32() uint32 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
}

// 绝对值
func AbsInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// 绝对值
func AbsInt64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}
