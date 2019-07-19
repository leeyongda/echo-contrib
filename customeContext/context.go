package customeContext

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	Ctx context.Context
	// 使用标准库context

}

// 使用滴滴开源库json.
// 因为go-echo 不支持使用tags 编译注入 第三方json.
// 所以只能扩展上下文，去替代标准库里JSON的.
// go-gin 支持注入的.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// echo 框架 自定义JSON 解析，不依赖标标准库。
func (c *Context) JSON(code int, i interface{}) error {
	enc := json.NewEncoder(c.Response())
	c.Response().WriteHeader(code)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return enc.Encode(i)
}

func (c *Context) CheckTimeout(r chan interface{}) bool {
	select {
	case <-c.Ctx.Done():
		return true
	case rs := <-r:
		r <- rs
		// 接受的数据，塞回通道里面去
		// 记得关闭通道
		close(r)
		return false
	}

}

type (
	// HandlerFunc ...
	HandlerFunc func(*Context) error
)

func Handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*Context)
		return h(ctx)
	}
}
