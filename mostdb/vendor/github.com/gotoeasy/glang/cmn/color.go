package cmn

import (
	"strconv"
)

// 颜色转换（HEX => RGB）
//
//	rgb := HexToRgb("ccc") // rgb: [204 204 204]
//	rgb := HexToRgb("aabbcc") // rgb: [170 187 204]
//	rgb := HexToRgb("#aabbcc") // rgb: [170 187 204]
//	rgb := HexToRgb("0xad99c0") // rgb: [170 187 204]
func HexToRgb(hex string) (rgb []int) {
	hex = Trim(hex)
	if hex == "" {
		return
	}

	// like from css. eg "#ccc" "#ad99c0"
	if hex[0] == '#' {
		hex = hex[1:]
	}

	hex = ToLower(hex)
	switch len(hex) {
	case 3: // "ccc"
		hex = BytesToString([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	case 8: // "0xad99c0"
		hex = TrimPrefix(hex, "0x")
	}

	// recheck
	if len(hex) != 6 {
		return
	}

	// convert string to int64
	if i64, err := strconv.ParseInt(hex, 16, 32); err == nil {
		color := int(i64)
		// parse int
		rgb = make([]int, 3)
		rgb[0] = color >> 16
		rgb[1] = (color & 0x00FF00) >> 8
		rgb[2] = color & 0x0000FF
	}
	return
}

// 颜色转换（RGB => HEX）
//
//	hex := RgbToHex(170, 187, 204) // hex: "#aabbcc"
func RgbToHex(r int, g int, b int) string {
	hex := "#"
	if r < 16 {
		hex += "0"
	}
	hex += strconv.FormatInt(int64(r), 16)

	if g < 16 {
		hex += "0"
	}
	hex += strconv.FormatInt(int64(g), 16)

	if b < 16 {
		hex += "0"
	}
	hex += strconv.FormatInt(int64(b), 16)
	return hex
}
