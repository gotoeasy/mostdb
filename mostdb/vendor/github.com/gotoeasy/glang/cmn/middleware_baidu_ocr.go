package cmn

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/nfnt/resize"
)

type baiduToken struct {
	ExpiresIn   int64  `json:"expires_in,omitempty"`   // 默认2592000秒，即30天
	AccessToken string `json:"access_token,omitempty"` // 令牌
	Expire      int64  `json:"expire,omitempty"`       // 最大有效时间（秒）
}

type BaiduOcr struct {
	apiKey    string
	secretKey string
	token     *baiduToken
}

// 创建百度OCR对象（参数apiKey和secretKey在百度注册应用后获取）
func NewBaiduOcr(apiKey string, secretKey string) *BaiduOcr {
	return &BaiduOcr{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

// 定额发票识别，支持png/jpeg/jpg/bpm/pdf类型（返回JSON识别结果）
/*
 支持对各类定额发票的发票代码、发票号码、金额、发票所在地、发票金额小写、省、市7个关键字段进行结构化识别
*/
func (b *BaiduOcr) QuotaInvoice(filebytes []byte, isPdf bool) (string, error) {
	// 接口文档 https://cloud.baidu.com/apiexplorer/index.html?Product=GWSE-DJAQ8YwekkQ&Api=GWAI-ZwLB9psds3b
	var host = "https://aip.baidubce.com/rest/2.0/ocr/v1/quota_invoice"
	uri, err := url.Parse(host)
	if err != nil {
		return "", err
	}
	query := uri.Query()
	query.Set("access_token", b.getAccessToken())
	uri.RawQuery = query.Encode()

	b64str := Base64(filebytes)
	if len(filebytes) > 4*1000*1000 || len(b64str) > 4000*1000 {
		Debug("尝试压缩图片，压缩前", len(filebytes), len(b64str))
		filebytes = b.compressImg(filebytes)
		b64str = Base64(filebytes)
		Debug("尝试压缩图片，压缩后", len(filebytes), len(b64str))
	}
	if len(b64str) > 4096*1024 {
		return "", errors.New("文件太大了")
	}

	formMap := make(map[string]string)
	if isPdf {
		formMap["pdf_file"] = b64str
	} else {
		formMap["image"] = b64str
	}
	return HttpPostForm(uri.String(), formMap)
}

// 增值税发票识别，支持png/jpeg/jpg/bpm/pdf类型（返回JSON识别结果）
/*
支持对增值税普票、专票、卷票、电子发票的所有字段进行结构化识别，
包括发票基本信息、销售方及购买方信息、商品信息、价税信息等，其中四要素识别准确率超过 99.9%；
同时，支持对增值税卷票的 21 个关键字段进行识别，包括发票类型、发票代码、发票号码、机打号码、
机器编号、收款人、销售方名称、销售方纳税人识别号、开票日期、购买方名称、购买方纳税人识别号、
项目、单价、数量、金额、税额、合计金额(小写)、合计金额(大写)、校验码、省、市，四要素平均识别准确率可达95%以上。
*/
func (b *BaiduOcr) VatInvoice(filebytes []byte, isPdf bool) (string, error) {
	// 文档借口 https://cloud.baidu.com/apiexplorer/index.html?Product=GWSE-DJAQ8YwekkQ&Api=GWAI-Cv8DjGvFoje
	var host = "https://aip.baidubce.com/rest/2.0/ocr/v1/vat_invoice"
	uri, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	query := uri.Query()
	query.Set("access_token", b.getAccessToken())
	uri.RawQuery = query.Encode()

	b64str := base64.StdEncoding.EncodeToString(filebytes)
	if len(filebytes) > 4*1000*1000 || len(b64str) > 4000*1000 {
		Debug("尝试压缩图片，压缩前", len(filebytes), len(b64str))
		filebytes = b.compressImg(filebytes)
		b64str = base64.StdEncoding.EncodeToString(filebytes)
		Debug("尝试压缩图片，压缩后", len(filebytes), len(b64str))
	}
	if len(b64str) > 4096*1024 {
		return "", errors.New("文件太大了")
	}

	formMap := make(map[string]string)
	if isPdf {
		formMap["pdf_file"] = b64str
	} else {
		formMap["image"] = b64str
	}
	return HttpPostForm(uri.String(), formMap)
}

// 取令牌
func (b *BaiduOcr) getAccessToken() string {
	if b.token == nil {
		b.initAccessToken()
	}

	if b.token == nil {
		return ""
	}

	if b.token.Expire > time.Now().Unix() {
		return b.token.AccessToken
	}

	b.token = nil
	return b.getAccessToken()
}

func (b *BaiduOcr) initAccessToken() {

	var host = "https://aip.baidubce.com/oauth/2.0/token"
	var param = map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     b.apiKey,
		"client_secret": b.secretKey,
	}

	uri, err := url.Parse(host)
	if err != nil {
		Error(err)
		return
	}
	query := uri.Query()
	for k, v := range param {
		query.Set(k, v)
	}
	uri.RawQuery = query.Encode()

	response, err := http.Get(uri.String())
	if err != nil {
		Error(err)
		return
	}
	bts, err := io.ReadAll(response.Body)
	if err != nil {
		Error(err)
		return
	}

	tk := &baiduToken{}
	err = json.Unmarshal(bts, tk)
	if err != nil {
		Error(err)
		return
	}

	tk.Expire = time.Now().Unix() + tk.ExpiresIn - 60
	b.token = tk
}

func (b *BaiduOcr) compressImg(buf []byte) []byte {

	// 文件压缩
	img, layout, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		Error(err)
		return buf
	}

	// 修改图片的大小(最大3072，百度接口base64编码后小于4M，分辨率不高于4096*4096)
	set := resize.Thumbnail(3072, 3072, img, resize.Lanczos3) // Lanczos3 算法文件最大，图片最清晰，NearestNeighbor 最差
	newBuf := bytes.Buffer{}
	switch layout {
	case "png":
		err = png.Encode(&newBuf, set)
	case "jpeg", "jpg":
		err = jpeg.Encode(&newBuf, set, &jpeg.Options{Quality: 80})
	case "bmp":
		jpgBuf := ImgBmpToJpg(buf, nil)
		return b.compressImg(jpgBuf)
	default:
		Error("暂不支持该文件压缩")
		return buf
	}
	if err != nil {
		Error(err)
		return buf
	}
	if newBuf.Len() < len(buf) {
		buf = newBuf.Bytes()
	}
	return buf
}
