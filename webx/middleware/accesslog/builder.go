package accesslog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/logx"
)

// 受制于泛型，这里只能使用包变量，如无任何实例赋予就用这个
var logg logx.Logger = logx.Log

// // 自定义的logger，建议实例化赋予
func SetLogger(l logx.Logger) {
	logg = l
}

type MiddlewareBuilder struct {
	logFunc       func(ctx context.Context, al AccessLog)
	allowReqBody  bool
	allowRespBody bool
}

// fn 的 ctx 其实是 gin.Context
func NewMiddlewareBuilder(fn func(ctx context.Context, al AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: fn,
		// 默认不打印
		allowReqBody:  false,
		allowRespBody: false,
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		al := AccessLog{
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
		}
		if b.allowReqBody && c.Request.Body != nil {
			// 直接忽略 error，不影响程序运行
			reqBodyBytes, err := c.GetRawData()
			if err != nil {
				logg.Warn("GetRawData reqBodyBytes ", err)
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
			b.logFunc(c, al)
		}()
		// 这里会执行到业务代码
		c.Next()
	}
}

// AccessLog 你可以打印很多的信息，根据需要自己加
type AccessLog struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	ReqBody    string `json:"req_body"`
	Duration   string `json:"duration"`
	StatusCode int    `json:"status_code"`
	RespBody   string `json:"resp_body"`
}

func (al AccessLog) String() string {
	b, err := json.Marshal(al)
	if err != nil {
		logg.Warn("AccessLog Marshal ", err)
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
