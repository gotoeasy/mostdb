package web

import (
	"mostdb/conf"
	"mostdb/most"

	"github.com/gotoeasy/glang/cmn"
)

var operationChan chan *OperationModel

func init() {
	operationChan = make(chan *OperationModel, 64)
	go func() {
		for {
			o := <-operationChan
			if o.KvParam.Key == "" {
				o.resultChan <- most.MostResultNg("参数有误")
				continue
			}

			if o.OpType == "set" {
				handleSetOpDataModel(o, false)
			} else if o.OpType == "apiset" {
				handleSetOpDataModel(o, true)
			} else if o.OpType == "get" {
				handleGetOpDataModel(o, false)
			} else if o.OpType == "apiget" {
				handleGetOpDataModel(o, true)
			} else if o.OpType == "del" {
				handleDelOpDataModel(o, false)
			} else if o.OpType == "apidel" {
				handleDelOpDataModel(o, true)
			} else {
				cmn.Error("不支持的操作", o.OpType)
				o.resultChan <- most.MostResultNg("不支持的操作 " + o.OpType)
			}
		}
	}()
}

func submitOperation(o *OperationModel) {
	operationChan <- o
}

// 插入
func handleSetOpDataModel(o *OperationModel, isApi bool) {

	db := most.NewDataStorageHandle("") // TODO 存储名管理

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
				bytes, err := cmn.HttpPostJson(url+conf.GetContextPath()+"/api/add", kvjson)
				if err != nil {
					cmn.Warn("远程调用失败", err)
					return most.MostResultNg(err.Error())
				}

				rs := &most.MostResult{}
				err = rs.LoadJson(cmn.BytesToString(bytes))
				if err != nil {
					cmn.Error("远程调用json结果解析失败", err)
					return most.MostResultNg(err.Error())
				}
				return rs
			})
		}
	}
	// 等待过半一致结果
	future.WaitForMostResult()
	o.resultChan <- future.Result
}

// 删除
func handleDelOpDataModel(o *OperationModel, isApi bool) {

	db := most.NewDataStorageHandle("") // TODO 存储名管理

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
					bytes, err := cmn.HttpPostJson(url+conf.GetContextPath()+"/api/del", kvjson)
					if err != nil {
						cmn.Warn("远程调用失败", err)
						return most.MostResultNg(err.Error())
					}

					rs := &most.MostResult{}
					err = rs.LoadJson(cmn.BytesToString(bytes))
					if err != nil {
						cmn.Error("远程调用json结果解析失败", err)
						return most.MostResultNg(err.Error())
					}
					return rs
				})
			}
		}
	}
	// 等待过半一致结果
	future.WaitForMostResult()
	o.resultChan <- future.Result
}

// 读取
func handleGetOpDataModel(o *OperationModel, isApi bool) {

	db := most.NewDataStorageHandle("") // TODO 存储名管理

	future := most.NewFuture()
	// 本机处理任务
	future.AddTask(func() *most.MostResult {
		bytes, err := db.Get(cmn.StringToBytes(o.KvParam.Key))
		if err != nil {
			cmn.Debug("指定KEY找不到", err)
			return most.MostResultNg(err.Error())
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
					bytes, err := cmn.HttpPostJson(url+conf.GetContextPath()+"/api/get", kvjson)
					if err != nil {
						cmn.Debug("远程调用失败", err)
						return most.MostResultNg(err.Error())
					}

					rs := &most.MostResult{}
					err = rs.LoadJson(cmn.BytesToString(bytes))
					if err != nil {
						cmn.Error("远程调用json结果解析失败", err)
						return most.MostResultNg(err.Error())
					}
					return rs
				})
			}
		}
	}
	// 等待过半一致结果
	future.WaitForMostResult()
	o.resultChan <- future.Result
}
