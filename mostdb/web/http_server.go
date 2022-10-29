package web

import (
	"fmt"
	"mostdb/conf"
	"mostdb/most"
	"net/http"
	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/gin-gonic/gin"
	"github.com/gotoeasy/glang/cmn"
	"github.com/valyala/fasthttp"
)

type IgnoreGinStdoutWritter struct{}

func (w *IgnoreGinStdoutWritter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func Run() {

	if cmn.EqualsIngoreCase("gin", conf.GetWebFramework()) {
		runWithGin()
	} else if cmn.EqualsIngoreCase("fasthttp", conf.GetWebFramework()) {
		runWithFasthttp()
	}

}

func runWithFasthttp() {

	router := fasthttprouter.New()

	// 增删改查
	router.POST(conf.GetContextPath()+"/get", FasthttpKvGetController)
	router.POST(conf.GetContextPath()+"/set", FasthttpKvSetController)
	router.POST(conf.GetContextPath()+"/del", FasthttpKvDelController)
	router.POST(conf.GetContextPath()+"/api/get", FasthttpKvApiGetController)
	router.POST(conf.GetContextPath()+"/api/set", FasthttpKvApiSetController)
	router.POST(conf.GetContextPath()+"/api/del", FasthttpKvApiDelController)

	// 启动数据库
	most.InitDb()

	// 启动Web服务
	cmn.Info("启动Web服务，端口", conf.GetServerPort())
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", conf.GetServerPort()), router.Handler) // :5379
	if err != nil && err != http.ErrServerClosed {
		cmn.Error("%s", err) // 启动失败的话打印错误信息后退出
		os.Exit(1)
	}
}

func runWithGin() {

	// 配置gin
	gin.DisableConsoleColor()                     // 关闭Gin的日志颜色
	gin.DefaultWriter = &IgnoreGinStdoutWritter{} // 关闭Gin的默认日志输出
	gin.SetMode(gin.ReleaseMode)                  // 开启Gin的Release模式

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
		os.Exit(1)
	}
}
