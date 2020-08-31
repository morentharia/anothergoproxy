package main

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/elazarl/goproxy"
	"github.com/fatih/color"
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
		return nil, err
	}
	proxy.OnRequest().DoFunc(cacheHandlers.requestHandler)
	proxy.OnResponse().DoFunc(cacheHandlers.responseHandler)

	proxy.Verbose = options.Verbose
	return &Proxy{proxy}, nil
}

var reqBodyColor = color.New(color.FgMagenta).SprintFunc()
var urlColor = color.New(color.FgYellow).SprintFunc()
var respBodyColor = color.New(color.FgBlue).SprintFunc()

type CacheHandlers struct {
	cache          sync.Map
	sessionStorage SessionStorage
}

func NewCacheHandlers() (*CacheHandlers, error) {
	var err error
	if _, err = os.Stat(options.OutputPath); os.IsNotExist(err) {
		err = os.Mkdir(options.OutputPath, 0700)
		if err != nil {
			logrus.WithError(err).Error("mkdir")
			return nil, err
		}
	}
	if _, err := os.Stat(filepath.Join(options.OutputPath, "data")); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(options.OutputPath, "data"), 0700)
		if err != nil {
			logrus.WithError(err).Error("mkdir")
			return nil, err
		}
	}
	return &CacheHandlers{
		sessionStorage: SessionStorage{
			mux:     &sync.Mutex{},
			storage: make(map[int64]*RequestDTO),
			arr:     make([]int64, 0),
		},
	}, nil
}

type SessionStorage struct {
	mux     *sync.Mutex
	storage map[int64]*RequestDTO
	arr     []int64
}

func (s *SessionStorage) Store(key int64, value *RequestDTO) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.arr = append(s.arr, key)
	if len(s.arr) > 100 {
		for _, k := range s.arr[:50] {
			delete(s.storage, k)
		}
		s.arr = s.arr[50:]
	}
	s.storage[key] = value
}
func (s *SessionStorage) Load(key int64) (*RequestDTO, bool) {
	s.mux.Lock()
	value, ok := s.storage[key]
	s.mux.Unlock()
	return value, ok
}

func (c *CacheHandlers) requestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	reqDTO := NewRequestDTO(req)
	c.sessionStorage.Store(ctx.Session, reqDTO)

	f := NewCacheFile(reqDTO)
	if resp, err := f.Load(); err == nil {
		logrus.Printf("[%d] --> %s %s body: %s", ctx.Session, req.Method, urlColor(req.URL), f.RequestBodyFilename)
		return req, resp.HttpResponse()
	}

	logrus.Printf("[%d] --> %s %s", ctx.Session, req.Method, urlColor(req.URL))

	return req, nil
}

func (c *CacheHandlers) responseHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	logrus.Printf("[%d] <-- %d %s", ctx.Session, resp.StatusCode, urlColor(ctx.Req.URL))
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

	if err = NewCacheFile(reqDTO).Save(respDTO); err != nil {
		logrus.WithError(err).Error("save file")
		return nil
	}

	return resp
}
