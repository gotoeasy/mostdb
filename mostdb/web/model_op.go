package web

import "mostdb/most"

type OperationModel struct {
	OpType     string                // 操作类型
	KvParam    *most.KvData          // 请求参数
	Result     *most.MostResult      // 处理结果
	resultChan chan *most.MostResult // 处理结果通道
	done       bool                  // 是否已处理完
}

func NewOperationModel(opType string, kvData *most.KvData) *OperationModel {
	return &OperationModel{
		OpType:     opType,
		KvParam:    kvData,
		resultChan: make(chan *most.MostResult, 1),
	}
}

func (o *OperationModel) WaitForOperationResult() *most.MostResult {
	if o.done {
		return o.Result
	}

	rs := <-o.resultChan
	o.Result = rs
	o.done = true

	defer close(o.resultChan)
	return rs
}

func (o *OperationModel) Submit() *OperationModel {
	submitOperation(o)
	return o
}
