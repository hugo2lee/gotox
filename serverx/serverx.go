package serverx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
)

type Server struct {
	configer *configx.Configx
	logger   logx.Logger
	Engine   *gin.Engine
	httpSrv  *http.Server
}

func New(conf *configx.Configx, log logx.Logger) *Server {
	engine := gin.Default()
	return &Server{
		configer: conf,
		logger:   log,
		Engine:   engine,
		httpSrv: &http.Server{
			Addr:    conf.Addr(),
			Handler: engine,
		},
	}
}

// 启用http服务
func (s *Server) Run() error {
	return s.Engine.Run(s.configer.Addr())
}

func (s *Server) GracefullyUp(notifyStop func()) {
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("http server error: %s\n", err)
		notifyStop()
		return
	}
}

func (s *Server) GracefullyDown(notifyCtx context.Context) error {
	return s.httpSrv.Shutdown(notifyCtx)
}
