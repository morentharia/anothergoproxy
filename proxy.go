package main

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	*goproxy.ProxyHttpServer
}

func NewProxy() (*Proxy, error) {
	proxy := goproxy.NewProxyHttpServer()
	if options.UpstreamProxyURL != "" {
		proxy.Tr = &http.Transport{Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(options.UpstreamProxyURL)
		}}
		proxy.ConnectDial = proxy.NewConnectDialToProxy(options.UpstreamProxyURL)
	}
	if options.URLMatch != "" {
		proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(options.URLMatch))).HandleConnect(goproxy.AlwaysMitm)
	}

	cacheHandlers, err := NewCacheHandlers()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	proxy.OnRequest().DoFunc(cacheHandlers.requestHandler)
	proxy.OnResponse().DoFunc(cacheHandlers.responseHandler)
	proxy.OnResponse().DoFunc(disableCSPHandler)

	proxy.Verbose = options.Verbose
	return &Proxy{proxy}, nil
}

var reqBodyColor = color.New(color.FgMagenta).SprintFunc()
var urlColor = color.New(color.FgYellow).SprintFunc()
var respBodyColor = color.New(color.FgBlue).SprintFunc()

type CacheHandlers struct {
	cache          ReqRespCacheI
	sessionStorage *SessionStorage
}

func NewCacheHandlers() (*CacheHandlers, error) {
	cache, err := NewCacheFile()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &CacheHandlers{
		cache:          cache,
		sessionStorage: NewSessionStorage(),
	}, nil
}

func (c *CacheHandlers) requestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	reqDTO := NewRequestDTO(req)
	c.sessionStorage.Store(ctx.Session, reqDTO)

	if resp, err := c.cache.Load(reqDTO); err == nil {
		logrus.Printf("[%d] --> %s %s", ctx.Session, req.Method, urlColor(req.URL))
		return req, resp.HttpResponse()
	}

	logrus.Printf("[%d] --> %s %s", ctx.Session, req.Method, urlColor(req.URL))

	return req, nil
}

func (c *CacheHandlers) responseHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	logrus.Printf("[%d] <-- %d %s", ctx.Session, resp.StatusCode, urlColor(ctx.Req.URL))
	//TODO: websocket support !!
	if resp.StatusCode == 101 {
		return resp
	}
	//TODO:
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return resp
	}
	location := resp.Header.Get("Location")
	if location != "" {
		logrus.Printf("Location: %s", location)
	}
	var reqDTO *RequestDTO
	var ok bool
	if reqDTO, ok = c.sessionStorage.Load(ctx.Session); !ok {
		return resp
	}

	respDTO, err := NewResponseDTO(resp)
	if err != nil {
		logrus.WithError(err).Error("NewResponseDTO")
		return resp
	}

	if err = c.cache.Store(reqDTO, respDTO); err != nil {
		logrus.WithError(err).Error("save file")
		return nil
	}

	return resp
}
