package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

func (g *Gateway) checkSignature() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		g.mu.RLock()
		secret := g.systemConfig.Secret
		g.mu.RUnlock()
		return func(ctx echo.Context) error {
			if secret == "" {
				return next(ctx)
			}
			req := ctx.Request()
			sig := req.Header.Get("X-Signature")
			if sig == "" {
				return ctx.NoContent(204)
			}
			mac := hmac.New(sha1.New, []byte(secret))
			byteData, _ := ioutil.ReadAll(req.Body)
			io.WriteString(mac, string(byteData))
			res := fmt.Sprintf("%x", mac.Sum(nil))
			if res != sig {
				g.dprintf("Message not from go-cqhttp, reject")
				ctx.NoContent(204)
				return nil
			}
			req.Body = ioutil.NopCloser(bytes.NewReader(byteData))
			ctx.SetRequest(req)
			return next(ctx)
		}
	}
}
