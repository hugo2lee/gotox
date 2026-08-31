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

type HashAlgorithm string

const (
	MD5    HashAlgorithm = "Md5"
	SHA1   HashAlgorithm = "Sha1"
	SHA256 HashAlgorithm = "Sha256"
)

type HashBuilder struct {
	algorithms map[HashAlgorithm]func() hash.Hash
	logger     logx.Logger
}

func NewBuilder() *HashBuilder {
	return &HashBuilder{
		algorithms: make(map[HashAlgorithm]func() hash.Hash),
		logger:     logx.NewNoOpLogger(),
	}
}

func (b *HashBuilder) WithLogger(logger logx.Logger) *HashBuilder {
	if logger != nil {
		b.logger = logger
	}
	return b
}

func (b *HashBuilder) WithMd5() *HashBuilder {
	b.algorithms[MD5] = md5.New
	return b
}

func (b *HashBuilder) WithSha1() *HashBuilder {
	b.algorithms[SHA1] = sha1.New
	return b
}

func (b *HashBuilder) WithSha256() *HashBuilder {
	b.algorithms[SHA256] = sha256.New
	return b
}

func (b *HashBuilder) WithAlgorithm(algorithm HashAlgorithm, hasherFunc func() hash.Hash) *HashBuilder {
	b.algorithms[algorithm] = hasherFunc
	return b
}

func (b *HashBuilder) SetHash(c *gin.Context, hv *bodyTemp) {
	if c.Writer.Status() != http.StatusOK {
		return
	}

	bodyStr := hv.body.String()
	if bodyStr == "" {
		b.logger.Error("HashResponse bodyStr is empty")
		return
	}
	for algorithm, hasherFunc := range b.algorithms {
		hasher := hasherFunc()
		_, err := io.WriteString(hasher, bodyStr)
		if err != nil {
			b.logger.Error("HashResponse WriteString %s %v", algorithm, err)
			continue
		}
		hash := hex.EncodeToString(hasher.Sum(nil))
		c.Header(fmt.Sprintf("Content-%s", algorithm), hash)
	}
}

func (b *HashBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		originalWriter := c.Writer
		hv := &bodyTemp{body: bytes.NewBuffer(nil)}
		c.Writer = &hashBodyWriter{ResponseWriter: originalWriter, bt: hv}
		c.Next()
		b.SetHash(c, hv)
		originalWriter.WriteHeaderNow()
		if hv.write {
			if _, err := originalWriter.Write(hv.body.Bytes()); err != nil {
				b.logger.Error("HashResponse Write %v", err)
			}
		}
		if hv.writeString {
			if _, err := originalWriter.WriteString(hv.body.String()); err != nil {
				b.logger.Error("HashResponse WriteString %v", err)
			}
		}
	}
}

type bodyTemp struct {
	body        *bytes.Buffer
	write       bool
	writeString bool
}

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
