package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"github.com/miRemid/yuki/message"
)

func (g *Gateway) ReverseProxy(ctx *gin.Context) {
	defer ctx.Status(204)

	var (
		data  bytes.Buffer
		err   error
		msg   string
		cmd   string
		param string
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
			return
		} else {
			splits := strings.Split(msg, " ")
			cmd = splits[0]
			param = strings.Join(splits[1:], " ")
		}
		g.dprintf("[Cmd] %s, [Param] %s", cmd, param)
	case message.Notice:
		g.dprintf("receive notice message")
	case message.MetaEvent:
		g.dprintf("receive meta_event message")
		return
	case message.Request:
		g.dprintf("receive request message")
	default:
		return
	}
	// g.mu.RLock()
	// target := g.nodes[cmd]
	// director := func(req *http.Request) {
	// 	req.URL.Scheme = "http"
	// 	req.URL.Host = target.RemoteAddr
	// 	req.Host = target.RemoteAddr
	// }
	// g.mu.RUnlock()
	// proxy := &httputil.ReverseProxy{Director: director}
	// proxy.ServeHTTP(ctx.Writer, ctx.Request)
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
