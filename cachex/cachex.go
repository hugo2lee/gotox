package cachex

import (
	"context"
	"sync"
	"time"

	"github.com/hugo2lee/gotox/resourcex"
	"github.com/patrickmn/go-cache"
)

var (
	_ resourcex.Resource = (*Cachex)(nil)
	_ Cachexer           = (*Cachex)(nil)
)

type Cachexer interface {
	Set(key string, value any)
	Get(key string) (any, bool)
}

type Cachex struct {
	*cache.Cache
}

func New() Cachexer {
	return &Cachex{cache.New(2*time.Minute, 3*time.Minute)}
}

func (c *Cachex) Name() string {
	return "cachex"
}

func (c *Cachex) Set(key string, value any) {
	c.Cache.Set(key, value, cache.DefaultExpiration)
}

func (c *Cachex) Close(ctx context.Context, wg *sync.WaitGroup) {
	c.Cache.Flush()
	wg.Done()
}
