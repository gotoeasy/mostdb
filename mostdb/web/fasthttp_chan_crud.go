package web

import (
	"mostdb/conf"
	"mostdb/most"
	"strings"
	"time"

	"github.com/gotoeasy/glang/cmn"
	"github.com/valyala/fasthttp"
)

func submitOperationFasthttp(o *OperationModel) {

	if o.KvParam.Key == "" {
		o.resultChan <- most.MostResultNg("参数有误")
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
		o.resultChan <- most.MostResultNg("不支持的操作 " + o.OpType)
	}
}

// 插入
func handleSetOpDataModelFasthttp(o *OperationModel, isApi bool) {

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
				bytes, err := FasthttpPostJson(url+conf.GetContextPath()+"/api/set", kvjson)
				if err != nil {
					cmn.Warn("远程调用失败", err)
					return most.MostResultNg(err.Error())
				}

				rs := &most.MostResult{}
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
	o.resultChan <- future.Result
}

// 删除
func handleDelOpDataModelFasthttp(o *OperationModel, isApi bool) {

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
					bytes, err := FasthttpPostJson(url+conf.GetContextPath()+"/api/del", kvjson)
					if err != nil {
						cmn.Warn("远程调用失败", err)
						return most.MostResultNg(err.Error())
					}

					rs := &most.MostResult{}
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
	o.resultChan <- future.Result
}

// 读取
func handleGetOpDataModelFasthttp(o *OperationModel, isApi bool) {

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
					bytes, err := FasthttpPostJson(url+conf.GetContextPath()+"/api/get", kvjson)
					if err != nil {
						cmn.Debug("远程调用失败", err)
						return most.MostResultNg(err.Error())
					}

					rs := &most.MostResult{}
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
	o.resultChan <- future.Result
}

func FasthttpPostJson(url string, jsondata string, headers ...string) ([]byte, error) {

	// req := &fasthttp.Request{}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.SetBody(cmn.StringToBytes(jsondata))

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json;charset=UTF-8")
	for i, max := 0, len(headers); i < max; i++ {
		strs := strings.Split(headers[i], ":")
		req.Header.Set(strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1]))
	}

	// res := &fasthttp.Response{}
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	client := &fasthttp.Client{
		ReadTimeout:        5 * time.Second,
		MaxConnWaitTimeout: 5 * time.Second,
	}
	err := client.Do(req, res)
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
