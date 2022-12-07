package cmn

import "time"

// 计算函数执行时间（毫秒）
func ExecTime(fn func()) int64 {
	start := time.Now()
	fn()
	return time.Since(start).Milliseconds()
}

// 执行回调函数，错误时重试
func Retry(callback func() error, retryTimes int, duration time.Duration) (err error) {
	for i := 1; i <= retryTimes; i++ {
		if err = callback(); err != nil {
			time.Sleep(duration)
			continue
		}
		return
	}
	return
}
