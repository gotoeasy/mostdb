package web

import (
	"mostdb/conf"

	"github.com/gotoeasy/glang/cmn"
)

func Run() {

	// 启动数据库
	//most.InitDb()

	httpserver := cmn.NewFasthttpServer()

	// 增删改查
	httpserver.
		HandlePost(conf.GetContextPath()+"/get", FasthttpKvGetController).
		HandlePost(conf.GetContextPath()+"/set", FasthttpKvSetController).
		HandlePost(conf.GetContextPath()+"/del", FasthttpKvDelController).
		HandlePost(conf.GetContextPath()+"/api/get", FasthttpKvApiGetController).
		HandlePost(conf.GetContextPath()+"/api/set", FasthttpKvApiSetController).
		HandlePost(conf.GetContextPath()+"/api/del", FasthttpKvApiDelController)

	// 打印日志
	cmn.Info("启动Web服务，端口", conf.GetServerPort())
	cmn.Info("ServerUrl", conf.GetServerUrl())
	cmn.Info("ClusterUrls", conf.GetClusterUrls())

	// 启动服务
	err := httpserver.SetPort(conf.GetServerPort()).Start()
	if err != nil {
		cmn.Fatalln("启动失败", err) // 启动失败的话打印错误信息后退出
	}
}
