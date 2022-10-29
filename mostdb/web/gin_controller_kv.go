package web

import (
	"mostdb/most"

	"github.com/gin-gonic/gin"
)

func KvSetController(c *gin.Context) {
	crudHandle("set", c)
}

func KvDelController(c *gin.Context) {
	crudHandle("del", c)
}

func KvGetController(c *gin.Context) {
	crudHandle("get", c)
}

func KvApiSetController(c *gin.Context) {
	crudHandle("apiset", c)
}

func KvApiDelController(c *gin.Context) {
	crudHandle("apidel", c)
}

func KvApiGetController(c *gin.Context) {
	crudHandle("apiget", c)
}

func crudHandle(opType string, c *gin.Context) {
	kv := &most.KvData{}
	err := c.BindJSON(kv)
	if err != nil {
		c.JSON(200, most.MostResultNg(err.Error()))
		return
	}

	if opType == "get" {
		// 读取本地及节点数据，不需要排队
		o := NewOperationModel(opType, kv)
		handleGetOpDataModel(o, false)
		rs := <-o.resultChan
		c.JSON(200, rs)
	} else if opType == "apiget" {
		// API读取本地数据，不需要排队
		o := NewOperationModel(opType, kv)
		handleGetOpDataModel(o, true)
		rs := <-o.resultChan
		c.JSON(200, rs)
	} else {
		// 提交排队等待异步处理结果
		rs := NewOperationModel(opType, kv).Submit().WaitForOperationResult()
		c.JSON(200, rs)
	}

}
