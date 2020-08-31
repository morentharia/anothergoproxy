package main

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func disableCSPHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	// Allow everything https://stackoverflow.com/questions/35978863/allow-all-content-security-policy
	resp.Header.Set("Content-Security-Policy", "default-src *  data: blob: filesystem: about: ws: wss: 'unsafe-inline' 'unsafe-eval' ; script-src * data: blob: 'unsafe-inline' 'unsafe-eval'; connect-src * data: blob: 'unsafe-inline'; img-src * data: blob: 'unsafe-inline'; frame-src * data: blob: ; style-src * data: blob: 'unsafe-inline'; font-src * data: blob: 'unsafe-inline';")
	return resp
}
