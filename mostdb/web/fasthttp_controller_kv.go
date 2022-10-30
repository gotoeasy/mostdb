package web

import (
	"encoding/json"
	"mostdb/most"

	"github.com/valyala/fasthttp"
)

func FasthttpKvGetController(c *fasthttp.RequestCtx) {
	crudHandleFasthttp("get", c)
}

func FasthttpKvSetController(c *fasthttp.RequestCtx) {
	crudHandleFasthttp("set", c)
}

func FasthttpKvDelController(c *fasthttp.RequestCtx) {
	crudHandleFasthttp("del", c)
}

func FasthttpKvApiGetController(c *fasthttp.RequestCtx) {
	crudHandleFasthttp("apiget", c)
}

func FasthttpKvApiSetController(c *fasthttp.RequestCtx) {
	crudHandleFasthttp("apiset", c)
}

func FasthttpKvApiDelController(c *fasthttp.RequestCtx) {
	crudHandleFasthttp("apidel", c)
}

func crudHandleFasthttp(opType string, c *fasthttp.RequestCtx) {
	kv := &most.KvData{}
	bts := c.PostBody()
	json.Unmarshal(bts, kv)

	// 提交排队等待异步处理结果
	rs := NewOperationModel(opType, kv).Submit().WaitForOperationResult()

	c.SetContentType("application/json")
	c.SetStatusCode(200)
	jsonBytes, _ := json.Marshal(rs)
	c.SetBody(jsonBytes)

}
