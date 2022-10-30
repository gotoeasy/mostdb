package cmn

import (
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func FasthttpPostJson(url string, jsondata string, headers ...string) ([]byte, error) {

	// req := &fasthttp.Request{}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.SetBody(StringToBytes(jsondata))

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json;charset=UTF-8")
	for i, max := 0, len(headers); i < max; i++ {
		strs := strings.Split(headers[i], ":")
		req.Header.Set(strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1]))
	}

	// res := &fasthttp.Response{}
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	// 默认5分钟超时，因为超过5分钟通常已没啥意义
	client := &fasthttp.Client{
		ReadTimeout:        5 * time.Minute,
		MaxConnWaitTimeout: 5 * time.Minute,
	}
	err := client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
