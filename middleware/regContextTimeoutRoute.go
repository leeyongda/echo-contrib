package middleware

import (
	"context"
	"github.com/leeyongda/echo-contrib/customeContext"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Config struct {
	TimeOut  time.Duration // 控制超时时间
	TimeUnit string        // 时间单位
}

// 注册路由超时，防止下游一直再等待。
func RegContextTimeOutRoute(c Config) echo.MiddlewareFunc {
	return regContextTimeOutRoute(c)
}

// 内部注册超时路由
func regContextTimeOutRoute(cfg Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.TimeOut > 0 {
				var timer time.Duration
				times := strings.ToLower(cfg.TimeUnit)
				switch times {
				case "s":
					timer = time.Second * cfg.TimeOut
				case "ms":
					timer = time.Microsecond * cfg.TimeOut
				default:
					timer = time.Second * cfg.TimeOut
				}
				ctx, cancel := context.WithTimeout(context.Background(), timer)
				defer cancel()
				cc := &customeContext.Context{c, ctx}
				return next(cc)
			}
			// 否则不做超时处理，超时默认时间就是http 设置的超时
			ctx := context.Background()
			cc := &customeContext.Context{c, ctx}
			return next(cc)

		}
	}
}
