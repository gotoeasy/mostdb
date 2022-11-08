package cmn

import (
	"time"
)

// 当日的yyyymmdd格式
func Today() string {
	return time.Now().Format("20060102")
}

// 当前日期加减天数后的yyyymmdd格式
func GetYyyymmdd(addDays int) string {
	return time.Now().AddDate(0, 0, addDays).Format("20060102")
}
