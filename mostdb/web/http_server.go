package web

import (
	"fmt"
	"mostdb/conf"
	"mostdb/most"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gotoeasy/glang/cmn"
)

type IgnoreGinStdoutWritter struct{}

func (w *IgnoreGinStdoutWritter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func Run() {

	// 配置gin
	gin.DisableConsoleColor()                     // 关闭Gin的日志颜色
	gin.DefaultWriter = &IgnoreGinStdoutWritter{} // 关闭Gin的默认日志输出
	gin.SetMode(gin.ReleaseMode)                  // 开启Gin的Release模式

	// 按配置判断启用GZIP压缩
	ginEngine := gin.Default()

	// 增删改查
	ginEngine.POST(conf.GetContextPath()+"/get", KvGetController)
	ginEngine.POST(conf.GetContextPath()+"/set", KvSetController)
	ginEngine.POST(conf.GetContextPath()+"/del", KvDelController)
	ginEngine.POST(conf.GetContextPath()+"/api/get", KvApiGetController)
	ginEngine.POST(conf.GetContextPath()+"/api/set", KvApiSetController)
	ginEngine.POST(conf.GetContextPath()+"/api/del", KvApiDelController)

	// 启动数据库
	most.InitDb()

	// 启动Web服务
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", conf.GetServerPort()), // :5379
		Handler: ginEngine,
	}
	cmn.Info("启动Web服务，端口", conf.GetServerPort())
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		cmn.Error("%s", err) // 启动失败的话打印错误信息后退出
	}
}
