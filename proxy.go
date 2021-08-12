package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"github.com/miRemid/yuki/message"
	"github.com/miRemid/yuki/tools"
)

func (g *Gateway) reverseProxy(ctx *gin.Context) {
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
			ctx.Status(204)
			return
		} else {
			splits := strings.Split(msg, " ")
			cmd = splits[0]
			param = strings.Join(splits[1:], " ")
		}
		g.dprintf("[Cmd] %s, [Param] %s", cmd, param)
	case message.Notice:
		cmd = "notice"
		g.dprintf("receive notice message")
	case message.MetaEvent:
		cmd = "meta"
		g.dprintf("receive meta_event message")
	case message.Request:
		cmd = "request"
		g.dprintf("receive request message")
	default:
		ctx.Status(204)
		return
	}

	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(data.Bytes()))
	g.mu.RLock()
	node, err := g.selector.Peek(ctx.ClientIP())
	g.mu.RUnlock()
	if err != nil {
		g.dprintf("peek node error: %v", err)
		ctx.Status(204)
		return
	}
	targetURL := node.RemoteAddr
	// check rules
	if rule, ok := g.rules[cmd]; ok {
		g.dprintf("command %s using rules remote addr: %v", cmd, rule.RemoteAddr)
		targetURL = rule.RemoteAddr
	}
	target, _ := tools.CheckValidURL(targetURL)
	// if url is about "http://127.0.0.1/"
	if strings.HasSuffix(target.Path, "/") {
		target.Path = target.Path + cmd
	} else {
		target.Path = target.Path + "/" + cmd
	}

	g.dprintf("Reverse Proxy to RemoteAddr: %v", target)
	proxy := httputil.ReverseProxy{
		Director:       tools.Director(target),
		ModifyResponse: tools.ModifyResponse(),
	}
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
