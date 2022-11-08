package cmn

import "time"

// 计算函数执行时间（毫秒）
func ExecTime(fn func()) int64 {
	start := time.Now()
	fn()
	return time.Since(start).Milliseconds()
}
