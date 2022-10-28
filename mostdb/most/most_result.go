package most

import (
	"encoding/json"

	"github.com/gotoeasy/glang/cmn"
)

type MostResult struct {
	Code    int     `json:"code,omitempty"`    // 编码
	Message string  `json:"message,omitempty"` // 消息
	Success bool    `json:"success,omitempty"` // 成功与否
	Result  *KvData `json:"result,omitempty"`  // 数据
}

type KvData struct {
	Key   string `json:"key,omitempty"`   // 键
	Value string `json:"value,omitempty"` // 值
}

func MostResultNg(msg string) *MostResult {
	return &MostResult{
		Code:    500,
		Message: msg,
		Success: false,
	}
}

func MostResultOk() *MostResult {
	return &MostResult{
		Code:    200,
		Success: true,
	}
}

func MostResultOK(d *KvData) *MostResult {
	return &MostResult{
		Code:    200,
		Success: true,
		Result:  d,
	}
}

func (m *MostResult) ToJson() string {
	if m == nil {
		return "{}"
	}
	bt, _ := json.Marshal(m)
	return cmn.BytesToString(bt)
}

func (m *MostResult) LoadJson(jsonstr string) error {
	return json.Unmarshal(cmn.StringToBytes(jsonstr), m)
}

func (k *KvData) ToJson() string {
	if k == nil {
		return "{}"
	}
	bt, _ := json.Marshal(k)
	return cmn.BytesToString(bt)
}
