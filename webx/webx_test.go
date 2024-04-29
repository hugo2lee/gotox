package webx_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hugo2lee/gotox/webx"
	"github.com/pkg/errors"
)

func Test_Wrap(t *testing.T) {
	svr := gin.Default()
	svr.GET("/ping", webx.Wrap(func(ctx *gin.Context) (webx.Response, error) {
		return webx.Response{
			Code:    200,
			Message: "pong",
			Data:    nil,
		}, errors.New("test error")
	}))

	recorder := httptest.NewRecorder()
	svr.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/ping", nil))

	log.Printf("resp %v \n", recorder.Body.String())
}

func Test_WrapBody(t *testing.T) {
	svr := gin.Default()
	svr.POST("/ping", webx.WrapBind(func(ctx *gin.Context, req struct {
		Name string `json:"name"`
	},
	) (webx.Response, error) {
		return webx.Response{
			Code:    200,
			Message: "pong",
			Data:    nil,
		}, errors.New("post error")
	}))

	req := httptest.NewRequest(http.MethodPost, "/ping", bytes.NewBuffer([]byte("hello")))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	svr.ServeHTTP(recorder, req)

	log.Printf("resp %v \n", recorder.Body.String())
}

func Test_Response(t *testing.T) {
	interErr := webx.ErrMsg{
		Code:    500,
		Message: "svc error",
	}

	succRes := webx.ResponseSuccess("success")
	errRes := webx.ResponseErr(interErr)
	log.Printf("resp %v \n", succRes)
	log.Printf("resp %v \n", errRes)
}
