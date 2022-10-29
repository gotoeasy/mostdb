/**
 * 系统配置
 * 1）统一管理系统的全部配置信息
 * 2）所有配置都有默认值以便直接使用
 * 3）所有配置都可以通过环境变量设定覆盖，方便自定义配置，方便容器化部署
 */
package conf

import (
	"sort"
	"strings"

	"github.com/gotoeasy/glang/cmn"
)

var storeRoot string
var serverPort string
var serverUrl string
var clusterUrls []string
var contextPath string
var enableSecurityKey bool
var securityKey string
var headerSecurityKey string
var webFramework string

func init() {
	UpdateConfigByEnv()
	cmn.SetLogLevel("INFO") // 日志级别
}

func UpdateConfigByEnv() {
	// 读取环境变量初始化配置，各配置都有默认值
	storeRoot = cmn.GetEnvStr("MOSTDB_STORE_ROOT", "/opt/mostdb")             // 存储根目录
	serverPort = cmn.GetEnvStr("MOSTDB_SERVER_PORT", "5379")                  // web服务端口，默认“5379”
	serverUrl = cmn.GetEnvStr("MOSTDB_SERVER_URL", "")                        // 服务URL，默认“”，集群配置时自动获取地址可能不对，可通过这个设定
	splitUrls(cmn.GetEnvStr("MOSTDB_CLUSTER_URLS", ""))                       // 从服务器地址，多个时逗号分开，默认“”
	contextPath = cmn.GetEnvStr("MOSTDB_CONTEXT_PATH", "/mostdb")             // web服务contextPath
	enableSecurityKey = cmn.GetEnvBool("MOSTDB_API_KEY_ENABLE", false)        // web服务是否开启API秘钥校验，默认false
	headerSecurityKey = cmn.GetEnvStr("MOSTDB_API_KEY_NAME", "X-MOSTDB-AUTH") // web服务API秘钥的header键名
	securityKey = cmn.GetEnvStr("MOSTDB_API_KEY_VAULE", "mostdb")             // web服务API秘钥
	webFramework = cmn.GetEnvStr("MOSTDB_WEB_GIN", "fasthttp")                // web服务框架(gin/fasthttp)
}

func GetWebFramework() string {
	return webFramework
}

func GetServerPort() string {
	return serverPort
}

func GetServerUrl() string {
	return serverUrl
}

func GetClusterUrls() []string {
	return clusterUrls
}

func splitUrls(str string) {
	hosts := strings.Split(str, ";")
	for i := 0; i < len(hosts); i++ {
		url := strings.TrimSpace(hosts[i])
		if url != "" && url != serverUrl {
			clusterUrls = append(clusterUrls, url)
		}
	}

	// 倒序
	sort.Slice(clusterUrls, func(i, j int) bool {
		return clusterUrls[i] > clusterUrls[j]
	})
}

func IsEnableSecurityKey() bool {
	return enableSecurityKey
}

func GetHeaderSecurityKey() string {
	return headerSecurityKey
}

func GetSecurityKey() string {
	return securityKey
}

func GetContextPath() string {
	return contextPath
}

func GetStorageRoot() string {
	return storeRoot
}
