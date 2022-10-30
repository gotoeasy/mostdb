package cmn

import (
	"io"
	"net/http"
	"strings"
)

// 固定Content-Type:application/json;charset=UTF-8，其他自定义headers格式为 K:V
func HttpGetJson(url string, headers ...string) ([]byte, error) {

	// 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	for i, max := 0, len(headers); i < max; i++ {
		strs := strings.Split(headers[i], ":")
		req.Header.Set(strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1]))
	}

	// 读取响应内容
	client := http.Client{}
	res, err := client.Do(req)
	defer req.Body.Close()

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
