package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/leeyongda/echo-contrib/customeContext"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRegTimeOutRoute(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	hc := func(c *customeContext.Context) error {
		var result = make(chan interface{}, 1)
		// 模拟业务超时请求，延迟5秒
		go func() {
			time.Sleep(time.Second * 5)
			result <- "123"
			close(result)
		}()
		//if c.CheckTimeout(result) {
		//	return &echo.HTTPError{
		//		Code:     http.StatusServiceUnavailable,
		//		Message:  "request TimeOut",
		//		Internal: c.Ctx.Err(),
		//	}
		//}
		//return c.JSON(http.StatusOK, <-result)
		select {
		case <-c.Ctx.Done():
			return &echo.HTTPError{
				Code:     http.StatusServiceUnavailable,
				Message:  "request TimeOut",
				Internal: c.Ctx.Err(),
			}
		case rs := <-result:
			return c.JSON(http.StatusOK, rs)
		}

	}
	cfg := Config{
		TimeOut:  6, // 超时时间，时间单位秒
		TimeUnit: "s",
	}
	// 允许超时
	h := RegContextTimeOutRoute(cfg)(customeContext.Handler(hc))
	asserts := assert.New(t)
	asserts.NoError(h(c))

	// 不允许超时
	cfg.TimeOut = 4
	h = RegContextTimeOutRoute(cfg)(customeContext.Handler(hc))
	he := h(c).(*echo.HTTPError)
	asserts.Equal(http.StatusServiceUnavailable, he.Code)
}
