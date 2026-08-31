package accesslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/internal/pkg"
	"github.com/hugo2lee/gotox/logx"
)

const (
	Auth             = "Authorization"
	TraceIdName      = "X-Request-Id"
	SpanIdName       = "X-Request-Spanid"
	ParentSpanIdName = "X-Request-Parentspanid"
	GinKeyTraceName  = "gotox-traceid"
	Plaform          = "plaform"
	Token            = "token"

	redactedValue = "[REDACTED]"
)

type AccesslogCtl struct {
	logFunc                   func(ctx context.Context, al AccessLog)
	logger                    logx.Logger
	allowStamp                bool
	allowTrace                bool
	allowQuery                bool
	allowReqBody              bool
	allowRespBody             bool
	allowSensitiveCredentials bool
}

func NewBuilder(fn func(ctx context.Context, al AccessLog)) *AccesslogCtl {
	return &AccesslogCtl{
		logFunc: fn,
		logger:  logx.NewNoOpLogger(),
	}
}

func (b *AccesslogCtl) WithLogger(logger logx.Logger) *AccesslogCtl {
	if logger != nil {
		b.logger = logger
	}
	return b
}

func (b *AccesslogCtl) AllowTrace() *AccesslogCtl { b.allowTrace = true; return b }
func (b *AccesslogCtl) AllowStamp() *AccesslogCtl { b.allowStamp = true; return b }
func (b *AccesslogCtl) AllowQuery() *AccesslogCtl { b.allowQuery = true; return b }
func (b *AccesslogCtl) AllowReqBody() *AccesslogCtl { b.allowReqBody = true; return b }
func (b *AccesslogCtl) AllowRespBody() *AccesslogCtl { b.allowRespBody = true; return b }

// AllowSensitiveCredentials opts into logging raw Authorization/token values.
// Avoid this in persistent logs. By default, credential values are redacted.
func (b *AccesslogCtl) AllowSensitiveCredentials() *AccesslogCtl {
	b.allowSensitiveCredentials = true
	return b
}

func redactCredential(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 && parts[0] != "" {
		return parts[0] + " " + redactedValue
	}
	return redactedValue
}

func (b *AccesslogCtl) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		al := AccessLog{Method: c.Request.Method, Path: c.Request.URL.Path}

		if b.allowStamp {
			al.TimeStamp = time.Now().UnixMilli()
			al.Ip = fmt.Sprintf("%v|%v", c.ClientIP(), c.RemoteIP())
		}

		if b.allowTrace {
			auth := c.Request.Header.Get(Auth)
			token := c.Request.Header.Get(Token)
			if b.allowSensitiveCredentials {
				al.Auth = auth
				al.Token = token
			} else {
				al.Auth = redactCredential(auth)
				al.Token = redactCredential(token)
			}
			al.Plaform = c.Request.Header.Get(Plaform)
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
		}

		if b.allowQuery {
			al.Query = c.Request.URL.RawQuery
		}

		if b.allowReqBody && c.Request.Body != nil {
			reqBodyBytes, err := c.GetRawData()
			if err != nil {
				b.logger.Warn("GetRawData reqBodyBytes %v", err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			al.ReqBody = string(reqBodyBytes)
		}

		if b.allowRespBody {
			c.Writer = responseWriter{ResponseWriter: c.Writer, al: &al}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			if b.allowStamp && c.Keys != nil {
				if sn, ok := c.Keys["sn"].(string); ok {
					al.Sn = sn
				}
				if guid, ok := c.Keys["guid"].(string); ok {
					al.Guid = guid
				}
			}
			b.logFunc(c, al)
		}()
		c.Next()
	}
}

type AccessLog struct {
	TraceId      string `json:"trace_id"`
	SpanId       string `json:"span_id"`
	ParentSpanId string `json:"parent_span_id"`
	Auth         string `json:"authorization"`
	TimeStamp    int64  `json:"time_stamp"`
	Ip           string `json:"ip"`
	Sn           string `json:"sn"`
	Guid         string `json:"guid"`
	Plaform      string `json:"plaform"`
	Token        string `json:"token"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Query        string `json:"query"`
	ReqBody      string `json:"req_body"`
	Duration     string `json:"duration"`
	StatusCode   int    `json:"status_code"`
	RespBody     string `json:"resp_body"`
	logger       logx.Logger
}

func (al AccessLog) String() string {
	b, err := json.Marshal(al)
	if err != nil {
		logger := al.logger
		if logger == nil {
			logger = logx.NewNoOpLogger()
		}
		logger.Warn("AccessLog Marshal Error %v", err)
	}
	return string(b)
}

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
