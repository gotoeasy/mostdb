package cmn

import (
	"io"
	"net/http"
	"strings"
)

// 使用标准包进行Post请求，固定Content-Type:application/x-www-form-urlencoded，其他自定义headers格式为 K:V
func HttpPostForm(url string, formMap map[string]string, headers ...string) (string, error) {

	sendBody := http.Request{}
	sendBody.ParseForm()

	for k, v := range formMap {
		sendBody.Form.Add(k, v)
	}
	sendData := sendBody.Form.Encode()

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, strings.NewReader(sendData))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for i, max := 0, len(headers); i < max; i++ {
		strs := Split(headers[i], ":")
		request.Header.Set(Trim(strs[0]), Trim(strs[1]))
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return BytesToString(result), nil
}
