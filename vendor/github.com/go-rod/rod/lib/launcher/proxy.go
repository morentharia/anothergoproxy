package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"

	"github.com/go-rod/rod/lib/cdp"
	"github.com/go-rod/rod/lib/utils"
)

// HeaderName for remote launch
const HeaderName = "Rod-Launcher"

// NewRemote create a Launcher instance from remote defaults. You must use it with launch.NewProxy or
// use the docker image mentioned from here: https://github.com/go-rod/rod/blob/master/lib/examples/remote-launch
func NewRemote(remoteURL string) *Launcher {
	u, err := url.Parse(remoteURL)
	utils.E(err)

	l := New()
	l.remote = true
	l.url = toWS(*u).String()
	l.Flags = nil

	res, err := http.Get(toHTTP(*u).String())
	utils.E(err)

	utils.E(json.Unmarshal(utils.MustReadBytes(res.Body), l))

	return l
}

// JSON serialization
func (l *Launcher) JSON() []byte {
	return utils.MustToJSONBytes(l)
}

// Client for launching browser remotely, such as browser from a docker container.
func (l *Launcher) Client() *cdp.Client {
	l.mustRemote()
	header := http.Header{}
	header.Add(HeaderName, utils.MustToJSON(l))
	return cdp.New(l.url).Header(header)
}

func (l *Launcher) mustRemote() {
	if !l.remote {
		panic("Must be used with launcher.NewRemote")
	}
}

// Proxy to help launch browser remotely.
// Any http request will return a default Launcher based on remote OS environment.
// Any websocket request will start a new browser and the request will be proxied to the browser.
// The websocket header "Rod-Launcher" holds the options to launch browser.
// If the websocket is closed, the browser will be killed.
type Proxy struct {
	Log       func(string)
	isWindows bool
}

var _ http.Handler = &Proxy{}

// NewProxy instance
func NewProxy() *Proxy {
	return &Proxy{
		Log:       func(s string) {},
		isWindows: runtime.GOOS == "windows",
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		p.launch(w, r)
		return
	}

	p.defaults(w, r)
}

func (p *Proxy) defaults(w http.ResponseWriter, _ *http.Request) {
	l := New()
	utils.E(w.Write(l.JSON()))
}

func (p *Proxy) launch(w http.ResponseWriter, r *http.Request) {
	l := New().Log(p.Log)

	options := r.Header.Get(HeaderName)
	if options != "" {
		l.Flags = nil
		utils.E(json.Unmarshal([]byte(options), l))
	}

	u := l.MustLaunch()
	defer func() {
		l.kill()

		l.Cleanup()
	}()

	parsedURL, err := url.Parse(u)
	utils.E(err)

	p.Log(fmt.Sprintln("Launch", u, l.FormatArgs()))
	defer func() { p.Log(fmt.Sprintln("Close", u)) }()

	parsedWS, err := url.Parse(u)
	utils.E(err)
	parsedURL.Path = parsedWS.Path

	httputil.NewSingleHostReverseProxy(toHTTP(*parsedURL)).ServeHTTP(w, r)
}
