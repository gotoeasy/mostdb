package cmn

import (
	"io"
	"net/http"
	"strings"
)

// 使用标准包进行Post请求，固定Content-Type:application/json;charset=UTF-8，其他自定义headers格式为 K:V
func HttpPostJson(url string, jsondata string, headers ...string) ([]byte, error) {

	// 请求
	req, err := http.NewRequest("POST", url, strings.NewReader(jsondata))
	if err != nil {
		return nil, err
	}

	// 请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for i, max := 0, len(headers); i < max; i++ {
		strs := Split(headers[i], ":")
		if len(strs) > 1 {
			req.Header.Set(Trim(strs[0]), Trim(strs[1]))
		}
	}

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
