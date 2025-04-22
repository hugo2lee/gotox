package hashresponse

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/logx"
)

// 受制于泛型，这里只能使用包变量，如无任何实例赋予就用这个
var logg logx.Logger = logx.Log

// // 自定义的logger，建议实例化赋予
func SetLogger(l logx.Logger) {
	logg = l
}

// HashAlgorithm 定义支持的哈希算法类型
type HashAlgorithm string

const (
	MD5    HashAlgorithm = "Md5"
	SHA1   HashAlgorithm = "Sha1"
	SHA256 HashAlgorithm = "Sha256"
)

// ResponseHashBuilder 是中间件构建器
type ResponseHashBuilder struct {
	algorithms map[HashAlgorithm]func() hash.Hash
}

// NewBuilder 创建一个新的构建器实例
func NewBuilder() *ResponseHashBuilder {
	return &ResponseHashBuilder{
		algorithms: make(map[HashAlgorithm]func() hash.Hash),
	}
}

func (b *ResponseHashBuilder) WithMd5() *ResponseHashBuilder {
	b.algorithms[MD5] = md5.New
	return b
}

func (b *ResponseHashBuilder) WithSha1() *ResponseHashBuilder {
	b.algorithms[SHA1] = sha1.New
	return b
}

func (b *ResponseHashBuilder) WithSha256() *ResponseHashBuilder {
	b.algorithms[SHA256] = sha256.New
	return b
}

// WithAlgorithm 添加需要计算的哈希算法
func (b *ResponseHashBuilder) WithAlgorithm(algorithm HashAlgorithm, hasherFunc func() hash.Hash) *ResponseHashBuilder {
	b.algorithms[algorithm] = hasherFunc
	return b
}

func (b *ResponseHashBuilder) SetHash(c *gin.Context, buf *bytes.Buffer) {
	// 获取响应体内容
	bodyStr := buf.String()
	if bodyStr == "" {
		logg.Error("HashResponse bodyStr is empty")
		return
	}
	for algorithm, hasherFunc := range b.algorithms {
		hasher := hasherFunc()
		_, err := io.WriteString(hasher, bodyStr)
		if err != nil {
			logg.Error("HashResponse WriteString %s %v", algorithm, err)
			continue
		}
		hash := hex.EncodeToString(hasher.Sum(nil))
		c.Header(fmt.Sprintf("Content-%s", algorithm), hash)
	}
}

// Build 构建 Gin 中间件
func (b *ResponseHashBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 捕获响应 body
		buf := bytes.NewBuffer(nil)
		c.Writer = &hashResponseBodyWriter{
			ResponseWriter: c.Writer,
			hashMiddleBody: buf,
		}
		c.Next()
		// 在请求处理结束后计算哈希值
		b.SetHash(c, buf)
	}
}

// hashResponseBodyWriter 是一个自定义的 ResponseWriter，用于捕获响应 body
type hashResponseBodyWriter struct {
	gin.ResponseWriter
	hashMiddleBody *bytes.Buffer
}

func (w hashResponseBodyWriter) Write(b []byte) (int, error) {
	_, err := w.hashMiddleBody.Write(b)
	if err != nil {
		logg.Error("hashResponseBodyWriter Write %v", err)
	}
	return w.ResponseWriter.Write(b)
}

func (w hashResponseBodyWriter) WriteString(s string) (int, error) {
	_, err := w.hashMiddleBody.WriteString(s)
	if err != nil {
		logg.Error("hashResponseBodyWriter WriteString %v", err)
	}
	return w.ResponseWriter.WriteString(s)
}
