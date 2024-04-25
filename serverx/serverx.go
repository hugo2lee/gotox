package serverx

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
)

// need to close resources
type Resource interface {
	Name() string
	Close(context.Context, *sync.WaitGroup)
}

type Server struct {
	configer  *configx.ConfigCli
	logger    logx.Logger
	Engine    *gin.Engine
	httpSrv   *http.Server
	resources []Resource
}

func New(conf *configx.ConfigCli, log logx.Logger) *Server {
	engine := gin.Default()
	return &Server{
		configer: conf,
		logger:   log,
		Engine:   engine,
		httpSrv: &http.Server{
			Addr:    conf.Addr(),
			Handler: engine,
		},
		resources: make([]Resource, 0),
	}
}

func (s *Server) AddResource(res ...Resource) {
	s.resources = append(s.resources, res...)
}

func (s *Server) CloseResource(ctx context.Context) {
	wg := new(sync.WaitGroup)
	wg.Add(len(s.resources))
	for _, f := range s.resources {
		go f.Close(ctx, wg)
	}
	wgChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgChan)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("close resource timeout")
	case <-wgChan:
		s.logger.Info("all resource close")
	}
}

// 启用http服务
func (s *Server) Run() error {
	return s.Engine.Run(s.configer.Addr())
}

func (s *Server) GracefullyUp() {
	// 监听中断信号
	// notifyCtx, notifyStop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGSYS, syscall.SIGTERM)
	notifyCtx, notifyStop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGTERM)
	defer notifyStop()

	{
		// TODO task 生命周期管理
	}

	{
		// http连接
		go func() {
			if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				s.logger.Error("server error: %s\n", err)
				notifyStop()
				return
			}
		}()
	}

	{
		// 等待中断信号以优雅地关闭服务器
		<-notifyCtx.Done()
		s.logger.Info("shutting down gracefully in 15 seconds..., press Ctrl+C again to force")
	}

	{
		// http带超时关闭
		httpSrvCtx, httpSrvCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer httpSrvCancel()
		if err := s.httpSrv.Shutdown(httpSrvCtx); err != nil {
			s.logger.Error("Server Shutdown error: %v \n", err)
		}
		s.logger.Info("httpSrv close")
	}

	{
		if len(s.resources) != 0 {
			// 资源带超时关闭
			resourceCtx, resourceCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer resourceCancel()
			s.CloseResource(resourceCtx)
		}
	}

	s.logger.Info("Server exiting")
	time.Sleep(2 * time.Second)
	return
}
