package cmn

import (
	"sync"

	"github.com/wangbin/jiebago"
)

type TokenizerJiebago struct {
	segmenter jiebago.Segmenter
}

var _segmenterJiebago *TokenizerJiebago
var _segmenterJiebagoMu sync.Mutex

// 分词结果忽略的单字
var _ingoreCharsJiebago = "`~!@# $%^&*()-_=+[{]}\\|;:'\",<.>/?，。《》；：‘　’“”、|】｝【｛＋－—（）×＆…％￥＃＠！～·\t\r\n"

// 创建中文分词器（jiebago）
// 参数dicFile为字典文件，传入空时默认为"data/dictionary.txt"
func NewTokenizerJiebago(dicFile string) *TokenizerJiebago {
	if _segmenterJiebago != nil {
		return _segmenterJiebago
	}
	_segmenterJiebagoMu.Lock()
	defer _segmenterJiebagoMu.Unlock()
	if _segmenterJiebago != nil {
		return _segmenterJiebago
	}

	// 载入词典
	if IsBlank(dicFile) {
		dicFile = "data/dictionary.txt"
	}
	var segmenter jiebago.Segmenter
	segmenter.LoadDictionary(dicFile)

	_segmenterJiebago = &TokenizerJiebago{
		segmenter: segmenter,
	}
	return _segmenterJiebago
}

// 按搜索引擎模式进行分词（自动去重、去标点符号、忽略大小写）
func (t *TokenizerJiebago) CutForSearch(str string) []string {
	sch := t.segmenter.CutForSearch(ToLower(str), true)

	var rs []string
	var mapStr = make(map[string]string)
	for w := range sch {
		if Contains(_ingoreCharsJiebago, w) {
			continue // 忽略标点符号
		}
		if _, has := mapStr[w]; has {
			continue // 去重
		}
		mapStr[w] = ""
		rs = append(rs, w)
	}
	return rs
}
