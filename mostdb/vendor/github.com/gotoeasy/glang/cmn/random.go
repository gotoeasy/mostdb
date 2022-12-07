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

// 随机数
func RandomUint64() uint64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
}

// 随机半角英数字符串
func RandomString(length int) string {
	// str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return BytesToString(result)
}
