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
	"net/http"

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

// HashBuilder 是中间件构建器
type HashBuilder struct {
	algorithms map[HashAlgorithm]func() hash.Hash
}

// NewBuilder 创建一个新的构建器实例
func NewBuilder() *HashBuilder {
	return &HashBuilder{
		algorithms: make(map[HashAlgorithm]func() hash.Hash),
	}
}

func (hashBuilder *HashBuilder) WithMd5() *HashBuilder {
	hashBuilder.algorithms[MD5] = md5.New
	return hashBuilder
}

func (hashBuilder *HashBuilder) WithSha1() *HashBuilder {
	hashBuilder.algorithms[SHA1] = sha1.New
	return hashBuilder
}

func (hashBuilder *HashBuilder) WithSha256() *HashBuilder {
	hashBuilder.algorithms[SHA256] = sha256.New
	return hashBuilder
}

// WithAlgorithm 添加需要计算的哈希算法
func (hashBuilder *HashBuilder) WithAlgorithm(algorithm HashAlgorithm, hasherFunc func() hash.Hash) *HashBuilder {
	hashBuilder.algorithms[algorithm] = hasherFunc
	return hashBuilder
}

func (hashBuilder *HashBuilder) SetHash(c *gin.Context, hv *bodyTemp) {
	if c.Writer.Status() != http.StatusOK {
		return
	}

	// 获取响应体内容
	bodyStr := hv.body.String()
	if bodyStr == "" {
		logg.Error("HashResponse bodyStr is empty")
		return
	}
	for algorithm, hasherFunc := range hashBuilder.algorithms {
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
func (b *HashBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 保存原始 Writer
		originalWriter := c.Writer
		// 捕获响应 body
		hv := &bodyTemp{
			body: bytes.NewBuffer(nil),
		}
		c.Writer = &hashBodyWriter{
			ResponseWriter: originalWriter,
			bt:             hv,
		}
		c.Next()
		// 在请求处理结束后计算哈希值
		b.SetHash(c, hv)
		originalWriter.WriteHeaderNow()
		if hv.write {
			if _, err := originalWriter.Write(hv.body.Bytes()); err != nil {
				logg.Error("HashResponse Write %v", err)
			}
		}
		if hv.writeString {
			if _, err := originalWriter.WriteString(hv.body.String()); err != nil {
				logg.Error("HashResponse WriteString %v", err)
			}
		}
	}
}

type bodyTemp struct {
	body        *bytes.Buffer
	write       bool
	writeString bool
}

// hashBodyWriter 是一个自定义的 ResponseWriter，用于捕获响应 body
type hashBodyWriter struct {
	gin.ResponseWriter
	bt *bodyTemp
}

func (w hashBodyWriter) Write(b []byte) (int, error) {
	w.bt.write = true
	return w.bt.body.Write(b)
}

func (w hashBodyWriter) WriteString(s string) (int, error) {
	w.bt.writeString = true
	return w.bt.body.WriteString(s)
}
