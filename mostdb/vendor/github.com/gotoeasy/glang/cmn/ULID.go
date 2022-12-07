package cmn

import "github.com/oklog/ulid/v2"

// ULID
//
// ULID常作为UUID的替代方案，固定26位长度（10位时间戳+16位随机数），适用于数据库ID。
// 主要特点：毫秒精度有序（仅同一毫秒内无序）、无特殊字符
func ULID() string {
	return ulid.Make().String()
}
