package cmn

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
