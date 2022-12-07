package cmn

import (
	"bytes"

	"github.com/anaskhan96/soup"
	"golang.org/x/net/html"
)

// 对html进行分词前的文本提取（提取html中的文本并进行反转义）
func GetHtmlText(strHtml string) string {
	doc := soup.HTMLParse(strHtml)
	// doc.FullText() 提取的文本是连续拼接不能满足分词需要，所以复制内容插入空格处理，顺便增加反转义

	var buf bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
			buf.WriteString(" ") // 插入空格
		}
		if n.Type == html.ElementNode {
			f(n.FirstChild)
		}
		if n.NextSibling != nil {
			f(n.NextSibling)
		}
	}

	f(doc.Pointer.FirstChild)
	return html.UnescapeString(buf.String()) // 反转义
}
