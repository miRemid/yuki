package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Director(target *url.URL) func(*http.Request) {
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
}

func ModifyResponse(format, command string) func(*http.Response) error {
	return func(res *http.Response) error {
		if command == "meta" || command == "request" || command == "notice" {
			return nil
		}
		// TODO: modify 404 request
		if res.StatusCode == 404 {
			log.Println("404")
			ioutil.ReadAll(res.Body)
			response := gin.H{
				"reply": fmt.Sprintf(format, command),
			}
			data, _ := json.Marshal(&response)
			res.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			res.Header.Set("Content-Length", fmt.Sprint(len(data)))
		}
		return nil
	}
}
