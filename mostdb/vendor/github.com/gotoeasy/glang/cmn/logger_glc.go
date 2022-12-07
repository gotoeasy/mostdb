package cmn

import (
	"fmt"
	"strings"
)

// 日志中心客户端结构体
//
// 日志中心见 https://github.com/gotoeasy/glogcenter
type GLogCenterClient struct {
	apiUrl   string
	system   string
	apiKey   string
	enable   bool
	logLevel int
	logChan  chan string // 用chan控制日志发送顺序
}

// 日志中心选项
type GlcOptions struct {
	ApiUrl   string // 日志中心的添加日志接口地址
	System   string // 系统名（对应日志中心检索页面的分类栏）
	ApiKey   string // 日志中心的ApiKey
	Enable   bool   // 是否开启发送到日志中心
	LogLevel string // 日志级别（trace/debug/info/warn/error/fatal）
}

var glc *GLogCenterClient

// 按环境编配配置初始化glc对象，方便开箱即用，外部使用时可通过SetLogCenterClient重新设定
func init() {
	SetLogCenterClient(NewGLogCenterClient(&GlcOptions{
		ApiUrl:   GetEnvStr("GLC_API_URL", ""),
		System:   GetEnvStr("GLC_SYSTEM", "glang"),
		ApiKey:   GetEnvStr("GLC_API_KEY", ""),
		Enable:   GetEnvBool("GLC_ENABLE", false),
		LogLevel: GetEnvStr("GLC_LOG_LEVEL", "debug"),
	}))
}

// 创建日志中心客户端对象
func NewGLogCenterClient(o *GlcOptions) *GLogCenterClient {
	if o == nil {
		o = &GlcOptions{}
	}

	glc := &GLogCenterClient{
		apiUrl:  o.ApiUrl,
		system:  o.System,
		apiKey:  o.ApiKey,
		enable:  o.Enable,
		logChan: make(chan string, 2048),
	}

	if EqualsIngoreCase("DEBUG", o.LogLevel) {
		glc.logLevel = 1
	} else if EqualsIngoreCase("INFO", o.LogLevel) {
		glc.logLevel = 2
	} else if EqualsIngoreCase("WARN", o.LogLevel) {
		glc.logLevel = 3
	} else if EqualsIngoreCase("ERROR", o.LogLevel) {
		glc.logLevel = 4
	} else if EqualsIngoreCase("FATAL", o.LogLevel) {
		logLevel = 5
	}

	go func() {
		for {
			text := <-glc.logChan
			FasthttpPostJson(glc.apiUrl, text, glc.apiKey)
		}
	}()

	return glc
}

// 设定GLC日志中心客户端
func SetLogCenterClient(glcClient *GLogCenterClient) {
	glc = glcClient
}

// 发送Trace级别日志到日志中心
func (g *GLogCenterClient) Trace(v ...any) {
	if glc.enable && glc.logLevel <= 0 {
		g.Println("TRACE " + fmt.Sprint(v...))
	}
}

// 发送指定系统名的Trace级别日志到日志中心
func (g *GLogCenterClient) TraceSys(system string, v ...any) {
	if glc.enable && glc.logLevel <= 0 {
		g.print(system, "TRACE "+fmt.Sprint(v...))
	}
}

// 发送Debug级别日志到日志中心
func (g *GLogCenterClient) Debug(v ...any) {
	if glc.enable && glc.logLevel <= 1 {
		g.Println("DEBUG " + fmt.Sprint(v...))
	}
}

// 发送Debug级别日志到日志中心
func (g *GLogCenterClient) DebugSys(system string, v ...any) {
	if glc.enable && glc.logLevel <= 1 {
		g.Println("DEBUG " + fmt.Sprint(v...))
	}
}

// 发送Info级别日志到日志中心
func (g *GLogCenterClient) Info(v ...any) {
	if glc.enable && glc.logLevel <= 2 {
		g.Println("INFO " + fmt.Sprint(v...))
	}
}

// 发送指定系统名的Info级别日志到日志中心
func (g *GLogCenterClient) InfoSys(system string, v ...any) {
	if glc.enable && glc.logLevel <= 2 {
		g.print(system, "INFO "+fmt.Sprint(v...))
	}
}

// 发送Warn级别日志到日志中心
func (g *GLogCenterClient) Warn(v ...any) {
	if glc.enable && glc.logLevel <= 3 {
		g.Println("WARN " + fmt.Sprint(v...))
	}
}

// 发送指定系统名的Warn级别日志到日志中心
func (g *GLogCenterClient) WarnSys(system string, v ...any) {
	if glc.enable && glc.logLevel <= 3 {
		g.print(system, "WARN "+fmt.Sprint(v...))
	}
}

// 发送Error级别日志到日志中心
func (g *GLogCenterClient) Error(v ...any) {
	if glc.enable && glc.logLevel <= 4 {
		g.Println("ERROR " + fmt.Sprint(v...))
	}
}

// 发送指定系统名的Error级别日志到日志中心
func (g *GLogCenterClient) ErrorSys(system string, v ...any) {
	if glc.enable && glc.logLevel <= 4 {
		g.print(system, "ERROR "+fmt.Sprint(v...))
	}
}

// 发送Fatal级别日志到日志中心
func (g *GLogCenterClient) Fatal(v ...any) {
	if glc.enable && glc.logLevel <= 5 {
		g.Println("FATAL " + fmt.Sprint(v...))
	}
}

// 发送指定系统名的Fatal级别日志到日志中心
func (g *GLogCenterClient) FatalSys(system string, v ...any) {
	if glc.enable && glc.logLevel <= 5 {
		g.print(system, "FATAL "+fmt.Sprint(v...))
	}
}

// 发送日志到日志中心
func (g *GLogCenterClient) Println(text string) {
	g.print(g.system, text)
}

func (g *GLogCenterClient) print(system string, text string) {
	if IsBlank(text) {
		return
	}
	var data strings.Builder
	data.WriteString("{")
	data.WriteString(`"system":"` + g.encodeGlcJsonValue(system) + `"`)
	data.WriteString(`,"date":"` + FormatSystemDate(FMT_YYYY_MM_DD_HH_MM_SS_SSS) + `"`)
	data.WriteString(`,"text":"` + g.encodeGlcJsonValue(text) + `"`)
	data.WriteString("}")

	g.logChan <- data.String()
}

func (g *GLogCenterClient) encodeGlcJsonValue(v string) string {
	v = ReplaceAll(v, `"`, `\"`)
	v = ReplaceAll(v, "\t", "\\\\t")
	v = ReplaceAll(v, "\r", "\\\\r")
	v = ReplaceAll(v, "\n", "\\\\n")
	return v
}
