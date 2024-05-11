package resourcex

import (
	"context"
	"sync"

	"github.com/hugo2lee/gotox/logx"
)

var logg logx.Logger = logx.Log

// // 自定义的logger，建议实例化赋予
func SetLogger(l logx.Logger) {
	logg = l
}

// need to close resources
type Resource interface {
	Name() string
	Close(context.Context, *sync.WaitGroup)
}

type ResourceCli struct {
	name    string
	closeFn func(ctx context.Context)
}

func NewResourceCli(name string, fn func(ctx context.Context)) *ResourceCli {
	return &ResourceCli{
		name:    name,
		closeFn: fn,
	}
}

func (r *ResourceCli) Name() string {
	return r.name
}

func (r *ResourceCli) Close(ctx context.Context, wg *sync.WaitGroup) {
	r.closeFn(ctx)
	wg.Done()
	logg.Info("resource \"%s\" closed", r.name)
}

type Resourcer struct {
	resources map[string]Resource
}

func NewResourcer() *Resourcer {
	return &Resourcer{
		resources: make(map[string]Resource),
	}
}

func (r *Resourcer) AddResource(res ...Resource) {
	for _, f := range res {
		r.resources[f.Name()] = f
	}
}

func (r *Resourcer) CloseAll(ctx context.Context) {
	wg := new(sync.WaitGroup)
	wg.Add(len(r.resources))
	for _, f := range r.resources {
		go f.Close(ctx, wg)
	}
	wgChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgChan)
	}()

	select {
	case <-ctx.Done():
		logg.Info("close resource timeout")
	case <-wgChan:
		logg.Info("all resource closed")
	}
}
