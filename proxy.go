package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"github.com/miRemid/yuki/message"
)

func (g *Gateway) reverseProxy(ctx *gin.Context) {
	var (
		data    bytes.Buffer
		err     error
		msg     string
		cmd     string
		param   string
		subpath string
	)
	body := ctx.Request.Body
	io.Copy(&data, body)
	body.Close()
	post_type := gjson.GetBytes(data.Bytes(), "post_type").String()

	switch post_type {
	case message.Message:
		g.dprintf("receive normal message")
		raw_message := gjson.GetBytes(data.Bytes(), "raw_message").String()
		// check prefix
		if msg, err = g.checkPrefix(raw_message); err != nil {
			// reject
			g.dprintf("check prefix error: %v", err)
			ctx.Status(204)
			return
		} else {
			splits := strings.Split(msg, " ")
			cmd = splits[0]
			param = strings.Join(splits[1:], " ")
		}
		subpath = "/" + cmd
		g.dprintf("[Cmd] %s, [Param] %s", cmd, param)
	case message.Notice:
		subpath = "/notice"
		g.dprintf("receive notice message")
	case message.MetaEvent:
		subpath = "/meta"
		g.dprintf("receive meta_event message")
	case message.Request:
		subpath = "/request"
		g.dprintf("receive request message")
	default:
		ctx.Status(204)
		return
	}

	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(data.Bytes()))
	g.mu.RLock()
	node, err := g.selector.Peek(ctx.ClientIP())
	if err != nil {
		g.dprintf("peek node error: %v", err)
		ctx.Status(204)
		return
	}
	target := node.RemoteAddr
	if err != nil {
		g.dprintf("url parse error: %v", err)
		ctx.Status(204)
		return
	}
	g.dprintf("Reverse Proxy to RemoteAddr: %v", target)
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.URL.Path = subpath
		req.Host = target
	}
	modifyReponse := func(res *http.Response) error {
		// TODO: if command 404, quick reply
		return nil
	}
	proxy := httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyReponse,
	}
	g.mu.RUnlock()
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func (g *Gateway) checkPrefix(message string) (string, error) {
	g.mu.RLock()
	prefix := g.systemConfig.Prefix
	g.mu.RUnlock()

	for _, p := range prefix {
		if matched, err := regexp.MatchString(fmt.Sprintf("^%s", p), message); err != nil {
			return "", err
		} else if matched {
			return message[len(p):], nil
		}
	}
	return "", errors.New("not matched")
}
