package cmn

import (
	"time"

	"github.com/valyala/fasthttp"
)

// 使用Fasthttp进行Post请求，固定Content-Type:application/json;charset=UTF-8，其他自定义headers格式为 K:V
func FasthttpPostJson(url string, jsondata string, headers ...string) ([]byte, error) {

	// req := &fasthttp.Request{}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.SetBody(StringToBytes(jsondata))

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json;charset=UTF-8")
	for i, max := 0, len(headers); i < max; i++ {
		strs := Split(headers[i], ":")
		req.Header.Set(Trim(strs[0]), Trim(strs[1]))
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
