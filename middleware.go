package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func (g *Gateway) checkSignature(ctx *gin.Context) {
	g.mu.RLock()
	secret := g.systemConfig.Secret
	g.mu.RUnlock()
	if secret != "" {
		sig := ctx.Request.Header.Get("X-Signature")
		if sig == "" {
			g.dprintf("X-Signature not found in request header")
			ctx.Status(204)
			return
		}
		sig = sig[len("sha1="):]
		mac := hmac.New(sha1.New, []byte(secret))
		byteData, _ := ioutil.ReadAll(ctx.Request.Body)
		io.WriteString(mac, string(byteData))
		res := fmt.Sprintf("%x", mac.Sum(nil))
		if res != sig {
			g.dprintf("Message not from go-cqhttp, reject")
			ctx.Status(204)
			return
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(byteData))
	}
	ctx.Next()
}
