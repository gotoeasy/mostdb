package web

import (
	"fmt"
	"mostdb/conf"
	"mostdb/most"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/gotoeasy/glang/cmn"
	"github.com/valyala/fasthttp"
)

func Run() {

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
	cmn.Info("ServerUrl", conf.GetServerUrl())
	cmn.Info("ClusterUrls", conf.GetClusterUrls())
	server := &fasthttp.Server{
		Handler:            router.Handler,
		MaxRequestBodySize: 10 * 1024 * 1024,
	}
	err := server.ListenAndServe(fmt.Sprintf(":%s", conf.GetServerPort())) // :5379
	// err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", conf.GetServerPort()), router.Handler) // :5379
	if err != nil && err != http.ErrServerClosed {
		cmn.Fatalln("%s", err) // 启动失败的话打印错误信息后退出
	}
}
