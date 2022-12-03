package web

import (
	"mostdb/conf"
	"mostdb/most"
	"path/filepath"

	"github.com/gotoeasy/glang/cmn"
)

func submitOperationFasthttp(o *OperationModel) {

	if o.KvParam.Key == "" {
		o.opResultChan <- most.MostResultNg("参数有误")
		return
	}

	if o.OpType == "set" {
		handleSetOpDataModelFasthttp(o, false)
	} else if o.OpType == "apiset" {
		handleSetOpDataModelFasthttp(o, true)
	} else if o.OpType == "get" {
		handleGetOpDataModelFasthttp(o, false)
	} else if o.OpType == "apiget" {
		handleGetOpDataModelFasthttp(o, true)
	} else if o.OpType == "del" {
		handleDelOpDataModelFasthttp(o, false)
	} else if o.OpType == "apidel" {
		handleDelOpDataModelFasthttp(o, true)
	} else {
		cmn.Error("不支持的操作", o.OpType)
		o.opResultChan <- most.MostResultNg("不支持的操作 " + o.OpType)
	}
}

// 保存
func handleSetOpDataModelFasthttp(o *OperationModel, isApi bool) {

	db := cmn.NewLevelDB(getDbPath(), nil) // TODO 存储名管理

	future := most.NewFuture()
	// 本机处理任务
	future.AddTask(func() *most.MostResult {
		err := db.Put(cmn.StringToBytes(o.KvParam.Key), cmn.StringToBytes(o.KvParam.Value))
		if err != nil {
			cmn.Error("本地调用失败", err)
			return most.MostResultNg(err.Error())
		}
		return most.MostResultOk()
	})
	// 节点处理任务
	size := len(conf.GetClusterUrls())
	if !isApi && size > 0 {
		kvjson := o.KvParam.ToJson()
		for i := 0; i < size; i++ {
			url := conf.GetClusterUrls()[i]
			future.AddTask(func() *most.MostResult {
				bytes, err := cmn.FasthttpPostJson(url+conf.GetContextPath()+"/api/set", kvjson)
				if err != nil {
					cmn.Error("远程调用失败", err)
					return most.MostResultNg(err.Error())
				}

				rs := &most.MostResult{}
				defer func() {
					if e := recover(); e != nil {
						cmn.Error("网络抖动导致接收数据不完整", cmn.BytesToString(bytes), "异常", e)
					}
				}()

				err = rs.LoadJson(cmn.BytesToString(bytes))
				if err != nil {
					cmn.Error("远程调用json结果解析失败", cmn.BytesToString(bytes), "异常：", err)
					return most.MostResultNg(err.Error())
				}
				return rs
			})
		}
	}
	// 等待过半一致结果
	future.WaitForMostResult()
	o.opResultChan <- future.Result
}

// 删除
func handleDelOpDataModelFasthttp(o *OperationModel, isApi bool) {

	db := cmn.NewLevelDB(getDbPath(), nil) // TODO 存储名管理

	future := most.NewFuture()
	// 本机处理任务
	future.AddTask(func() *most.MostResult {
		err := db.Del(cmn.StringToBytes(o.KvParam.Key))
		if err != nil {
			cmn.Error("本地调用失败", err)
			return most.MostResultNg(err.Error())
		}
		return most.MostResultOk()
	})
	// 节点处理任务
	if !isApi {
		kvjson := o.KvParam.ToJson()
		for i := 0; i < len(conf.GetClusterUrls()); i++ {
			url := conf.GetClusterUrls()[i]
			if conf.GetServerUrl() != url {
				future.AddTask(func() *most.MostResult {
					bytes, err := cmn.FasthttpPostJson(url+conf.GetContextPath()+"/api/del", kvjson)
					if err != nil {
						cmn.Error("远程调用失败", err)
						return most.MostResultNg(err.Error())
					}

					rs := &most.MostResult{}
					defer func() {
						if e := recover(); e != nil {
							cmn.Error("网络抖动导致接收数据不完整", cmn.BytesToString(bytes), "异常", e)
						}
					}()

					err = rs.LoadJson(cmn.BytesToString(bytes))
					if err != nil {
						cmn.Error("远程调用json结果解析失败", cmn.BytesToString(bytes), "异常：", err)
						return most.MostResultNg(err.Error())
					}
					return rs
				})
			}
		}
	}
	// 等待过半一致结果
	future.WaitForMostResult()
	o.opResultChan <- future.Result
}

// 读取
func handleGetOpDataModelFasthttp(o *OperationModel, isApi bool) {

	db := cmn.NewLevelDB(getDbPath(), nil) // TODO 存储名管理

	future := most.NewFuture()
	// 本机处理任务
	future.AddTask(func() *most.MostResult {
		bytes, err := db.Get(cmn.StringToBytes(o.KvParam.Key))
		if err != nil || bytes == nil {
			cmn.Error("指定KEY找不到或值为nil", err)
			return most.MostResultNg("NotFound")
		}
		return most.MostResultOK(&most.KvData{
			Key:   o.KvParam.Key,
			Value: cmn.BytesToString(bytes),
		})
	})
	// 节点处理任务
	if !isApi {
		kvjson := o.KvParam.ToJson()
		for i := 0; i < len(conf.GetClusterUrls()); i++ {
			url := conf.GetClusterUrls()[i]
			if conf.GetServerUrl() != url {
				future.AddTask(func() *most.MostResult {
					bytes, err := cmn.FasthttpPostJson(url+conf.GetContextPath()+"/api/get", kvjson)
					if err != nil {
						cmn.Error("远程调用失败", err)
						return most.MostResultNg(err.Error())
					}

					rs := &most.MostResult{}
					defer func() {
						if e := recover(); e != nil {
							cmn.Error("网络抖动导致接收数据不完整", cmn.BytesToString(bytes), "异常", e)
						}
					}()

					err = rs.LoadJson(cmn.BytesToString(bytes))
					if err != nil {
						cmn.Error("远程调用json结果解析失败", cmn.BytesToString(bytes), "异常：", err)
						return most.MostResultNg(err.Error())
					}

					return rs
				})
			}
		}
	}
	// 等待过半一致结果
	future.WaitForMostResult()
	o.opResultChan <- future.Result
}

func getDbPath() string {
	return filepath.Join(conf.GetStorageRoot(), "store")
}
