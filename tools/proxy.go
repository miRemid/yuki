package tools

import (
	"net/http"
	"net/url"
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

func ModifyResponse() func(*http.Response) error {
	return func(res *http.Response) error {
		// TODO: modify 404 request
		return nil
	}
}
