package accesslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/internal/pkg"
	"github.com/hugo2lee/gotox/logx"
)

// 受制于泛型，这里只能使用包变量，如无任何实例赋予就用这个
var logg logx.Logger = logx.Log

// // 自定义的logger，建议实例化赋予
func SetLogger(l logx.Logger) {
	logg = l
}

const (
	Auth             = "Authorization"
	TraceIdName      = "X-Request-Id"
	SpanIdName       = "X-Request-Spanid"
	ParentSpanIdName = "X-Request-Parentspanid"
	GinKeyTraceName  = "gotox-traceid"
	Plaform          = "plaform" // 平台，web、android、ios、pc等
	Token            = "token"   // token 字段，通常是 jwt token 或者其他的 token
)

type AccesslogCtl struct {
	logFunc       func(ctx context.Context, al AccessLog)
	allowStamp    bool
	allowTrace    bool
	allowQuery    bool
	allowReqBody  bool
	allowRespBody bool
}

// fn 的 ctx 其实是 gin.Context
func NewBuilder(fn func(ctx context.Context, al AccessLog)) *AccesslogCtl {
	return &AccesslogCtl{
		logFunc: fn,
		// 默认不打印
		allowTrace:    false,
		allowStamp:    false,
		allowQuery:    false,
		allowReqBody:  false,
		allowRespBody: false,
	}
}

func (b *AccesslogCtl) AllowTrace() *AccesslogCtl {
	b.allowTrace = true
	return b
}

func (b *AccesslogCtl) AllowStamp() *AccesslogCtl {
	b.allowStamp = true
	return b
}

func (b *AccesslogCtl) AllowQuery() *AccesslogCtl {
	b.allowQuery = true
	return b
}

func (b *AccesslogCtl) AllowReqBody() *AccesslogCtl {
	b.allowReqBody = true
	return b
}

func (b *AccesslogCtl) AllowRespBody() *AccesslogCtl {
	b.allowRespBody = true
	return b
}

func (b *AccesslogCtl) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		al := AccessLog{
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
		}

		if b.allowStamp {
			al.TimeStamp = time.Now().UnixMilli()
			al.Ip = fmt.Sprintf("%v|%v", c.ClientIP(), c.RemoteIP())
		}

		if b.allowTrace {
			al.Auth = c.Request.Header.Get(Auth)
			al.Plaform = c.Request.Header.Get(Plaform)
			al.Token = c.Request.Header.Get(Token)
			al.TraceId = c.Request.Header.Get(TraceIdName)
			if al.TraceId == "" {
				al.TraceId = pkg.GenUuid()
				al.ParentSpanId = ""
				al.SpanId = al.TraceId
			} else {
				al.ParentSpanId = c.Request.Header.Get(SpanIdName)
				al.SpanId = pkg.GenUuid()
			}
			if c.Keys == nil {
				c.Keys = make(map[string]any)
			}
			c.Keys[GinKeyTraceName] = al.TraceId
			// c.Keys["spanid"] = al.SpanId
			// c.Keys["parentspanid"] = al.ParentSpanId
		}

		if b.allowQuery {
			al.Query = c.Request.URL.RawQuery
		}

		if b.allowReqBody && c.Request.Body != nil {
			// 直接忽略 error，不影响程序运行
			reqBodyBytes, err := c.GetRawData()
			if err != nil {
				logg.Warn("GetRawData reqBodyBytes %v", err)
			}
			// Request.Body 是一个 Stream（流）对象，所以是只能读取一次的
			// 因此读完之后要放回去，不然后续步骤是读不到的
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			al.ReqBody = string(reqBodyBytes)
		}

		if b.allowRespBody {
			c.Writer = responseWriter{
				ResponseWriter: c.Writer,
				al:             &al,
			}
		}

		defer func() {
			duration := time.Since(start)
			al.Duration = duration.String()

			if b.allowStamp {
				if c.Keys != nil {
					if sn, ok := c.Keys["sn"].(string); ok {
						al.Sn = sn
					}
					if guid, ok := c.Keys["guid"].(string); ok {
						al.Guid = guid
					}
				}
			}

			b.logFunc(c, al)
		}()
		// 这里会执行到业务代码
		c.Next()
	}
}

// AccessLog 你可以打印很多的信息，根据需要自己加
type AccessLog struct {
	// 链路追踪
	TraceId      string `json:"trace_id"`
	SpanId       string `json:"span_id"`
	ParentSpanId string `json:"parent_span_id"`
	Auth         string `json:"authorization"`

	// 业务特征
	TimeStamp int64  `json:"time_stamp"`
	Ip        string `json:"ip"`
	Sn        string `json:"sn"`
	Guid      string `json:"guid"`
	Plaform   string `json:"plaform"` // 平台，web、android、ios、pc等
	Token     string `json:"token"`

	// 业务信息
	Method     string `json:"method"`
	Path       string `json:"path"`
	Query      string `json:"query"`
	ReqBody    string `json:"req_body"`
	Duration   string `json:"duration"`
	StatusCode int    `json:"status_code"`
	RespBody   string `json:"resp_body"`
}

func (al AccessLog) String() string {
	b, err := json.Marshal(al)
	if err != nil {
		logg.Warn("AccessLog Marshal Error %v", err)
	}
	return string(b)
}

// responseWriter 包装 gin.ResponseWriter，实现 gin response 时一并执行方法
type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (r responseWriter) WriteHeader(statusCode int) {
	r.al.StatusCode = statusCode
	r.ResponseWriter.Header().Set(TraceIdName, r.al.TraceId)
	r.ResponseWriter.Header().Set(SpanIdName, r.al.SpanId)
	r.ResponseWriter.Header().Set(ParentSpanIdName, r.al.ParentSpanId)
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r responseWriter) Write(data []byte) (int, error) {
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r responseWriter) WriteString(data string) (int, error) {
	r.al.RespBody = data
	return r.ResponseWriter.WriteString(data)
}
